package mtgcards

import "database/sql"
import "strings"

type dbGetQueries struct {
    GetSetHashQuery *sql.Stmt
    GetCardHashQuery *sql.Stmt
}

type dbInsertQueries struct {
    InsertSetQuery *sql.Stmt
    InsertSetTranslationQuery *sql.Stmt
	InsertCardQuery *sql.Stmt
    InsertAltLangDataQuery *sql.Stmt
	InsertBaseTypeQuery *sql.Stmt
	InsertCardPrintingQuery *sql.Stmt
	InsertCardSubtypeQuery *sql.Stmt
	InsertCardSupertypeQuery *sql.Stmt
	InsertFrameEffectQuery *sql.Stmt
	InsertLeadershipSkillQuery *sql.Stmt
	InsertLegalityQuery *sql.Stmt
	InsertOtherFaceQuery *sql.Stmt
	InsertPurchaseURLQuery *sql.Stmt
	InsertRulingQuery *sql.Stmt
	InsertVariationQuery *sql.Stmt
}

type dbUpdateQueries struct {
    UpdateSetQuery *sql.Stmt
	UpdateCardQuery *sql.Stmt
}

type dbDeleteQueries struct {
    DeleteSetTranslationsQuery *sql.Stmt
    DeleteAltLangDataQuery *sql.Stmt
    DeleteBaseTypesQuery *sql.Stmt
    DeleteCardPrintingsQuery *sql.Stmt
    DeleteCardSubtypesQuery *sql.Stmt
    DeleteCardSupertypesQuery *sql.Stmt
    DeleteFrameEffectsQuery *sql.Stmt
    DeleteLeadershipSkillsQuery *sql.Stmt
    DeleteLegalitiesQuery *sql.Stmt
	DeleteOtherFaceQuery *sql.Stmt
    DeletePurchaseURLsQuery *sql.Stmt
    DeleteRulingsQuery *sql.Stmt
	DeleteVariationsQuery *sql.Stmt
}

func (queries *dbGetQueries) ForTx(tx *sql.Tx) *dbGetQueries {
    var txQueries dbGetQueries

    txQueries.GetSetHashQuery = tx.Stmt(queries.GetSetHashQuery)
    txQueries.GetCardHashQuery = tx.Stmt(queries.GetCardHashQuery)

    return &txQueries
}

func (queries *dbInsertQueries) ForTx(tx *sql.Tx) *dbInsertQueries {
    var txQueries dbInsertQueries

    txQueries.InsertSetQuery = tx.Stmt(queries.InsertSetQuery)
    txQueries.InsertSetTranslationQuery = tx.Stmt(queries.InsertSetTranslationQuery)
	txQueries.InsertCardQuery = tx.Stmt(queries.InsertCardQuery)
    txQueries.InsertAltLangDataQuery = tx.Stmt(queries.InsertAltLangDataQuery)
	txQueries.InsertBaseTypeQuery = tx.Stmt(queries.InsertBaseTypeQuery)
	txQueries.InsertCardPrintingQuery = tx.Stmt(queries.InsertCardPrintingQuery)
	txQueries.InsertCardSubtypeQuery = tx.Stmt(queries.InsertCardSubtypeQuery)
	txQueries.InsertCardSupertypeQuery = tx.Stmt(queries.InsertCardSupertypeQuery)
	txQueries.InsertFrameEffectQuery = tx.Stmt(queries.InsertFrameEffectQuery)
	txQueries.InsertLeadershipSkillQuery = tx.Stmt(queries.InsertLeadershipSkillQuery)
	txQueries.InsertLegalityQuery = tx.Stmt(queries.InsertLegalityQuery)
	txQueries.InsertOtherFaceQuery = tx.Stmt(queries.InsertOtherFaceQuery)
	txQueries.InsertPurchaseURLQuery = tx.Stmt(queries.InsertPurchaseURLQuery)
	txQueries.InsertRulingQuery = tx.Stmt(queries.InsertRulingQuery)
	txQueries.InsertVariationQuery = tx.Stmt(queries.InsertVariationQuery)

    return &txQueries
}

