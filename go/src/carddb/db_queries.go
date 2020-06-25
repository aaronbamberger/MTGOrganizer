package carddb

import "database/sql"

type DBGetQueries struct {
    GetSetHashQuery *sql.Stmt
    GetCardHashQuery *sql.Stmt
    GetTokenHashQuery *sql.Stmt
}

type DBInsertQueries struct {
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
    InsertTokenQuery *sql.Stmt
    InsertTokenBaseTypeQuery *sql.Stmt
    InsertTokenSubtypeQuery *sql.Stmt
    InsertTokenSupertypeQuery *sql.Stmt
    InsertTokenReverseRelatedQuery *sql.Stmt
}

type DBUpdateQueries struct {
    UpdateSetQuery *sql.Stmt
	UpdateCardQuery *sql.Stmt
    UpdateTokenQuery *sql.Stmt
}

type DBDeleteQueries struct {
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
    DeleteTokenBaseTypesQuery *sql.Stmt
    DeleteTokenSubtypesQuery *sql.Stmt
    DeleteTokenSupertypesQuery *sql.Stmt
    DeleteTokenReverseRelatedQuery *sql.Stmt
}

func (queries *DBGetQueries) ForTx(tx *sql.Tx) *DBGetQueries {
    var txQueries DBGetQueries

    txQueries.GetSetHashQuery = tx.Stmt(queries.GetSetHashQuery)
    txQueries.GetCardHashQuery = tx.Stmt(queries.GetCardHashQuery)
    txQueries.GetTokenHashQuery = tx.Stmt(queries.GetTokenHashQuery)

    return &txQueries
}

func (queries *DBInsertQueries) ForTx(tx *sql.Tx) *DBInsertQueries {
    var txQueries DBInsertQueries

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
    txQueries.InsertTokenQuery = tx.Stmt(queries.InsertTokenQuery)
    txQueries.InsertTokenBaseTypeQuery = tx.Stmt(queries.InsertTokenBaseTypeQuery)
    txQueries.InsertTokenSubtypeQuery = tx.Stmt(queries.InsertTokenSubtypeQuery)
    txQueries.InsertTokenSupertypeQuery = tx.Stmt(queries.InsertTokenSupertypeQuery)
    txQueries.InsertTokenReverseRelatedQuery = tx.Stmt(queries.InsertTokenReverseRelatedQuery)

    return &txQueries
}

func (queries *DBUpdateQueries) ForTx(tx *sql.Tx) *DBUpdateQueries {
    var txQueries DBUpdateQueries

    txQueries.UpdateSetQuery = tx.Stmt(queries.UpdateSetQuery)
	txQueries.UpdateCardQuery = tx.Stmt(queries.UpdateCardQuery)
    txQueries.UpdateTokenQuery = tx.Stmt(queries.UpdateTokenQuery)

    return &txQueries
}

func (queries *DBDeleteQueries) ForTx(tx *sql.Tx) *DBDeleteQueries {
    var txQueries DBDeleteQueries

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
    txQueries.DeleteTokenBaseTypesQuery = tx.Stmt(queries.DeleteTokenBaseTypesQuery)
    txQueries.DeleteTokenSubtypesQuery = tx.Stmt(queries.DeleteTokenSubtypesQuery)
    txQueries.DeleteTokenSupertypesQuery = tx.Stmt(queries.DeleteTokenSupertypesQuery)
    txQueries.DeleteTokenReverseRelatedQuery = tx.Stmt(queries.DeleteTokenReverseRelatedQuery)

    return &txQueries
}

func (queries *DBGetQueries) Prepare(db *sql.DB) error {
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

    queries.GetTokenHashQuery, err = db.Prepare(`SELECT token_hash, token_id
        FROM all_tokens
        WHERE uuid = ?`)
    if err != nil {
        return err
    }

    return nil
}

func (queries *DBGetQueries) Cleanup() {
    if queries.GetSetHashQuery != nil {
        queries.GetSetHashQuery.Close()
    }

    if queries.GetCardHashQuery != nil {
        queries.GetCardHashQuery.Close()
    }

    if queries.GetTokenHashQuery != nil {
        queries.GetTokenHashQuery.Close()
    }
}

