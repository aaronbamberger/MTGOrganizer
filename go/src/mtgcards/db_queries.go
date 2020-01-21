package mtgcards

import "database/sql"
import "encoding/hex"
import "fmt"
import "hash"
import "strings"
import "sync"

var atomicPropertiesIdQuery *sql.Stmt
var insertAtomicPropertiesQuery *sql.Stmt
var numAtomicPropertiesQuery *sql.Stmt
var insertCardQuery *sql.Stmt
var insertAltLangDataQuery *sql.Stmt
var insertCardPrintingQuery *sql.Stmt
var insertCardSubtypeQuery *sql.Stmt
var insertCardSupertypeQuery *sql.Stmt
var insertFrameEffectQuery *sql.Stmt
var insertLeadershipSkillQuery *sql.Stmt
var insertLegalityQuery *sql.Stmt
var insertOtherFaceIdQuery *sql.Stmt
var insertPurchaseUrlQuery *sql.Stmt
var insertRulingQuery *sql.Stmt
var insertSetTranslationQuery *sql.Stmt
var insertVariationQuery *sql.Stmt
var setHashQuery *sql.Stmt
var insertSetQuery *sql.Stmt
var insertGameFormatQuery *sql.Stmt
var insertLegalityOptionQuery *sql.Stmt
var insertPurchaseSiteQuery *sql.Stmt
var insertLeadershipFormatQuery *sql.Stmt
var insertSetTranslationLanguageQuery *sql.Stmt

var gameFormats sync.Map
var legalityOptions sync.Map
var purchaseSites sync.Map
var leadershipFormats sync.Map
var setTranslationLanguages sync.Map
//var frameEffectOptions sync.Map

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
	fmt.Printf("Prepare query 1\n")
	numAtomicPropertiesQuery, err = db.Prepare(`SELECT COUNT(scryfall_oracle_id)
		FROM atomic_card_data
		WHERE card_data_hash = ?`)
	if err != nil {
		return err
	}

	fmt.Printf("Prepare query 2\n")
	atomicPropertiesIdQuery, err = db.Prepare(`SELECT atomic_card_data_id, scryfall_oracle_id
		FROM atomic_card_data
		WHERE card_data_hash = ?`)
	if err != nil {
		return err
	}

	fmt.Printf("Prepare query 3\n")
	insertAtomicPropertiesQuery, err = db.Prepare(`INSERT INTO atomic_card_data
		(card_data_hash, color_identity, color_indicator, colors, converted_mana_cost,
		edhrec_rank, face_converted_mana_cost, hand, is_reserved, layout, life,
		loyalty, mana_cost, mtgstocks_id, name, card_power, scryfall_oracle_id,
		side, text, toughness, card_type)
		VALUES
		(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return err
	}

	fmt.Printf("Prepare query 4\n")
	insertCardQuery, err = db.Prepare(`INSERT INTO all_cards
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
		return err
	}

	fmt.Printf("Prepare query 5\n")
	insertAltLangDataQuery, err = db.Prepare(`INSERT INTO alternate_language_data
		(atomic_card_data_id, flavor_text, language, multiverse_id, name, text, card_type)
		VALUES
		(?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return err
	}

	fmt.Printf("Prepare query 6\n")
	insertCardPrintingQuery, err = db.Prepare(`INSERT INTO card_printings
		(atomic_card_data_id, set_code)
		VALUES
		(?, ?)`)
	if err != nil {
		return err
	}

	fmt.Printf("Prepare query 7\n")
	insertCardSubtypeQuery, err = db.Prepare(`INSERT INTO card_subtypes
		(atomic_card_data_id, card_subtype)
		VALUES
		(?, ?)`)
	if err != nil {
		return err
	}

	fmt.Printf("Prepare query 8\n")
	insertCardSupertypeQuery, err = db.Prepare(`INSERT INTO card_supertypes
		(atomic_card_data_id, card_supertype)
		VALUES
		(?, ?)`)
	if err != nil {
		return err
	}

	fmt.Printf("Prepare query 9\n")
	insertFrameEffectQuery, err = db.Prepare(`INSERT INTO frame_effects
		(card_uuid, frame_effect)
		VALUES
		(?, ?)`)
	if err != nil {
		return err
	}

	fmt.Printf("Prepare query 10\n")
	insertLeadershipSkillQuery, err = db.Prepare(`INSERT INTO leadership_skills
		(atomic_card_data_id, leadership_format_id, leader_legal)
		VALUES
		(?, ?, ?)`)
	if err != nil {
		return err
	}

	fmt.Printf("Prepare query 11\n")
	insertLegalityQuery, err = db.Prepare(`INSERT INTO legalities
		(atomic_card_data_id, game_format_id, legality_option_id)
		VALUES
		(?, ?, ?)`)
	if err != nil {
		return err
	}

	fmt.Printf("Prepare query 12\n")
	insertOtherFaceIdQuery, err = db.Prepare(`INSERT INTO other_faces
		(card_uuid, other_face_uuid)
		VALUES
		(?, ?)`)
	if err != nil {
		return err
	}

	fmt.Printf("Prepare query 13\n")
	insertPurchaseUrlQuery, err = db.Prepare(`INSERT INTO purchase_urls
		(atomic_card_data_id, purchase_site_id, purchase_url)
		VALUES
		(?, ?, ?)`)
	if err != nil {
		return err
	}

	fmt.Printf("Prepare query 14\n")
	insertRulingQuery, err = db.Prepare(`INSERT INTO rulings
		(atomic_card_data_id, ruling_date, ruling_text)
		VALUES
		(?, ?, ?)`)
	if err != nil {
		return err
	}

	fmt.Printf("Prepare query 15\n")
	insertSetTranslationQuery, err = db.Prepare(`INSERT INTO set_translations
		(set_id, set_translation_language_id, set_translated_name)
		VALUES
		(?, ?, ?)`)
	if err != nil {
		return err
	}

	fmt.Printf("Prepare query 16\n")
	insertVariationQuery, err = db.Prepare(`INSERT INTO variations
		(card_uuid, variation_uuid)
		VALUES
		(?, ?)`)
	if err != nil {
		return err
	}

	fmt.Printf("Prepare query 17\n")
	setHashQuery, err = db.Prepare("SELECT set_hash FROM sets WHERE code = ?")
	if err != nil {
		return err
	}

	fmt.Printf("Prepare query 18\n")
	insertSetQuery, err = db.Prepare(`INSERT INTO sets
		(set_hash, base_size, block_name, code, is_foreign_only, is_foil_only,
		is_online_only, is_partial_preview, keyrune_code, mcm_name, mcm_id,
		mtgo_code, name, parent_code, release_date, tcgplayer_group_id,
		total_set_size, set_type)
		VALUES
		(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return err
	}

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

	return nil
}

