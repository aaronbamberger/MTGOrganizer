package backend

import "context"
import "fmt"
import "log"
import "database/sql"
import "strings"

import "mtgcards"

func sendError(done <-chan interface{}, respChan chan<- ResponseMessage, err error) {
    select {
    case <-done:
    case respChan <- ResponseMessage{Type: ErrorResponse, Value: err}:
    }
}

type CardSearchResult struct {
    Name string `json:"name"`
    UUID string `json:"uuid"`
    SetName string `json:"setName"`
    SetKeyruneCode string `json:"setKeyruneCode"`
}

func cardSearch(db *sql.DB,
        request string,
        done <-chan interface{},
        respChan chan<- ResponseMessage) {

    dbConn, err := db.Conn(context.Background())
    if err != nil {
        sendError(done, respChan, err)
        return
    }
    defer dbConn.Close()

    processedName := strings.ReplaceAll(request, "s", "'?s")
    res, err := dbConn.QueryContext(
        context.Background(),
        `SELECT DISTINCT
        all_cards.name, all_cards.uuid, sets.name, sets.keyrune_code
        FROM
        all_cards INNER JOIN sets ON all_cards.set_id = sets.set_id
        WHERE all_cards.name RLIKE ?
        ORDER BY all_cards.name ASC`,
        fmt.Sprintf("\\b%s", processedName))
    if err != nil {
        sendError(done, respChan, err)
        return
    }
    defer res.Close()

    cards := make([]CardSearchResult, 0)

    for res.Next() {
        var cardName string
        var cardUUID string
        var setName string
        var setKeyruneCode string
        err = res.Scan(&cardName, &cardUUID, &setName, &setKeyruneCode)
        if err != nil {
            sendError(done, respChan, err)
            return
        }
        cards = append(cards,
            CardSearchResult{
                Name: cardName,
                UUID: cardUUID,
                SetName: setName,
                SetKeyruneCode: setKeyruneCode})
    }
    if err = res.Err(); err != nil {
        sendError(done, respChan, err)
        return
    }

    select {
    case <-done:
    case respChan <- ResponseMessage{Type:CardSearchResponse, Value: cards}:
    }
}

type CardDetail struct {
    mtgcards.MTGCard
    CardId int `json:"card_id"`
    SetId int `json:"set_id"`
}

