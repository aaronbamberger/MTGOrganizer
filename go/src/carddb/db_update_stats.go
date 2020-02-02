package carddb

import "sync"

type DBUpdateStats struct {
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

func (stats *DBUpdateStats) AddToTotalSets(delta int) {
    stats.mutex.Lock()
    defer stats.mutex.Unlock()
    stats.totalSets += delta
}

func (stats *DBUpdateStats) AddToTotalNewSets(delta int) {
    stats.mutex.Lock()
    defer stats.mutex.Unlock()
    stats.totalNewSets += delta
}

func (stats *DBUpdateStats) AddToTotalExistingSets(delta int) {
    stats.mutex.Lock()
    defer stats.mutex.Unlock()
    stats.totalExistingSets += delta
}

func (stats *DBUpdateStats) AddToExistingSetsSkipped(delta int) {
    stats.mutex.Lock()
    defer stats.mutex.Unlock()
    stats.existingSetsSkipped += delta
}

func (stats *DBUpdateStats) AddToExistingSetsUpdated(delta int) {
    stats.mutex.Lock()
    defer stats.mutex.Unlock()
    stats.existingSetsUpdated += delta
}

func (stats *DBUpdateStats) AddToTotalCards(delta int) {
    stats.mutex.Lock()
    defer stats.mutex.Unlock()
    stats.totalCards += delta
}

func (stats *DBUpdateStats) AddToTotalNewCards(delta int) {
    stats.mutex.Lock()
    defer stats.mutex.Unlock()
    stats.totalNewCards += delta
}

func (stats *DBUpdateStats) AddToTotalExistingCards(delta int) {
    stats.mutex.Lock()
    defer stats.mutex.Unlock()
    stats.totalExistingCards += delta
}

func (stats *DBUpdateStats) AddToTotalNewCardsInNewSets(delta int) {
    stats.mutex.Lock()
    defer stats.mutex.Unlock()
    stats.totalNewCardsInNewSets += delta
}

func (stats *DBUpdateStats) AddToTotalNewCardsInExistingSets(delta int) {
    stats.mutex.Lock()
    defer stats.mutex.Unlock()
    stats.totalNewCardsInExistingSets += delta
}

func (stats *DBUpdateStats) AddToExistingCardsSkipped(delta int) {
    stats.mutex.Lock()
    defer stats.mutex.Unlock()
    stats.existingCardsSkipped += delta
}

func (stats *DBUpdateStats) AddToExistingCardsUpdated(delta int) {
    stats.mutex.Lock()
    defer stats.mutex.Unlock()
    stats.existingCardsUpdated += delta
}

func (stats *DBUpdateStats) AddToTotalTokens(delta int) {
    stats.mutex.Lock()
    defer stats.mutex.Unlock()
    stats.totalTokens += delta
}

func (stats *DBUpdateStats) AddToTotalNewTokens(delta int) {
    stats.mutex.Lock()
    defer stats.mutex.Unlock()
    stats.totalNewTokens += delta
}

func (stats *DBUpdateStats) AddToTotalExistingTokens(delta int) {
    stats.mutex.Lock()
    defer stats.mutex.Unlock()
    stats.totalExistingTokens += delta
}

func (stats *DBUpdateStats) AddToTotalNewTokensInNewSets(delta int) {
    stats.mutex.Lock()
    defer stats.mutex.Unlock()
    stats.totalNewTokensInNewSets += delta
}

func (stats *DBUpdateStats) AddToTotalNewTokensInExistingSets(delta int) {
    stats.mutex.Lock()
    defer stats.mutex.Unlock()
    stats.totalNewTokensInExistingSets += delta
}

func (stats *DBUpdateStats) AddToExistingTokensSkipped(delta int) {
    stats.mutex.Lock()
    defer stats.mutex.Unlock()
    stats.existingTokensSkipped += delta
}

func (stats *DBUpdateStats) AddToExistingTokensUpdated(delta int) {
    stats.mutex.Lock()
    defer stats.mutex.Unlock()
    stats.existingTokensUpdated += delta
}

func (stats *DBUpdateStats) TotalSets() int {
    stats.mutex.RLock()
    defer stats.mutex.RUnlock()
    return stats.totalSets
}

func (stats *DBUpdateStats) TotalNewSets() int {
    stats.mutex.RLock()
    defer stats.mutex.RUnlock()
    return stats.totalNewSets
}

func (stats *DBUpdateStats) TotalExistingSets() int {
    stats.mutex.RLock()
    defer stats.mutex.RUnlock()
    return stats.totalExistingSets
}

func (stats *DBUpdateStats) ExistingSetsSkipped() int {
    stats.mutex.RLock()
    defer stats.mutex.RUnlock()
    return stats.existingSetsSkipped
}

func (stats *DBUpdateStats) ExistingSetsUpdated() int {
    stats.mutex.RLock()
    defer stats.mutex.RUnlock()
    return stats.existingSetsUpdated
}

func (stats *DBUpdateStats) TotalCards() int {
    stats.mutex.RLock()
    defer stats.mutex.RUnlock()
    return stats.totalCards
}

func (stats *DBUpdateStats) TotalNewCards() int {
    stats.mutex.RLock()
    defer stats.mutex.RUnlock()
    return stats.totalNewCards
}

func (stats *DBUpdateStats) TotalExistingCards() int {
    stats.mutex.RLock()
    defer stats.mutex.RUnlock()
    return stats.totalExistingCards
}

func (stats *DBUpdateStats) TotalNewCardsInNewSets() int {
    stats.mutex.RLock()
    defer stats.mutex.RUnlock()
    return stats.totalNewCardsInNewSets
}

func (stats *DBUpdateStats) TotalNewCardsInExistingSets() int {
    stats.mutex.RLock()
    defer stats.mutex.RUnlock()
    return stats.totalNewCardsInExistingSets
}

func (stats *DBUpdateStats) ExistingCardsSkipped() int {
    stats.mutex.RLock()
    defer stats.mutex.RUnlock()
    return stats.existingCardsSkipped
}

func (stats *DBUpdateStats) ExistingCardsUpdated() int {
    stats.mutex.RLock()
    defer stats.mutex.RUnlock()
    return stats.existingCardsUpdated
}

func (stats *DBUpdateStats) TotalTokens() int {
    stats.mutex.RLock()
    defer stats.mutex.RUnlock()
    return stats.totalTokens
}

func (stats *DBUpdateStats) TotalNewTokens() int {
    stats.mutex.RLock()
    defer stats.mutex.RUnlock()
    return stats.totalNewTokens
}

func (stats *DBUpdateStats) TotalExistingTokens() int {
    stats.mutex.RLock()
    defer stats.mutex.RUnlock()
    return stats.totalExistingTokens
}

func (stats *DBUpdateStats) TotalNewTokensInNewSets() int {
    stats.mutex.RLock()
    defer stats.mutex.RUnlock()
    return stats.totalNewTokensInNewSets
}

func (stats *DBUpdateStats) TotalNewTokensInExistingSets() int {
    stats.mutex.RLock()
    defer stats.mutex.RUnlock()
    return stats.totalNewTokensInExistingSets
}

func (stats *DBUpdateStats) ExistingTokensSkipped() int {
    stats.mutex.RLock()
    defer stats.mutex.RUnlock()
    return stats.existingTokensSkipped
}

func (stats *DBUpdateStats) ExistingTokensUpdated() int {
    stats.mutex.RLock()
    defer stats.mutex.RUnlock()
    return stats.existingTokensUpdated
}
