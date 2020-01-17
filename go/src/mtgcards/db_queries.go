package mtgcards

import "database/sql"
import "fmt"
import "strings"

var atomicPropertiesIdQuery *sql.Stmt
var insertAtomicPropertiesQuery *sql.Stmt
var numAtomicPropertiesQuery *sql.Stmt

func CreateDbQueries(db *sql.DB) error {
	var err error
	numAtomicPropertiesQuery, err = db.Prepare(`SELECT COUNT(scryfall_oracle_id)
		FROM atomic_card_data
		WHERE card_data_hash = ?`)
	if err != nil {
		return err
	}

	atomicPropertiesIdQuery, err = db.Prepare(`SELECT atomic_card_data_id, scryfall_oracle_id
		FROM atomic_card_data
		WHERE card_data_hash = ?`)
	if err != nil {
		return err
	}

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
}

func (card MTGCard) InsertAtomicPropertiesToDb(atomicPropertiesHash string) (int64, error) {
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

	if err != nil {
		return 0, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	if rowsAffected != 1 {
		return 0, fmt.Errorf("Insert of atomic card data affected an unexpected num rows: %d", rowsAffected)
	}

	lastInsertId, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return lastInsertId, nil
}

func GetGameFormats(db *sql.DB) (map[string]int, error) {
	retrieveGameFormatsQuery, err := db.Prepare(`SELECT game_format_id, game_format_name
		FROM game_formats`)
	defer retrieveGameFormatsQuery.Close()
	if err != nil {
		return nil, err
	}

	gameFormats := make(map[string]int)

	gameFormatRows, err := retrieveGameFormatsQuery.Query()
	if err != nil {
		return nil, err
	}
	defer gameFormatRows.Close()

	for gameFormatRows.Next() {
		if err := gameFormatRows.Err(); err != nil {
			return nil, err
		}

		var gameFormatId int
		var gameFormatName string
		err := gameFormatRows.Scan(&gameFormatId, &gameFormatName)
		if err != nil {
			return nil, err
		}
		gameFormats[gameFormatName] = gameFormatId
	}

	return gameFormats, nil
}

func GetLegalityOptions(db *sql.DB) (map[string]int, error) {
	retrieveLegalityOptionsQuery, err := db.Prepare(`SELECT legality_option_id, legality_option_name
		FROM legality_options`)
	defer retrieveLegalityOptionsQuery.Close()
	if err != nil {
		return nil, err
	}

	legalityOptions := make(map[string]int)

	legalityOptionsRows, err := retrieveLegalityOptionsQuery.Query()
	if err != nil {
		return nil, err
	}
	defer legalityOptionsRows.Close()

	for legalityOptionsRows.Next() {
		if err := legalityOptionsRows.Err(); err != nil {
			return nil, err
		}

		var legalityOptionId int
		var legalityOptionName string
		err := legalityOptionsRows.Scan(&legalityOptionId, &legalityOptionName)
		if err != nil {
			return nil, err
		}
		legalityOptions[legalityOptionName] = legalityOptionId
	}

	return legalityOptions, nil
}

func GetPurchaseSites(db *sql.DB) (map[string]int, error) {
	retrievePurchaseSitesQuery, err := db.Prepare(`SELECT purchase_site_id, purchase_site_name
		FROM purchase_sites`)
	defer retrievePurchaseSitesQuery.Close()
	if err != nil {
		return nil, err
	}

	purchaseSites := make(map[string]int)

	purchaseSitesRows, err := retrievePurchaseSitesQuery.Query()
	if err != nil {
		return nil, err
	}
	defer purchaseSitesRows.Close()

	for purchaseSitesRows.Next() {
		if err := purchaseSitesRows.Err(); err != nil {
			return nil, err
		}

		var purchaseSiteId int
		var purchaseSiteName string
		err := purchaseSitesRows.Scan(&purchaseSiteId, &purchaseSiteName)
		if err != nil {
			return nil, err
		}
		purchaseSites[purchaseSiteName] = purchaseSiteId
	}

	return purchaseSites, nil
}

func GetSetTranslationLanguages(db *sql.DB) (map[string]int, error) {
	retrieveSetTranslationLanguagesQuery, err := db.Prepare(`SELECT set_translation_language_id,
		set_translation_language FROM set_translation_languages`)
	defer retrieveSetTranslationLanguagesQuery.Close()
	if err != nil {
		return nil, err
	}

	setTranslationLanguages := make(map[string]int)

	setTranslationLanguagesRows, err := retrieveSetTranslationLanguagesQuery.Query()
	if err != nil {
		return nil, err
	}
	defer setTranslationLanguagesRows.Close()

	for setTranslationLanguagesRows.Next() {
		if err := setTranslationLanguagesRows.Err(); err != nil {
			return nil, err
		}

		var setTranslationLanguageId int
		var setTranslationLanguage string
		err := setTranslationLanguagesRows.Scan(&setTranslationLanguageId, &setTranslationLanguage)
		if err != nil {
			return nil, err
		}
		setTranslationLanguages[setTranslationLanguage] = setTranslationLanguageId
	}

	return setTranslationLanguages, nil
}

func GetAtomicPropertiesId(atomicPropertiesHash string, card *MTGCard) (int64, bool, error) {
	// First, check how many entries are already in the db with this card hash
	// If it's 0, this atomic data isn't in the db, so we can return without getting the id
	// If it's 1, we can just return the retrieved ID
	// If it's more than 1, we have a hash collision, so we use the scryfall_oracle_id to disambiguate

	var count int
	countResult := numAtomicPropertiesQuery.QueryRow(atomicPropertiesHash)
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
		idResult := atomicPropertiesIdQuery.QueryRow(atomicPropertiesHash)
		if err := idResult.Scan(&atomicPropertiesId, &scryfallOracleId); err != nil {
			return 0, false, err
		}
		return atomicPropertiesId, true, nil
	} else {
		// Hash collision, so need to iterate and check the scryfall_oracle_id
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