func cardDetail(db *sql.DB,
        uuid string,
        done <-chan interface{},
        respChan chan<- ResponseMessage) {

    dbConn, err := db.Conn(context.Background())
    if err != nil {
        log.Print(err)
        return
    }
    defer dbConn.Close()

    // Get the basic card info
    cardInfo := dbConn.QueryRowContext(
        context.Background(),
        `SELECT
        card_id,
        artist,
        border_color,
        card_number,
        card_power,
        card_type,
        converted_mana_cost,
        duel_deck,
        edhrec_rank,
        face_converted_mana_cost,
        flavor_text,
        frame_version,
        hand,
        has_foil,
        has_non_foil,
        is_alternative,
        is_arena,
        is_full_art,
        is_mtgo,
        is_online_only,
        is_oversized,
        is_paper,
        is_promo,
        is_reprint,
        is_reserved,
        is_starter,
        is_story_spotlight,
        is_textless,
        is_timeshifted,
        layout,
        life,
        loyalty,
        mana_cost,
        mcm_id,
        mcm_meta_id,
        mtg_arena_id,
        mtgo_foil_id,
        mtgo_id,
        mtgstocks_id,
        multiverse_id,
        name,
        original_text,
        original_type,
        rarity,
        scryfall_id,
        scryfall_illustration_id,
        scryfall_oracle_id,
        set_id,
        side,
        tcgplayer_product_id,
        text,
        toughness,
        watermark
        FROM
        all_cards
        WHERE uuid = ?`,
        uuid)

    var card CardDetail
    var duelDeck sql.NullString
    var edhrecRank sql.NullInt64
    var flavorText sql.NullString
    var hand sql.NullString
    var life sql.NullString
    var loyalty sql.NullString
    var mtgArenaId sql.NullInt64
    var mtgoFoilId sql.NullInt64
    var mtgoId sql.NullInt64
    var name sql.NullString
    var scryfallIllustrationId sql.NullString
    var side sql.NullString
    err = cardInfo.Scan(&card.CardId,
        &card.Artist,
        &card.BorderColor,
        &card.Number,
        &card.Power,
        &card.Type,
        //&card.ColorIdentity,
        //&card.ColorIndicator,
        //&card.Colors,
        &card.ConvertedManaCost,
        &duelDeck,
        &edhrecRank,
        &card.FaceConvertedManaCost,
        &flavorText,
        &card.FrameVersion,
        &hand,
        &card.HasFoil,
        &card.HasNonFoil,
        &card.IsAlternative,
        &card.IsArena,
        &card.IsFullArt,
        &card.IsMTGO,
        &card.IsOnlineOnly,
        &card.IsOversized,
        &card.IsPaper,
        &card.IsPromo,
        &card.IsReprint,
        &card.IsReserved,
        &card.IsStarter,
        &card.IsStorySpotlight,
        &card.IsTextless,
        &card.IsTimeshifted,
        &card.Layout,
        &life,
        &loyalty,
        &card.ManaCost,
        &card.MCMId,
        &card.MCMMetaId,
        &mtgArenaId,
        &mtgoFoilId,
        &mtgoId,
        &card.MTGStocksId,
        &card.MultiverseId,
        &name,
        &card.OriginalText,
        &card.OriginalType,
        &card.Rarity,
        &card.ScryfallId,
        &scryfallIllustrationId,
        &card.ScryfallOracleId,
        &card.SetId,
        &side,
        &card.TCGPlayerProductId,
        &card.Text,
        &card.Toughness,
        &card.Watermark)
    if err != nil {
        log.Printf("Error parsing basic card info: %s", err)
        sendError(done, respChan, err)
        return
    }

    // Populate the optional fields into the card if they're not null
    if duelDeck.Valid {
        card.DuelDeck = duelDeck.String
    }
    if edhrecRank.Valid {
        card.EDHRecRank = int(edhrecRank.Int64)
    }
    if flavorText.Valid {
        card.FlavorText = flavorText.String
    }
    if hand.Valid {
        card.Hand = hand.String
    }
    if life.Valid {
        card.Life = life.String
    }
    if loyalty.Valid {
        card.Loyalty = loyalty.String
    }
    if mtgArenaId.Valid {
        card.MTGArenaId = int(mtgArenaId.Int64)
    }
    if mtgoFoilId.Valid {
        card.MTGOFoilId = int(mtgoFoilId.Int64)
    }
    if mtgoId.Valid {
        card.MTGOId = int(mtgoId.Int64)
    }
    if name.Valid {
        card.Name = name.String
    }
    if scryfallIllustrationId.Valid {
        card.ScryfallIllustrationId = scryfallIllustrationId.String
    }
    if side.Valid {
        card.Side = side.String
    }

    // Get the card printings
    printings, err := dbConn.QueryContext(
        context.Background(),
        `SELECT set_code
        FROM card_printings
        WHERE card_id = ?`,
        card.CardId)
    if err != nil {
        log.Printf("Error getting card printings: %s", err)
        sendError(done, respChan, err)
        return
    }
    card.Printings = make([]string, 0)
    for printings.Next() {
        var setCode string
        err = printings.Scan(&setCode)
        if err != nil {
            sendError(done, respChan, err)
            printings.Close()
            return
        }
        card.Printings = append(card.Printings, setCode)
    }
    printings.Close()
    if err = printings.Err(); err != nil {
        sendError(done, respChan, err)
        return
    }

    // Get the card variations
    variations, err := dbConn.QueryContext(
        context.Background(),
        `SELECT variation_uuid
        FROM variations
        WHERE card_id = ?`,
        card.CardId)
    if err != nil {
        log.Printf("Error getting card variations: %s", err)
        sendError(done, respChan, err)
        return
    }
    card.Variations = make([]string, 0)
    for variations.Next() {
        var variationUUID string
        err = variations.Scan(&variationUUID)
        if err != nil {
            sendError(done, respChan, err)
            variations.Close()
            return
        }
        card.Variations = append(card.Variations, variationUUID)
    }
    variations.Close()
    if err = variations.Err(); err != nil {
        sendError(done, respChan, err)
        return
    }

    // Get the format legalities
    legalities, err := dbConn.QueryContext(
        context.Background(),
        `SELECT game_formats.game_format_name, legality_options.legality_option_name
        FROM
        legalities INNER JOIN game_formats
        ON legalities.game_format_id = game_formats.game_format_id
        INNER JOIN legality_options
        ON legalities.legality_option_id = legality_options.legality_option_id
        WHERE legalities.card_id = ?`,
        card.CardId)
    if err != nil {
        sendError(done, respChan, err)
        return
    }
    card.Legalities = make(map[string]string)
    for legalities.Next() {
        var gameFormat string
        var legalityOption string
        err = legalities.Scan(&gameFormat, &legalityOption)
        if err != nil {
            sendError(done, respChan, err)
            legalities.Close()
            return
        }
        card.Legalities[gameFormat] = legalityOption
    }
    legalities.Close()
    if err = legalities.Err(); err != nil {
        sendError(done, respChan, err)
        return
    }

    // Get the leadership skills
    leadershipSkills, err := dbConn.QueryContext(
        context.Background(),
        `SELECT leadership_formats.leadership_format_name, leadership_skills.leader_legal
        FROM
        leadership_skills INNER JOIN leadership_formats
        ON leadership_skills.leadership_format_id = leadership_formats.leadership_format_id
        WHERE leadership_skills.card_id = ?`,
        card.CardId)
    if err != nil {
        log.Printf("Error getting leadership skills: %s", err)
        sendError(done, respChan, err)
        return
    }
    card.LeadershipSkills = make(map[string]bool)
    for leadershipSkills.Next() {
        var leadershipFormat string
        var leaderLegal bool
        err = leadershipSkills.Scan(&leadershipFormat, &leaderLegal)
        if err != nil {
            sendError(done, respChan, err)
            leadershipSkills.Close()
            return
        }
        card.LeadershipSkills[leadershipFormat] = leaderLegal
    }
    leadershipSkills.Close()
    if err = leadershipSkills.Err(); err != nil {
        sendError(done, respChan, err)
        return
    }

    select {
    case <-done:
    case respChan <- ResponseMessage{Type: CardDetailResponse, Value: card}:
    }
}
