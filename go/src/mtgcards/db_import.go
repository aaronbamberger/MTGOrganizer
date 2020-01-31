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

    var getQueries dbGetQueries
    var insertQueries dbInsertQueries
    var updateQueries dbUpdateQueries
    var deleteQueries dbDeleteQueries

    defer getQueries.Cleanup()
    err = getQueries.Prepare(db)
    if err != nil {
        return nil, err
    }
    defer insertQueries.Cleanup()
    err = insertQueries.Prepare(db)
    if err != nil {
        return nil, err
    }
    defer updateQueries.Cleanup()
    err = updateQueries.Prepare(db)
    if err != nil {
        return nil, err
    }
    defer deleteQueries.Cleanup()
    err = deleteQueries.Prepare(db)
    if err != nil {
        return nil, err
    }

	for _, set := range sets {
		setImportWg.Add(1)
		go maybeInsertSetToDb(
            db,
            &getQueries,
            &insertQueries,
            &updateQueries,
            &deleteQueries,
            &stats,
            &setImportWg,
            set)
	}

	setImportWg.Wait()
	return &stats, nil
}

func maybeInsertSetToDb(db *sql.DB,
        getQueries *dbGetQueries,
        insertQueries *dbInsertQueries,
        updateQueries *dbUpdateQueries,
        deleteQueries *dbDeleteQueries,
        stats *DbUpdateStats,
		wg *sync.WaitGroup,
        set MTGSet) {
	defer wg.Done()
	ctx := context.Background()

	// Open a DB connection
	conn, err := db.Conn(ctx)
	if err != nil {
		log.Print(err)
		return
	}
	defer conn.Close()

	// Transaction for inserting the set itself
	setTx, err := conn.BeginTx(ctx, nil)
	if err != nil {
		log.Print(err)
		return
	}

	// Hash the set for later use
	set.Canonicalize()
	setHash := HashToHexString(set.Hash())

    setGetQueries := getQueries.ForTx(setTx)

	// First, check to see if this set is in the DB at all
	setExists, setDbHash, setId, err := set.GetHashAndId(setGetQueries)
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

            // TODO: Update the set

			setTx.Commit()
            stats.AddToExistingSetsUpdated(1)

			// For each card, check if the card exists, and if so, if the hash
			// matches
			for idx := range set.Cards {
                // Need to access by index here to get a pointer to the card,
                // not a copy
                card := &set.Cards[idx]

				// Transaction for each card
				cardTx, err := conn.BeginTx(ctx, nil)
				if err != nil {
					log.Print(err)
					continue
				}

				cardGetQueries := getQueries.ForTx(cardTx)

				cardExists, cardDbHash, _, err := card.GetHashAndId(cardGetQueries)
				if err != nil {
					log.Print(err)
					cardTx.Rollback()
					continue
				}

				if !cardExists {
                    cardInsertQueries := insertQueries.ForTx(cardTx)
                    err := card.InsertToDb(cardInsertQueries, setId)
					if err != nil {
						log.Print(err)
						cardTx.Rollback()
						continue
					}

                    cardTx.Commit()
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

                        // TODO: Update the card
                        /*
						newAtomicPropetiesAdded, err := card.UpdateCardDataInDb(cardQueries,
                            atomicCardDataId, setId)
						if err != nil {
							log.Print(err)
							cardTx.Rollback()
							continue
						}
                        */
                        cardTx.Commit()
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

        setInsertQueries := insertQueries.ForTx(setTx)
		setId, err := set.InsertToDb(setInsertQueries)
		if err != nil {
			log.Print(err)
			setTx.Rollback()
			return
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
			cardTx, err := conn.BeginTx(ctx, nil)
			if err != nil {
				log.Print(err)
				continue
			}

            cardInsertQueries := insertQueries.ForTx(cardTx)

			err = card.InsertToDb(cardInsertQueries, setId)
			if err != nil {
				log.Print(err)
				cardTx.Rollback()
				continue
			}
			cardTx.Commit()

            totalCards += 1
            totalNewCards += 1
			totalNewCardsInNewSets += 1
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

