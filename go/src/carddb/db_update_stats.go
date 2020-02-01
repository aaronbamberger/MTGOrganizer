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
