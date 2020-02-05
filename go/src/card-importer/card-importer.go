package main

//import "database/sql"
import _ "github.com/go-sql-driver/mysql"
import "log"
import "mtgcards"
//import "carddb"

func main() {
    allPrices, err := mtgcards.DownloadAllPrices(true)
    if err != nil {
        log.Fatal(err)
    }

    log.Printf("Num prices: %d\n", len(allPrices))

    count := 0
    for uuid, prices := range allPrices {
        log.Printf("UUID: %s, Prices: %s", uuid, prices)
        count += 1
        if count > 5 {
            break
        }
    }

    /*
	allSets, err := mtgcards.DownloadAllPrintings(true)
	if err != nil {
		log.Fatal(err)
    }

	// Connect to the database
	db, err := sql.Open("mysql", "app_user:app_db_password@tcp(172.18.0.2)/mtg_cards")
	defer db.Close()
	if err != nil {
		log.Fatal(err)
	}

	db.SetMaxIdleConns(1000)

	stats, err := carddb.ImportSetsToDB(db, allSets)
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
	log.Printf("Total existing cards: %d\n", stats.TotalExistingCards())
	log.Printf("Existing cards skipped due to hash match: %d\n",
        stats.ExistingCardsSkipped())
	log.Printf("Existing cards updated due to hash mismatch: %d\n",
        stats.ExistingCardsUpdated())
    log.Printf("Total tokens in update: %d\n", stats.TotalTokens())
    log.Printf("Total new tokens: %d\n", stats.TotalNewTokens())
    log.Printf("Total new tokens in new sets: %d\n", stats.TotalNewTokensInNewSets())
    log.Printf("Total new tokens in existing sets: %d\n",
        stats.TotalNewTokensInExistingSets())
    log.Printf("Total existing tokens: %d\n", stats.TotalExistingTokens())
    log.Printf("Existing tokens skipped due to hash match: %d\n",
        stats.ExistingTokensSkipped())
    log.Printf("Existing tokens updated due to hash mismatch: %d\n",
        stats.ExistingTokensUpdated())
    */
}
