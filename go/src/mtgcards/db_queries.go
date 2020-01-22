package mtgcards

import "context"
import "database/sql"
import "encoding/hex"
import "fmt"
import "hash"
import "log"
import "strings"
import "sync"

var insertGameFormatQuery *sql.Stmt
var insertLegalityOptionQuery *sql.Stmt
var insertPurchaseSiteQuery *sql.Stmt
var insertLeadershipFormatQuery *sql.Stmt
var insertSetTranslationLanguageQuery *sql.Stmt
var insertBaseTypeOptionQuery *sql.Stmt

var gameFormats sync.Map
var legalityOptions sync.Map
var purchaseSites sync.Map
var leadershipFormats sync.Map
var setTranslationLanguages sync.Map
var baseTypeOptions sync.Map
//var frameEffectOptions sync.Map

const cardsPerThread = 50

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

func CreateDbQueries(db *sql.DB) error {
	var err error

	insertGameFormatQuery, err = db.Prepare(`INSERT INTO game_formats
		(game_format_name)
		VALUES
		(?)`)
	if err != nil {
		return err
	}

	insertLegalityOptionQuery, err = db.Prepare(`INSERT INTO legality_options
		(legality_option_name)
		VALUES
		(?)`)
	if err != nil {
		return err
	}

	insertPurchaseSiteQuery, err = db.Prepare(`INSERT INTO purchase_sites
		(purchase_site_name)
		VALUES
		(?)`)
	if err != nil {
		return err
	}

	insertLeadershipFormatQuery, err = db.Prepare(`INSERT INTO leadership_formats
		(leadership_format_name)
		VALUES
		(?)`)
	if err != nil {
		return err
	}

	insertSetTranslationLanguageQuery, err = db.Prepare(`INSERT INTO set_translation_languages
		(set_translation_language)
		VALUES
		(?)`)
	if err != nil {
		return err
	}

	insertBaseTypeOptionQuery, err = db.Prepare(`INSERT INTO base_type_options
		(base_type_option)
		VALUES
		(?)`)
	if err != nil {
		return err
	}

	err = populateGameFormatsCache(db)
	if err != nil {
		return err
	}

	err = populateLegalityOptionsCache(db)
	if err != nil {
		return err
	}

	err = populatePurchaseSitesCache(db)
	if err != nil {
		return err
	}

	err = populateLeadershipFormatsCache(db)
	if err != nil {
		return err
	}

	err = populateSetTranslationLanguagesCache(db)
	if err != nil {
		return err
	}

	err = populateBaseTypeOptionsCache(db)
	if err != nil {
		return err
	}

	return nil
}

func CloseDbQueries() {
	if insertGameFormatQuery != nil {
		insertGameFormatQuery.Close()
	}

	if insertLegalityOptionQuery != nil {
		insertLegalityOptionQuery.Close()
	}

	if insertPurchaseSiteQuery != nil {
		insertPurchaseSiteQuery.Close()
	}

	if insertLeadershipFormatQuery != nil {
		insertLeadershipFormatQuery.Close()
	}

	if insertSetTranslationLanguageQuery != nil {
		insertSetTranslationLanguageQuery.Close()
	}

	if insertBaseTypeOptionQuery != nil {
		insertBaseTypeOptionQuery.Close()
	}
}

