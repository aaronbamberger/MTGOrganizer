package main

import "compress/gzip"
import "database/sql"
import _ "github.com/go-sql-driver/mysql"
import "encoding/json"
import "encoding/hex"
import "fmt"
import "hash"
import "mtgcards"
import "net/http"

func HashToHexString(hashVal hash.Hash) string {
	hashBytes := make([]byte, 0, hashVal.Size())
	hashBytes = hashVal.Sum(hashBytes)
	return hex.EncodeToString(hashBytes)
}

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
	/*
	insertCardQuery, err := db.Prepare(`INSERT INTO all_cards
		(uuid, full_card_hash, atomic_card_data_hash, artist, border_color,
		card_number, scryfall_id, watermark, frame_version, mcm_id, mcm_meta_id,
		multiverse_id, original_text, original_type, rarity, tcgplayer_product_id,
		duel_deck, flavor_text, has_foil, has_non_foil, is_alternative, is_arena,
		is_full_art, is_mtgo, is_online_only, is_oversized, is_paper, is_promo,
		is_reprint, is_starter, is_story_spotlight, is_textless, is_timeshifted,
		mtg_arena_id, mtgo_foil_id, mtgo_id, scryfall_illustration_id)
		VALUES
		(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?,
		?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	defer insertCardQuery.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	insertAltLangDataQuery, err := db.Prepare(`INSERT INTO alternate_language_data
		(card_data_hash, flavor_text, language, multiverse_id, name, text, card_type)
		VALUES
		(?, ?, ?, ?, ?, ?, ?)`)
	defer insertAltLangDataQuery.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	insertCardPrintingQuery, err := db.Prepare(`INSERT INTO card_printings
		(card_data_hash, set_id)
		VALUES
		(?, ?)`)
	defer insertCardPrintingQuery.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	insertCardSubtypeQuery, err := db.Prepare(`INSERT INTO card_subtypes
		(card_data_hash, card_subtype)
		VALUES
		(?, ?)`)
	defer insertCardSubtypeQuery.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	insertCardSupertypeQuery, err := db.Prepare(`INSERT INTO card_supertypes
		(card_data_hash, card_supertype)
		VALUES
		(?, ?)`)
	defer insertCardSupertypeQuery.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	insertFrameEffectQuery, err := db.Prepare(`INSERT INTO frame_effects
		(card_uuid, frame_effect)
		VALUES
		(?, ?)`)
	defer insertFrameEffectQuery.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	insertLeadershipSkillQuery, err := db.Prepare(`INSERT INTO leadership_skills
		(card_data_hash, brawl_leader_legal, commander_leader_legal, oathbreaker_leader_legal)
		VALUES
		(?, ?, ?, ?)`)
	defer insertLeadershipSkillQuery.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	insertLegalityQuery, err := db.Prepare(`INSERT INTO legalities
		(card_data_hash, game_format_id, legality_option_id)
		VALUES
		(?, ?, ?, ?)`)
	defer insertLegalityQuery.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	insertOtherFaceIdQuery, err := db.Prepare(`INSERT INTO other_faces
		(card_uuid, other_face_id)
		VALUES
		(?, ?)`)
	if err != nil {
		fmt.Println(err)
		return
	}
	insertPurchaseUrlQuery, err := db.Prepare(`INSERT INTO purchase_urls
		(purchase_site_id, purchase_url)
		VALUES
		(?, ?)`)
	defer insertPurchaseUrlQuery.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	insertRulingQuery, err := db.Prepare(`INSERT INTO rulings
		(card_data_hash, ruling_date, ruling_text)
		VALUES
		(?, ?, ?)`)
	defer insertRulingQuery.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	insertSetTranslationQuery, err := db.Prepare(`INSERT INTO set_translations
		(set_id, set_translation_language_id, set_translated_name)
		VALUES
		(?, ?, ?)`)
	defer insertSetTranslationQuery.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	insertVariationQuery, err := db.Prepare(`INSERT INTO variations
		(card_uuid, variation_uuid)
		VALUES
		(?, ?)`)
	defer insertVariationQuery.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	*/

	err = mtgcards.CreateDbQueries(db)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Fetch some things from the db for future use
	/*
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
	setTranslationLanguages, err := mtgcards.GetSetTranslationLanguages(db)
	if err != nil {
		fmt.Println(err)
		return
	}
	*/

	for setCode, set := range allSets {
		fmt.Printf("Processing set with code %s\n", setCode)

		// Hash the set for later use
		set.Canonicalize()
		setHash := HashToHexString(set.Hash())

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
				}
			}
		} else {
			// This set does not already exist in the db
			fmt.Printf("Set %s does not exist in the db, inserting the set\n", setCode)
			_, err := insertSetQuery.Exec(setHash, set.BaseSetSize, set.Block, set.Code, set.IsForeignOnly,
				set.IsFoilOnly, set.IsOnlineOnly, set.IsPartialPreview, set.KeyruneCode, set.MCMName,
				set.MCMId, set.MTGOCode, set.Name, set.ParentCode, set.ReleaseDate, set.TCGPlayerGroupId,
				set.TotalSetSize, set.Type)
			if err != nil {
				fmt.Println(err)
				continue
			}
			// Insert all of the cards in the set.  No need to check the full card hash, since we're bulk
			// inserting the entire set
			fmt.Printf("Processing cards in set %s\n", setCode)
			for _, card := range set.Cards {
				card.Canonicalize()
				// First, calculate the atomic properties hash, so we can see if this card
				// shares its atomic properties with an existing card in the db
				atomicPropertiesHash := HashToHexString(card.AtomicPropertiesHash())
				atomicPropertiesExist, err := mtgcards.CheckAtomicPropertiesDataExistence(atomicPropertiesHash)
				if err != nil {
					fmt.Println(err)
					continue
				}

				// If this is a newly seen set of atomic properties, we need to insert
				// a new record
				if !atomicPropertiesExist {
					card.InsertAtomicPropertiesToDb(atomicPropertiesHash)
				}

				// Now, insert the rest of the card data
			}
		}
	}
}
