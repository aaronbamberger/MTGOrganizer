package carddb

import influx "github.com/influxdata/influxdb1-client/v2"
import "sync"
import "time"

type PricesUpdateStats struct {
    totalCardRecords int
    totalPriceRecords int
    mtgoPriceRecords int
    mtgoFoilPriceRecords int
    paperPriceRecords int
    paperFoilPriceRecords int
}

func (stats *PricesUpdateStats) AddToDb(client influx.Client) error {
    bpConfig := influx.BatchPointsConfig{Database: "mtg_cards"}
    bp, err := influx.NewBatchPoints(bpConfig)
    if err != nil {
        return err
    }

    fields := map[string]interface{} {
        "total_card_records": stats.totalCardRecords,
        "total_price_records": stats.totalPriceRecords,
        "mtgo_price_records": stats.mtgoPriceRecords,
        "mtgo_foil_price_records": stats.mtgoFoilPriceRecords,
        "paper_price_records": stats.paperPriceRecords,
        "paper_foil_price_records": stats.paperFoilPriceRecords}

    point, err := influx.NewPoint("price_updates",
        nil,
        fields,
        time.Now())

    bp.AddPoint(point)

    err = client.Write(bp)
    if err != nil {
        return err
    }

    return nil
}

func (stats *PricesUpdateStats) AddToTotalCardRecords(delta int) {
    stats.totalCardRecords += delta
}

func (stats *PricesUpdateStats) AddToTotalPriceRecords(delta int) {
    stats.totalPriceRecords += delta
}

func (stats *PricesUpdateStats) AddToMTGOPriceRecords(delta int) {
    stats.mtgoPriceRecords += delta
}

func (stats *PricesUpdateStats) AddToMTGOFoilPriceRecords(delta int) {
    stats.mtgoFoilPriceRecords += delta
}

func (stats *PricesUpdateStats) AddToPaperPriceRecords(delta int) {
    stats.paperPriceRecords += delta
}

func (stats *PricesUpdateStats) AddToPaperFoilPriceRecords(delta int) {
    stats.paperFoilPriceRecords += delta
}

func (stats *PricesUpdateStats) TotalCardRecords() int {
    return stats.totalCardRecords
}

func (stats *PricesUpdateStats) TotalPriceRecords() int {
    return stats.totalPriceRecords
}

func (stats *PricesUpdateStats) MTGOPriceRecords() int {
    return stats.mtgoPriceRecords
}

func (stats *PricesUpdateStats) MTGOFoilPriceRecords() int {
    return stats.mtgoFoilPriceRecords
}

func (stats *PricesUpdateStats) PaperPriceRecords() int {
    return stats.paperPriceRecords
}

func (stats *PricesUpdateStats) PaperFoilPriceRecords() int {
    return stats.paperFoilPriceRecords
}

type CardUpdateStats struct {
	mutex sync.RWMutex

    // Stats about sets in the update
	totalSets int

    // Stats about new sets in the update
	totalNewSets int

    // Stats about existing sets in the update
	totalExistingSets int
	existingSetsSkipped int
	existingSetsUpdated int

    // Stats about cards in the update
	totalCards int

    // Stats about new cards in the update
	totalNewCards int
	totalNewCardsInNewSets int
	totalNewCardsInExistingSets int

    // Stats about existing cards in the update
    totalExistingCards int
	existingCardsSkipped int
	existingCardsUpdated int

    // Stats about tokens in the update
    totalTokens int

    // Stats about new tokens in the update
    totalNewTokens int
    totalNewTokensInNewSets int
    totalNewTokensInExistingSets int

    // Stats about existing tokens in the update
    totalExistingTokens int
    existingTokensSkipped int
    existingTokensUpdated int
}

