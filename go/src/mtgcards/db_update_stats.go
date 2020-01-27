package mtgcards

import "sync"

type DbUpdateStats struct {
	mutex sync.RWMutex

	totalSets int
	totalNewSets int
	totalExistingSets int

	existingSetsSkipped int
	existingSetsUpdated int

	totalCards int
	totalNewCards int
	totalNewAtomicCards int
    totalExistingCards int

	totalNewCardsInNewSets int
	totalNewCardsInExistingSets int

	existingCardsSkipped int
	existingCardsUpdated int
}

func (stats *DbUpdateStats) AddToTotalSets(delta int) {
    stats.mutex.Lock()
    defer stats.mutex.Unlock()
    stats.totalSets += delta
}

func (stats *DbUpdateStats) AddToTotalNewSets(delta int) {
    stats.mutex.Lock()
    defer stats.mutex.Unlock()
    stats.totalNewSets += delta
}

func (stats *DbUpdateStats) AddToTotalExistingSets(delta int) {
    stats.mutex.Lock()
    defer stats.mutex.Unlock()
    stats.totalExistingSets += delta
}

func (stats *DbUpdateStats) AddToExistingSetsSkipped(delta int) {
    stats.mutex.Lock()
    defer stats.mutex.Unlock()
    stats.existingSetsSkipped += delta
}

func (stats *DbUpdateStats) AddToExistingSetsUpdated(delta int) {
    stats.mutex.Lock()
    defer stats.mutex.Unlock()
    stats.existingSetsUpdated += delta
}

func (stats *DbUpdateStats) AddToTotalCards(delta int) {
    stats.mutex.Lock()
    defer stats.mutex.Unlock()
    stats.totalCards += delta
}

func (stats *DbUpdateStats) AddToTotalNewCards(delta int) {
    stats.mutex.Lock()
    defer stats.mutex.Unlock()
    stats.totalNewCards += delta
}

func (stats *DbUpdateStats) AddToTotalNewAtomicCards(delta int) {
    stats.mutex.Lock()
    defer stats.mutex.Unlock()
    stats.totalNewAtomicCards += delta
}

func (stats *DbUpdateStats) AddToTotalExistingCards(delta int) {
    stats.mutex.Lock()
    defer stats.mutex.Unlock()
    stats.totalExistingCards += delta
}

func (stats *DbUpdateStats) AddToTotalNewCardsInNewSets(delta int) {
    stats.mutex.Lock()
    defer stats.mutex.Unlock()
    stats.totalNewCardsInNewSets += delta
}

func (stats *DbUpdateStats) AddToTotalNewCardsInExistingSets(delta int) {
    stats.mutex.Lock()
    defer stats.mutex.Unlock()
    stats.totalNewCardsInExistingSets += delta
}

func (stats *DbUpdateStats) AddToExistingCardsSkipped(delta int) {
    stats.mutex.Lock()
    defer stats.mutex.Unlock()
    stats.existingCardsSkipped += delta
}

func (stats *DbUpdateStats) AddToExistingCardsUpdated(delta int) {
    stats.mutex.Lock()
    defer stats.mutex.Unlock()
    stats.existingCardsUpdated += delta
}

func (stats *DbUpdateStats) TotalSets() int {
    stats.mutex.RLock()
    defer stats.mutex.RUnlock()
    return stats.totalSets
}

func (stats *DbUpdateStats) TotalNewSets() int {
    stats.mutex.RLock()
    defer stats.mutex.RUnlock()
    return stats.totalNewSets
}

func (stats *DbUpdateStats) TotalExistingSets() int {
    stats.mutex.RLock()
    defer stats.mutex.RUnlock()
    return stats.totalExistingSets
}

func (stats *DbUpdateStats) ExistingSetsSkipped() int {
    stats.mutex.RLock()
    defer stats.mutex.RUnlock()
    return stats.existingSetsSkipped
}

func (stats *DbUpdateStats) ExistingSetsUpdated() int {
    stats.mutex.RLock()
    defer stats.mutex.RUnlock()
    return stats.existingSetsUpdated
}

func (stats *DbUpdateStats) TotalCards() int {
    stats.mutex.RLock()
    defer stats.mutex.RUnlock()
    return stats.totalCards
}

func (stats *DbUpdateStats) TotalNewCards() int {
    stats.mutex.RLock()
    defer stats.mutex.RUnlock()
    return stats.totalNewCards
}

func (stats *DbUpdateStats) TotalNewAtomicCards() int {
    stats.mutex.RLock()
    defer stats.mutex.RUnlock()
    return stats.totalNewAtomicCards
}

func (stats *DbUpdateStats) TotalExistingCards() int {
    stats.mutex.RLock()
    defer stats.mutex.RUnlock()
    return stats.totalExistingCards
}

func (stats *DbUpdateStats) TotalNewCardsInNewSets() int {
    stats.mutex.RLock()
    defer stats.mutex.RUnlock()
    return stats.totalNewCardsInNewSets
}

func (stats *DbUpdateStats) TotalNewCardsInExistingSets() int {
    stats.mutex.RLock()
    defer stats.mutex.RUnlock()
    return stats.totalNewCardsInExistingSets
}

func (stats *DbUpdateStats) ExistingCardsSkipped() int {
    stats.mutex.RLock()
    defer stats.mutex.RUnlock()
    return stats.existingCardsSkipped
}

func (stats *DbUpdateStats) ExistingCardsUpdated() int {
    stats.mutex.RLock()
    defer stats.mutex.RUnlock()
    return stats.existingCardsUpdated
}
