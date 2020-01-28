package mtgcards

import "database/sql"

type dbQueries struct {
    SetHashQuery *sql.Stmt
    CardHashQuery *sql.Stmt
    InsertSetQuery *sql.Stmt
    UpdateSetQuery *sql.Stmt
    InsertSetTranslationQuery *sql.Stmt
	NumAtomicPropertiesQuery *sql.Stmt
	AtomicPropertiesIdQuery *sql.Stmt
	AtomicPropertiesHashQuery *sql.Stmt
	InsertAtomicPropertiesQuery *sql.Stmt
	InsertCardQuery *sql.Stmt
	UpdateCardQuery *sql.Stmt
	InsertAltLangDataQuery *sql.Stmt
	InsertCardPrintingQuery *sql.Stmt
	InsertCardSubtypeQuery *sql.Stmt
	InsertCardSupertypeQuery *sql.Stmt
	InsertFrameEffectQuery *sql.Stmt
	InsertLeadershipSkillQuery *sql.Stmt
	InsertLegalityQuery *sql.Stmt
	InsertOtherFaceIdQuery *sql.Stmt
	InsertPurchaseURLQuery *sql.Stmt
	InsertRulingQuery *sql.Stmt
	InsertVariationQuery *sql.Stmt
	InsertBaseTypeQuery *sql.Stmt
	DeleteFrameEffectQuery *sql.Stmt
	DeleteOtherFaceQuery *sql.Stmt
	DeleteVariationQuery *sql.Stmt
	UpdateRefCntQuery *sql.Stmt
}

func (queries *dbQueries) queriesForTx(tx *sql.Tx) *dbQueries {
    var txQueries dbQueries
    txQueries.SetHashQuery = tx.Stmt(queries.SetHashQuery)
    txQueries.CardHashQuery = tx.Stmt(queries.CardHashQuery)
    txQueries.InsertSetQuery = tx.Stmt(queries.InsertSetQuery)
    txQueries.UpdateSetQuery = tx.Stmt(queries.UpdateSetQuery)
    txQueries.InsertSetTranslationQuery = tx.Stmt(queries.InsertSetTranslationQuery)
    txQueries.NumAtomicPropertiesQuery = tx.Stmt(queries.NumAtomicPropertiesQuery)
    txQueries.AtomicPropertiesIdQuery = tx.Stmt(queries.AtomicPropertiesIdQuery)
    txQueries.AtomicPropertiesHashQuery = tx.Stmt(queries.AtomicPropertiesHashQuery)
    txQueries.InsertAtomicPropertiesQuery = tx.Stmt(queries.InsertAtomicPropertiesQuery)
    txQueries.InsertCardQuery = tx.Stmt(queries.InsertCardQuery)
    txQueries.UpdateCardQuery = tx.Stmt(queries.UpdateCardQuery)
    txQueries.InsertAltLangDataQuery = tx.Stmt(queries.InsertAltLangDataQuery)
    txQueries.InsertCardPrintingQuery = tx.Stmt(queries.InsertCardPrintingQuery)
    txQueries.InsertCardSubtypeQuery = tx.Stmt(queries.InsertCardSubtypeQuery)
    txQueries.InsertCardSupertypeQuery = tx.Stmt(queries.InsertCardSupertypeQuery)
    txQueries.InsertFrameEffectQuery = tx.Stmt(queries.InsertFrameEffectQuery)
    txQueries.InsertLeadershipSkillQuery = tx.Stmt(queries.InsertLeadershipSkillQuery)
    txQueries.InsertLegalityQuery = tx.Stmt(queries.InsertLegalityQuery)
    txQueries.InsertOtherFaceIdQuery = tx.Stmt(queries.InsertOtherFaceIdQuery)
    txQueries.InsertPurchaseURLQuery = tx.Stmt(queries.InsertPurchaseURLQuery)
    txQueries.InsertRulingQuery = tx.Stmt(queries.InsertRulingQuery)
    txQueries.InsertVariationQuery = tx.Stmt(queries.InsertVariationQuery)
    txQueries.InsertBaseTypeQuery = tx.Stmt(queries.InsertBaseTypeQuery)
    txQueries.DeleteFrameEffectQuery = tx.Stmt(queries.DeleteFrameEffectQuery)
    txQueries.DeleteOtherFaceQuery = tx.Stmt(queries.DeleteOtherFaceQuery)
    txQueries.DeleteVariationQuery = tx.Stmt(queries.DeleteVariationQuery)
    txQueries.UpdateRefCntQuery = tx.Stmt(queries.UpdateRefCntQuery)

    return &txQueries
}