func (stats *CardUpdateStats) AddToDb(client influx.Client) error {
    stats.mutex.RLock()
    defer stats.mutex.RUnlock()

    bpConfig := influx.BatchPointsConfig{Database: "mtg_cards"}
    bp, err := influx.NewBatchPoints(bpConfig)
    if err != nil {
        return err
    }

    fields := map[string]interface{} {
        "total_sets": stats.totalSets,
        "total_new_sets": stats.totalNewSets,
        "total_existing_sets": stats.totalExistingSets,
        "existing_sets_skipped": stats.existingSetsSkipped,
        "existing_sets_updated": stats.existingSetsUpdated,
        "total_cards": stats.totalCards,
        "total_new_cards": stats.totalNewCards,
        "total_new_cards_in_new_sets": stats.totalNewCardsInNewSets,
        "total_new_cards_in_existing_sets": stats.totalNewCardsInExistingSets,
        "total_existing_cards": stats.totalExistingCards,
        "existing_cards_skipped": stats.existingCardsSkipped,
        "existing_cards_updated": stats.existingCardsUpdated,
        "total_tokens": stats.totalTokens,
        "total_new_tokens": stats.totalNewTokens,
        "total_new_tokens_in_new_sets": stats.totalNewTokensInNewSets,
        "total_new_tokens_in_existing_sets": stats.totalNewTokensInExistingSets,
        "total_existing_tokens": stats.totalExistingTokens,
        "existing_tokens_skipped": stats.existingTokensSkipped,
        "existing_tokens_updated": stats.existingTokensUpdated}

    point, err := influx.NewPoint("card_updates",
        nil,
        fields,
        time.Now())

    bp.AddPoint(point)

    err = client.Write(bp)
    if err != nil {
        return err
    }

    return nil
}

func (stats *CardUpdateStats) AddToTotalSets(delta int) {
    stats.mutex.Lock()
    defer stats.mutex.Unlock()
    stats.totalSets += delta
}

func (stats *CardUpdateStats) AddToTotalNewSets(delta int) {
    stats.mutex.Lock()
    defer stats.mutex.Unlock()
    stats.totalNewSets += delta
}

func (stats *CardUpdateStats) AddToTotalExistingSets(delta int) {
    stats.mutex.Lock()
    defer stats.mutex.Unlock()
    stats.totalExistingSets += delta
}

func (stats *CardUpdateStats) AddToExistingSetsSkipped(delta int) {
    stats.mutex.Lock()
    defer stats.mutex.Unlock()
    stats.existingSetsSkipped += delta
}

func (stats *CardUpdateStats) AddToExistingSetsUpdated(delta int) {
    stats.mutex.Lock()
    defer stats.mutex.Unlock()
    stats.existingSetsUpdated += delta
}

func (stats *CardUpdateStats) AddToTotalCards(delta int) {
    stats.mutex.Lock()
    defer stats.mutex.Unlock()
    stats.totalCards += delta
}

func (stats *CardUpdateStats) AddToTotalNewCards(delta int) {
    stats.mutex.Lock()
    defer stats.mutex.Unlock()
    stats.totalNewCards += delta
}

func (stats *CardUpdateStats) AddToTotalExistingCards(delta int) {
    stats.mutex.Lock()
    defer stats.mutex.Unlock()
    stats.totalExistingCards += delta
}

func (stats *CardUpdateStats) AddToTotalNewCardsInNewSets(delta int) {
    stats.mutex.Lock()
    defer stats.mutex.Unlock()
    stats.totalNewCardsInNewSets += delta
}

func (stats *CardUpdateStats) AddToTotalNewCardsInExistingSets(delta int) {
    stats.mutex.Lock()
    defer stats.mutex.Unlock()
    stats.totalNewCardsInExistingSets += delta
}

func (stats *CardUpdateStats) AddToExistingCardsSkipped(delta int) {
    stats.mutex.Lock()
    defer stats.mutex.Unlock()
    stats.existingCardsSkipped += delta
}

func (stats *CardUpdateStats) AddToExistingCardsUpdated(delta int) {
    stats.mutex.Lock()
    defer stats.mutex.Unlock()
    stats.existingCardsUpdated += delta
}

func (stats *CardUpdateStats) AddToTotalTokens(delta int) {
    stats.mutex.Lock()
    defer stats.mutex.Unlock()
    stats.totalTokens += delta
}

func (stats *CardUpdateStats) AddToTotalNewTokens(delta int) {
    stats.mutex.Lock()
    defer stats.mutex.Unlock()
    stats.totalNewTokens += delta
}

func (stats *CardUpdateStats) AddToTotalExistingTokens(delta int) {
    stats.mutex.Lock()
    defer stats.mutex.Unlock()
    stats.totalExistingTokens += delta
}

func (stats *CardUpdateStats) AddToTotalNewTokensInNewSets(delta int) {
    stats.mutex.Lock()
    defer stats.mutex.Unlock()
    stats.totalNewTokensInNewSets += delta
}

func (stats *CardUpdateStats) AddToTotalNewTokensInExistingSets(delta int) {
    stats.mutex.Lock()
    defer stats.mutex.Unlock()
    stats.totalNewTokensInExistingSets += delta
}

func (stats *CardUpdateStats) AddToExistingTokensSkipped(delta int) {
    stats.mutex.Lock()
    defer stats.mutex.Unlock()
    stats.existingTokensSkipped += delta
}