func MaybeInsertSetToDb(db *sql.DB, wg *sync.WaitGroup, set MTGSet) {
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
	setHashQuery, err := tx.Prepare("SELECT set_hash FROM sets WHERE code = ?")
	if err != nil {
		log.Print(err)
		tx.Rollback()
		return
	}
	defer setHashQuery.Close()

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

	// Prepare the various statements needed for card operations
	numAtomicPropertiesQuery, err := tx.Prepare(`SELECT COUNT(scryfall_oracle_id)
		FROM atomic_card_data
		WHERE card_data_hash = ?`)
	if err != nil {
		log.Print(err)
		return
	}
	defer numAtomicPropertiesQuery.Close()

	atomicPropertiesIdQuery, err := tx.Prepare(`SELECT atomic_card_data_id,
		scryfall_oracle_id
		FROM atomic_card_data
		WHERE card_data_hash = ?`)
	if err != nil {
		log.Print(err)
		return
	}
	defer atomicPropertiesIdQuery.Close()

	insertAtomicPropertiesQuery, err := tx.Prepare(`INSERT INTO atomic_card_data
		(card_data_hash, color_identity, color_indicator, colors, converted_mana_cost,
		edhrec_rank, face_converted_mana_cost, hand, is_reserved, layout, life,
		loyalty, mana_cost, mtgstocks_id, name, card_power, scryfall_oracle_id,
		side, text, toughness, card_type)
		VALUES
		(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		log.Print(err)
		return
	}
	defer insertAtomicPropertiesQuery.Close()

	insertCardQuery, err := tx.Prepare(`INSERT INTO all_cards
		(uuid, full_card_hash, atomic_card_data_id, artist, border_color,
		card_number, scryfall_id, watermark, frame_version, mcm_id, mcm_meta_id,
		multiverse_id, original_text, original_type, rarity, tcgplayer_product_id,
		duel_deck, flavor_text, has_foil, has_non_foil, is_alternative, is_arena,
		is_full_art, is_mtgo, is_online_only, is_oversized, is_paper, is_promo,
		is_reprint, is_starter, is_story_spotlight, is_textless, is_timeshifted,
		mtg_arena_id, mtgo_foil_id, mtgo_id, scryfall_illustration_id)
		VALUES
		(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?,
		?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		log.Print(err)
		return
	}
	defer insertCardQuery.Close()

	insertAltLangDataQuery, err := tx.Prepare(`INSERT INTO alternate_language_data
		(atomic_card_data_id, flavor_text, language, multiverse_id, name, text, card_type)
		VALUES
		(?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		log.Print(err)
		return
	}
	defer insertAltLangDataQuery.Close()

	insertCardPrintingQuery, err := tx.Prepare(`INSERT INTO card_printings
		(atomic_card_data_id, set_code)
		VALUES
		(?, ?)`)
	if err != nil {
		log.Print(err)
		return
	}
	defer insertCardPrintingQuery.Close()

	insertCardSubtypeQuery, err := tx.Prepare(`INSERT INTO card_subtypes
		(atomic_card_data_id, card_subtype)
		VALUES
		(?, ?)`)
	if err != nil {
		log.Print(err)
		return
	}
	defer insertCardSubtypeQuery.Close()

	insertCardSupertypeQuery, err := tx.Prepare(`INSERT INTO card_supertypes
		(atomic_card_data_id, card_supertype)
		VALUES
		(?, ?)`)
	if err != nil {
		log.Print(err)
		return
	}
	defer insertCardSupertypeQuery.Close()

	insertFrameEffectQuery, err := tx.Prepare(`INSERT INTO frame_effects
		(card_uuid, frame_effect)
		VALUES
		(?, ?)`)
	if err != nil {
		log.Print(err)
		return
	}
	defer insertFrameEffectQuery.Close()

	insertLeadershipSkillQuery, err := tx.Prepare(`INSERT INTO leadership_skills
		(atomic_card_data_id, leadership_format_id, leader_legal)
		VALUES
		(?, ?, ?)`)
	if err != nil {
		log.Print(err)
		return
	}
	defer insertLeadershipSkillQuery.Close()

	insertLegalityQuery, err := tx.Prepare(`INSERT INTO legalities
		(atomic_card_data_id, game_format_id, legality_option_id)
		VALUES
		(?, ?, ?)`)
	if err != nil {
		log.Print(err)
		return
	}
	defer insertLegalityQuery.Close()

	insertOtherFaceIdQuery, err := tx.Prepare(`INSERT INTO other_faces
		(card_uuid, other_face_uuid)
		VALUES
		(?, ?)`)
	if err != nil {
		log.Print(err)
		return
	}
	defer insertOtherFaceIdQuery.Close()

	insertPurchaseUrlQuery, err := tx.Prepare(`INSERT INTO purchase_urls
		(atomic_card_data_id, purchase_site_id, purchase_url)
		VALUES
		(?, ?, ?)`)
	if err != nil {
		log.Print(err)
		return
	}
	defer insertPurchaseUrlQuery.Close()

	insertRulingQuery, err := tx.Prepare(`INSERT INTO rulings
		(atomic_card_data_id, ruling_date, ruling_text)
		VALUES
		(?, ?, ?)`)
	if err != nil {
		log.Print(err)
		return
	}
	defer insertRulingQuery.Close()

	insertVariationQuery, err := tx.Prepare(`INSERT INTO variations
		(card_uuid, variation_uuid)
		VALUES
		(?, ?)`)
	if err != nil {
		log.Print(err)
		return
	}
	defer insertVariationQuery.Close()

	insertSetTranslationQuery, err := tx.Prepare(`INSERT INTO set_translations
		(set_id, set_translation_language_id, set_translated_name)
		VALUES
		(?, ?, ?)`)
	if err != nil {
		log.Print(err)
		return
	}
	defer insertSetTranslationQuery.Close()

	insertBaseTypeQuery, err := tx.Prepare(`INSERT INTO base_types
		(atomic_card_data_id, base_type_option_id)
		VALUES
		(?, ?)`)
	if err != nil {
		log.Print(err)
		return
	}

	// Hash the set for later use
	set.Canonicalize()
	setHash := HashToHexString(set.Hash())

	// First, check to see if this set is in the DB at all
	setExists, setDbHash, err := set.CheckIfSetExists(setHashQuery)
	if err != nil {
		log.Print(err)
		tx.Rollback()
		return
	}

	if setExists {
		log.Printf("Set %s already exists in the database\n", set.Code)
		// This set already exists in the db
		// Check to see if the hash matcdbhes what's already in the db
		if setDbHash == setHash {
			// Hashes match, so we can skip updating this set in the db
			log.Printf("Set %s in db matches hash %s, skipping update...\n", set.Code, setDbHash)
		} else {
			// Hashes don't match, so we need to look at each card in the set to update
			log.Printf("Set %s hashes don't match (db: %s, json: %s), updating set...\n",
				set.Code, setDbHash, setHash)
			//TODO: Maybe update cards in set
		}
	} else {
		// This set does not already exist in the db
		setId, err := set.InsertSetToDb(insertSetQuery, setHash)
		if err != nil {
			log.Print(err)
			tx.Rollback()
			return
		}

		// Insert the set translations
		for lang, name := range set.Translations {
			err := InsertSetTranslationToDb(insertSetTranslationLanguageQuery, insertSetTranslationQuery,
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
			card.Canonicalize()
			// First, calculate the atomic properties hash, so we can see if this card
			// shares its atomic properties with an existing card in the db
			var atomicPropId int64
			var exists bool
			atomicPropHash := HashToHexString(card.AtomicPropertiesHash())
			atomicPropId, exists, err = card.GetAtomicPropertiesId(numAtomicPropertiesQuery,
					atomicPropertiesIdQuery, atomicPropHash)
			if err != nil {
				log.Print(err)
				return
			}

			if !exists {
				// If the atomic properties don't exist already, we need to insert
				// a new record
				atomicPropId, err = card.InsertAtomicPropertiesToDb(insertAtomicPropertiesQuery,
					atomicPropHash)
				if err != nil {
					log.Print(err)
					return
				}
			}

			// Now, insert the rest of the card data
			err = card.InsertCardToDb(insertCardQuery, atomicPropId)
			if err != nil {
				log.Print(err)
				return
			}

			// Alternate language data
			for _, altLangData := range card.AlternateLanguageData {
				err = altLangData.InsertAltLangDataToDb(insertAltLangDataQuery, atomicPropId)
				if err != nil {
					log.Print(err)
				}
			}

			// Frame effects
			for _, frameEffect := range card.FrameEffects {
				err = card.InsertFrameEffectToDb(insertFrameEffectQuery, frameEffect)
				if err != nil {
					log.Print(err)
				}
			}

			// Leadership skills
			for leadershipFormat, leaderValid := range card.LeadershipSkills {
				err = InsertLeadershipSkillToDb(insertLeadershipFormatQuery,
					insertLeadershipSkillQuery, atomicPropId, leadershipFormat,
					leaderValid)
				if err != nil {
					log.Print(err)
				}
			}

			// Legalities
			for format, legality := range card.Legalities {
				err = InsertLegalityToDb(insertGameFormatQuery,
					insertLegalityOptionQuery, insertLegalityQuery, atomicPropId,
					format, legality)
				if err != nil {
					log.Print(err)
				}
			}

			// Other face IDs
			for _, otherFaceId := range card.OtherFaceIds {
				err = card.InsertOtherFaceIdToDb(insertOtherFaceIdQuery, otherFaceId)
				if err != nil {
					log.Print(err)
				}
			}

			// Printings
			for _, setCode := range card.Printings {
				err = InsertCardPrintingToDb(insertCardPrintingQuery, atomicPropId, setCode)
				if err != nil {
					log.Print(err)
				}
			}

			// Purchase URLs
			for site, url := range card.PurchaseURLs {
				err = InsertPurchaseURLToDb(insertPurchaseSiteQuery,
					insertPurchaseUrlQuery, atomicPropId, site, url)
			}
			if err != nil {
				log.Print(err)
			}

			// Rulings
			for _, ruling := range card.Rulings {
				err = ruling.InsertRulingToDb(insertRulingQuery, atomicPropId)
				if err != nil {
					log.Print(err)
				}
			}

			// Subtypes
			for _, subtype := range card.Subtypes {
				err = InsertCardSubtypeToDb(insertCardSubtypeQuery, atomicPropId, subtype)
				if err != nil {
					log.Print(err)
				}
			}

			// Supertypes
			for _, supertype := range card.Supertypes {
				err = InsertCardSupertypeToDb(insertCardSupertypeQuery, atomicPropId, supertype)
				if err != nil {
					log.Print(err)
				}
			}

			// Calculate the set of "base" types, which I'm defining as the set
			// subtraction of card.Types - (card.Subtypes + card.Supertypes)
			cardBaseTypes := make(map[string]bool)
			for _, cardType := range card.Types {
				var inSubtype, inSupertype bool
				for _, subtype := range card.Subtypes {
					if subtype == cardType {
						inSubtype = true
						break
					}
				}
				for _, supertype := range card.Supertypes {
					if supertype == cardType {
						inSupertype = true
						break
					}
				}
				if !inSubtype && !inSupertype {
					cardBaseTypes[cardType] = true
				}
			}
			for baseType, _ := range cardBaseTypes {
				err = InsertBaseTypeToDb(insertBaseTypeOptionQuery,
					insertBaseTypeQuery, atomicPropId, baseType)
				if err != nil {
					log.Print(err)
				}
			}

			// Variations
			for _, variation := range card.Variations {
				err = card.InsertVariationToDb(insertVariationQuery, variation)
				if err != nil {
					log.Print(err)
				}
			}
		}
	}

	tx.Commit()
	log.Printf("Done processing set %s\n", set.Code)
}

func (set *MTGSet) CheckIfSetExists(query *sql.Stmt) (bool, string, error) {
	// First, check to see if this set is in the DB at all
	setRows, err := query.Query(set.Code)
	if err != nil {
		return false, "", err
	}
	defer setRows.Close()

	if setRows.Next() {
		// This set already exists in the db
		// Get the hash associated with the existing set
		var dbSetHash string
		err := setRows.Scan(&dbSetHash)
		if err != nil {
			return false, "", err
		}

		return true, dbSetHash, nil
	} else {
		// This set doesn't exist in the db
		return false, "", nil
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


func (card *MTGCard) InsertFrameEffectToDb(query *sql.Stmt, frameEffect string) error {
	res, err := query.Exec(card.UUID, frameEffect)
	if err != nil {
		return err
	}

	return checkRowsAffected(res, 1, "insert frame effect")
}

func InsertLeadershipSkillToDb(insertFormatQuery *sql.Stmt, insertSkillQuery *sql.Stmt,
		atomicPropertiesId int64, leadershipFormat string, leaderLegal bool) error {
	// Get the leadership format id from the cache
	var leadershipFormatId int64
	leadershipFormatIdTemp, loaded := leadershipFormats.Load(leadershipFormat)
	if !loaded {
		// This is the unlikely case where we have a new value that isn't pre-populated in the db
		res, err := insertFormatQuery.Exec(leadershipFormat)
		if err != nil {
			return err
		}

		leadershipFormatId, err = res.LastInsertId()
		if err != nil {
			return err
		}
		leadershipFormats.Store(leadershipFormat, leadershipFormatId)
	} else {
		leadershipFormatId = leadershipFormatIdTemp.(int64)
	}

	res, err := insertSkillQuery.Exec(atomicPropertiesId, leadershipFormatId, leaderLegal)
	if err != nil {
		return err
	}

	return checkRowsAffected(res, 1, "insert leadership skill")
}

func InsertLegalityToDb(insertFormatQuery *sql.Stmt, insertOptionQuery *sql.Stmt,
		insertLegalityQuery *sql.Stmt, atomicPropertiesId int64, gameFormat string,
		legalityOption string) error {
	// Get the game format id from the cache
	var gameFormatId int64
	gameFormatIdTemp, loaded := gameFormats.Load(gameFormat)
	if !loaded {
		// This is the unlikely case where we have a new value that isn't pre-populated in the db
		res, err := insertFormatQuery.Exec(gameFormat)
		if err != nil {
			return err
		}

		gameFormatId, err = res.LastInsertId()
		if err != nil {
			return err
		}
		gameFormats.Store(gameFormat, gameFormatId)
	} else {
		gameFormatId = gameFormatIdTemp.(int64)
	}

	// Get the legality option id from the cache
	var legalityOptionId int64
	legalityOptionIdTemp, loaded := legalityOptions.Load(legalityOption)
	if !loaded {
		// This is the unlikely case where we have a new value that isn't pre-populated in the db
		res, err := insertOptionQuery.Exec(legalityOption)
		if err != nil {
			return err
		}

		legalityOptionId, err = res.LastInsertId()
		if err != nil {
			return err
		}
		legalityOptions.Store(legalityOption, legalityOptionId)
	} else {
		legalityOptionId = legalityOptionIdTemp.(int64)
	}

	res, err := insertLegalityQuery.Exec(atomicPropertiesId, gameFormatId, legalityOptionId)
	if err != nil {
		return err
	}

	return checkRowsAffected(res, 1, "insert legality")
}

func InsertCardPrintingToDb(query *sql.Stmt, atomicPropertiesId int64, setCode string) error {
	res, err := query.Exec(atomicPropertiesId, setCode)
	if err != nil {
		return err
	}

	return checkRowsAffected(res, 1, "insert card printing")
}

func InsertPurchaseURLToDb(insertSiteQuery *sql.Stmt, insertUrlQuery *sql.Stmt,
		atomicPropertiesId int64, purchaseSite string, purchaseURL string) error {
	// Get the purchase site id from the cache
	var purchaseSiteId int64
	purchaseSiteIdTemp, loaded := purchaseSites.Load(purchaseSite)
	if !loaded {
		// This is the unlikely case where we have a new value that isn't pre-populated in the db
		res, err := insertSiteQuery.Exec(purchaseSite)
		if err != nil {
			return err
		}

		purchaseSiteId, err = res.LastInsertId()
		if err != nil {
			return err
		}
		purchaseSites.Store(purchaseSite, purchaseSiteId)
	} else {
		purchaseSiteId = purchaseSiteIdTemp.(int64)
	}

	res, err := insertUrlQuery.Exec(atomicPropertiesId, purchaseSiteId, purchaseURL)
	if err != nil {
		return err
	}

	return checkRowsAffected(res, 1, "insert purchase url")
}

func InsertSetTranslationToDb(insertLangQuery *sql.Stmt, insertTranslationQuery *sql.Stmt,
		setId int64, translationLang string, translatedName string) error {
	// Get the language id from the cache
	var languageId int64
	languageIdTemp, loaded := setTranslationLanguages.Load(translationLang)
	if !loaded {
		// This is the unlikely case where we have a new value that isn't pre-populated in the db
		res, err := insertLangQuery.Exec(translationLang)
		if err != nil {
			return err
		}

		languageId, err = res.LastInsertId()
		if err != nil {
			return err
		}
		setTranslationLanguages.Store(translationLang, languageId)
	} else {
		languageId = languageIdTemp.(int64)
	}

	res, err := insertTranslationQuery.Exec(setId, languageId, translatedName)
	if err != nil {
		return err
	}

	return checkRowsAffected(res, 1, "insert set name translation")
}

func InsertBaseTypeToDb(insertBaseTypeOptionQuery *sql.Stmt, insertBaseTypeQuery *sql.Stmt,
		atomicPropertiesId int64, baseTypeOption string) error {
	// Get the base type id from the cache
	var baseTypeOptionId int64
	baseTypeOptionIdTemp, loaded := baseTypeOptions.Load(baseTypeOption)
	if !loaded {
		// This is the unlikely case where we have a new value that isn't pre-populated in the db
		res, err := insertBaseTypeOptionQuery.Exec(baseTypeOption)
		if err != nil {
			return err
		}

		baseTypeOptionId, err = res.LastInsertId()
		if err != nil {
			return err
		}
		baseTypeOptions.Store(baseTypeOption, baseTypeOptionId)
	} else {
		baseTypeOptionId = baseTypeOptionIdTemp.(int64)
	}

	res, err := insertBaseTypeQuery.Exec(atomicPropertiesId, baseTypeOptionId)
	if err != nil {
		return err
	}

	return checkRowsAffected(res, 1, "insert base type")
}

func (altLangInfo *MTGCardAlternateLanguageInfo) InsertAltLangDataToDb(query *sql.Stmt,
		atomicPropertiesId int64) error {
	res, err := query.Exec(atomicPropertiesId, altLangInfo.FlavorText,
		altLangInfo.Language, altLangInfo.MultiverseId, altLangInfo.Name,
		altLangInfo.Text, altLangInfo.Type)
	if err != nil {
		return err
	}

	return checkRowsAffected(res, 1, "insert alt lang info")
}

func (ruling *MTGCardRuling) InsertRulingToDb(query *sql.Stmt, atomicPropertiesId int64) error {
	res, err := query.Exec(atomicPropertiesId, ruling.Date, ruling.Text)
	if err != nil {
		return err
	}

	return checkRowsAffected(res, 1, "insert ruling")
}

func InsertCardSubtypeToDb(query *sql.Stmt, atomicPropertiesId int64, subtype string) error {
	res, err := query.Exec(atomicPropertiesId, subtype)
	if err != nil {
		return err
	}

	return checkRowsAffected(res, 1, "insert subtype")
}

func InsertCardSupertypeToDb(query *sql.Stmt, atomicPropertiesId int64, supertype string) error {
	res, err := query.Exec(atomicPropertiesId, supertype)
	if err != nil {
		return err
	}

	return checkRowsAffected(res, 1, "insert supertype")
}

func (card *MTGCard) InsertOtherFaceIdToDb(query *sql.Stmt, otherFaceUUID string) error {
	res, err := query.Exec(card.UUID, otherFaceUUID)
	if err != nil {
		return err
	}

	return checkRowsAffected(res, 1, "insert other face ID")
}

func (card *MTGCard) InsertVariationToDb(query *sql.Stmt, variationUUID string) error {
	res, err := query.Exec(card.UUID, variationUUID)
	if err != nil {
		return err
	}

	return checkRowsAffected(res, 1, "insert variation")
}

func (card *MTGCard) InsertCardToDb(query *sql.Stmt, atomicPropertiesId int64) error {
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

	res, err := query.Exec(card.UUID, cardHash, atomicPropertiesId,
		card.Artist, card.BorderColor, card.Number, card.ScryfallId,
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

func (card *MTGCard) InsertAtomicPropertiesToDb(query *sql.Stmt,
		atomicPropertiesHash string) (int64, error) {
	// Build the set values needed for color_identity, color_indicator, and colors
	var colorIdentity sql.NullString
	var colorIndicator sql.NullString
	var colors sql.NullString
	var edhrecRank sql.NullInt32
	var hand sql.NullString
	var life sql.NullString
	var loyalty sql.NullString
	var name sql.NullString
	var side sql.NullString

	if len(card.ColorIdentity) > 0 {
		colorIdentity.String = strings.Join(card.ColorIdentity, ",")
		colorIdentity.Valid = true
	}

	if len(card.ColorIndicator) > 0 {
		colorIndicator.String = strings.Join(card.ColorIndicator, ",")
		colorIndicator.Valid = true
	}

	if len(card.Colors) > 0 {
		colors.String = strings.Join(card.Colors, ",")
		colors.Valid = true
	}

	if card.EDHRecRank != 0 {
		edhrecRank.Int32 = int32(card.EDHRecRank)
		edhrecRank.Valid = true
	}

	if len(card.Hand) > 0 {
		hand.String = card.Hand
		hand.Valid = true
	}

	if len(card.Life) > 0 {
		life.String = card.Life
		life.Valid = true
	}

	if len(card.Loyalty) > 0 {
		loyalty.String = card.Loyalty
		loyalty.Valid = true
	}

	if len(card.Name) > 0 {
		name.String = card.Name
		name.Valid = true
	}

	if len(card.Side) > 0 {
		side.String = card.Side
		side.Valid = true
	}

	res, err := query.Exec(atomicPropertiesHash,
		colorIdentity,
		colorIndicator,
		colors,
		card.ConvertedManaCost,
		edhrecRank,
		card.FaceConvertedManaCost,
		hand,
		card.IsReserved,
		card.Layout,
		life,
		loyalty,
		card.ManaCost,
		card.MTGStocksId,
		name,
		card.Power,
		card.ScryfallOracleId,
		side,
		card.Text,
		card.Toughness,
		card.Type)

	if err != nil {
		return 0, err
	}

	err = checkRowsAffected(res, 1, "insert atomic card data")
	if err != nil {
		return 0, err
	}

	lastInsertId, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return lastInsertId, nil
}

func populateGameFormatsCache(db *sql.DB) error {
	retrieveGameFormatsQuery, err := db.Prepare(`SELECT game_format_id, game_format_name
		FROM game_formats`)
	defer retrieveGameFormatsQuery.Close()
	if err != nil {
		return err
	}

	gameFormatRows, err := retrieveGameFormatsQuery.Query()
	if err != nil {
		return err
	}
	defer gameFormatRows.Close()

	for gameFormatRows.Next() {
		if err := gameFormatRows.Err(); err != nil {
			return err
		}

		var gameFormatId int64
		var gameFormatName string
		err := gameFormatRows.Scan(&gameFormatId, &gameFormatName)
		if err != nil {
			return err
		}
		gameFormats.Store(gameFormatName, gameFormatId)
	}

	return nil
}

func populateLegalityOptionsCache(db *sql.DB) error {
	retrieveLegalityOptionsQuery, err := db.Prepare(`SELECT legality_option_id, legality_option_name
		FROM legality_options`)
	defer retrieveLegalityOptionsQuery.Close()
	if err != nil {
		return err
	}

	legalityOptionsRows, err := retrieveLegalityOptionsQuery.Query()
	if err != nil {
		return err
	}
	defer legalityOptionsRows.Close()

	for legalityOptionsRows.Next() {
		if err := legalityOptionsRows.Err(); err != nil {
			return err
		}

		var legalityOptionId int64
		var legalityOptionName string
		err := legalityOptionsRows.Scan(&legalityOptionId, &legalityOptionName)
		if err != nil {
			return err
		}
		legalityOptions.Store(legalityOptionName, legalityOptionId)
	}

	return nil
}

func populatePurchaseSitesCache(db *sql.DB) error {
	retrievePurchaseSitesQuery, err := db.Prepare(`SELECT purchase_site_id, purchase_site_name
		FROM purchase_sites`)
	defer retrievePurchaseSitesQuery.Close()
	if err != nil {
		return err
	}

	purchaseSitesRows, err := retrievePurchaseSitesQuery.Query()
	if err != nil {
		return err
	}
	defer purchaseSitesRows.Close()

	for purchaseSitesRows.Next() {
		if err := purchaseSitesRows.Err(); err != nil {
			return err
		}

		var purchaseSiteId int64
		var purchaseSiteName string
		err := purchaseSitesRows.Scan(&purchaseSiteId, &purchaseSiteName)
		if err != nil {
			return err
		}
		purchaseSites.Store(purchaseSiteName, purchaseSiteId)
	}

	return nil
}

func populateSetTranslationLanguagesCache(db *sql.DB) error {
	retrieveSetTranslationLanguagesQuery, err := db.Prepare(`SELECT set_translation_language_id,
		set_translation_language FROM set_translation_languages`)
	defer retrieveSetTranslationLanguagesQuery.Close()
	if err != nil {
		return err
	}

	setTranslationLanguagesRows, err := retrieveSetTranslationLanguagesQuery.Query()
	if err != nil {
		return err
	}
	defer setTranslationLanguagesRows.Close()

	for setTranslationLanguagesRows.Next() {
		if err := setTranslationLanguagesRows.Err(); err != nil {
			return err
		}

		var setTranslationLanguageId int64
		var setTranslationLanguage string
		err := setTranslationLanguagesRows.Scan(&setTranslationLanguageId, &setTranslationLanguage)
		if err != nil {
			return err
		}
		setTranslationLanguages.Store(setTranslationLanguage, setTranslationLanguageId)
	}

	return nil
}

func populateBaseTypeOptionsCache(db *sql.DB) error {
	retrieveBaseTypeOptionsQuery, err := db.Prepare(`SELECT base_type_option_id,
		base_type_option FROM base_type_options`)
	defer retrieveBaseTypeOptionsQuery.Close()
	if err != nil {
		return err
	}

	baseTypeOptionsRows, err := retrieveBaseTypeOptionsQuery.Query()
	if err != nil {
		return err
	}
	defer baseTypeOptionsRows.Close()

	for baseTypeOptionsRows.Next() {
		if err := baseTypeOptionsRows.Err(); err != nil {
			return err
		}

		var baseTypeOptionId int64
		var baseTypeOption string
		err := baseTypeOptionsRows.Scan(&baseTypeOptionId, &baseTypeOption)
		if err != nil {
			return err
		}
		baseTypeOptions.Store(baseTypeOption, baseTypeOptionId)
	}

	return nil
}

func populateLeadershipFormatsCache(db *sql.DB) error {
	retrieveLeadershipFormatsQuery, err := db.Prepare(`SELECT leadership_format_id,
		leadership_format_name FROM leadership_formats`)
	defer retrieveLeadershipFormatsQuery.Close()

	leadershipFormatsRows, err := retrieveLeadershipFormatsQuery.Query()
	if err != nil {
		return  err
	}
	defer leadershipFormatsRows.Close()

	for leadershipFormatsRows.Next() {
		if err := leadershipFormatsRows.Err(); err != nil {
			return err
		}

		var leadershipFormatId int64
		var leadershipFormatName string
		err := leadershipFormatsRows.Scan(&leadershipFormatId, &leadershipFormatName)
		if err != nil {
			return err
		}
		leadershipFormats.Store(leadershipFormatName, leadershipFormatId)
	}

	return nil
}

func (card *MTGCard) GetAtomicPropertiesId(numPropsQuery *sql.Stmt, propIdQuery *sql.Stmt,
		atomicPropertiesHash string) (int64, bool, error) {
	// First, check how many entries are already in the db with this card hash
	// If it's 0, this atomic data isn't in the db, so we can return without getting the id
	// If it's 1, we can just return the retrieved ID
	// If it's more than 1, we have a hash collision, so we use the scryfall_oracle_id to disambiguate

	var count int
	countResult := numPropsQuery.QueryRow(atomicPropertiesHash)
	if err := countResult.Scan(&count); err != nil {
		return 0, false, err
	}

	if count == 0 {
		return 0, false, nil
	}

	// Since count is at least 1, we need to query the actual ID
	var atomicPropertiesId int64
	var scryfallOracleId string
	if count == 1 {
		// Only need to query the Id
		idResult := propIdQuery.QueryRow(atomicPropertiesHash)
		if err := idResult.Scan(&atomicPropertiesId, &scryfallOracleId); err != nil {
			return 0, false, err
		}
		return atomicPropertiesId, true, nil
	} else {
		// Hash collision, so need to iterate and check the scryfall_oracle_id
		results, err := propIdQuery.Query(atomicPropertiesHash)
		if err != nil {
			return 0, false, err
		}
		defer results.Close()
		for results.Next() {
			if err := results.Err(); err != nil {
				return 0, false, err
			}
			if err := results.Scan(&atomicPropertiesId, &scryfallOracleId); err != nil {
				return 0, false, err
			}
			if card.ScryfallOracleId == scryfallOracleId {
				return atomicPropertiesId, true, nil
			}
		}

		// We shouldn't get here, since it means there are multiple entries with the correct
		// hash, but none that match the scryfall_oracle_id, so return an error
		return 0, false, fmt.Errorf("Multiple atomic data with proper hash, but no matches")
	}
}
