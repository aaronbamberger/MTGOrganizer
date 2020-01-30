package mtgcards

import "context"
import "database/sql"
import "encoding/hex"
import "fmt"
import "hash"
import "log"
import "sync"

func checkRowsAffected(res sql.Result, expectedAffected int64, errString string) error {
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected != expectedAffected {
		return fmt.Errorf("Query %s affected an unexpected number of rows, expected %d, got %d\n",
			errString, expectedAffected, rowsAffected)
	}

	return nil
}

func HashToHexString(hashVal hash.Hash) string {
	hashBytes := make([]byte, 0, hashVal.Size())
	hashBytes = hashVal.Sum(hashBytes)
	return hex.EncodeToString(hashBytes)
}

func ImportSetsToDb(db *sql.DB, sets map[string]MTGSet) (*DbUpdateStats, error) {
	var setImportWg sync.WaitGroup
	var stats DbUpdateStats

	// We defer the cleanup before calling the setup function, because the setup
	// function might get partway through the initialization, and then error out,
	// leaving some things to be cleaned up.  The cleanup function will only
	// cleanup things that have been actually set up, so it's safe to call it
	// regardless of where the setup might fail
	defer cleanupOptionTables()
	err := prepareOptionTables(db)
	if err != nil {
		return nil, err
	}

	dbQueries, err  := prepareDbQueries(db)
	if err != nil {
		dbQueries.cleanupDbQueries()
		return nil, err
	}
	defer dbQueries.cleanupDbQueries()

    // Create the thread that serializes all access to the atomic_card_data
    // db table
    atomicPropertiesRequest := make(chan atomicPropRequest)
    resChannel := make(chan error)
    quitChannel := make(chan interface{})
    go atomicPropDbThread(db, atomicPropertiesRequest, resChannel, quitChannel)

    err = <-resChannel
    if err != nil {
        quitChannel <- nil
        <-resChannel
        return nil, err
    }

	for _, set := range sets {
		setImportWg.Add(1)
		go maybeInsertSetToDb(db, dbQueries, &stats, &setImportWg, atomicPropertiesRequest, set)
	}

	setImportWg.Wait()
    quitChannel <- nil
    <-resChannel
	return &stats, nil
}