func (stats *CardUpdateStats) AddToExistingTokensUpdated(delta int) {
    stats.mutex.Lock()
    defer stats.mutex.Unlock()
    stats.existingTokensUpdated += delta
}

func (stats *CardUpdateStats) TotalSets() int {
    stats.mutex.RLock()
    defer stats.mutex.RUnlock()
    return stats.totalSets
}

func (stats *CardUpdateStats) TotalNewSets() int {
    stats.mutex.RLock()
    defer stats.mutex.RUnlock()
    return stats.totalNewSets
}

func (stats *CardUpdateStats) TotalExistingSets() int {
    stats.mutex.RLock()
    defer stats.mutex.RUnlock()
    return stats.totalExistingSets
}

func (stats *CardUpdateStats) ExistingSetsSkipped() int {
    stats.mutex.RLock()
    defer stats.mutex.RUnlock()
    return stats.existingSetsSkipped
}

func (stats *CardUpdateStats) ExistingSetsUpdated() int {
    stats.mutex.RLock()
    defer stats.mutex.RUnlock()
    return stats.existingSetsUpdated
}

func (stats *CardUpdateStats) TotalCards() int {
    stats.mutex.RLock()
    defer stats.mutex.RUnlock()
    return stats.totalCards
}

func (stats *CardUpdateStats) TotalNewCards() int {
    stats.mutex.RLock()
    defer stats.mutex.RUnlock()
    return stats.totalNewCards
}

func (stats *CardUpdateStats) TotalExistingCards() int {
    stats.mutex.RLock()
    defer stats.mutex.RUnlock()
    return stats.totalExistingCards
}

func (stats *CardUpdateStats) TotalNewCardsInNewSets() int {
    stats.mutex.RLock()
    defer stats.mutex.RUnlock()
    return stats.totalNewCardsInNewSets
}

func (stats *CardUpdateStats) TotalNewCardsInExistingSets() int {
    stats.mutex.RLock()
    defer stats.mutex.RUnlock()
    return stats.totalNewCardsInExistingSets
}

func (stats *CardUpdateStats) ExistingCardsSkipped() int {
    stats.mutex.RLock()
    defer stats.mutex.RUnlock()
    return stats.existingCardsSkipped
}

func (stats *CardUpdateStats) ExistingCardsUpdated() int {
    stats.mutex.RLock()
    defer stats.mutex.RUnlock()
    return stats.existingCardsUpdated
}

func (stats *CardUpdateStats) TotalTokens() int {
    stats.mutex.RLock()
    defer stats.mutex.RUnlock()
    return stats.totalTokens
}

func (stats *CardUpdateStats) TotalNewTokens() int {
    stats.mutex.RLock()
    defer stats.mutex.RUnlock()
    return stats.totalNewTokens
}

func (stats *CardUpdateStats) TotalExistingTokens() int {
    stats.mutex.RLock()
    defer stats.mutex.RUnlock()
    return stats.totalExistingTokens
}

func (stats *CardUpdateStats) TotalNewTokensInNewSets() int {
    stats.mutex.RLock()
    defer stats.mutex.RUnlock()
    return stats.totalNewTokensInNewSets
}

func (stats *CardUpdateStats) TotalNewTokensInExistingSets() int {
    stats.mutex.RLock()
    defer stats.mutex.RUnlock()
    return stats.totalNewTokensInExistingSets
}

func (stats *CardUpdateStats) ExistingTokensSkipped() int {
    stats.mutex.RLock()
    defer stats.mutex.RUnlock()
    return stats.existingTokensSkipped
}

func (stats *CardUpdateStats) ExistingTokensUpdated() int {
    stats.mutex.RLock()
    defer stats.mutex.RUnlock()
    return stats.existingTokensUpdated
}

func AddSingleUpdateStatsToDb(client influx.Client,
        cardsUpdated bool,
        pricesUpdated bool,
        cardUpdateDuration time.Duration,
        priceUpdateDuration time.Duration) error {

    bpConfig := influx.BatchPointsConfig{Database: "mtg_cards"}
    bp, err := influx.NewBatchPoints(bpConfig)
    if err != nil {
        return err
    }

    fields := map[string]interface{} {
        "cards_updated": cardsUpdated,
        "prices_updated": pricesUpdated,
        "card_update_duration": cardUpdateDuration.Seconds(),
        "price_update_duration": priceUpdateDuration.Seconds()}

    point, err := influx.NewPoint("updates",
        nil,
        fields,
        time.Now())

    bp.AddPoint(point)

    err = client.Write(bp)
    if err != nil {
        return err
    }

    return nil
}

