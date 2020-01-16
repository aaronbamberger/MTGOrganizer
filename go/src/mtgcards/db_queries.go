package mtgcards

import "database/sql"
import "fmt"
import "strings"

var atomicPropertiesHashQuery *sql.Stmt
var insertAtomicPropertiesQuery *sql.Stmt

func CreateDbQueries(db *sql.DB) error {
	var err error
	atomicPropertiesHashQuery, err = db.Prepare(`SELECT card_data_hash FROM atomic_card_data
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
	if atomicPropertiesHashQuery != nil {
		atomicPropertiesHashQuery.Close()
	}
	if insertAtomicPropertiesQuery != nil {
		insertAtomicPropertiesQuery.Close()
	}
}

func (card MTGCard) InsertAtomicPropertiesToDb(atomicPropertiesHash string) error {
	// Build the set values needed for color_identity, color_indicator, and colors
	var colorIdentity string
	var colorIndicator string
	var colors string

	if len(card.ColorIdentity) > 0 {
		colorIdentity = "'" + strings.Join(card.ColorIdentity, ",") + "'"
	}

	if len(card.ColorIndicator) > 0 {
		colorIndicator = "'" + strings.Join(card.ColorIndicator, ",") + "'"
	}

	if len(card.Colors) > 0 {
		colors = "'" + strings.Join(card.Colors, ",") + "'"
	}

	fmt.Printf("For card %s\n", card.Name)
	fmt.Printf("Color Ident: %v\n", card.ColorIdentity)
	fmt.Printf("Color Ident String: %s\n", colorIdentity)
	fmt.Printf("Color Ind: %v\n", card.ColorIndicator)
	fmt.Printf("Color Ind String: %s\n", colorIndicator)
	fmt.Printf("Colors: %v\n", card.Colors)
	fmt.Printf("Colors String: %s\n", colors)

	// TODO: Figure out why inserts of set field aren't working
	res, err := insertAtomicPropertiesQuery.Exec(atomicPropertiesHash,
		"",
		"",
		"",
		card.ConvertedManaCost,
		card.EDHRecRank,
		card.FaceConvertedManaCost,
		card.Hand,
		card.IsReserved,
		card.Layout,
		card.Life,
		card.Loyalty,
		card.ManaCost,
		card.MTGStocksId,
		card.Name,
		card.Power,
		card.ScryfallOracleId,
		card.Side,
		card.Text,
		card.Toughness,
		card.Type)

	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected != 1 {
		fmt.Printf("Insert of atomic card data affected an unexpected number of rows: %d\n", rowsAffected)
	}

	return nil
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

func CheckAtomicPropertiesDataExistence(atomicPropertiesHash string) (bool, error) {
	resultRows, err := atomicPropertiesHashQuery.Query(atomicPropertiesHash)
	if err != nil {
		return false, err
	}

	defer resultRows.Close()

	// Only need to check if we've returned a row, if we have, we already
	// know the hash
	if resultRows.Next() {
		return true, nil
	} else {
		return false, resultRows.Err()
	}
}