func (queries *DBInsertQueries) Prepare(db *sql.DB) error {
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
		(uuid, card_hash, artist, ascii_name, border_color, card_number, card_power,
        card_type, color_identity, color_indicator, colors, converted_mana_cost,
        duel_deck, edhrec_rank, face_converted_mana_cost, flavor_name, flavor_text,
        frame_version, hand, has_foil, has_non_foil, is_alternative,
        is_arena, is_buy_a_box, is_date_stamped, is_full_art, is_mtgo, is_online_only,
        is_oversized, is_paper, is_promo, is_reprint, is_reserved, is_starter,
        is_story_spotlight, is_textless, is_timeshifted, layout, life, loyalty,
        mana_cost, mcm_id, mcm_meta_id, mtg_arena_id, mtgo_foil_id, mtgo_id,
        mtgstocks_id, multiverse_id, name, original_text, original_type, rarity,
        scryfall_id, scryfall_illustration_id, scryfall_oracle_id, set_id,
        side, tcgplayer_product_id, text, toughness, watermark)
		VALUES
		(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?,
		?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?,
        ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)
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

	queries.InsertBaseTypeQuery, err = db.Prepare(`INSERT INTO card_base_types
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

    // Since multiple sets might have the same token, we treat trying to insert
    // an existing token as expected, and just do nothing if we try and insert
    // a duplicate token
    queries.InsertTokenQuery, err = db.Prepare(`INSERT INTO all_tokens
        (uuid, token_hash, artist, border_color, card_number, card_power,
        card_type, color_identity, color_indicator, colors, is_online_only,
        layout, loyalty, name, scryfall_id, scryfall_illustration_id,
        scryfall_oracle_id, set_id, side, text, toughness, watermark)
        VALUES
        (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
        ON DUPLICATE KEY UPDATE uuid=uuid`)
    if err != nil {
        return err
    }

    queries.InsertTokenBaseTypeQuery, err = db.Prepare(`INSERT INTO token_base_types
        (token_id, base_type_option_id)
        VALUES
        (?, ?)`)
    if err != nil {
        return err
    }

    queries.InsertTokenSubtypeQuery, err = db.Prepare(`INSERT INTO token_subtypes
        (token_id, subtype_option_id)
        VALUES
        (?, ?)`)
    if err != nil {
        return err
    }

    queries.InsertTokenSupertypeQuery, err = db.Prepare(`INSERT INTO token_supertypes
        (token_id, supertype_option_id)
        VALUES
        (?, ?)`)
    if err != nil {
        return err
    }

    queries.InsertTokenReverseRelatedQuery, err = db.Prepare(`INSERT INTO token_reverse_related
        (token_id, reverse_related_card)
        VALUES
        (?, ?)`)
    if err != nil {
        return err
    }

    return nil
}

func (queries *DBInsertQueries) Cleanup() {
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

    if queries.InsertTokenQuery != nil {
        queries.InsertTokenQuery.Close()
    }

    if queries.InsertTokenBaseTypeQuery != nil {
        queries.InsertTokenBaseTypeQuery.Close()
    }

    if queries.InsertTokenSubtypeQuery != nil {
        queries.InsertTokenSubtypeQuery.Close()
    }

    if queries.InsertTokenSupertypeQuery != nil {
        queries.InsertTokenSupertypeQuery.Close()
    }

    if queries.InsertTokenReverseRelatedQuery != nil {
        queries.InsertTokenReverseRelatedQuery.Close()
    }
}

func (queries *DBUpdateQueries) Prepare(db *sql.DB) error {
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
        ascii_name = ?,
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
        flavor_name = ?,
        flavor_text = ?,
        frame_version = ?,
        hand = ?,
        has_foil = ?,
        has_non_foil = ?,
        is_alternative = ?,
        is_arena = ?,
        is_buy_a_box = ?,
        is_date_stamped = ?,
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

    queries.UpdateTokenQuery, err = db.Prepare(`UPDATE all_tokens SET
        token_hash = ?,
        artist = ?,
        border_color = ?,
        card_number = ?,
        card_power = ?,
        card_type = ?,
        color_identity = ?,
        color_indicator = ?,
        colors = ?,
        is_online_only = ?,
        layout = ?,
        loyalty = ?,
        name = ?,
        scryfall_id = ?,
        scryfall_illustration_id = ?,
        scryfall_oracle_id = ?,
        set_id = ?,
        side = ?,
        text = ?,
        toughness = ?,
        watermark = ?
        WHERE uuid = ?`)
    if err != nil {
        return err
    }

    return nil
}

func (queries *DBUpdateQueries) Cleanup() {
    if queries.UpdateSetQuery != nil {
        queries.UpdateSetQuery.Close()
    }

	if queries.UpdateCardQuery != nil {
        queries.UpdateCardQuery.Close()
    }

    if queries.UpdateTokenQuery != nil {
        queries.UpdateTokenQuery.Close()
    }
}

func (queries *DBDeleteQueries) Prepare(db *sql.DB) error {
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
        FROM card_base_types
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

    queries.DeleteTokenBaseTypesQuery, err = db.Prepare(`DELETE
        FROM token_base_types
        WHERE token_id = ?`)
    if err != nil {
        return err
    }

    queries.DeleteTokenSubtypesQuery, err = db.Prepare(`DELETE
        FROM token_subtypes
        WHERE token_id = ?`)
    if err != nil {
        return err
    }

    queries.DeleteTokenSupertypesQuery, err = db.Prepare(`DELETE
        FROM token_supertypes
        WHERE token_id = ?`)
    if err != nil {
        return err
    }

    queries.DeleteTokenReverseRelatedQuery, err = db.Prepare(`DELETE
        FROM token_reverse_related
        WHERE token_id = ?`)
    if err != nil {
        return err
    }

    return nil
}

func (queries *DBDeleteQueries) Cleanup() {
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

    if queries.DeleteTokenBaseTypesQuery != nil {
        queries.DeleteTokenBaseTypesQuery.Close()
    }

    if queries.DeleteTokenSubtypesQuery != nil {
        queries.DeleteTokenSubtypesQuery.Close()
    }

    if queries.DeleteTokenSupertypesQuery != nil {
        queries.DeleteTokenSupertypesQuery.Close()
    }

    if queries.DeleteTokenReverseRelatedQuery != nil {
        queries.DeleteTokenReverseRelatedQuery.Close()
    }
}
