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

type DbUpdateStats struct {
	Mutex sync.Mutex
	TotalSetsInUpdate int
	TotalExistingSetsInDb int
	TotalNewSets int
	ExistingSetsSkippedDueToHashMatch int
	ExistingSetsUpdatedDueToHashMismatch int
	TotalCardsInUpdate int
	TotalNewCards int
	TotalNewAtomicCards int
	TotalNewCardsInNewSets int
	TotalNewCardsInExistingSets int
	ExistingCardsSkippedDueToHashMatch int
	ExistingCardsUpdatedDueToHashMismatch int
}

func ImportSetsToDb(db *sql.DB, sets map[string]MTGSet) (*DbUpdateStats, error) {
	var setImportWg sync.WaitGroup
	var updateStats DbUpdateStats

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

	for _, set := range sets {
		setImportWg.Add(1)
		go maybeInsertSetToDb(db, &updateStats, &setImportWg, set)
	}

	setImportWg.Wait()
	return &updateStats, nil
}

type cardDbQueries struct {
	NumAtomicPropertiesQuery *sql.Stmt
	AtomicPropertiesIdQuery *sql.Stmt
	AtomicPropertiesHashQuery *sql.Stmt
	InsertAtomicPropertiesQuery *sql.Stmt
	InsertCardQuery *sql.Stmt
	UpdateCardQuery *sql.Stmt
	InsertAltLangDataQuery *sql.Stmt
	InsertCardPrintingQuery *sql.Stmt
	InsertCardSubtypeQuery *sql.Stmt
	InsertCardSupertypeQuery *sql.Stmt
	InsertFrameEffectQuery *sql.Stmt
	InsertLeadershipSkillQuery *sql.Stmt
	InsertLegalityQuery *sql.Stmt
	InsertOtherFaceIdQuery *sql.Stmt
	InsertPurchaseURLQuery *sql.Stmt
	InsertRulingQuery *sql.Stmt
	InsertVariationQuery *sql.Stmt
	InsertBaseTypeQuery *sql.Stmt
	DeleteFrameEffectQuery *sql.Stmt
	DeleteOtherFaceQuery *sql.Stmt
	DeleteVariationQuery *sql.Stmt
	UpdateRefCntQuery *sql.Stmt
}

