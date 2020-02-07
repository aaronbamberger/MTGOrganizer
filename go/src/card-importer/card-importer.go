package main

import "database/sql"
import _ "github.com/go-sql-driver/mysql"
import influx "github.com/influxdata/influxdb1-client/v2"
import "log"
import "mtgcards"
import "carddb"

func main() {
	// Connect to the mariadb database
	db, err := sql.Open("mysql", "app_user:app_db_password@tcp(172.18.0.2)/mtg_cards?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	db.SetMaxIdleConns(1000)

    // Connect to the influxdb database
    influxClientConfig := influx.HTTPConfig{
        Addr: "http://172.18.0.4:8086",
        Username: "app_user",
        Password: "app_db_password"}

    influxClient, err := influx.NewHTTPClient(influxClientConfig)
    if err != nil {
        log.Fatal(err)
    }
    defer influxClient.Close()

    // First, get the latest version data to see if we even need to bother downloading
    // any additional data
    onlineVersion, err := mtgcards.DownloadVersion()
    if err != nil {
        log.Fatal(err)
    }
    dbVersion, err := carddb.GetDbLastUpdate(db)
    if err != nil {
        log.Fatal(err)
    }

    log.Printf("Online version:\n%s", onlineVersion)
    log.Printf("DB version:\n%s", dbVersion)

    // Update cards if necessary
    if onlineVersion.BuildDate.After(dbVersion.BuildDate) {
        log.Printf("Cards are out of date, updating...\n")

        log.Printf("Downloading new cards...\n")
	    allSets, err := mtgcards.DownloadAllPrintings(false)
        if err != nil {
            log.Fatal(err)
        }

        log.Printf("Updating cards in db...\n")
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
    } else {
        log.Printf("Already have latest version of cards, skipping update...\n")
    }

    // Update prices if necessary
    if onlineVersion.PricesDate.After(dbVersion.PricesDate) {
        log.Printf("Prices are out of date, updating...\n")

        log.Printf("Downloading new prices...\n")
        allPrices, err := mtgcards.DownloadAllPrices(false)
        if err != nil {
            log.Fatal(err)
        }

        log.Printf("Updating prices in db...\n")
        stats, err := carddb.ImportPricesToDb(influxClient, dbVersion.PricesDate, allPrices)
        if err != nil {
            log.Fatal(err)
        }
        log.Printf("New price records added for %d cards\n", stats.CardRecordsAdded)
        log.Printf("New MTGO price records added: %d\n", stats.MTGOPriceRecordsAdded)
        log.Printf("New MTGO foil price records added: %d\n", stats.MTGOFoilPriceRecordsAdded)
        log.Printf("New paper price records added: %d\n", stats.PaperPriceRecordsAdded)
        log.Printf("New paper foil price records added: %d\n", stats.PaperFoilPriceRecordsAdded)
    } else {
        log.Printf("Already have latest version of prices, skipping update...\n")
    }

    // Update the update times in the db
    log.Printf("Updating last update time in DB\n")
    err = carddb.UpdateDbLastUpdate(db, onlineVersion)
    if err != nil {
        log.Fatal(err)
    }
}