func maybeInsertSetToDb(
        db *sql.DB,
        queries *dbQueries,
        stats *DbUpdateStats,
		wg *sync.WaitGroup,
        atomicPropertiesRequest chan atomicPropRequest,
        set MTGSet) {
	defer wg.Done()
	ctx := context.Background()

	// Open a DB connection
	dbConn, err := db.Conn(ctx)
	if err != nil {
		log.Print(err)
		return
	}
	defer dbConn.Close()

	// Transaction for inserting the set itself
	setTx, err := dbConn.BeginTx(ctx, nil)
	if err != nil {
		log.Print(err)
		return
	}

	setQueries := queries.queriesForTx(setTx)

	// Hash the set for later use
	set.Canonicalize()
	setHash := HashToHexString(set.Hash())

	// First, check to see if this set is in the DB at all
	setExists, setDbHash, setId, err := set.CheckIfSetExists(setQueries)
	if err != nil {
		log.Print(err)
		setTx.Rollback()
		return
	}

	stats.AddToTotalSets(1)

    totalCards := 0
	totalNewCards := 0
	totalNewCardsInNewSets := 0
	totalNewCardsInExistingSets := 0
	totalNewAtomicRecordsForNewCards := 0
    totalExistingAtomicRecordsForNewCards := 0
	totalExistingCards := 0
	totalExistingCardsHashSkipped := 0
	totalExistingCardsUpdated := 0
    totalNewAtomicRecordsForExistingCards := 0
    totalExistingAtomicRecordsForExistingCards := 0

	if setExists {
		log.Printf("Set %s already exists in the database\n", set.Code)
		// This set already exists in the db
		// Check to see if the hash matcdbhes what's already in the db
		if setDbHash == setHash {
			// Hashes match, so we can skip updating this set in the db
			log.Printf("Set %s in db matches hash %s, skipping update...\n", set.Code, setDbHash)
			setTx.Commit()
			stats.AddToExistingSetsSkipped(1)
		} else {
			// Hashes don't match, so we need to first update the set itself, and then
			// look at each card in the set to see if it needs to be updated
			log.Printf("Set %s hashes don't match (db: %s, json: %s), updating set...\n",
				set.Code, setDbHash, setHash)
			err := set.UpdateSetInDb(setQueries, setHash, setId)
			if err != nil {
				log.Print(err)
				setTx.Rollback()
				return
			}
			setTx.Commit()
            stats.AddToExistingSetsUpdated(1)

			// For each card, check if the card exists, and if so, if the hash
			// matches
			for idx := range set.Cards {
                // Need to access by index here to get a pointer to the card,
                // not a copy
                card := &set.Cards[idx]

				// Transaction for each card
				cardTx, err := dbConn.BeginTx(ctx, nil)
				if err != nil {
					log.Print(err)
					continue
				}

				cardQueries := queries.queriesForTx(cardTx)

				cardExists, cardDbHash, atomicCardDataId, err := card.CheckIfCardExists(cardQueries)
				if err != nil {
					log.Print(err)
					cardTx.Rollback()
					continue
				}

				if !cardExists {
					newAtomicPropertiesAdded, err := card.InsertAllCardDataToDb(
                        cardQueries,
                        setId,
                        atomicPropertiesRequest)
					if err != nil {
						log.Print(err)
						cardTx.Rollback()
						continue
					}
                    cardTx.Commit()
					if newAtomicPropertiesAdded {
						totalNewAtomicRecordsForNewCards += 1
					} else {
                        totalExistingAtomicRecordsForNewCards += 1
                    }
                    totalCards += 1
					totalNewCards += 1
					totalNewCardsInExistingSets += 1
				} else {
					// Check if the stored hash matches
					cardHash := HashToHexString(card.Hash())
					if cardHash == cardDbHash {
						// Can skip
						log.Printf("Card %s hash matches in db (%s), skipping", card.Name, cardHash)
                        cardTx.Commit()
                        totalCards += 1
                        totalExistingCards += 1
						totalExistingCardsHashSkipped += 1
					} else {
						// Need to update card
						log.Printf("Card %s hash doesn't match (db: %s, card: %s), updating",
							card.Name, cardDbHash, cardHash)
						newAtomicPropetiesAdded, err := card.UpdateCardDataInDb(
                            cardQueries,
                            atomicCardDataId,
                            setId,
                            atomicPropertiesRequest)
						if err != nil {
							log.Print(err)
							cardTx.Rollback()
							continue
						}
                        cardTx.Commit()
                        if newAtomicPropetiesAdded {
                            totalNewAtomicRecordsForExistingCards += 1
                        } else {
                            totalExistingAtomicRecordsForExistingCards += 1
                        }
                        totalCards += 1
                        totalExistingCards += 1
                        totalExistingCardsUpdated += 1
					}
				}
			}
		}
		stats.AddToTotalExistingSets(1)
	} else {
		// This set does not already exist in the db
		setId, err := set.InsertSetToDb(setQueries, setHash)
		if err != nil {
			log.Print(err)
			setTx.Rollback()
			return
		}

		// Insert the set translations
		for lang, name := range set.Translations {
			err := InsertSetTranslationToDb(setQueries, setId, lang, name)
			if err != nil {
				log.Print(err)
				setTx.Rollback()
				return
			}
		}

		setTx.Commit()
		stats.AddToTotalNewSets(1)

		// Insert all of the cards in the set.  No need to check the full card hash, since we're bulk
		// inserting the entire set
		log.Printf("Processing cards in set %s\n", set.Code)
		for idx := range set.Cards {
            // Need to access by index here to get a pointer to the card,
            // not a copy
            card := &set.Cards[idx]

			// Transaction for each card
			cardTx, err := dbConn.BeginTx(ctx, nil)
			if err != nil {
				log.Print(err)
				continue
			}

			cardQueries := queries.queriesForTx(cardTx)

			newAtomicPropertiesAdded, err := card.InsertAllCardDataToDb(
                cardQueries,
                setId,
                atomicPropertiesRequest)
			if err != nil {
				log.Print(err)
				cardTx.Rollback()
				continue
			}
			cardTx.Commit()

            totalCards += 1
            totalNewCards += 1
			totalNewCardsInNewSets += 1
			if newAtomicPropertiesAdded {
				totalNewAtomicRecordsForNewCards += 1
			} else {
                totalExistingAtomicRecordsForNewCards += 1
            }
		}
	}

	stats.AddToTotalCards(totalCards)
	stats.AddToTotalNewCards(totalNewCards)
	stats.AddToTotalNewCardsInNewSets(totalNewCardsInNewSets)
	stats.AddToTotalNewCardsInExistingSets(totalNewCardsInExistingSets)
	stats.AddToTotalNewAtomicRecordsForNewCards(totalNewAtomicRecordsForNewCards)
    stats.AddToTotalNewAtomicRecordsForExistingCards(totalNewAtomicRecordsForExistingCards)
    stats.AddToTotalExistingAtomicRecordsForNewCards(totalExistingAtomicRecordsForNewCards)
    stats.AddToTotalExistingAtomicRecordsForExistingCards(totalExistingAtomicRecordsForExistingCards)
	stats.AddToTotalExistingCards(totalExistingCards)
	stats.AddToExistingCardsSkipped(totalExistingCardsHashSkipped)
	stats.AddToExistingCardsUpdated(totalExistingCardsUpdated)
	log.Printf("Done processing set %s\n", set.Code)
}

