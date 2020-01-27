package mtgcards

import "database/sql"
import "fmt"
import "strings"

func (card *MTGCard) InsertAtomicPropertiesToDb(queries *cardDbQueries,
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

	res, err := queries.InsertAtomicPropertiesQuery.Exec(atomicPropertiesHash,
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

	atomicPropId, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

    // Now, insert all of the atomic properties that live in other tables

	// Alternate language data
	for _, altLangData := range card.AlternateLanguageData {
		err = altLangData.InsertAltLangDataToDb(queries, atomicPropId)
		if err != nil {
			return 0, err
		}
	}

	// Leadership skills
	for leadershipFormat, leaderValid := range card.LeadershipSkills {
		err = InsertLeadershipSkillToDb(queries, atomicPropId, leadershipFormat, leaderValid)
		if err != nil {
			return 0, err
		}
	}

	// Legalities
	for format, legality := range card.Legalities {
		err = InsertLegalityToDb(queries, atomicPropId, format, legality)
		if err != nil {
			return 0, err
		}
	}

    // Printings
	for _, setCode := range card.Printings {
		err = InsertCardPrintingToDb(queries, atomicPropId, setCode)
		if err != nil {
			return 0, err
		}
	}

	// Purchase URLs
	for site, url := range card.PurchaseURLs {
		err = InsertPurchaseURLToDb(queries, atomicPropId, site, url)
		if err != nil {
			return 0, err
		}
	}

	// Rulings
	for _, ruling := range card.Rulings {
		err = ruling.InsertRulingToDb(queries, atomicPropId)
		if err != nil {
			return 0, err
		}
	}

	// Subtypes
	for _, subtype := range card.Subtypes {
		err = InsertCardSubtypeToDb(queries, atomicPropId, subtype)
		if err != nil {
			return 0, err
		}
	}

	// Supertypes
	for _, supertype := range card.Supertypes {
		err = InsertCardSupertypeToDb(queries, atomicPropId, supertype)
		if err != nil {
			return 0, err
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
		err = InsertBaseTypeToDb(queries, atomicPropId, baseType)
		if err != nil {
			return 0, err
		}
	}

	return atomicPropId, nil
}

func InsertCardSubtypeToDb(queries *cardDbQueries, atomicPropertiesId int64, subtype string) error {
	subtypeId, err := getSubtypeOptionId(subtype)
	if err != nil {
		return err
	}

	res, err := queries.InsertCardSubtypeQuery.Exec(atomicPropertiesId, subtypeId)
	if err != nil {
		return err
	}

	return checkRowsAffected(res, 1, "insert card subtype")
}

func InsertCardSupertypeToDb(queries *cardDbQueries, atomicPropertiesId int64, supertype string) error {
	supertypeId, err := getSupertypeOptionId(supertype)
	if err != nil {
		return err
	}

	res, err := queries.InsertCardSupertypeQuery.Exec(atomicPropertiesId, supertypeId)
	if err != nil {
		return err
	}

	return checkRowsAffected(res, 1, "insert card supertype")
}

func (altLangInfo *MTGCardAlternateLanguageInfo) InsertAltLangDataToDb(queries *cardDbQueries,
		atomicPropertiesId int64) error {
	res, err := queries.InsertAltLangDataQuery.Exec(atomicPropertiesId, altLangInfo.FlavorText,
		altLangInfo.Language, altLangInfo.MultiverseId, altLangInfo.Name,
		altLangInfo.Text, altLangInfo.Type)
	if err != nil {
		return err
	}

	return checkRowsAffected(res, 1, "insert alt lang info")
}

func (ruling *MTGCardRuling) InsertRulingToDb(queries *cardDbQueries, atomicPropertiesId int64) error {
	res, err := queries.InsertRulingQuery.Exec(atomicPropertiesId, ruling.Date, ruling.Text)
	if err != nil {
		return err
	}

	return checkRowsAffected(res, 1, "insert ruling")
}

func InsertBaseTypeToDb(queries *cardDbQueries, atomicPropertiesId int64,
		baseTypeOption string) error {
	baseTypeOptionId, err := getBaseTypeOptionId(baseTypeOption)
	if err != nil {
		return err
	}

	res, err := queries.InsertBaseTypeQuery.Exec(atomicPropertiesId, baseTypeOptionId)
	if err != nil {
		return err
	}

	return checkRowsAffected(res, 1, "insert base type")
}

func InsertLeadershipSkillToDb(queries *cardDbQueries, atomicPropertiesId int64,
		leadershipFormat string, leaderLegal bool) error {

	leadershipFormatId, err := getLeadershipFormatId(leadershipFormat)
	if err != nil {
		return err
	}

	res, err := queries.InsertLeadershipSkillQuery.Exec(atomicPropertiesId,
        leadershipFormatId, leaderLegal)
	if err != nil {
		return err
	}

	return checkRowsAffected(res, 1, "insert leadership skill")
}

func InsertLegalityToDb(queries *cardDbQueries, atomicPropertiesId int64, gameFormat string,
		legalityOption string) error {
	gameFormatId, err := getGameFormatId(gameFormat)
	if err != nil {
		return err
	}

	legalityOptionId, err := getLegalityOptionId(legalityOption)
	if err != nil {
		return err
	}

	res, err := queries.InsertLegalityQuery.Exec(atomicPropertiesId, gameFormatId, legalityOptionId)
	if err != nil {
		return err
	}

	return checkRowsAffected(res, 1, "insert legality")
}

func InsertCardPrintingToDb(queries *cardDbQueries, atomicPropertiesId int64, setCode string) error {
	res, err := queries.InsertCardPrintingQuery.Exec(atomicPropertiesId, setCode)
	if err != nil {
		return err
	}

	return checkRowsAffected(res, 1, "insert card printing")
}

func InsertPurchaseURLToDb(queries *cardDbQueries, atomicPropertiesId int64,
		purchaseSite string, purchaseURL string) error {
	purchaseSiteId, err := getPurchaseSiteId(purchaseSite)
	if err != nil {
		return err
	}

	res, err := queries.InsertPurchaseURLQuery.Exec(atomicPropertiesId, purchaseSiteId, purchaseURL)
	if err != nil {
		return err
	}

	return checkRowsAffected(res, 1, "insert purchase url")
}

func (card *MTGCard) GetAtomicPropertiesId(queries *cardDbQueries,
        atomicPropertiesHash string) (int64, int64, bool, error) {
	// First, check how many entries are already in the db with this card hash
	// If it's 0, this atomic data isn't in the db, so we can return without getting the id
	// If it's 1, we can just return the retrieved ID
	// If it's more than 1, we have a hash collision, so we use the scryfall_oracle_id to disambiguate

	var count int
	countResult := queries.NumAtomicPropertiesQuery.QueryRow(atomicPropertiesHash)
	if err := countResult.Scan(&count); err != nil {
		return 0, 0, false, err
	}

	if count == 0 {
		return 0, 0, false, nil
	}

	// Since count is at least 1, we need to query the actual ID
	var atomicPropertiesId int64
    var refCnt int64
	var scryfallOracleId string
	if count == 1 {
		// Only need to query the Id
		idResult := queries.AtomicPropertiesIdQuery.QueryRow(atomicPropertiesHash)
		if err := idResult.Scan(&atomicPropertiesId, &refCnt, &scryfallOracleId); err != nil {
			return 0, 0, false, err
		}
		return atomicPropertiesId, refCnt, true, nil
	} else {
		// Hash collision, so need to iterate and check the scryfall_oracle_id
		results, err := queries.AtomicPropertiesIdQuery.Query(atomicPropertiesHash)
		if err != nil {
			return 0, 0, false, err
		}
		defer results.Close()
		for results.Next() {
			if err := results.Err(); err != nil {
				return 0, 0, false, err
			}
			if err := results.Scan(&atomicPropertiesId, &refCnt, &scryfallOracleId); err != nil {
				return 0, 0, false, err
			}
			if card.ScryfallOracleId == scryfallOracleId {
				return atomicPropertiesId, refCnt, true, nil
			}
		}

		// We shouldn't get here, since it means there are multiple entries with the correct
		// hash, but none that match the scryfall_oracle_id, so return an error
		return 0, 0, false, fmt.Errorf("Multiple atomic data with proper hash, but no matches")
	}
}
