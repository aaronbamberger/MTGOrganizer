package main

import "database/sql"
import _ "github.com/go-sql-driver/mysql"
import "log"
import "mtgcards"

func main() {
	allSets, err := mtgcards.DownloadAllPrintings(true)
	if err != nil {
		log.Fatal(err)
	}

	// Connect to the database
	db, err := sql.Open("mysql", "app_user:app_db_password@tcp(172.18.0.3)/mtg_cards")
	defer db.Close()
	if err != nil {
		log.Fatal(err)
	}

	err = mtgcards.CreateDbQueries(db)
	if err != nil {
		log.Print(err)
		return
	}
	defer mtgcards.CloseDbQueries()

	// Fetch some things from the db for future use
	gameFormats, err := mtgcards.GetGameFormats(db)
	if err != nil {
		log.Print(err)
		return
	}
	legalityOptions, err := mtgcards.GetLegalityOptions(db)
	if err != nil {
		log.Print(err)
		return
	}
	purchaseSites, err := mtgcards.GetPurchaseSites(db)
	if err != nil {
		log.Print(err)
		return
	}
	/*
	setTranslationLanguages, err := mtgcards.GetSetTranslationLanguages(db)
	if err != nil {
		fmt.Println(err)
		return
	}
	*/
	leadershipFormats, err := mtgcards.GetLeadershipFormats(db)
	if err != nil {
		log.Print(err)
		return
	}

	totalSets := len(allSets)
	currentSet := 1
	set := allSets["7ED"]
	//for _, set := range allSets {
		log.Printf("Processing set with code %s (%d of %d)\n", set.Code, currentSet, totalSets)
		currentSet += 1

		// Hash the set for later use
		set.Canonicalize()
		setHash := mtgcards.HashToHexString(set.Hash())

		// First, check to see if this set is in the DB at all
		setExists, setDbHash, err := set.CheckIfSetExists()
		if err != nil {
			log.Print(err)
			//continue
		}

		if setExists {
			log.Printf("Set %s already exists in the database\n", set.Code)
			// This set already exists in the db
			// Check to see if the hash matches what's already in the db
			if setDbHash == setHash {
				// Hashes match, so we can skip updating this set in the db
				log.Printf("Set %s in db matches hash %s, skipping update...\n", set.Code, setDbHash)
				//continue
			} else {
				// Hashes don't match, so we need to look at each card in the set to update
				log.Printf("Set %s hashes don't match (db: %s, json: %s), updating set...\n",
					set.Code, setDbHash, setHash)
				//TODO: Maybe update cards in set
			}
		} else {
			// This set does not already exist in the db
			err := set.InsertSetToDb(setHash)
			if err != nil {
				log.Print(err)
				//continue
			}

			// Insert the set translations
			//TODO: Figure out a better way to do this

			// Insert all of the cards in the set.  No need to check the full card hash, since we're bulk
			// inserting the entire set
			log.Printf("Processing cards in set %s\n", set.Code)
			for _, card := range set.Cards {
				card.Canonicalize()
				// First, calculate the atomic properties hash, so we can see if this card
				// shares its atomic properties with an existing card in the db
				var atomicPropId int64
				var exists bool
				atomicPropHash := mtgcards.HashToHexString(card.AtomicPropertiesHash())
				atomicPropId, exists, err = mtgcards.GetAtomicPropertiesId(atomicPropHash, &card)
				if err != nil {
					log.Print(err)
					continue
				}

				if !exists {
					// If the atomic properties don't exist already, we need to insert
					// a new record
					atomicPropId, err = card.InsertAtomicPropertiesToDb(atomicPropHash)
					if err != nil {
						log.Print(err)
						continue
					}
				}

				// Now, insert the rest of the card data
				err = card.InsertCardToDb(atomicPropId)
				if err != nil {
					log.Print(err)
					continue
				}

				// Alternate language data
				for _, altLangData := range card.AlternateLanguageData {
					err = altLangData.InsertAltLangDataToDb(atomicPropId)
					if err != nil {
						log.Print(err)
					}
				}

				// Frame effects
				for _, frameEffect := range card.FrameEffects {
					err = card.InsertFrameEffectToDb(frameEffect)
					if err != nil {
						log.Print(err)
					}
				}

				// Leadership skills
				for leadershipFormat, leaderValid := range card.LeadershipSkills {
					err = mtgcards.InsertLeadershipSkillToDb(atomicPropId,
						leadershipFormats[leadershipFormat], leaderValid)
					if err != nil {
						log.Print(err)
					}
				}

				// Legalities
				for format, legality := range card.Legalities {
					err = mtgcards.InsertLegalityToDb(atomicPropId,
						gameFormats[format], legalityOptions[legality])
					if err != nil {
						log.Print(err)
					}
				}

				// Other face IDs
				for _, otherFaceId := range card.OtherFaceIds {
					err = card.InsertOtherFaceIdToDb(otherFaceId)
					if err != nil {
						log.Print(err)
					}
				}

				// Printings
				for _, setCode := range card.Printings {
					err = mtgcards.InsertCardPrintingToDb(atomicPropId, setCode)
					if err != nil {
						log.Print(err)
					}
				}

				// Purchase URLs
				err = mtgcards.InsertPurchaseURLToDb(atomicPropId,
					purchaseSites["cardmarket"], card.PurchaseURLs.Cardmarket)
				if err != nil {
					log.Print(err)
				}
				err = mtgcards.InsertPurchaseURLToDb(atomicPropId,
					purchaseSites["tcgplayer"], card.PurchaseURLs.TCGPlayer)
				if err != nil {
					log.Print(err)
				}
				err = mtgcards.InsertPurchaseURLToDb(atomicPropId,
					purchaseSites["mtgstocks"], card.PurchaseURLs.MTGStocks)
				if err != nil {
					log.Print(err)
				}

				// Rulings
				for _, ruling := range card.Rulings {
					err = ruling.InsertRulingToDb(atomicPropId)
					if err != nil {
						log.Print(err)
					}
				}

				// Subtypes
				for _, subtype := range card.Subtypes {
					err = mtgcards.InsertCardSubtypeToDb(atomicPropId, subtype)
					if err != nil {
						log.Print(err)
					}
				}

				// Supertypes
				for _, supertype := range card.Supertypes {
					err = mtgcards.InsertCardSupertypeToDb(atomicPropId, supertype)
					if err != nil {
						log.Print(err)
					}
				}

				// Variations
				for _, variation := range card.Variations {
					err = card.InsertVariationToDb(variation)
					if err != nil {
						log.Print(err)
					}
				}
			}
		}
	//}
}