func (card *MTGCard) UpdateCardDataInDb(
    queries *dbQueries,
	atomicPropertiesId int64,
    setId int64,
    atomicPropertiesRequest chan atomicPropRequest) (bool, error) {
	// First, check to see if the atomic properties hash still matches.  If it does,
	// we just need to update the rest of the card data, and can leave it pointing
	// to the same atomic properties record.
    atomicPropHash := HashToHexString(card.AtomicPropertiesHash())
    responseChan := make(chan atomicPropResponse)

    request := atomicPropRequest{
        RequestType: AtomicPropRequestGetHash,
        ResponseChan: responseChan,
        AtomicPropertiesId: atomicPropertiesId}
    atomicPropertiesRequest <- request
    response := <-responseChan

    if response.Error != nil {
        return false, response.Error
    }

    if atomicPropHash != response.AtomicPropertiesHash {
        // The atomic properties don't match, so release our reference to the old
        // atomic properties and grab a reference to the new atomic properties
        request = atomicPropRequest{
            RequestType: AtomicPropRequestRemoveRef,
            ResponseChan: responseChan,
            AtomicPropertiesId: atomicPropertiesId}
        atomicPropertiesRequest <- request
        response = <-responseChan

        if response.Error != nil {
            return false, response.Error
        }

        request = atomicPropRequest{
            RequestType: AtomicPropRequestGetRef,
            AtomicPropertiesHash: atomicPropHash,
            ColorIdentity: card.ColorIdentity,
            ColorIndicator: card.ColorIndicator,
            Colors: card.Colors,
            ConvertedManaCost: card.ConvertedManaCost,
            EDHRecRank: card.EDHRecRank,
            FaceConvertedManaCost: card.FaceConvertedManaCost,
            Hand: card.Hand,
            IsReserved: card.IsReserved,
            Layout: card.Layout,
            Life: card.Life,
            Loyalty: card.Loyalty,
            ManaCost: card.ManaCost,
            MTGStocksId: card.MTGStocksId,
            Name: card.Name,
            Power: card.Power,
            ScryfallOracleId: card.ScryfallOracleId,
            Side: card.Side,
            Text: card.Text,
            Toughness: card.Toughness,
            Type: card.Type,
            ResponseChan: responseChan}
        atomicPropertiesRequest <- request
        response = <-responseChan

        if response.Error != nil {
            return false, response.Error
        }

        atomicPropertiesId = response.AtomicPropertiesId
    }

	// Now, update the card record, clear any entries from auxilliary tables belonging
	// to the old card record, and insert new auxilliary entries for the updated card record
    err := card.UpdateCardInDb(queries, atomicPropertiesId, setId)
	if err != nil {
		return false, err
	}

	err = card.DeleteOtherTableCardData(queries)
	if err != nil {
		return false, err
	}

	err = card.InsertOtherTableCardData(queries)
	if err != nil {
		return false, err
	}

	return response.NewRecordAdded, nil
}

