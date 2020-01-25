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
	log.Printf("Total sets in update: %d\n", dbUpdateStats.TotalSetsInUpdate)
	log.Printf("Total existing sets in the database: %d\n", dbUpdateStats.TotalExistingSetsInDb)
	log.Printf("Total new sets: %d\n", dbUpdateStats.TotalNewSets)
	log.Printf("Existing sets skipped due to hash match: %d\n",
		dbUpdateStats.ExistingSetsSkippedDueToHashMatch)
	log.Printf("Existing sets updated due to hash mismatch: %d\n",
		dbUpdateStats.ExistingSetsUpdatedDueToHashMismatch)
	log.Printf("Total cards in update: %d\n", dbUpdateStats.TotalCardsInUpdate)
	log.Printf("Total new cards: %d\n", dbUpdateStats.TotalNewCards)
	log.Printf("Total new cards in new sets: %d\n", dbUpdateStats.TotalNewCardsInNewSets)
	log.Printf("Total new cards in existing sets: %d\n", dbUpdateStats.TotalNewCardsInExistingSets)
	log.Printf("Total new atomic cards: %d\n", dbUpdateStats.TotalNewAtomicCards)
	log.Printf("Existing cards skipped due to hash match: %d\n",
		dbUpdateStats.ExistingCardsSkippedDueToHashMatch)
	log.Printf("Existing cards updated due to hash mismatch: %d\n",
		dbUpdateStats.ExistingCardsUpdatedDueToHashMismatch)
}
