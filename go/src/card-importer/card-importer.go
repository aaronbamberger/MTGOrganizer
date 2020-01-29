package main

import "database/sql"
import _ "github.com/go-sql-driver/mysql"
import "log"
import "mtgcards"

func main() {
	allSets, err := mtgcards.DownloadAllPrintings(true)
	if err != nil {
		log.Fatal(err)
	}

	// Connect to the database
	db, err := sql.Open("mysql", "app_user:app_db_password@tcp(172.18.0.3)/mtg_cards")
	defer db.Close()
	if err != nil {
		log.Fatal(err)
	}

	db.SetMaxIdleConns(1000)

	stats, err := mtgcards.ImportSetsToDb(db, allSets)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Finished DB Update\n")
	log.Printf("Total sets in update: %d\n", stats.TotalSets())
	log.Printf("Total existing sets in the database: %d\n",
        stats.TotalExistingSets())
	log.Printf("Total new sets: %d\n", stats.TotalNewSets())
	log.Printf("Existing sets skipped due to hash match: %d\n",
        stats.ExistingSetsSkipped())
	log.Printf("Existing sets updated due to hash mismatch: %d\n",
        stats.ExistingSetsUpdated())
	log.Printf("Total cards in update: %d\n", stats.TotalCards())
	log.Printf("Total new cards: %d\n", stats.TotalNewCards())
	log.Printf("Total new cards in new sets: %d\n", stats.TotalNewCardsInNewSets())
	log.Printf("Total new cards in existing sets: %d\n",
        stats.TotalNewCardsInExistingSets())
	log.Printf("New atomic records for new cards: %d\n",
        stats.TotalNewAtomicRecordsForNewCards())
    log.Printf("New atomic records for existing cards: %d\n",
        stats.TotalNewAtomicRecordsForExistingCards())
    log.Printf("Existing atomic records for new cards: %d\n",
        stats.TotalExistingAtomicRecordsForNewCards())
    log.Printf("Existing atomic records for existing cards: %d\n",
        stats.TotalExistingAtomicRecordsForExistingCards())
	log.Printf("Total existing cards: %d\n", stats.TotalExistingCards())
	log.Printf("Existing cards skipped due to hash match: %d\n",
        stats.ExistingCardsSkipped())
	log.Printf("Existing cards updated due to hash mismatch: %d\n",
        stats.ExistingCardsUpdated())
}