func CloseDbQueries() {
	if numAtomicPropertiesQuery != nil {
		numAtomicPropertiesQuery.Close()
	}

	if atomicPropertiesIdQuery != nil {
		atomicPropertiesIdQuery.Close()
	}

	if insertAtomicPropertiesQuery != nil {
		insertAtomicPropertiesQuery.Close()
	}

	if insertCardQuery != nil {
		insertCardQuery.Close()
	}

	if insertAltLangDataQuery != nil {
		insertAltLangDataQuery.Close()
	}

	if insertCardPrintingQuery != nil {
		insertCardPrintingQuery.Close()
	}

	if insertCardSubtypeQuery != nil {
		insertCardSubtypeQuery.Close()
	}

	if insertCardSupertypeQuery != nil {
		insertCardSupertypeQuery.Close()
	}

	if insertFrameEffectQuery != nil {
		insertFrameEffectQuery.Close()
	}

	if insertLeadershipSkillQuery != nil {
		insertLeadershipSkillQuery.Close()
	}

	if insertLegalityQuery != nil {
		insertLegalityQuery.Close()
	}

	if insertOtherFaceIdQuery != nil {
		insertOtherFaceIdQuery.Close()
	}

	if insertPurchaseUrlQuery != nil {
		insertPurchaseUrlQuery.Close()
	}

	if insertRulingQuery != nil {
		insertRulingQuery.Close()
	}

	if insertSetTranslationQuery != nil {
		insertSetTranslationQuery.Close()
	}

	if insertVariationQuery != nil {
		insertVariationQuery.Close()
	}

	if setHashQuery != nil {
		setHashQuery.Close()
	}

	if insertSetQuery != nil {
		insertSetQuery.Close()
	}

	if insertGameFormatQuery != nil {
		insertGameFormatQuery.Close()
	}

	if insertLegalityOptionQuery != nil {
		insertLegalityQuery.Close()
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
}

func (set *MTGSet) CheckIfSetExists(lock *sync.RWMutex) (bool, string, error) {
	lock.Lock()
	defer lock.Unlock()

	// First, check to see if this set is in the DB at all
	setRows, err := setHashQuery.Query(set.Code)
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

func (set *MTGSet) InsertSetToDb(setHash string, lock *sync.RWMutex) error {
	lock.RLock()
	res, err := insertSetQuery.Exec(setHash, set.BaseSetSize, set.Block, set.Code, set.IsForeignOnly,
		set.IsFoilOnly, set.IsOnlineOnly, set.IsPartialPreview, set.KeyruneCode, set.MCMName,
		set.MCMId, set.MTGOCode, set.Name, set.ParentCode, set.ReleaseDate, set.TCGPlayerGroupId,
		set.TotalSetSize, set.Type)
	lock.RUnlock()
	if err != nil {
		return err
	}

	return checkRowsAffected(res, 1, "insert set")
}


func (card *MTGCard) InsertFrameEffectToDb(frameEffect string, lock *sync.RWMutex) error {
	lock.RLock()
	res, err := insertFrameEffectQuery.Exec(card.UUID, frameEffect)
	lock.RUnlock()
	if err != nil {
		return err
	}

	return checkRowsAffected(res, 1, "insert frame effect")
}

func InsertLeadershipSkillToDb(atomicPropertiesId int64, leadershipFormat string, leaderLegal bool,
		lock *sync.RWMutex) error {
	// Get the leadership format id from the cache
	var leadershipFormatId int64
	leadershipFormatIdTemp, loaded := leadershipFormats.Load(leadershipFormat)
	if !loaded {
		// This is the unlikely case where we have a new value that isn't pre-populated in the db
		lock.RLock()
		res, err := insertLeadershipFormatQuery.Exec(leadershipFormat)
		lock.RUnlock()
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


	lock.RLock()
	res, err := insertLeadershipSkillQuery.Exec(atomicPropertiesId, leadershipFormatId, leaderLegal)
	lock.RUnlock()
	if err != nil {
		return err
	}

	return checkRowsAffected(res, 1, "insert leadership skill")
}

func InsertLegalityToDb(atomicPropertiesId int64, gameFormat string, legalityOption string,
		lock *sync.RWMutex) error {
	// Get the game format id from the cache
	var gameFormatId int64
	gameFormatIdTemp, loaded := gameFormats.Load(gameFormat)
	if !loaded {
		// This is the unlikely case where we have a new value that isn't pre-populated in the db
		lock.RLock()
		res, err := insertGameFormatQuery.Exec(gameFormat)
		lock.RUnlock()
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
		lock.RLock()
		res, err := insertLegalityOptionQuery.Exec(legalityOption)
		lock.RUnlock()
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

	lock.RLock()
	res, err := insertLegalityQuery.Exec(atomicPropertiesId, gameFormatId, legalityOptionId)
	lock.RUnlock()
	if err != nil {
		return err
	}

	return checkRowsAffected(res, 1, "insert legality")
}

func InsertCardPrintingToDb(atomicPropertiesId int64, setCode string, lock *sync.RWMutex) error {
	lock.RLock()
	res, err := insertCardPrintingQuery.Exec(atomicPropertiesId, setCode)
	lock.RUnlock()
	if err != nil {
		return err
	}

	return checkRowsAffected(res, 1, "insert card printing")
}

func InsertPurchaseURLToDb(atomicPropertiesId int64, purchaseSite string, purchaseURL string,
		lock *sync.RWMutex) error {
	// Get the purchase site id from the cache
	var purchaseSiteId int64
	purchaseSiteIdTemp, loaded := purchaseSites.Load(purchaseSite)
	if !loaded {
		// This is the unlikely case where we have a new value that isn't pre-populated in the db
		lock.RLock()
		res, err := insertPurchaseSiteQuery.Exec(purchaseSite)
		lock.RUnlock()
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

	lock.RLock()
	res, err := insertPurchaseUrlQuery.Exec(atomicPropertiesId, purchaseSiteId, purchaseURL)
	lock.RUnlock()
	if err != nil {
		return err
	}

	return checkRowsAffected(res, 1, "insert purchase url")
}

func (altLangInfo *MTGCardAlternateLanguageInfo) InsertAltLangDataToDb(atomicPropertiesId int64,
		lock *sync.RWMutex) error {
	lock.RLock()
	res, err := insertAltLangDataQuery.Exec(atomicPropertiesId, altLangInfo.FlavorText,
		altLangInfo.Language, altLangInfo.MultiverseId, altLangInfo.Name,
		altLangInfo.Text, altLangInfo.Type)
	lock.RUnlock()
	if err != nil {
		return err
	}

	return checkRowsAffected(res, 1, "insert alt lang info")
}

func (ruling *MTGCardRuling) InsertRulingToDb(atomicPropertiesId int64, lock *sync.RWMutex) error {
	lock.RLock()
	res, err := insertRulingQuery.Exec(atomicPropertiesId, ruling.Date, ruling.Text)
	lock.RUnlock()
	if err != nil {
		return err
	}

	return checkRowsAffected(res, 1, "insert ruling")
}

func InsertCardSubtypeToDb(atomicPropertiesId int64, subtype string, lock *sync.RWMutex) error {
	lock.RLock()
	res, err := insertCardSubtypeQuery.Exec(atomicPropertiesId, subtype)
	lock.RUnlock()
	if err != nil {
		return err
	}

	return checkRowsAffected(res, 1, "insert subtype")
}

func InsertCardSupertypeToDb(atomicPropertiesId int64, supertype string, lock *sync.RWMutex) error {
	lock.RLock()
	res, err := insertCardSupertypeQuery.Exec(atomicPropertiesId, supertype)
	lock.RUnlock()
	if err != nil {
		return err
	}

	return checkRowsAffected(res, 1, "insert supertype")
}

func InsertSetTranslationToDb(setId int, translationLangId int, translatedName string,
		lock *sync.RWMutex) error {
	lock.RLock()
	res, err := insertSetTranslationQuery.Exec(setId, translationLangId, translatedName)
	lock.RUnlock()
	if err != nil {
		return err
	}

	return checkRowsAffected(res, 1, "insert set name translation")
}

func (card *MTGCard) InsertOtherFaceIdToDb(otherFaceUUID string, lock *sync.RWMutex) error {
	lock.RLock()
	res, err := insertOtherFaceIdQuery.Exec(card.UUID, otherFaceUUID)
	lock.RUnlock()
	if err != nil {
		return err
	}

	return checkRowsAffected(res, 1, "insert other face ID")
}

func (card *MTGCard) InsertVariationToDb(variationUUID string, lock *sync.RWMutex) error {
	lock.RLock()
	res, err := insertVariationQuery.Exec(card.UUID, variationUUID)
	lock.RUnlock()
	if err != nil {
		return err
	}

	return checkRowsAffected(res, 1, "insert variation")
}

func (card *MTGCard) InsertCardToDb(atomicPropertiesId int64, lock *sync.RWMutex) error {
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

	lock.RLock()
	res, err := insertCardQuery.Exec(card.UUID, cardHash, atomicPropertiesId,
		card.Artist, card.BorderColor, card.Number, card.ScryfallId,
		card.Watermark, card.FrameVersion, card.MCMId, card.MCMMetaId,
		card.MultiverseId, card.OriginalText, card.OriginalType,
		card.Rarity, card.TCGPlayerProductId, duelDeck, flavorText,
		card.HasFoil, card.HasNonFoil, card.IsAlternative, card.IsArena,
		card.IsFullArt, card.IsMTGO, card.IsOnlineOnly, card.IsOversized,
		card.IsPaper, card.IsPromo, card.IsReprint, card.IsStarter,
		card.IsStorySpotlight, card.IsTextless, card.IsTimeshifted,
		mtgArenaId, mtgoFoilId, mtgoId, scryfallIllustrationId)
	lock.RUnlock()

	if err != nil {
		return err
	}

	return checkRowsAffected(res, 1, "insert card data")
}

func (card *MTGCard) InsertAtomicPropertiesToDb(atomicPropertiesHash string,
		lock *sync.RWMutex) (int64, error) {
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

	lock.RLock()
	res, err := insertAtomicPropertiesQuery.Exec(atomicPropertiesHash,
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
	lock.RUnlock()

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

func GetAtomicPropertiesId(atomicPropertiesHash string, card *MTGCard,
		lock *sync.RWMutex) (int64, bool, error) {
	// First, check how many entries are already in the db with this card hash
	// If it's 0, this atomic data isn't in the db, so we can return without getting the id
	// If it's 1, we can just return the retrieved ID
	// If it's more than 1, we have a hash collision, so we use the scryfall_oracle_id to disambiguate

	var count int
	lock.Lock()
	countResult := numAtomicPropertiesQuery.QueryRow(atomicPropertiesHash)
	lock.Unlock()
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
		lock.Lock()
		idResult := atomicPropertiesIdQuery.QueryRow(atomicPropertiesHash)
		lock.Unlock()
		if err := idResult.Scan(&atomicPropertiesId, &scryfallOracleId); err != nil {
			return 0, false, err
		}
		return atomicPropertiesId, true, nil
	} else {
		// Hash collision, so need to iterate and check the scryfall_oracle_id
		lock.Lock()
		defer lock.Unlock()
		results, err := atomicPropertiesIdQuery.Query(atomicPropertiesHash)
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