func (card *MTGCard) InsertAllCardDataToDb(
        queries *dbQueries,
        setId int64,
        atomicPropertiesRequest chan atomicPropRequest) (bool, error) {

    // Get a reference to an atomic card properties record for this card
	atomicPropHash := HashToHexString(card.AtomicPropertiesHash())
    responseChan := make(chan atomicPropResponse)

    request := atomicPropRequest{
        RequestType: AtomicPropRequestGetRef,
        AtomicPropertiesHash: atomicPropHash,
        ColorIdentity: card.ColorIdentity,
        ColorIndicator: card.ColorIndicator,
        Colors: card.Colors,
        ConvertedManaCost: card.ConvertedManaCost,
        EDHRecRank: card.EDHRecRank,
        FaceConvertedManaCost: card.FaceConvertedManaCost,
        Hand: card.Hand,
        IsReserved: card.IsReserved,
        Layout: card.Layout,
        Life: card.Life,
        Loyalty: card.Loyalty,
        ManaCost: card.ManaCost,
        MTGStocksId: card.MTGStocksId,
        Name: card.Name,
        Power: card.Power,
        ScryfallOracleId: card.ScryfallOracleId,
        Side: card.Side,
        Text: card.Text,
        Toughness: card.Toughness,
        Type: card.Type,
        ResponseChan: responseChan}

    atomicPropertiesRequest <- request
    response := <-responseChan
	if response.Error != nil {
		return false, response.Error
	}

    err := card.InsertRemainingAtomicPropertiesToDb(queries, response.AtomicPropertiesId)
    if err != nil {
        return false, err
    }

	// Insert the main card record in the all_cards table
	err = card.InsertCardToDb(queries, response.AtomicPropertiesId, setId)
	if err != nil {
		return false, err
	}

	// Insert the rest of the card data
	err = card.InsertOtherTableCardData(queries)
	if err != nil {
		return false, nil
	}

	return response.NewRecordAdded, nil
}

func (card *MTGCard) InsertOtherTableCardData(queries *dbQueries) error {
	// Insert the card data that doesn't live in the all_cards table

	// Frame effects
	for _, frameEffect := range card.FrameEffects {
		err := card.InsertFrameEffectToDb(queries, frameEffect)
		if err != nil {
			return err
		}
	}

	// Other face IDs
	for _, otherFaceId := range card.OtherFaceIds {
		err := card.InsertOtherFaceIdToDb(queries, otherFaceId)
		if err != nil {
			return err
		}
	}

	// Variations
	for _, variation := range card.Variations {
		err := card.InsertVariationToDb(queries, variation)
		if err != nil {
			return err
		}
	}

	return nil
}

func (card *MTGCard) DeleteOtherTableCardData(queries *dbQueries) error {
	_, err := queries.DeleteFrameEffectQuery.Exec(card.UUID)
	if err != nil {
		return err
	}

	_, err = queries.DeleteOtherFaceQuery.Exec(card.UUID)
	if err != nil {
		return err
	}

	_, err = queries.DeleteVariationQuery.Exec(card.UUID)
	if err != nil {
		return err
	}

	return nil
}

func (card *MTGCard) CheckIfCardExists(queries *dbQueries) (bool, string, int64, error) {
	resultRow := queries.CardHashQuery.QueryRow(card.UUID)

	var cardHash string
	var atomicCardDataId int64
	err := resultRow.Scan(&cardHash, &atomicCardDataId)
	if err != nil {
		if err == sql.ErrNoRows {
			// This card isn't in the database
			return false, "", 0, nil
		} else {
			return false, "", 0, err
		}
	} else {
		return true, cardHash, atomicCardDataId, nil
	}
}

