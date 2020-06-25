package carddb

import "database/sql"
import "mtgcards"
import "strings"

func GetCardHashAndIdFromDB(
        cardUUID string,
        queries *DBGetQueries) (bool, string, int64, error) {
	res := queries.GetCardHashQuery.QueryRow(cardUUID)

	var cardHash string
	var cardId int64
	err := res.Scan(&cardHash, &cardId)
	if err != nil {
		if err == sql.ErrNoRows {
			// This card isn't in the database
			return false, "", 0, nil
		} else {
			return false, "", 0, err
		}
	} else {
		return true, cardHash, cardId, nil
	}
}

func InsertCardToDB(
        card *mtgcards.MTGCard,
        setId int64,
        queries *DBInsertQueries) error {
	// Build all of the values that can be null
    var asciiName sql.NullString
    var colorIdentity sql.NullString
	var colorIndicator sql.NullString
	var colors sql.NullString
    var duelDeck sql.NullString
	var edhrecRank sql.NullInt32
    var flavorName sql.NullString
	var flavorText sql.NullString
	var hand sql.NullString
    var life sql.NullString
	var loyalty sql.NullString
	var mtgArenaId sql.NullInt32
	var mtgoFoilId sql.NullInt32
	var mtgoId sql.NullInt32
    var name sql.NullString
	var scryfallIllustrationId sql.NullString
	var side sql.NullString

    if len(card.AsciiName) > 0 {
        asciiName.String = card.AsciiName
        asciiName.Valid = true
    }

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

    if len(card.DuelDeck) > 0 {
		duelDeck.String = card.DuelDeck
		duelDeck.Valid = true
	}

	if card.EDHRecRank != 0 {
		edhrecRank.Int32 = int32(card.EDHRecRank)
		edhrecRank.Valid = true
	}

    if len(card.FlavorName) > 0 {
        flavorName.String = card.FlavorName
        flavorName.Valid = true
    }

    if len(card.FlavorText) > 0 {
		flavorText.String = card.FlavorText
		flavorText.Valid = true
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

	if len(card.Name) > 0 {
		name.String = card.Name
		name.Valid = true
	}

	if len(card.ScryfallIllustrationId) > 0 {
		scryfallIllustrationId.String = card.ScryfallIllustrationId
		scryfallIllustrationId.Valid = true
	}

	if len(card.Side) > 0 {
		side.String = card.Side
		side.Valid = true
	}

	res, err := queries.InsertCardQuery.Exec(
        card.UUID,
        card.Hash(),
        card.Artist,
        asciiName,
        card.BorderColor,
        card.Number,
        card.Power,
        card.Type,
		colorIdentity,
		colorIndicator,
		colors,
		card.ConvertedManaCost,
        duelDeck,
		edhrecRank,
		card.FaceConvertedManaCost,
        flavorName,
        flavorText,
        card.FrameVersion,
		hand,
        card.HasFoil,
        card.HasNonFoil,
        card.IsAlternative,
        card.IsArena,
        card.IsBuyABox,
        card.IsDateStamped,
        card.IsFullArt,
        card.IsMTGO,
        card.IsOnlineOnly,
        card.IsOversized,
        card.IsPaper,
        card.IsPromo,
        card.IsReprint,
		card.IsReserved,
        card.IsStarter,
        card.IsStorySpotlight,
        card.IsTextless,
        card.IsTimeshifted,
		card.Layout,
		life,
		loyalty,
		card.ManaCost,
        card.MCMId,
        card.MCMMetaId,
        mtgArenaId,
        mtgoFoilId,
        mtgoId,
		card.MTGStocksId,
        card.MultiverseId,
		name,
        card.OriginalText,
        card.OriginalType,
        card.Rarity,
        card.ScryfallId,
        scryfallIllustrationId,
		card.ScryfallOracleId,
        setId,
		side,
        card.TCGPlayerProductId,
		card.Text,
		card.Toughness,
		card.Watermark)

	if err != nil {
		return err
	}

	cardId, err := res.LastInsertId()
	if err != nil {
		return err
	}

    // Now, insert all of card data that doesn't live in the all_cards table
    err = InsertOtherCardDataToDB(cardId, card, queries)
    if err != nil {
        return nil
    }

	return nil
}

func UpdateCardInDB(
        cardId int64,
        setId int64,
        card *mtgcards.MTGCard,
        updateQueries *DBUpdateQueries,
        deleteQueries *DBDeleteQueries,
        insertQueries *DBInsertQueries) error {
    // Build all of the values that can be null
    var asciiName sql.NullString
    var colorIdentity sql.NullString
	var colorIndicator sql.NullString
	var colors sql.NullString
    var duelDeck sql.NullString
	var edhrecRank sql.NullInt32
    var flavorName sql.NullString
	var flavorText sql.NullString
	var hand sql.NullString
    var life sql.NullString
	var loyalty sql.NullString
	var mtgArenaId sql.NullInt32
	var mtgoFoilId sql.NullInt32
	var mtgoId sql.NullInt32
    var name sql.NullString
	var scryfallIllustrationId sql.NullString
	var side sql.NullString

    if len(card.AsciiName) > 0 {
        asciiName.String = card.AsciiName
        asciiName.Valid = true
    }

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

    if len(card.DuelDeck) > 0 {
		duelDeck.String = card.DuelDeck
		duelDeck.Valid = true
	}

	if card.EDHRecRank != 0 {
		edhrecRank.Int32 = int32(card.EDHRecRank)
		edhrecRank.Valid = true
	}

    if len(card.FlavorName) > 0 {
        flavorName.String = card.FlavorName
        flavorName.Valid = true
    }

    if len(card.FlavorText) > 0 {
		flavorText.String = card.FlavorText
		flavorText.Valid = true
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

	if len(card.Name) > 0 {
		name.String = card.Name
		name.Valid = true
	}

	if len(card.ScryfallIllustrationId) > 0 {
		scryfallIllustrationId.String = card.ScryfallIllustrationId
		scryfallIllustrationId.Valid = true
	}

	if len(card.Side) > 0 {
		side.String = card.Side
		side.Valid = true
	}

    // First, update the main card record
	_, err := updateQueries.UpdateCardQuery.Exec(
        card.Hash(),
        card.Artist,
        asciiName,
        card.BorderColor,
        card.Number,
        card.Power,
        card.Type,
		colorIdentity,
		colorIndicator,
		colors,
		card.ConvertedManaCost,
        duelDeck,
		edhrecRank,
		card.FaceConvertedManaCost,
        flavorName,
        flavorText,
        card.FrameVersion,
		hand,
        card.HasFoil,
        card.HasNonFoil,
        card.IsAlternative,
        card.IsArena,
        card.IsBuyABox,
        card.IsDateStamped,
        card.IsFullArt,
        card.IsMTGO,
        card.IsOnlineOnly,
        card.IsOversized,
        card.IsPaper,
        card.IsPromo,
        card.IsReprint,
		card.IsReserved,
        card.IsStarter,
        card.IsStorySpotlight,
        card.IsTextless,
        card.IsTimeshifted,
		card.Layout,
		life,
		loyalty,
		card.ManaCost,
        card.MCMId,
        card.MCMMetaId,
        mtgArenaId,
        mtgoFoilId,
        mtgoId,
		card.MTGStocksId,
        card.MultiverseId,
		name,
        card.OriginalText,
        card.OriginalType,
        card.Rarity,
        card.ScryfallId,
        scryfallIllustrationId,
		card.ScryfallOracleId,
        setId,
		side,
        card.TCGPlayerProductId,
		card.Text,
		card.Toughness,
		card.Watermark,
        card.UUID)
	if err != nil {
		return err
	}

    // Next, delete the rest of the old card data
    err = DeleteOtherCardDataFromDB(cardId, deleteQueries)
    if err != nil {
        return err
    }

    // Finally, insert the rest of the new card data
    err = InsertOtherCardDataToDB(cardId, card, insertQueries)
    if err != nil {
        return err
    }

    return nil
}

func InsertOtherCardDataToDB(
        cardId int64,
        card *mtgcards.MTGCard,
        queries *DBInsertQueries) error {
	// Alternate language data
	for _, altLangData := range card.AlternateLanguageData {
        err := InsertAltLangInfoToDB(cardId, &altLangData, queries)
		if err != nil {
			return err
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
        err := InsertBaseTypeToDB(cardId, baseType, queries)
		if err != nil {
			return err
		}
	}

    // Frame effects
	for _, frameEffect := range card.FrameEffects {
		err := InsertFrameEffectToDB(cardId, frameEffect, queries)
		if err != nil {
			return err
		}
	}

	// Leadership skills
	for leadershipFormat, leaderValid := range card.LeadershipSkills {
        err := InsertLeadershipSkillToDB(cardId, leadershipFormat, leaderValid, queries)
		if err != nil {
			return err
		}
	}

	// Legalities
	for format, legality := range card.Legalities {
        err := InsertLegalityToDB(cardId, format, legality, queries)
		if err != nil {
			return err
		}
	}

	// Other face IDs
	for _, otherFaceId := range card.OtherFaceIds {
		err := InsertOtherFaceIdToDB(cardId, otherFaceId, queries)
		if err != nil {
			return err
		}
	}

    // Printings
	for _, setCode := range card.Printings {
        err := InsertPrintingToDB(cardId, setCode, queries)
		if err != nil {
			return err
		}
	}

	// Purchase URLs
	for site, url := range card.PurchaseURLs {
        err := InsertPurchaseURLToDB(cardId, site, url, queries)
		if err != nil {
			return err
		}
	}

	// Rulings
	for _, ruling := range card.Rulings {
        err := InsertRulingToDB(cardId, &ruling, queries)
		if err != nil {
			return err
		}
	}

	// Subtypes
	for _, subtype := range card.Subtypes {
        err := InsertSubtypeToDB(cardId, subtype, queries)
		if err != nil {
			return err
		}
	}

	// Supertypes
	for _, supertype := range card.Supertypes {
        err := InsertSupertypeToDB(cardId, supertype, queries)
		if err != nil {
			return err
		}
	}

	// Variations
	for _, variation := range card.Variations {
		err := InsertVariationToDB(cardId, variation, queries)
		if err != nil {
			return err
		}
	}

    return nil
}

func InsertSubtypeToDB(
        cardId int64,
        subtype string,
        queries *DBInsertQueries) error {
	subtypeId, err := getSubtypeOptionId(subtype)
	if err != nil {
		return err
	}

	_, err = queries.InsertCardSubtypeQuery.Exec(cardId, subtypeId)
	if err != nil {
		return err
	}

    return nil
}

func InsertSupertypeToDB(
        cardId int64,
        supertype string,
        queries *DBInsertQueries) error {
	supertypeId, err := getSupertypeOptionId(supertype)
	if err != nil {
		return err
	}

	_, err = queries.InsertCardSupertypeQuery.Exec(cardId, supertypeId)
	if err != nil {
		return err
	}

	return nil
}

func InsertAltLangInfoToDB(
        cardId int64,
        altLangInfo *mtgcards.MTGCardAlternateLanguageInfo,
        queries *DBInsertQueries) error {
	_, err := queries.InsertAltLangDataQuery.Exec(
        cardId,
        altLangInfo.FlavorText,
		altLangInfo.Language,
        altLangInfo.MultiverseId,
        altLangInfo.Name,
		altLangInfo.Text,
        altLangInfo.Type)

	if err != nil {
		return err
	}

	return nil
}

func InsertRulingToDB(
        cardId int64,
        ruling *mtgcards.MTGCardRuling,
        queries *DBInsertQueries) error {
	_, err := queries.InsertRulingQuery.Exec(cardId, ruling.Date, ruling.Text)
	if err != nil {
		return err
	}

	return nil
}

func InsertBaseTypeToDB(
        cardId int64,
        baseTypeOption string,
        queries *DBInsertQueries) error {
	baseTypeOptionId, err := getBaseTypeOptionId(baseTypeOption)
	if err != nil {
		return err
	}

	_, err = queries.InsertBaseTypeQuery.Exec(cardId, baseTypeOptionId)
	if err != nil {
		return err
	}

	return nil
}

func InsertLeadershipSkillToDB(
        cardId int64,
        leadershipFormat string,
        leaderLegal bool,
        queries *DBInsertQueries) error {
	leadershipFormatId, err := getLeadershipFormatId(leadershipFormat)
	if err != nil {
		return err
	}

	_, err = queries.InsertLeadershipSkillQuery.Exec(
        cardId,
        leadershipFormatId,
        leaderLegal)
	if err != nil {
		return err
	}

	return nil
}

func InsertLegalityToDB(
        cardId int64,
        gameFormat string,
		legalityOption string,
        queries *DBInsertQueries) error {
	gameFormatId, err := getGameFormatId(gameFormat)
	if err != nil {
		return err
	}

	legalityOptionId, err := getLegalityOptionId(legalityOption)
	if err != nil {
		return err
	}

	_, err = queries.InsertLegalityQuery.Exec(cardId, gameFormatId, legalityOptionId)
	if err != nil {
		return err
	}

	return nil
}

func InsertPrintingToDB(
        cardId int64,
        setCode string,
        queries *DBInsertQueries) error {
	_, err := queries.InsertCardPrintingQuery.Exec(cardId, setCode)
	if err != nil {
		return err
	}

	return nil
}

func InsertPurchaseURLToDB(
        cardId int64,
        purchaseSite string,
        purchaseURL string,
        queries *DBInsertQueries) error {
	purchaseSiteId, err := getPurchaseSiteId(purchaseSite)
	if err != nil {
		return err
	}

	_, err = queries.InsertPurchaseURLQuery.Exec(cardId, purchaseSiteId, purchaseURL)
	if err != nil {
		return err
	}

	return nil
}

func InsertFrameEffectToDB(
        cardId int64,
        frameEffect string,
        queries *DBInsertQueries) error {
	frameEffectId, err := getFrameEffectId(frameEffect)
	if err != nil {
		return err
	}

	_, err = queries.InsertFrameEffectQuery.Exec(cardId, frameEffectId)
	if err != nil {
		return err
	}

	return nil
}

func InsertOtherFaceIdToDB(
        cardId int64,
        otherFaceUUID string,
        queries *DBInsertQueries) error {
	_, err := queries.InsertOtherFaceQuery.Exec(cardId, otherFaceUUID)
	if err != nil {
		return err
	}

	return nil
}

func InsertVariationToDB(
        cardId int64,
        variationUUID string,
        queries *DBInsertQueries) error {
	_, err := queries.InsertVariationQuery.Exec(cardId, variationUUID)
	if err != nil {
		return err
	}

	return nil
}

func DeleteOtherCardDataFromDB(
        cardId int64,
        queries *DBDeleteQueries) error {
    // Alternate language data
    _, err := queries.DeleteAltLangDataQuery.Exec(cardId)
    if err != nil {
        return err
    }

    // Base types
    _, err = queries.DeleteBaseTypesQuery.Exec(cardId)
    if err != nil {
        return err
    }

    // Frame effects
    _, err = queries.DeleteFrameEffectsQuery.Exec(cardId)
    if err != nil {
        return err
    }

	// Leadership skills
    _, err = queries.DeleteLeadershipSkillsQuery.Exec(cardId)
    if err != nil {
        return err
    }

	// Legalities
    _, err = queries.DeleteLegalitiesQuery.Exec(cardId)
    if err != nil {
        return err
    }

	// Other face IDs
    _, err = queries.DeleteOtherFaceQuery.Exec(cardId)
    if err != nil {
        return err
    }

    // Printings
    _, err = queries.DeleteCardPrintingsQuery.Exec(cardId)
    if err != nil {
        return err
    }

	// Purchase URLs
    _, err = queries.DeletePurchaseURLsQuery.Exec(cardId)
    if err != nil {
        return err
    }

	// Rulings
    _, err = queries.DeleteRulingsQuery.Exec(cardId)
    if err != nil {
        return err
    }

	// Subtypes
    _, err = queries.DeleteCardSubtypesQuery.Exec(cardId)
    if err != nil {
        return err
    }

	// Supertypes
    _, err = queries.DeleteCardSupertypesQuery.Exec(cardId)
    if err != nil {
        return err
    }

	// Variations
    _, err = queries.DeleteVariationsQuery.Exec(cardId)
    if err != nil {
        return err
    }

    return nil
}
