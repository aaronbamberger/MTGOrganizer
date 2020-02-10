package carddb

import "context"
import "database/sql"
import "log"
import "mtgcards"
import "sync"

func ImportSetsToDB(
        db *sql.DB,
        sets map[string]mtgcards.MTGSet) (*CardUpdateStats, error) {
	var setImportWg sync.WaitGroup
	var stats CardUpdateStats

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

    var getQueries DBGetQueries
    var insertQueries DBInsertQueries
    var updateQueries DBUpdateQueries
    var deleteQueries DBDeleteQueries

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

func maybeInsertSetToDb(
        db *sql.DB,
        getQueries *DBGetQueries,
        insertQueries *DBInsertQueries,
        updateQueries *DBUpdateQueries,
        deleteQueries *DBDeleteQueries,
        stats *CardUpdateStats,
		wg *sync.WaitGroup,
        set mtgcards.MTGSet) {
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
	setHash := set.Hash()

    setGetQueries := getQueries.ForTx(setTx)

	// First, check to see if this set is in the DB at all
	setExists, setDbHash, setId, err := GetSetHashAndIdFromDB(set.Code, setGetQueries)
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
	totalExistingCards := 0
	totalExistingCardsHashSkipped := 0
	totalExistingCardsUpdated := 0

    totalTokens := 0
    totalNewTokens := 0
    totalNewTokensInNewSets := 0
    totalNewTokensInExistingSets := 0
    totalExistingTokens := 0
    totalExistingTokensHashSkipped := 0
    totalExistingTokensUpdated := 0

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

            setUpdateQueries := updateQueries.ForTx(setTx)
            setDeleteQueries := deleteQueries.ForTx(setTx)
            setInsertQueries := insertQueries.ForTx(setTx)
            err := UpdateSetInDB(setId, &set, setUpdateQueries, setDeleteQueries, setInsertQueries)
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
				cardTx, err := conn.BeginTx(ctx, nil)
				if err != nil {
					log.Print(err)
					continue
				}

				cardGetQueries := getQueries.ForTx(cardTx)

				cardExists, cardDbHash, cardId, err := GetCardHashAndIdFromDB(card.UUID, cardGetQueries)
				if err != nil {
					log.Print(err)
					cardTx.Rollback()
					continue
				}

				if !cardExists {
                    cardInsertQueries := insertQueries.ForTx(cardTx)
                    err := InsertCardToDB(card, setId, cardInsertQueries)
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
					cardHash := card.Hash()
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

                        cardUpdateQueries := updateQueries.ForTx(cardTx)
                        cardDeleteQueries := deleteQueries.ForTx(cardTx)
                        cardInsertQueries := insertQueries.ForTx(cardTx)
                        err := UpdateCardInDB(
                            cardId,
                            setId,
                            card,
                            cardUpdateQueries,
                            cardDeleteQueries,
                            cardInsertQueries)
                        if err != nil {
                            log.Print(err)
                            cardTx.Rollback()
                            continue
                        }
                        cardTx.Commit()
                        totalCards += 1
                        totalExistingCards += 1
                        totalExistingCardsUpdated += 1
					}
				}
			}

			// For each token, check if the token exists, and if so, if the hash
			// matches
			for idx := range set.Tokens {
                // Need to access by index here to get a pointer to the token,
                // not a copy
                token := &set.Tokens[idx]

				// Transaction for each token
				tokenTx, err := conn.BeginTx(ctx, nil)
				if err != nil {
					log.Print(err)
					continue
				}

				tokenGetQueries := getQueries.ForTx(tokenTx)

				tokenExists, tokenDbHash, tokenId, err := GetTokenHashAndIdFromDB(
                    token.UUID,
                    tokenGetQueries)
				if err != nil {
					log.Print(err)
					tokenTx.Rollback()
					continue
				}

				if !tokenExists {
                    tokenInsertQueries := insertQueries.ForTx(tokenTx)
                    inserted, err := InsertTokenToDB(token, setId, tokenInsertQueries)
					if err != nil {
						log.Print(err)
						tokenTx.Rollback()
						continue
					}

                    tokenTx.Commit()
                    if inserted {
                        totalTokens += 1
					    totalNewTokens += 1
					    totalNewTokensInExistingSets += 1
                    }
				} else {
					// Check if the stored hash matches
					tokenHash := token.Hash()
					if tokenHash == tokenDbHash {
						// Can skip
						log.Printf("Token %s hash matches in db (%s), skipping", token.Name, tokenHash)
                        tokenTx.Commit()
                        totalTokens += 1
                        totalExistingTokens += 1
						totalExistingTokensHashSkipped += 1
					} else {
						// Need to update token
						log.Printf("Token %s hash doesn't match (db: %s, token: %s), updating",
							token.Name, tokenDbHash, tokenHash)

                        tokenUpdateQueries := updateQueries.ForTx(tokenTx)
                        tokenDeleteQueries := deleteQueries.ForTx(tokenTx)
                        tokenInsertQueries := insertQueries.ForTx(tokenTx)
                        err := UpdateTokenInDB(
                            tokenId,
                            setId,
                            token,
                            tokenUpdateQueries,
                            tokenDeleteQueries,
                            tokenInsertQueries)
                        if err != nil {
                            log.Print(err)
                            tokenTx.Rollback()
                            continue
                        }
                        tokenTx.Commit()
                        totalTokens += 1
                        totalExistingTokens += 1
                        totalExistingTokensUpdated += 1
					}
				}
			}

		}
		stats.AddToTotalExistingSets(1)
	} else {
		// This set does not already exist in the db

        setInsertQueries := insertQueries.ForTx(setTx)
		setId, err := InsertSetToDB(&set, setInsertQueries)
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

			err = InsertCardToDB(card, setId, cardInsertQueries)
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

		// Insert all of the tokens in the set.  No need to check the full token hash, since we're bulk
		// inserting the entire set
		log.Printf("Processing tokens in set %s\n", set.Code)
		for idx := range set.Tokens {
            // Need to access by index here to get a pointer to the token,
            // not a copy
            token := &set.Tokens[idx]

			// Transaction for each token
			tokenTx, err := conn.BeginTx(ctx, nil)
			if err != nil {
				log.Print(err)
				continue
			}

            tokenInsertQueries := insertQueries.ForTx(tokenTx)

            inserted, err := InsertTokenToDB(token, setId, tokenInsertQueries)
			if err != nil {
				log.Print(err)
				tokenTx.Rollback()
				continue
			}
			tokenTx.Commit()

            if inserted {
                totalTokens += 1
                totalNewTokens += 1
			    totalNewTokensInNewSets += 1
            }
		}

	}

	stats.AddToTotalCards(totalCards)
	stats.AddToTotalNewCards(totalNewCards)
	stats.AddToTotalNewCardsInNewSets(totalNewCardsInNewSets)
	stats.AddToTotalNewCardsInExistingSets(totalNewCardsInExistingSets)
	stats.AddToTotalExistingCards(totalExistingCards)
	stats.AddToExistingCardsSkipped(totalExistingCardsHashSkipped)
	stats.AddToExistingCardsUpdated(totalExistingCardsUpdated)
    stats.AddToTotalTokens(totalTokens)
    stats.AddToTotalNewTokens(totalNewTokens)
    stats.AddToTotalNewTokensInNewSets(totalNewTokensInNewSets)
    stats.AddToTotalNewTokensInExistingSets(totalNewTokensInExistingSets)
    stats.AddToTotalExistingTokens(totalExistingTokens)
    stats.AddToExistingTokensSkipped(totalExistingTokensHashSkipped)
    stats.AddToExistingTokensUpdated(totalExistingTokensUpdated)
	log.Printf("Done processing set %s\n", set.Code)
}