func (queries *dbUpdateQueries) ForTx(tx *sql.Tx) *dbUpdateQueries {
    var txQueries dbUpdateQueries

    txQueries.UpdateSetQuery = tx.Stmt(queries.UpdateSetQuery)
	txQueries.UpdateCardQuery = tx.Stmt(queries.UpdateCardQuery)

    return &txQueries
}

func (queries *dbDeleteQueries) ForTx(tx *sql.Tx) *dbDeleteQueries {
    var txQueries dbDeleteQueries

    txQueries.DeleteSetTranslationsQuery = tx.Stmt(queries.DeleteSetTranslationsQuery)
    txQueries.DeleteAltLangDataQuery = tx.Stmt(queries.DeleteAltLangDataQuery)
    txQueries.DeleteBaseTypesQuery = tx.Stmt(queries.DeleteBaseTypesQuery)
    txQueries.DeleteCardPrintingsQuery = tx.Stmt(queries.DeleteCardPrintingsQuery)
    txQueries.DeleteCardSubtypesQuery = tx.Stmt(queries.DeleteCardSubtypesQuery)
    txQueries.DeleteCardSupertypesQuery = tx.Stmt(queries.DeleteCardSupertypesQuery)
    txQueries.DeleteFrameEffectsQuery = tx.Stmt(queries.DeleteFrameEffectsQuery)
    txQueries.DeleteLeadershipSkillsQuery = tx.Stmt(queries.DeleteLeadershipSkillsQuery)
    txQueries.DeleteLegalitiesQuery = tx.Stmt(queries.DeleteLegalitiesQuery)
	txQueries.DeleteOtherFaceQuery = tx.Stmt(queries.DeleteOtherFaceQuery)
    txQueries.DeletePurchaseURLsQuery = tx.Stmt(queries.DeletePurchaseURLsQuery)
    txQueries.DeleteRulingsQuery = tx.Stmt(queries.DeleteRulingsQuery)
	txQueries.DeleteVariationsQuery = tx.Stmt(queries.DeleteVariationsQuery)

    return &txQueries
}

func (queries *dbGetQueries) Prepare(db *sql.DB) error {
    var err error

    queries.GetSetHashQuery, err = db.Prepare(`SELECT set_hash, set_id 
        FROM sets 
        WHERE code = ?`)
	if err != nil {
        return err
	}

	queries.GetCardHashQuery, err = db.Prepare(`SELECT card_hash, card_id
        FROM all_cards
        WHERE uuid = ?`)
	if err != nil {
        return err
	}

    return nil
}

func (queries *dbGetQueries) Cleanup() {
    if queries.GetSetHashQuery != nil {
        queries.GetSetHashQuery.Close()
    }

    if queries.GetCardHashQuery != nil {
        queries.GetCardHashQuery.Close()
    }
}