func (set *MTGSet) CheckIfSetExists(queries *dbQueries) (bool, string, int64, error) {
	resultRow := queries.SetHashQuery.QueryRow(set.Code)

	var setHash string
	var setId int64
	err := resultRow.Scan(&setHash, &setId)
	if err != nil {
		if err == sql.ErrNoRows {
			// This set isn't in the database
			return false, "", 0, nil
		} else {
			return false, "", 0, err
		}
	} else {
		return true, setHash, setId, nil
	}
}

func (set *MTGSet) InsertSetToDb(queries *dbQueries, setHash string) (int64, error) {
	res, err := queries.InsertSetQuery.Exec(setHash, set.BaseSetSize, set.Block,
		set.Code, set.IsForeignOnly, set.IsFoilOnly, set.IsOnlineOnly,
		set.IsPartialPreview, set.KeyruneCode, set.MCMName, set.MCMId, set.MTGOCode,
		set.Name, set.ParentCode, set.ReleaseDate, set.TCGPlayerGroupId,
		set.TotalSetSize, set.Type)
	if err != nil {
		return 0, err
	}

	err = checkRowsAffected(res, 1, "insert set")
	if err != nil {
		return 0, err
	}

	setId, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return setId, nil
}

func (set *MTGSet) UpdateSetInDb(queries *dbQueries, setHash string, setId int64) error {
	res, err := queries.UpdateSetQuery.Exec(setHash, set.BaseSetSize, set.Block,
		set.Code, set.IsForeignOnly, set.IsFoilOnly, set.IsOnlineOnly,
		set.IsPartialPreview, set.KeyruneCode, set.MCMName, set.MCMId, set.MTGOCode,
		set.Name, set.ParentCode, set.ReleaseDate, set.TCGPlayerGroupId,
		set.TotalSetSize, set.Type, setId)
	if err != nil {
		return err
	}

	return checkRowsAffected(res, 1, "update set")
}

func InsertSetTranslationToDb(queries *dbQueries, setId int64, translationLang string,
		translatedName string) error {
	languageId, err := getSetTranslationLanguageId(translationLang)
	if err != nil {
		return err
	}

	res, err := queries.InsertSetTranslationQuery.Exec(setId, languageId, translatedName)
	if err != nil {
		return err
	}

	return checkRowsAffected(res, 1, "insert set name translation")
}

func (card *MTGCard) InsertFrameEffectToDb(queries *dbQueries, frameEffect string) error {
	frameEffectId, err := getFrameEffectId(frameEffect)
	if err != nil {
		return err
	}

	res, err := queries.InsertFrameEffectQuery.Exec(card.UUID, frameEffectId)
	if err != nil {
		return err
	}

	return checkRowsAffected(res, 1, "insert frame effect")
}

func (card *MTGCard) InsertOtherFaceIdToDb(queries *dbQueries, otherFaceUUID string) error {
	res, err := queries.InsertOtherFaceIdQuery.Exec(card.UUID, otherFaceUUID)
	if err != nil {
		return err
	}

	return checkRowsAffected(res, 1, "insert other face ID")
}

func (card *MTGCard) InsertVariationToDb(queries *dbQueries, variationUUID string) error {
	res, err := queries.InsertVariationQuery.Exec(card.UUID, variationUUID)
	if err != nil {
		return err
	}

	return checkRowsAffected(res, 1, "insert variation")
}