func prepareDbQueries(db *sql.DB) (*dbQueries, error) {
	var err error
    var queries dbQueries

    queries.SetHashQuery, err = db.Prepare(`SELECT set_hash, set_id FROM sets WHERE code = ?`)
	if err != nil {
        return &queries, err
	}

	queries.CardHashQuery, err = db.Prepare(`SELECT full_card_hash, atomic_card_data_id
        FROM all_cards WHERE uuid = ?`)
	if err != nil {
        return &queries, err
	}

	queries.InsertSetQuery, err = db.Prepare(`INSERT INTO sets
		(set_hash, base_size, block_name, code, is_foreign_only, is_foil_only,
		is_online_only, is_partial_preview, keyrune_code, mcm_name, mcm_id,
		mtgo_code, name, parent_code, release_date, tcgplayer_group_id,
		total_set_size, set_type)
		VALUES
		(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
        return &queries, err
	}

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
        return &queries, err
	}

	queries.InsertSetTranslationQuery, err = db.Prepare(`INSERT INTO set_translations
		(set_id, set_translation_language_id, set_translated_name)
		VALUES
		(?, ?, ?)`)
	if err != nil {
		return &queries, err
	}

	queries.NumAtomicPropertiesQuery, err = db.Prepare(`SELECT COUNT(scryfall_oracle_id)
		FROM atomic_card_data
		WHERE card_data_hash = ?`)
	if err != nil {
		return &queries, err
	}

	queries.AtomicPropertiesIdQuery, err = db.Prepare(`SELECT atomic_card_data_id,
		ref_cnt,
		scryfall_oracle_id
		FROM atomic_card_data
		WHERE card_data_hash = ?`)
	if err != nil {
		return &queries, err
	}

	queries.AtomicPropertiesHashQuery, err = db.Prepare(`SELECT card_data_hash
		FROM atomic_card_data
		WHERE atomic_card_data_id = ?`)
	if err != nil {
		return &queries, err
	}

	queries.InsertAtomicPropertiesQuery, err = db.Prepare(`INSERT INTO atomic_card_data
		(card_data_hash, color_identity, color_indicator, colors, converted_mana_cost,
		edhrec_rank, face_converted_mana_cost, hand, is_reserved, layout, life,
		loyalty, mana_cost, mtgstocks_id, name, card_power, scryfall_oracle_id,
		side, text, toughness, card_type)
		VALUES
		(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return &queries, err
	}

	queries.InsertCardQuery, err = db.Prepare(`INSERT INTO all_cards
		(uuid, full_card_hash, atomic_card_data_id, set_id, artist, border_color,
		card_number, scryfall_id, watermark, frame_version, mcm_id, mcm_meta_id,
		multiverse_id, original_text, original_type, rarity, tcgplayer_product_id,
		duel_deck, flavor_text, has_foil, has_non_foil, is_alternative, is_arena,
		is_full_art, is_mtgo, is_online_only, is_oversized, is_paper, is_promo,
		is_reprint, is_starter, is_story_spotlight, is_textless, is_timeshifted,
		mtg_arena_id, mtgo_foil_id, mtgo_id, scryfall_illustration_id)
		VALUES
		(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?,
		?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return &queries, err
	}

	queries.UpdateCardQuery, err = db.Prepare(`UPDATE all_cards SET
		full_card_hash = ?,
		atomic_card_data_id = ?,
		set_id = ?,
		artist = ?,
		border_color = ?,
		card_number = ?,
		scryfall_id = ?,
		watermark = ?,
		frame_version = ?,
		mcm_id = ?,
		mcm_meta_id = ?,
		multiverse_id = ?,
		original_text = ?,
		original_type = ?,
		rarity = ?,
		tcgplayer_product_id = ?,
		duel_deck = ?,
		flavor_text = ?,
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
		is_starter = ?,
		is_story_spotlight = ?,
		is_textless = ?,
		is_timeshifted = ?,
		mtg_arena_id = ?,
		mtgo_foil_id = ?,
		mtgo_id = ?,
		scryfall_illustration_id = ?
		WHERE uuid = ?`)
    if err != nil {
        return &queries, err
    }

	queries.InsertAltLangDataQuery, err = db.Prepare(`INSERT INTO alternate_language_data
		(atomic_card_data_id, flavor_text, language, multiverse_id, name, text, card_type)
		VALUES
		(?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return &queries, err
	}

	queries.InsertCardPrintingQuery, err = db.Prepare(`INSERT INTO card_printings
		(atomic_card_data_id, set_code)
		VALUES
		(?, ?)`)
	if err != nil {
		return &queries, err
	}

	queries.InsertCardSubtypeQuery, err = db.Prepare(`INSERT INTO card_subtypes
		(atomic_card_data_id, subtype_option_id)
		VALUES
		(?, ?)`)
	if err != nil {
		return &queries, err
	}

	queries.InsertCardSupertypeQuery, err = db.Prepare(`INSERT INTO card_supertypes
		(atomic_card_data_id, supertype_option_id)
		VALUES
		(?, ?)`)
	if err != nil {
		return &queries, err
	}

	queries.InsertFrameEffectQuery, err = db.Prepare(`INSERT INTO frame_effects
		(card_uuid, frame_effect_option_id)
		VALUES
		(?, ?)`)
	if err != nil {
		return &queries, err
	}

	queries.InsertLeadershipSkillQuery, err = db.Prepare(`INSERT INTO leadership_skills
		(atomic_card_data_id, leadership_format_id, leader_legal)
		VALUES
		(?, ?, ?)`)
	if err != nil {
		return &queries, err
	}

	queries.InsertLegalityQuery, err = db.Prepare(`INSERT INTO legalities
		(atomic_card_data_id, game_format_id, legality_option_id)
		VALUES
		(?, ?, ?)`)
	if err != nil {
		return &queries, err
	}

	queries.InsertOtherFaceIdQuery, err = db.Prepare(`INSERT INTO other_faces
		(card_uuid, other_face_uuid)
		VALUES
		(?, ?)`)
	if err != nil {
		return &queries, err
	}

	queries.InsertPurchaseURLQuery, err = db.Prepare(`INSERT INTO purchase_urls
		(atomic_card_data_id, purchase_site_id, purchase_url)
		VALUES
		(?, ?, ?)`)
	if err != nil {
		return &queries, err
	}

	queries.InsertRulingQuery, err = db.Prepare(`INSERT INTO rulings
		(atomic_card_data_id, ruling_date, ruling_text)
		VALUES
		(?, ?, ?)`)
	if err != nil {
		return &queries, err
	}

	queries.InsertVariationQuery, err = db.Prepare(`INSERT INTO variations
		(card_uuid, variation_uuid)
		VALUES
		(?, ?)`)
	if err != nil {
		return &queries, err
	}

	queries.InsertBaseTypeQuery, err = db.Prepare(`INSERT INTO base_types
		(atomic_card_data_id, base_type_option_id)
		VALUES
		(?, ?)`)
	if err != nil {
		return &queries, err
	}

	queries.DeleteFrameEffectQuery, err = db.Prepare(`DELETE FROM frame_effects
		WHERE card_uuid = ?`)
	if err != nil {
		return &queries, err
	}

	queries.DeleteOtherFaceQuery, err = db.Prepare(`DELETE FROM other_faces
		WHERE card_uuid = ?`)
	if err != nil {
		return &queries, err
	}

	queries.DeleteVariationQuery, err = db.Prepare(`DELETE FROM variations
		WHERE card_uuid = ?`)
	if err != nil {
		return &queries, err
	}

	queries.UpdateRefCntQuery, err = db.Prepare(`UPDATE atomic_card_data
		SET ref_cnt = ?
		WHERE atomic_card_data_id = ?`)
    if err != nil {
        return &queries, nil
    }

	return &queries, nil
}

func (queries *dbQueries) cleanupDbQueries() {
    if queries.SetHashQuery != nil {
        queries.SetHashQuery.Close()
    }

    if queries.CardHashQuery != nil {
        queries.CardHashQuery.Close()
    }

    if queries.InsertSetQuery != nil {
        queries.InsertSetQuery.Close()
    }

    if queries.UpdateSetQuery != nil {
        queries.UpdateSetQuery.Close()
    }

    if queries.InsertSetTranslationQuery != nil {
        queries.InsertSetTranslationQuery.Close()
    }

	if queries.NumAtomicPropertiesQuery != nil {
		queries.NumAtomicPropertiesQuery.Close()
	}

	if queries.AtomicPropertiesIdQuery != nil {
		queries.AtomicPropertiesIdQuery.Close()
	}

	if queries.AtomicPropertiesHashQuery != nil {
		queries.AtomicPropertiesHashQuery.Close()
	}

	if queries.InsertAtomicPropertiesQuery != nil {
		queries.InsertAtomicPropertiesQuery.Close()
	}

	if queries.InsertCardQuery != nil {
		queries.InsertCardQuery.Close()
	}

	if queries.UpdateCardQuery != nil {
		queries.UpdateCardQuery.Close()
	}

	if queries.InsertAltLangDataQuery != nil {
		queries.InsertAltLangDataQuery.Close()
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

	if queries.InsertOtherFaceIdQuery != nil {
		queries.InsertOtherFaceIdQuery.Close()
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

	if queries.InsertBaseTypeQuery != nil {
		queries.InsertBaseTypeQuery.Close()
	}

	if queries.DeleteFrameEffectQuery != nil {
		queries.DeleteFrameEffectQuery.Close()
	}

	if queries.DeleteOtherFaceQuery != nil {
		queries.DeleteOtherFaceQuery.Close()
	}

	if queries.DeleteVariationQuery != nil {
		queries.DeleteVariationQuery.Close()
	}

	if queries.UpdateRefCntQuery != nil {
		queries.UpdateRefCntQuery.Close()
	}
}