func (queries *cardDbQueries) prepareCardDbQueries(tx *sql.Tx) error {
	var err error

	queries.NumAtomicPropertiesQuery, err = tx.Prepare(`SELECT COUNT(scryfall_oracle_id)
		FROM atomic_card_data
		WHERE card_data_hash = ?`)
	if err != nil {
		return err
	}

	queries.AtomicPropertiesIdQuery, err = tx.Prepare(`SELECT atomic_card_data_id,
		ref_cnt,
		scryfall_oracle_id
		FROM atomic_card_data
		WHERE card_data_hash = ?`)
	if err != nil {
		return err
	}

	queries.AtomicPropertiesHashQuery, err = tx.Prepare(`SELECT card_data_hash
		FROM atomic_card_data
		WHERE atomic_card_data_id = ?`)
	if err != nil {
		return err
	}

	queries.InsertAtomicPropertiesQuery, err = tx.Prepare(`INSERT INTO atomic_card_data
		(card_data_hash, color_identity, color_indicator, colors, converted_mana_cost,
		edhrec_rank, face_converted_mana_cost, hand, is_reserved, layout, life,
		loyalty, mana_cost, mtgstocks_id, name, card_power, scryfall_oracle_id,
		side, text, toughness, card_type)
		VALUES
		(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return err
	}

	queries.InsertCardQuery, err = tx.Prepare(`INSERT INTO all_cards
		(uuid, full_card_hash, atomic_card_data_id, set_id, artist, border_color,
		card_number, scryfall_id, watermark, frame_version, mcm_id, mcm_meta_id,
		multiverse_id, original_text, original_type, rarity, tcgplayer_product_id,
		duel_deck, flavor_text, has_foil, has_non_foil, is_alternative, is_arena,
		is_full_art, is_mtgo, is_online_only, is_oversized, is_paper, is_promo,
		is_reprint, is_starter, is_story_spotlight, is_textless, is_timeshifted,
		mtg_arena_id, mtgo_foil_id, mtgo_id, scryfall_illustration_id)
		VALUES
		(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?,
		?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return err
	}

	queries.UpdateCardQuery, err = tx.Prepare(`UPDATE all_cards SET
		full_card_hash = ?,
		atomic_card_data_id = ?,
		set_id = ?,
		artist = ?,
		border_color = ?,
		card_number = ?,
		scryfall_id = ?,
		watermark = ?,
		frame_version = ?,
		mcm_id = ?,
		mcm_meta_id = ?,
		multiverse_id = ?,
		original_text = ?,
		original_type = ?,
		rarity = ?,
		tcgplayer_product_id = ?,
		duel_deck = ?,
		flavor_text = ?,
		has_foil = ?,
		has_non_foil = ?,
		is_alternative = ?,
		is_arena = ?,
		is_full_art = ?,
		is_mtgo = ?,
		is_online_only = ?,
		is_oversized = ?,
		is_paper = ?,
		is_promo = ?,
		is_reprint = ?,
		is_starter = ?,
		is_story_spotlight = ?,
		is_textless = ?,
		is_timeshifted = ?,
		mtg_arena_id = ?,
		mtgo_foil_id = ?,
		mtgo_id = ?,
		scryfall_illustration_id = ?
		WHERE uuid = ?`)

	queries.InsertAltLangDataQuery, err = tx.Prepare(`INSERT INTO alternate_language_data
		(atomic_card_data_id, flavor_text, language, multiverse_id, name, text, card_type)
		VALUES
		(?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return err
	}

	queries.InsertCardPrintingQuery, err = tx.Prepare(`INSERT INTO card_printings
		(atomic_card_data_id, set_code)
		VALUES
		(?, ?)`)
	if err != nil {
		return err
	}

	queries.InsertCardSubtypeQuery, err = tx.Prepare(`INSERT INTO card_subtypes
		(atomic_card_data_id, subtype_option_id)
		VALUES
		(?, ?)`)
	if err != nil {
		return err
	}

	queries.InsertCardSupertypeQuery, err = tx.Prepare(`INSERT INTO card_supertypes
		(atomic_card_data_id, supertype_option_id)
		VALUES
		(?, ?)`)
	if err != nil {
		return err
	}

	queries.InsertFrameEffectQuery, err = tx.Prepare(`INSERT INTO frame_effects
		(card_uuid, frame_effect_option_id)
		VALUES
		(?, ?)`)
	if err != nil {
		return err
	}

	queries.InsertLeadershipSkillQuery, err = tx.Prepare(`INSERT INTO leadership_skills
		(atomic_card_data_id, leadership_format_id, leader_legal)
		VALUES
		(?, ?, ?)`)
	if err != nil {
		return err
	}

	queries.InsertLegalityQuery, err = tx.Prepare(`INSERT INTO legalities
		(atomic_card_data_id, game_format_id, legality_option_id)
		VALUES
		(?, ?, ?)`)
	if err != nil {
		return err
	}

	queries.InsertOtherFaceIdQuery, err = tx.Prepare(`INSERT INTO other_faces
		(card_uuid, other_face_uuid)
		VALUES
		(?, ?)`)
	if err != nil {
		return err
	}

	queries.InsertPurchaseURLQuery, err = tx.Prepare(`INSERT INTO purchase_urls
		(atomic_card_data_id, purchase_site_id, purchase_url)
		VALUES
		(?, ?, ?)`)
	if err != nil {
		return err
	}

	queries.InsertRulingQuery, err = tx.Prepare(`INSERT INTO rulings
		(atomic_card_data_id, ruling_date, ruling_text)
		VALUES
		(?, ?, ?)`)
	if err != nil {
		return err
	}

	queries.InsertVariationQuery, err = tx.Prepare(`INSERT INTO variations
		(card_uuid, variation_uuid)
		VALUES
		(?, ?)`)
	if err != nil {
		return err
	}

	queries.InsertBaseTypeQuery, err = tx.Prepare(`INSERT INTO base_types
		(atomic_card_data_id, base_type_option_id)
		VALUES
		(?, ?)`)
	if err != nil {
		return err
	}

	queries.DeleteFrameEffectQuery, err = tx.Prepare(`DELETE FROM frame_effects
		WHERE card_uuid = ?`)
	if err != nil {
		return err
	}

	queries.DeleteOtherFaceQuery, err = tx.Prepare(`DELETE FROM other_faces
		WHERE card_uuid = ?`)
	if err != nil {
		return err
	}

	queries.DeleteVariationQuery, err = tx.Prepare(`DELETE FROM variations
		WHERE card_uuid = ?`)
	if err != nil {
		return err
	}

	queries.UpdateRefCntQuery, err = tx.Prepare(`UPDATE atomic_card_data
		SET ref_cnt = ?
		WHERE atomic_card_data_id = ?`)

	return nil
}

func (queries *cardDbQueries) cleanupCardDbQueries() {
	if queries.NumAtomicPropertiesQuery != nil {
		queries.NumAtomicPropertiesQuery.Close()
	}

	if queries.AtomicPropertiesIdQuery != nil {
		queries.AtomicPropertiesIdQuery.Close()
	}

	if queries.AtomicPropertiesHashQuery != nil {
		queries.AtomicPropertiesHashQuery.Close()
	}

	if queries.InsertAtomicPropertiesQuery != nil {
		queries.InsertAtomicPropertiesQuery.Close()
	}

	if queries.InsertCardQuery != nil {
		queries.InsertCardQuery.Close()
	}

	if queries.UpdateCardQuery != nil {
		queries.UpdateCardQuery.Close()
	}

	if queries.InsertAltLangDataQuery != nil {
		queries.InsertAltLangDataQuery.Close()
	}

	if queries.InsertCardPrintingQuery != nil {
		queries.InsertCardPrintingQuery.Close()
	}

	if queries.InsertCardSubtypeQuery != nil {
		queries.InsertCardSubtypeQuery.Close()
	}

	if queries.InsertCardSupertypeQuery != nil {
		queries.InsertCardSupertypeQuery.Close()
	}

	if queries.InsertFrameEffectQuery != nil {
		queries.InsertFrameEffectQuery.Close()
	}

	if queries.InsertLeadershipSkillQuery != nil {
		queries.InsertLeadershipSkillQuery.Close()
	}

	if queries.InsertLegalityQuery != nil {
		queries.InsertLegalityQuery.Close()
	}

	if queries.InsertOtherFaceIdQuery != nil {
		queries.InsertOtherFaceIdQuery.Close()
	}

	if queries.InsertPurchaseURLQuery != nil {
		queries.InsertPurchaseURLQuery.Close()
	}

	if queries.InsertRulingQuery != nil {
		queries.InsertRulingQuery.Close()
	}

	if queries.InsertVariationQuery != nil {
		queries.InsertVariationQuery.Close()
	}

	if queries.InsertBaseTypeQuery != nil {
		queries.InsertBaseTypeQuery.Close()
	}

	if queries.DeleteFrameEffectQuery != nil {
		queries.DeleteFrameEffectQuery.Close()
	}

	if queries.DeleteOtherFaceQuery != nil {
		queries.DeleteOtherFaceQuery.Close()
	}

	if queries.DeleteVariationQuery != nil {
		queries.DeleteVariationQuery.Close()
	}

	if queries.UpdateRefCntQuery != nil {
		queries.UpdateRefCntQuery.Close()
	}
}

func maybeInsertSetToDb(db *sql.DB, updateStats *DbUpdateStats, wg *sync.WaitGroup, set MTGSet) {
	defer wg.Done()
	ctx := context.Background()

	// Open a DB connection
	dbConn, err := db.Conn(ctx)
	if err != nil {
		log.Print(err)
		return
	}
	defer dbConn.Close()

	tx, err := dbConn.BeginTx(ctx, nil)
	if err != nil {
		log.Print(err)
		return
	}

	// Prepare the various statements needed for set operations
	setHashQuery, err := tx.Prepare(`SELECT set_hash, set_id FROM sets WHERE code = ?`)
	if err != nil {
		log.Print(err)
		tx.Rollback()
		return
	}
	defer setHashQuery.Close()

	cardHashQuery, err := tx.Prepare(`SELECT full_card_hash, atomic_card_data_id FROM all_cards
		WHERE uuid = ?`)
	if err != nil {
		log.Print(err)
		tx.Rollback()
		return
	}
	defer cardHashQuery.Close()

	insertSetQuery, err := tx.Prepare(`INSERT INTO sets
		(set_hash, base_size, block_name, code, is_foreign_only, is_foil_only,
		is_online_only, is_partial_preview, keyrune_code, mcm_name, mcm_id,
		mtgo_code, name, parent_code, release_date, tcgplayer_group_id,
		total_set_size, set_type)
		VALUES
		(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		log.Print(err)
		tx.Rollback()
		return
	}
	defer insertSetQuery.Close()

	insertSetTranslationQuery, err := tx.Prepare(`INSERT INTO set_translations
		(set_id, set_translation_language_id, set_translated_name)
		VALUES
		(?, ?, ?)`)
	if err != nil {
		log.Print(err)
		return
	}
	defer insertSetTranslationQuery.Close()

	// Prepare the various statements needed for card operations
	var cardQueries cardDbQueries
	defer cardQueries.cleanupCardDbQueries()
	err = cardQueries.prepareCardDbQueries(tx)
	if err != nil {
		log.Print(err)
		return
	}

	// Hash the set for later use
	set.Canonicalize()
	setHash := HashToHexString(set.Hash())

	// First, check to see if this set is in the DB at all
	setExists, setDbHash, setId, err := set.CheckIfSetExists(setHashQuery)
	if err != nil {
		log.Print(err)
		tx.Rollback()
		return
	}

	updateStats.Mutex.Lock()
	updateStats.TotalSetsInUpdate += 1
	updateStats.Mutex.Unlock()

	totalNewCards := 0
	totalNewCardsInNewSets := 0
	totalNewCardsInExistingSets := 0
	totalNewAtomicCards := 0
	totalExistingCardsHashSkipped := 0
	totalExistingCardsUpdated := 0

	if setExists {
		log.Printf("Set %s already exists in the database\n", set.Code)
		updateStats.Mutex.Lock()
		updateStats.TotalExistingSetsInDb += 1
		updateStats.Mutex.Unlock()
		// This set already exists in the db
		// Check to see if the hash matcdbhes what's already in the db
		if setDbHash == setHash {
			// Hashes match, so we can skip updating this set in the db
			log.Printf("Set %s in db matches hash %s, skipping update...\n", set.Code, setDbHash)
			updateStats.Mutex.Lock()
			updateStats.ExistingSetsSkippedDueToHashMatch += 1
			updateStats.Mutex.Unlock()
		} else {
			// Hashes don't match, so we need to look at each card in the set to update
			log.Printf("Set %s hashes don't match (db: %s, json: %s), updating set...\n",
				set.Code, setDbHash, setHash)
			updateStats.Mutex.Lock()
			updateStats.ExistingSetsUpdatedDueToHashMismatch += 1
			updateStats.Mutex.Unlock()

			// For each card, check if the card exists, and if so, if the hash
			// matches
			for _, card := range set.Cards {
				cardExists, cardDbHash, atomicCardDataId, err := card.CheckIfCardExists(cardHashQuery)
				if err != nil {
					log.Print(err)
					tx.Rollback()
					return
				}

				if !cardExists {
					newAtomicPropertiesAdded, err := card.InsertAllCardDataToDb(&cardQueries, setId)
					if err != nil {
						log.Print(err)
						tx.Rollback()
						return
					}
					if newAtomicPropertiesAdded {
						totalNewAtomicCards += 1
					}
					totalNewCardsInExistingSets += 1
				} else {
					// Check if the stored hash matches
					cardHash := HashToHexString(card.Hash())
					if cardHash == cardDbHash {
						// Can skip
						log.Printf("Card %s hash matches in db (%s), skipping", card.Name, cardHash)
						totalExistingCardsHashSkipped += 1
					} else {
						// Need to update card
						log.Printf("Card %s hash doesn't match (db: %s, card: %s), updating",
							card.Name, cardDbHash, cardHash)
						totalExistingCardsUpdated += 1
						err := card.UpdateCardDataInDb(&cardQueries, atomicCardDataId, setId)
						if err != nil {
							log.Print(err)
							tx.Rollback()
							return
						}
					}
				}
			}
		}
	} else {
		// This set does not already exist in the db
		setId, err := set.InsertSetToDb(insertSetQuery, setHash)
		if err != nil {
			log.Print(err)
			tx.Rollback()
			return
		}
		updateStats.Mutex.Lock()
		updateStats.TotalNewSets += 1
		updateStats.Mutex.Unlock()

		// Insert the set translations
		for lang, name := range set.Translations {
			err := InsertSetTranslationToDb(insertSetTranslationQuery,
				setId, lang, name)
			if err != nil {
				log.Print(err)
				tx.Rollback()
				return
			}
		}

		// Insert all of the cards in the set.  No need to check the full card hash, since we're bulk
		// inserting the entire set
		log.Printf("Processing cards in set %s\n", set.Code)
		for _, card := range set.Cards {
			totalNewCards += 1
			totalNewCardsInNewSets += 1
			card.Canonicalize()

			newAtomicPropertiesAdded, err := card.InsertAllCardDataToDb(&cardQueries, setId)
			if err != nil {
				log.Print(err)
				tx.Rollback()
				return
			}
			if newAtomicPropertiesAdded {
				totalNewAtomicCards += 1
			}
		}
	}

	tx.Commit()
	updateStats.Mutex.Lock()
	updateStats.TotalCardsInUpdate += totalNewCards
	updateStats.TotalNewCards += totalNewCards
	updateStats.TotalNewCardsInNewSets += totalNewCardsInNewSets
	updateStats.TotalNewAtomicCards += totalNewAtomicCards
	updateStats.TotalNewCardsInExistingSets += totalNewCardsInExistingSets
	updateStats.ExistingCardsSkippedDueToHashMatch += totalExistingCardsHashSkipped
	updateStats.ExistingCardsUpdatedDueToHashMismatch += totalExistingCardsUpdated
	updateStats.Mutex.Unlock()
	log.Printf("Done processing set %s\n", set.Code)
}

func (card *MTGCard) UpdateCardDataInDb(queries *cardDbQueries,
		atomicPropertiesId int64, setId int64) (error) {
	// First, check to see if the atomic properties hash still matches.  If it does,
	// we just need to update the rest of the card data, and can leave it pointing
	// to the same atomic properties record.
	var err error
	var dbHash string
	res := queries.AtomicPropertiesHashQuery.QueryRow(atomicPropertiesId)
	if err = res.Scan(&dbHash); err != nil {
		return err
	}
	atomicPropHash := HashToHexString(card.AtomicPropertiesHash())

	if dbHash != atomicPropHash {
		// The atomic properties hash doesn't match, so insert a new atomic properties
		// record
		atomicPropertiesId, err = card.InsertAtomicPropertiesToDb(queries, atomicPropHash)
		if err != nil {
			return err
		}
	}

	// Now, update the card record, clear any entries from auxilliary tables belonging
	// to the old card record, and insert new auxilliary entries for the updated card record
	err = card.UpdateCardInDb(queries, atomicPropertiesId, setId)
	if err != nil {
		return err
	}

	err = card.DeleteOtherTableCardData(queries)
	if err != nil {
		return err
	}

	err = card.InsertOtherTableCardData(queries)
	if err != nil {
		return err
	}

	return nil
}

func (card *MTGCard) InsertAllCardDataToDb(queries *cardDbQueries, setId int64) (bool, error) {
	newAtomicPropertiesAdded := false

	// First, calculate the atomic properties hash, so we can see if this card
	// shares its atomic properties with an existing card in the db
	atomicPropHash := HashToHexString(card.AtomicPropertiesHash())
	atomicPropId, refCnt, exists, err := card.GetAtomicPropertiesId(queries, atomicPropHash)
	if err != nil {
		return false, err
	}

	if !exists {
		// If the atomic properties don't exist already, we need to insert
		// a new record
		atomicPropId, err = card.InsertAtomicPropertiesToDb(queries, atomicPropHash)
		if err != nil {
			return false, err
		}
		newAtomicPropertiesAdded = true
	} else {
		// Otherwise, update the reference count of this atomic properties record
		_, err := queries.UpdateRefCntQuery.Exec(refCnt + 1, atomicPropId)
		if err != nil {
			return false, err
		}
	}

	// Insert the main card record in the all_cards table
	err = card.InsertCardToDb(queries, atomicPropId, setId)
	if err != nil {
		return false, err
	}

	// Insert the rest of the card data
	err = card.InsertOtherTableCardData(queries)
	if err != nil {
		return false, nil
	}

	return newAtomicPropertiesAdded, nil
}

func (card *MTGCard) InsertOtherTableCardData(queries *cardDbQueries) error {
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

func (card *MTGCard) DeleteOtherTableCardData(queries *cardDbQueries) error {
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

func (card *MTGCard) CheckIfCardExists(query *sql.Stmt) (bool, string, int64, error) {
	resultRow := query.QueryRow(card.UUID)

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

func (set *MTGSet) CheckIfSetExists(query *sql.Stmt) (bool, string, int64, error) {
	resultRow := query.QueryRow(set.Code)

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

func (set *MTGSet) InsertSetToDb(query *sql.Stmt, setHash string) (int64, error) {
	res, err := query.Exec(setHash, set.BaseSetSize, set.Block, set.Code, set.IsForeignOnly,
		set.IsFoilOnly, set.IsOnlineOnly, set.IsPartialPreview, set.KeyruneCode, set.MCMName,
		set.MCMId, set.MTGOCode, set.Name, set.ParentCode, set.ReleaseDate, set.TCGPlayerGroupId,
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

func InsertSetTranslationToDb(query *sql.Stmt, setId int64, translationLang string,
		translatedName string) error {
	languageId, err := getSetTranslationLanguageId(translationLang)
	if err != nil {
		return err
	}

	res, err := query.Exec(setId, languageId, translatedName)
	if err != nil {
		return err
	}

	return checkRowsAffected(res, 1, "insert set name translation")
}

func (card *MTGCard) InsertFrameEffectToDb(queries *cardDbQueries, frameEffect string) error {
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

func (card *MTGCard) InsertOtherFaceIdToDb(queries *cardDbQueries, otherFaceUUID string) error {
	res, err := queries.InsertOtherFaceIdQuery.Exec(card.UUID, otherFaceUUID)
	if err != nil {
		return err
	}

	return checkRowsAffected(res, 1, "insert other face ID")
}

func (card *MTGCard) InsertVariationToDb(queries *cardDbQueries, variationUUID string) error {
	res, err := queries.InsertVariationQuery.Exec(card.UUID, variationUUID)
	if err != nil {
		return err
	}

	return checkRowsAffected(res, 1, "insert variation")
}

func (card *MTGCard) InsertCardToDb(queries *cardDbQueries, atomicPropertiesId int64,
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

func (card *MTGCard) UpdateCardInDb(queries *cardDbQueries, atomicPropertiesId int64,
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

