package main

import "compress/gzip"
import "database/sql"
import _ "github.com/go-sql-driver/mysql"
import "encoding/json"
import "fmt"
import "mtgcards"
import "net/http"

func main() {
	resp, err := http.Get("https://www.mtgjson.com/files/AllPrintings.json.gz")
	if err != nil {
		fmt.Println("Error while downloading: %s\n", err)
	}
	defer resp.Body.Close()
	fmt.Printf("Response status: %s\n", resp.Status)
	fmt.Printf("Response proto: %s\n", resp.Proto)
	fmt.Printf("Response length: %d\n", resp.ContentLength)
	fmt.Printf("Response encodings: %v\n", resp.TransferEncoding)

	fmt.Printf("Parsing JSON response\n")
	decompressor, err := gzip.NewReader(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	decoder := json.NewDecoder(decompressor)
	var allSets map[string]mtgcards.MTGSet
	if err := decoder.Decode(&allSets); err != nil {
		fmt.Println(err)
		return
	}

	// Connect to the database
	db, err := sql.Open("mysql", "app_user:app_db_password@tcp(172.18.0.3)/mtg_cards")
	defer db.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	// Prepare all of the queries we'll want to use multiple times
	fmt.Printf("Preparing queries\n")
	setHashQuery, err := db.Prepare("SELECT set_hash FROM sets WHERE code = ?")
	defer setHashQuery.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Prepare insert set query\n")
	insertSetQuery, err := db.Prepare(`INSERT INTO sets
		(set_hash, base_size, block_name, code, is_foreign_only, is_foil_only,
		is_online_only, is_partial_preview, keyrune_code, mcm_name, mcm_id,
		mtgo_code, name, parent_code, release_date, tcgplayer_group_id,
		total_set_size, set_type)
		VALUES
		(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	defer insertSetQuery.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	err = mtgcards.CreateDbQueries(db)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer mtgcards.CloseDbQueries()

	// Fetch some things from the db for future use
	gameFormats, err := mtgcards.GetGameFormats(db)
	if err != nil {
		fmt.Println(err)
		return
	}
	legalityOptions, err := mtgcards.GetLegalityOptions(db)
	if err != nil {
		fmt.Println(err)
		return
	}
	purchaseSites, err := mtgcards.GetPurchaseSites(db)
	if err != nil {
		fmt.Println(err)
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
		fmt.Println(err)
		return
	}

	totalSets := len(allSets)
	currentSet := 1
	for setCode, set := range allSets {
		fmt.Printf("Processing set with code %s (%d of %d)\n", setCode, currentSet, totalSets)
		currentSet += 1

		// Hash the set for later use
		set.Canonicalize()
		setHash := mtgcards.HashToHexString(set.Hash())

		// First, check to see if this set is in the DB at all
		setRows, err := setHashQuery.Query(setCode)
		if err != nil {
			if setRows != nil {
				setRows.Close()
			}
			fmt.Println(err)
			continue
		}
		if setRows.Next() {
			fmt.Printf("Set %s already exists in the database\n", setCode)
			// This set already exists in the db
			// Check to see if the hash matches what's already in the db
			var dbSetHash string
			err := setRows.Scan(&dbSetHash)
			setRows.Close()
			if err != nil {
				fmt.Println(err)
				continue
			} else {
				if dbSetHash == setHash {
					// Hashes match, so we can skip updating this set in the db
					fmt.Printf("Set %s in db matches hash %s, skipping update...\n", setCode, dbSetHash)
					continue
				} else {
					// Hashes don't match, so we need to look at each card in the set to update
					fmt.Printf("Set %s hashes don't match (db: %s, json: %s), updating set...\n",
						setCode, dbSetHash, setHash)
					//TODO: Maybe update cards in set
				}
			}
		} else {
			// This set does not already exist in the db
			fmt.Printf("Set %s does not exist in the db, inserting the set\n", setCode)
			res, err := insertSetQuery.Exec(setHash, set.BaseSetSize, set.Block, set.Code, set.IsForeignOnly,
				set.IsFoilOnly, set.IsOnlineOnly, set.IsPartialPreview, set.KeyruneCode, set.MCMName,
				set.MCMId, set.MTGOCode, set.Name, set.ParentCode, set.ReleaseDate, set.TCGPlayerGroupId,
				set.TotalSetSize, set.Type)
			if err != nil {
				fmt.Println(err)
				continue
			}

			setId, err := res.LastInsertId()
			if err != nil {
				fmt.Println(err)
				continue
			}

			// Insert the set translations
			fmt.Printf("Set id: %d\n", setId)
			//TODO: Figure out a better way to do this

			// Insert all of the cards in the set.  No need to check the full card hash, since we're bulk
			// inserting the entire set
			fmt.Printf("Processing cards in set %s\n", setCode)
			for _, card := range set.Cards {
				card.Canonicalize()
				// First, calculate the atomic properties hash, so we can see if this card
				// shares its atomic properties with an existing card in the db
				var atomicPropId int64
				var exists bool
				atomicPropHash := mtgcards.HashToHexString(card.AtomicPropertiesHash())
				atomicPropId, exists, err = mtgcards.GetAtomicPropertiesId(atomicPropHash, &card)
				if err != nil {
					fmt.Println(err)
					continue
				}

				if !exists {
					// If the atomic properties don't exist already, we need to insert
					// a new record
					atomicPropId, err = card.InsertAtomicPropertiesToDb(atomicPropHash)
					if err != nil {
						fmt.Println(err)
						continue
					}
				}

				// Now, insert the rest of the card data
				err = card.InsertCardToDb(atomicPropId)
				if err != nil {
					fmt.Println(err)
					continue
				}

				// Alternate language data
				for _, altLangData := range card.AlternateLanguageData {
					err = altLangData.InsertAltLangDataToDb(atomicPropId)
					if err != nil {
						fmt.Println(err)
					}
				}

				// Leadership skills
				for leadershipFormat, leaderValid := range card.LeadershipSkills {
					err = mtgcards.InsertLeadershipSkillToDb(atomicPropId,
						leadershipFormats[leadershipFormat], leaderValid)
					if err != nil {
						fmt.Println(err)
					}
				}

				// Legalities
				for format, legality := range card.Legalities {
					err = mtgcards.InsertLegalityToDb(atomicPropId,
						gameFormats[format], legalityOptions[legality])
					if err != nil {
						fmt.Println(err)
					}
				}

				// Other face IDs
				for _, otherFaceId := range card.OtherFaceIds {
					err = card.InsertOtherFaceIdToDb(otherFaceId)
					if err != nil {
						fmt.Println(err)
					}
				}

				// Printings
				for _, setCode := range card.Printings {
					err = mtgcards.InsertCardPrintingToDb(atomicPropId, setCode)
					if err != nil {
						fmt.Println(err)
					}
				}

				// Purchase URLs
				err = mtgcards.InsertPurchaseURLToDb(atomicPropId,
					purchaseSites["cardmarket"], card.PurchaseURLs.Cardmarket)
				if err != nil {
					fmt.Println(err)
				}
				err = mtgcards.InsertPurchaseURLToDb(atomicPropId,
					purchaseSites["tcgplayer"], card.PurchaseURLs.TCGPlayer)
				if err != nil {
					fmt.Println(err)
				}
				err = mtgcards.InsertPurchaseURLToDb(atomicPropId,
					purchaseSites["mtgstocks"], card.PurchaseURLs.MTGStocks)
				if err != nil {
					fmt.Println(err)
				}

				// Rulings
				for _, ruling := range card.Rulings {
					err = ruling.InsertRulingToDb(atomicPropId)
					if err != nil {
						fmt.Println(err)
					}
				}

				// Subtypes
				for _, subtype := range card.Subtypes {
					err = mtgcards.InsertCardSubtypeToDb(atomicPropId, subtype)
					if err != nil {
						fmt.Println(err)
					}
				}

				// Supertypes
				for _, supertype := range card.Supertypes {
					err = mtgcards.InsertCardSupertypeToDb(atomicPropId, supertype)
					if err != nil {
						fmt.Println(err)
					}
				}

				// Variations
				for _, variation := range card.Variations {
					err = card.InsertVariationToDb(variation)
					if err != nil {
						fmt.Println(err)
					}
				}
			}
		}
	}
}