func (queries *dbInsertQueries) Prepare(db *sql.DB) error {
    var err error

	queries.InsertSetQuery, err = db.Prepare(`INSERT INTO sets
		(set_hash, base_size, block_name, code, is_foreign_only, is_foil_only,
		is_online_only, is_partial_preview, keyrune_code, mcm_name, mcm_id,
		mtgo_code, name, parent_code, release_date, tcgplayer_group_id,
		total_set_size, set_type)
		VALUES
		(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
        return err
	}

    queries.InsertSetTranslationQuery, err = db.Prepare(`INSERT INTO set_translations
		(set_id, set_translation_language_id, set_translated_name)
		VALUES
		(?, ?, ?)`)
	if err != nil {
		return err
	}

    queries.InsertCardQuery, err = db.Prepare(`INSERT INTO all_cards
		(uuid, card_hash, artist, border_color, card_number, card_power,
        card_type, color_identity, color_indicator, colors, converted_mana_cost,
        duel_deck, edhrec_rank, face_converted_mana_cost, flavor_text,
        frame_version, hand, has_foil, has_non_foil, is_alternative,
        is_arena, is_full_art, is_mtgo, is_online_only, is_oversized,
        is_paper, is_promo, is_reprint, is_reserved, is_starter,
        is_story_spotlight, is_textless, is_timeshifted, layout, life, loyalty,
        mana_cost, mcm_id, mcm_meta_id, mtg_arena_id, mtgo_foil_id, mtgo_id,
        mtgstocks_id, multiverse_id, name, original_text, original_type, rarity,
        scryfall_id, scryfall_illustration_id, scryfall_oracle_id, set_id,
        side, tcgplayer_product_id, text, toughness, watermark)
		VALUES
		(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?,
		?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?,
        ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return err
	}

	queries.InsertAltLangDataQuery, err = db.Prepare(`INSERT INTO alternate_language_data
		(card_id, flavor_text, language, multiverse_id, name, text, card_type)
		VALUES
		(?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return err
	}

	queries.InsertBaseTypeQuery, err = db.Prepare(`INSERT INTO base_types
		(card_id, base_type_option_id)
		VALUES
		(?, ?)`)
	if err != nil {
		return err
	}

	queries.InsertCardPrintingQuery, err = db.Prepare(`INSERT INTO card_printings
		(card_id, set_code)
		VALUES
		(?, ?)`)
	if err != nil {
		return err
	}

	queries.InsertCardSubtypeQuery, err = db.Prepare(`INSERT INTO card_subtypes
		(card_id, subtype_option_id)
		VALUES
		(?, ?)`)
	if err != nil {
		return err
	}

	queries.InsertCardSupertypeQuery, err = db.Prepare(`INSERT INTO card_supertypes
		(card_id, supertype_option_id)
		VALUES
		(?, ?)`)
	if err != nil {
		return err
	}

	queries.InsertFrameEffectQuery, err = db.Prepare(`INSERT INTO frame_effects
		(card_id, frame_effect_option_id)
		VALUES
		(?, ?)`)
	if err != nil {
		return err
	}

	queries.InsertLeadershipSkillQuery, err = db.Prepare(`INSERT INTO leadership_skills
		(card_id, leadership_format_id, leader_legal)
		VALUES
		(?, ?, ?)`)
	if err != nil {
		return err
	}

	queries.InsertLegalityQuery, err = db.Prepare(`INSERT INTO legalities
		(card_id, game_format_id, legality_option_id)
		VALUES
		(?, ?, ?)`)
	if err != nil {
		return err
	}

	queries.InsertOtherFaceQuery, err = db.Prepare(`INSERT INTO other_faces
		(card_id, other_face_uuid)
		VALUES
		(?, ?)`)
	if err != nil {
		return err
	}

	queries.InsertPurchaseURLQuery, err = db.Prepare(`INSERT INTO purchase_urls
		(card_id, purchase_site_id, purchase_url)
		VALUES
		(?, ?, ?)`)
	if err != nil {
		return err
	}

	queries.InsertRulingQuery, err = db.Prepare(`INSERT INTO rulings
		(card_id, ruling_date, ruling_text)
		VALUES
		(?, ?, ?)`)
	if err != nil {
		return err
	}

	queries.InsertVariationQuery, err = db.Prepare(`INSERT INTO variations
		(card_id, variation_uuid)
		VALUES
		(?, ?)`)
	if err != nil {
		return err
	}

    return nil
}

func (queries *dbInsertQueries) Cleanup() {
    if queries.InsertSetQuery != nil {
        queries.InsertSetQuery.Close()
    }

    if queries.InsertSetTranslationQuery != nil {
        queries.InsertSetTranslationQuery.Close()
    }

    if queries.InsertCardQuery != nil {
        queries.InsertCardQuery.Close()
    }

    if queries.InsertAltLangDataQuery != nil {
        queries.InsertAltLangDataQuery.Close()
    }

	if queries.InsertBaseTypeQuery != nil {
        queries.InsertBaseTypeQuery.Close()
    }

	if queries.InsertCardPrintingQuery != nil {
        queries.InsertCardPrintingQuery.Close()
    }

	if queries.InsertCardSubtypeQuery != nil {
        queries.InsertCardSubtypeQuery.Close()
    }

	if queries.InsertCardSupertypeQuery != nil {
        queries.InsertCardSupertypeQuery.Close()
    }

	if queries.InsertFrameEffectQuery != nil {
        queries.InsertFrameEffectQuery.Close()
    }

	if queries.InsertLeadershipSkillQuery != nil {
        queries.InsertLeadershipSkillQuery.Close()
    }

	if queries.InsertLegalityQuery != nil {
        queries.InsertLegalityQuery.Close()
    }

	if queries.InsertOtherFaceQuery != nil {
        queries.InsertOtherFaceQuery.Close()
    }

	if queries.InsertPurchaseURLQuery != nil {
        queries.InsertPurchaseURLQuery.Close()
    }

	if queries.InsertRulingQuery != nil {
        queries.InsertRulingQuery.Close()
    }

	if queries.InsertVariationQuery != nil {
        queries.InsertVariationQuery.Close()
    }
}

func (queries *dbUpdateQueries) Prepare(db *sql.DB) error {
    var err error

	queries.UpdateSetQuery, err = db.Prepare(`UPDATE sets SET
		set_hash = ?,
		base_size = ?,
		block_name = ?,
		code = ?,
		is_foreign_only = ?,
		is_foil_only = ?,
		is_online_only = ?,
		is_partial_preview = ?,
		keyrune_code = ?,
		mcm_name = ?,
		mcm_id = ?,
		mtgo_code = ?,
		name = ?,
		parent_code = ?,
		release_date = ?,
		tcgplayer_group_id = ?,
		total_set_size = ?,
		set_type = ?
		WHERE set_id = ?`)
	if err != nil {
        return err
	}

	queries.UpdateCardQuery, err = db.Prepare(`UPDATE all_cards SET
        card_hash = ?,
        artist = ?,
        border_color = ?,
        card_number = ?,
        card_power = ?,
        card_type = ?,
        color_identity = ?,
        color_indicator = ?,
        colors = ?,
        converted_mana_cost = ?,
        duel_deck = ?,
        edhrec_rank = ?,
        face_converted_mana_cost = ?,
        flavor_text = ?,
        frame_version = ?,
        hand = ?,
        has_foil = ?,
        has_non_foil = ?,
        is_alternative = ?,
        is_arena = ?,
        is_full_art = ?,
        is_mtgo = ?,
        is_online_only = ?,
        is_oversized = ?,
        is_paper = ?,
        is_promo = ?,
        is_reprint = ?,
        is_reserved = ?,
        is_starter = ?,
        is_story_spotlight = ?,
        is_textless = ?,
        is_timeshifted = ?,
        layout = ?,
        life = ?,
        loyalty = ?,
        mana_cost = ?,
        mcm_id = ?,
        mcm_meta_id = ?,
        mtg_arena_id = ?,
        mtgo_foil_id = ?,
        mtgo_id = ?,
        mtgstocks_id = ?,
        multiverse_id = ?,
        name = ?,
        original_text = ?,
        original_type = ?,
        rarity = ?,
        scryfall_id = ?,
        scryfall_illustration_id = ?,
        scryfall_oracle_id = ?,
        set_id = ?,
        side = ?,
        tcgplayer_product_id = ?,
        text = ?,
        toughness = ?,
        watermark = ?
		WHERE uuid = ?`)
    if err != nil {
        return err
    }

    return nil
}

func (queries *dbUpdateQueries) Cleanup() {
    if queries.UpdateSetQuery != nil {
        queries.UpdateSetQuery.Close()
    }

	if queries.UpdateCardQuery != nil {
        queries.UpdateCardQuery.Close()
    }
}

func (queries *dbDeleteQueries) Prepare(db *sql.DB) error {
    var err error

    queries.DeleteSetTranslationsQuery, err = db.Prepare(`DELETE
        FROM set_translations
        WHERE set_id = ?`)
    if err != nil {
        return err
    }

    queries.DeleteAltLangDataQuery, err = db.Prepare(`DELETE
        FROM alternate_language_data
        WHERE card_id = ?`)
    if err != nil {
        return err
    }

    queries.DeleteBaseTypesQuery, err = db.Prepare(`DELETE
        FROM base_types
        WHERE card_id = ?`)
    if err != nil {
        return err
    }

    queries.DeleteCardPrintingsQuery, err = db.Prepare(`DELETE
        FROM card_printings
        WHERE card_id = ?`)
    if err != nil {
        return err
    }

    queries.DeleteCardSubtypesQuery, err = db.Prepare(`DELETE
        FROM card_subtypes
        WHERE card_id = ?`)
    if err != nil {
        return err
    }

    queries.DeleteCardSupertypesQuery, err = db.Prepare(`DELETE
        FROM card_supertypes
        WHERE card_id = ?`)
    if err != nil {
        return err
    }

	queries.DeleteFrameEffectsQuery, err = db.Prepare(`DELETE
        FROM frame_effects
		WHERE card_id = ?`)
	if err != nil {
		return err
	}

    queries.DeleteLeadershipSkillsQuery, err = db.Prepare(`DELETE
        FROM leadership_skills
        WHERE card_id = ?`)
    if err != nil {
        return err
    }

    queries.DeleteLegalitiesQuery, err = db.Prepare(`DELETE
        FROM legalities
        WHERE card_id = ?`)
    if err != nil {
        return err
    }

	queries.DeleteOtherFaceQuery, err = db.Prepare(`DELETE
        FROM other_faces
		WHERE card_id = ?`)
	if err != nil {
		return err
	}

    queries.DeletePurchaseURLsQuery, err = db.Prepare(`DELETE
        FROM purchase_urls
        WHERE card_id = ?`)
    if err != nil {
        return err
    }

    queries.DeleteRulingsQuery, err = db.Prepare(`DELETE
        FROM rulings
        WHERE card_id = ?`)
    if err != nil {
        return err
    }

	queries.DeleteVariationsQuery, err = db.Prepare(`DELETE
        FROM variations
		WHERE card_id = ?`)
	if err != nil {
		return err
	}

    return nil
}

func (queries *dbDeleteQueries) Cleanup() {
    if queries.DeleteSetTranslationsQuery != nil {
        queries.DeleteSetTranslationsQuery.Close()
    }

    if queries.DeleteAltLangDataQuery != nil {
        queries.DeleteAltLangDataQuery.Close()
    }

    if queries.DeleteBaseTypesQuery != nil {
        queries.DeleteBaseTypesQuery.Close()
    }

    if queries.DeleteCardPrintingsQuery != nil {
        queries.DeleteCardPrintingsQuery.Close()
    }

    if queries.DeleteCardSubtypesQuery != nil {
        queries.DeleteCardSubtypesQuery.Close()
    }

    if queries.DeleteCardSupertypesQuery != nil {
        queries.DeleteCardSupertypesQuery.Close()
    }
    if queries.DeleteFrameEffectsQuery != nil {
        queries.DeleteFrameEffectsQuery.Close()
    }

    if queries.DeleteLeadershipSkillsQuery != nil {
        queries.DeleteLeadershipSkillsQuery.Close()
    }

    if queries.DeleteLegalitiesQuery != nil {
        queries.DeleteLegalitiesQuery.Close()
    }

	if queries.DeleteOtherFaceQuery != nil {
        queries.DeleteOtherFaceQuery.Close()
    }

    if queries.DeletePurchaseURLsQuery != nil {
        queries.DeletePurchaseURLsQuery.Close()
    }

    if queries.DeleteRulingsQuery != nil {
        queries.DeleteRulingsQuery.Close()
    }

	if queries.DeleteVariationsQuery != nil {
        queries.DeleteVariationsQuery.Close()
    }
}

func (card *MTGCard) GetHashAndId(queries *dbGetQueries) (bool, string, int64, error) {
	res := queries.GetCardHashQuery.QueryRow(card.UUID)

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

func (set *MTGSet) GetHashAndId(queries *dbGetQueries) (bool, string, int64, error) {
	res := queries.GetSetHashQuery.QueryRow(set.Code)

	var setHash string
	var setId int64
	err := res.Scan(&setHash, &setId)
	if err != nil {
		if err == sql.ErrNoRows {
			// This set isn't in the database
			return false, "", 0, nil
		} else {
			return false, "", 0, err
		}
	} else {
		return true, setHash, setId, nil
	}
}

func (set *MTGSet) InsertToDb(queries *dbInsertQueries) (int64, error) {
    setHash := HashToHexString(set.Hash())

    // Insert the set itself
	res, err := queries.InsertSetQuery.Exec(
        setHash,
        set.BaseSetSize,
        set.Block,
		set.Code,
        set.IsForeignOnly,
        set.IsFoilOnly,
        set.IsOnlineOnly,
		set.IsPartialPreview,
        set.KeyruneCode,
        set.MCMName,
        set.MCMId,
        set.MTGOCode,
		set.Name,
        set.ParentCode,
        set.ReleaseDate,
        set.TCGPlayerGroupId,
		set.TotalSetSize,
        set.Type)
	if err != nil {
		return 0, err
	}

    setId, err := res.LastInsertId()
    if err != nil {
        return 0, err
    }

	// Insert the set translations
	for lang, name := range set.Translations {
		err := InsertSetTranslationToDb(queries, setId, lang, name)
		if err != nil {
            return 0, err
		}
	}

	return setId, nil
}

func (card *MTGCard) InsertToDb(queries *dbInsertQueries, setId int64) error {
	// Build all of the values that can be null
    var colorIdentity sql.NullString
	var colorIndicator sql.NullString
	var colors sql.NullString
    var duelDeck sql.NullString
	var edhrecRank sql.NullInt32
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

    cardHash := HashToHexString(card.Hash())

	res, err := queries.InsertCardQuery.Exec(
        card.UUID,
        cardHash,
        card.Artist,
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
        flavorText,
        card.FrameVersion,
		hand,
        card.HasFoil,
        card.HasNonFoil,
        card.IsAlternative,
        card.IsArena,
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
    err = card.InsertOtherTableData(queries, cardId)
    if err != nil {
        return nil
    }

	return nil
}

func InsertSetTranslationToDb(queries *dbInsertQueries, setId int64, translationLang string,
		translatedName string) error {
	languageId, err := getSetTranslationLanguageId(translationLang)
	if err != nil {
		return err
	}

	_, err = queries.InsertSetTranslationQuery.Exec(setId, languageId, translatedName)
	if err != nil {
		return err
	}

	return nil
}

func InsertCardSubtypeToDb(queries *dbInsertQueries, cardId int64, subtype string) error {
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

func InsertCardSupertypeToDb(queries *dbInsertQueries, cardId int64, supertype string) error {
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

func (altLangInfo *MTGCardAlternateLanguageInfo) InsertToDb(queries *dbInsertQueries,
		cardId int64) error {
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

func (ruling *MTGCardRuling) InsertToDb(queries *dbInsertQueries, cardId int64) error {
	_, err := queries.InsertRulingQuery.Exec(cardId, ruling.Date, ruling.Text)
	if err != nil {
		return err
	}

	return nil
}

func InsertBaseTypeToDb(queries *dbInsertQueries, cardId int64, baseTypeOption string) error {
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

func InsertLeadershipSkillToDb(queries *dbInsertQueries, cardId int64,
		leadershipFormat string, leaderLegal bool) error {

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

func InsertLegalityToDb(queries *dbInsertQueries, cardId int64, gameFormat string,
		legalityOption string) error {
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

func InsertCardPrintingToDb(queries *dbInsertQueries, cardId int64, setCode string) error {
	_, err := queries.InsertCardPrintingQuery.Exec(cardId, setCode)
	if err != nil {
		return err
	}

	return nil
}

func InsertPurchaseURLToDb(queries *dbInsertQueries, cardId int64,
		purchaseSite string, purchaseURL string) error {
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

func InsertFrameEffectToDb(queries *dbInsertQueries, cardId int64, frameEffect string) error {
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

func InsertOtherFaceIdToDb(queries *dbInsertQueries, cardId int64, otherFaceUUID string) error {
	_, err := queries.InsertOtherFaceQuery.Exec(cardId, otherFaceUUID)
	if err != nil {
		return err
	}

	return nil
}

func InsertVariationToDb(queries *dbInsertQueries, cardId int64, variationUUID string) error {
	_, err := queries.InsertVariationQuery.Exec(cardId, variationUUID)
	if err != nil {
		return err
	}

	return nil
}

func (set *MTGSet) UpdateInDb(
        updateQueries *dbUpdateQueries,
        deleteQueries *dbDeleteQueries,
        insertQueries *dbInsertQueries,
        setId int64) error {

    // First, update the main set record
    setHash := HashToHexString(set.Hash())
	_, err := updateQueries.UpdateSetQuery.Exec(setHash, set.BaseSetSize, set.Block,
		set.Code, set.IsForeignOnly, set.IsFoilOnly, set.IsOnlineOnly,
		set.IsPartialPreview, set.KeyruneCode, set.MCMName, set.MCMId, set.MTGOCode,
		set.Name, set.ParentCode, set.ReleaseDate, set.TCGPlayerGroupId,
		set.TotalSetSize, set.Type, setId)
	if err != nil {
		return err
	}

    // Next, delete the old set translations
    _, err = deleteQueries.DeleteSetTranslationsQuery.Exec(setId)
    if err != nil {
        return err
    }

    // Finally, insert the new set translations
	for lang, name := range set.Translations {
		err := InsertSetTranslationToDb(insertQueries, setId, lang, name)
		if err != nil {
            return err
		}
	}

    return nil
}

func (card *MTGCard) InsertOtherTableData(queries *dbInsertQueries, cardId int64) error {
	// Alternate language data
	for _, altLangData := range card.AlternateLanguageData {
        err := altLangData.InsertToDb(queries, cardId)
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
        err := InsertBaseTypeToDb(queries, cardId, baseType)
		if err != nil {
			return err
		}
	}

    // Frame effects
	for _, frameEffect := range card.FrameEffects {
		err := InsertFrameEffectToDb(queries, cardId, frameEffect)
		if err != nil {
			return err
		}
	}

	// Leadership skills
	for leadershipFormat, leaderValid := range card.LeadershipSkills {
        err := InsertLeadershipSkillToDb(queries, cardId, leadershipFormat, leaderValid)
		if err != nil {
			return err
		}
	}

	// Legalities
	for format, legality := range card.Legalities {
        err := InsertLegalityToDb(queries, cardId, format, legality)
		if err != nil {
			return err
		}
	}

	// Other face IDs
	for _, otherFaceId := range card.OtherFaceIds {
		err := InsertOtherFaceIdToDb(queries, cardId, otherFaceId)
		if err != nil {
			return err
		}
	}

    // Printings
	for _, setCode := range card.Printings {
        err := InsertCardPrintingToDb(queries, cardId, setCode)
		if err != nil {
			return err
		}
	}

	// Purchase URLs
	for site, url := range card.PurchaseURLs {
        err := InsertPurchaseURLToDb(queries, cardId, site, url)
		if err != nil {
			return err
		}
	}

	// Rulings
	for _, ruling := range card.Rulings {
        err := ruling.InsertToDb(queries, cardId)
		if err != nil {
			return err
		}
	}

	// Subtypes
	for _, subtype := range card.Subtypes {
        err := InsertCardSubtypeToDb(queries, cardId, subtype)
		if err != nil {
			return err
		}
	}

	// Supertypes
	for _, supertype := range card.Supertypes {
        err := InsertCardSupertypeToDb(queries, cardId, supertype)
		if err != nil {
			return err
		}
	}

	// Variations
	for _, variation := range card.Variations {
		err := InsertVariationToDb(queries, cardId, variation)
		if err != nil {
			return err
		}
	}

    return nil
}

func DeleteOtherCardData(queries *dbDeleteQueries, cardId int64) error {
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

func (card *MTGCard) UpdateInDb(
        updateQueries *dbUpdateQueries,
        deleteQueries *dbDeleteQueries,
        insertQueries *dbInsertQueries,
        cardId int64,
        setId int64) error {
    // Build all of the values that can be null
    var colorIdentity sql.NullString
	var colorIndicator sql.NullString
	var colors sql.NullString
    var duelDeck sql.NullString
	var edhrecRank sql.NullInt32
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
    cardHash := HashToHexString(card.Hash())
	_, err := updateQueries.UpdateCardQuery.Exec(
        cardHash,
        card.Artist,
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
        flavorText,
        card.FrameVersion,
		hand,
        card.HasFoil,
        card.HasNonFoil,
        card.IsAlternative,
        card.IsArena,
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
    err = DeleteOtherCardData(deleteQueries, cardId)
    if err != nil {
        return err
    }

    // Finally, insert the rest of the new card data
    err = card.InsertOtherTableData(insertQueries, cardId)
    if err != nil {
        return err
    }

    return nil
}