func (card *MTGCard) InsertCardToDb(queries *dbQueries, atomicPropertiesId int64,
		setId int64) error {
	var duelDeck sql.NullString
	var flavorText sql.NullString
	var mtgArenaId sql.NullInt32
	var mtgoFoilId sql.NullInt32
	var mtgoId sql.NullInt32
	var scryfallIllustrationId sql.NullString

	if len(card.DuelDeck) > 0 {
		duelDeck.String = card.DuelDeck
		duelDeck.Valid = true
	}

	if len(card.FlavorText) > 0 {
		flavorText.String = card.FlavorText
		flavorText.Valid = true
	}

	if card.MTGArenaId > 0 {
		mtgArenaId.Int32 = int32(card.MTGArenaId)
		mtgArenaId.Valid = true
	}

	if card.MTGOFoilId > 0 {
		mtgoFoilId.Int32 = int32(card.MTGOFoilId)
		mtgoFoilId.Valid = true
	}

	if card.MTGOId > 0 {
		mtgoId.Int32 = int32(card.MTGOId)
		mtgoId.Valid = true
	}

	if len(card.ScryfallIllustrationId) > 0 {
		scryfallIllustrationId.String = card.ScryfallIllustrationId
		scryfallIllustrationId.Valid = true
	}

	cardHash := HashToHexString(card.Hash())

	res, err := queries.InsertCardQuery.Exec(card.UUID, cardHash, atomicPropertiesId,
		setId, card.Artist, card.BorderColor, card.Number, card.ScryfallId,
		card.Watermark, card.FrameVersion, card.MCMId, card.MCMMetaId,
		card.MultiverseId, card.OriginalText, card.OriginalType,
		card.Rarity, card.TCGPlayerProductId, duelDeck, flavorText,
		card.HasFoil, card.HasNonFoil, card.IsAlternative, card.IsArena,
		card.IsFullArt, card.IsMTGO, card.IsOnlineOnly, card.IsOversized,
		card.IsPaper, card.IsPromo, card.IsReprint, card.IsStarter,
		card.IsStorySpotlight, card.IsTextless, card.IsTimeshifted,
		mtgArenaId, mtgoFoilId, mtgoId, scryfallIllustrationId)

	if err != nil {
		return err
	}

	return checkRowsAffected(res, 1, "insert card data")
}

func (card *MTGCard) UpdateCardInDb(queries *dbQueries, atomicPropertiesId int64,
		setId int64) error {
	var duelDeck sql.NullString
	var flavorText sql.NullString
	var mtgArenaId sql.NullInt32
	var mtgoFoilId sql.NullInt32
	var mtgoId sql.NullInt32
	var scryfallIllustrationId sql.NullString

	if len(card.DuelDeck) > 0 {
		duelDeck.String = card.DuelDeck
		duelDeck.Valid = true
	}

	if len(card.FlavorText) > 0 {
		flavorText.String = card.FlavorText
		flavorText.Valid = true
	}

	if card.MTGArenaId > 0 {
		mtgArenaId.Int32 = int32(card.MTGArenaId)
		mtgArenaId.Valid = true
	}

	if card.MTGOFoilId > 0 {
		mtgoFoilId.Int32 = int32(card.MTGOFoilId)
		mtgoFoilId.Valid = true
	}

	if card.MTGOId > 0 {
		mtgoId.Int32 = int32(card.MTGOId)
		mtgoId.Valid = true
	}

	if len(card.ScryfallIllustrationId) > 0 {
		scryfallIllustrationId.String = card.ScryfallIllustrationId
		scryfallIllustrationId.Valid = true
	}

	cardHash := HashToHexString(card.Hash())

	res, err := queries.UpdateCardQuery.Exec(cardHash, atomicPropertiesId,
		setId, card.Artist, card.BorderColor, card.Number, card.ScryfallId,
		card.Watermark, card.FrameVersion, card.MCMId, card.MCMMetaId,
		card.MultiverseId, card.OriginalText, card.OriginalType,
		card.Rarity, card.TCGPlayerProductId, duelDeck, flavorText,
		card.HasFoil, card.HasNonFoil, card.IsAlternative, card.IsArena,
		card.IsFullArt, card.IsMTGO, card.IsOnlineOnly, card.IsOversized,
		card.IsPaper, card.IsPromo, card.IsReprint, card.IsStarter,
		card.IsStorySpotlight, card.IsTextless, card.IsTimeshifted,
		mtgArenaId, mtgoFoilId, mtgoId, scryfallIllustrationId, card.UUID)

	if err != nil {
		return err
	}

	return checkRowsAffected(res, 1, "update card data")
}

