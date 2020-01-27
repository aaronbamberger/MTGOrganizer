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

	dbUpdateStats, err := mtgcards.ImportSetsToDb(db, allSets)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Finished DB Update\n")
	log.Printf("Total sets in update: %d\n", dbUpdateStats.TotalSets())
	log.Printf("Total existing sets in the database: %d\n", dbUpdateStats.TotalExistingSets())
	log.Printf("Total new sets: %d\n", dbUpdateStats.TotalNewSets())
	log.Printf("Existing sets skipped due to hash match: %d\n", dbUpdateStats.ExistingSetsSkipped())
	log.Printf("Existing sets updated due to hash mismatch: %d\n", dbUpdateStats.ExistingSetsUpdated())
	log.Printf("Total cards in update: %d\n", dbUpdateStats.TotalCards())
	log.Printf("Total new cards: %d\n", dbUpdateStats.TotalNewCards())
	log.Printf("Total new cards in new sets: %d\n", dbUpdateStats.TotalNewCardsInNewSets())
	log.Printf("Total new cards in existing sets: %d\n", dbUpdateStats.TotalNewCardsInExistingSets())
	log.Printf("Total new atomic cards: %d\n", dbUpdateStats.TotalNewAtomicCards())
	log.Printf("Total existing cards: %d\n", dbUpdateStats.TotalExistingCards())
	log.Printf("Existing cards skipped due to hash match: %d\n", dbUpdateStats.ExistingCardsSkipped())
	log.Printf("Existing cards updated due to hash mismatch: %d\n", dbUpdateStats.ExistingCardsUpdated())
}
