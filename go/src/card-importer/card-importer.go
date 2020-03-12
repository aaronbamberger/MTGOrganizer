package main

import "database/sql"
import "fmt"
import _ "github.com/go-sql-driver/mysql"
import influx "github.com/influxdata/influxdb1-client/v2"
import "log"
import "mtgcards"
import "carddb"
import "os"
import "os/signal"
import "syscall"
import "time"

const (
    cardDBUrl = "card_db"
    cardPricesDBUrl = "card_prices_db"
    appUsername = "app_user"
    appPassword = "app_db_password"
)

func main() {
    // Create a ticker that will signal the program to check for updates
    // every 12 hours
    tickDuration, err := time.ParseDuration("12h")
    if err != nil {
        log.Fatal(err)
    }
    updateTicker := time.NewTicker(tickDuration)

    // Register a signal handler for SIGUSR1, so that updates can be requested
    // outside of the normal update period
    updateRequest := make(chan os.Signal)
    signal.Notify(updateRequest, syscall.SIGUSR1)

    // Register a signal handler for SIGINT, so that we can exit cleanly
    // in case an update is currently in progress (will wait for the update
    // to finish before we exit)
    closeRequest := make(chan os.Signal)
    signal.Notify(closeRequest, syscall.SIGINT)

    log.Printf("Card importer started with PID %d...\n", os.Getpid())
    for {
        select {
        // Handle this one first, if it's available, so that we exit
        // as soon as possible if there's a request
        case <-closeRequest:
            log.Printf("Close request received, exiting...\n")
            os.Exit(0)

        // Run the update checker either when the ticker goes off,
        // or on-demand, in response to SIGUSR1
        case <-updateTicker.C:
            log.Printf("Timer tick, checking for updates...\n")
            CheckForAndMaybeRunUpdate()
        case <-updateRequest:
            log.Printf("Manual update request, checking for updates...\n")
            CheckForAndMaybeRunUpdate()
        }
    }
}

func CheckForAndMaybeRunUpdate() {
	// Connect to the mariadb database
    cardDBConnStr := fmt.Sprintf("%s:%s@tcp(%s)/mtg_cards?parseTime=true",
        appUsername, appPassword, cardDBUrl)
	cardDB, err := sql.Open("mysql", cardDBConnStr)
	if err != nil {
		log.Print(err)
        return
	}
	defer cardDB.Close()
	cardDB.SetMaxIdleConns(1000)

    // Connect to the influxdb database
    cardPricesDBAddr := fmt.Sprintf("http://%s:8086", cardPricesDBUrl)
    influxClientConfig := influx.HTTPConfig{
        Addr: cardPricesDBAddr,
        Username: appUsername,
        Password: appPassword}

    pricesAndStatsDB, err := influx.NewHTTPClient(influxClientConfig)
    if err != nil {
        log.Print(err)
        return
    }
    defer pricesAndStatsDB.Close()

    // First, get the latest version data to see if we even need to bother downloading
    // any additional data
    onlineVersion, err := mtgcards.DownloadVersion()
    if err != nil {
        log.Print(err)
        return
    }
    dbVersion, err := carddb.GetDbLastUpdate(cardDB)
    if err != nil {
        log.Print(err)
        return
    }

    log.Printf("Online version:\n%s", onlineVersion)
    log.Printf("DB version:\n%s", dbVersion)

    cardsUpdated := false
    pricesUpdated := false
    imagesUpdated := false
    cardsUpdateDuration := time.Duration(0)
    pricesUpdateDuration := time.Duration(0)
    imagesUpdateDuration := time.Duration(0)

    // Update cards if necessary
    if onlineVersion.BuildDate.After(dbVersion.BuildDate) {
        cardsUpdateDuration, err = UpdateCards(cardDB, pricesAndStatsDB)
        if err != nil {
            log.Print(err)
        } else {
            cardsUpdated = true
        }
    } else {
        log.Printf("Already have latest version of cards, skipping update...\n")
    }

    // Update prices if necessary
    if onlineVersion.PricesDate.After(dbVersion.PricesDate) {
        pricesUpdateDuration, err = UpdatePrices(pricesAndStatsDB, onlineVersion.PricesDate)
        if err != nil {
            log.Print(err)
        } else {
            pricesUpdated = true
        }
    } else {
        log.Printf("Already have latest version of prices, skipping update...\n")
    }

    // Update the update times in the db, but only if something actually got updated
    if cardsUpdated || pricesUpdated {
        log.Printf("Updating last update time in DB\n")
        err = carddb.UpdateDbLastUpdate(cardDB, onlineVersion)
        if err != nil {
            log.Print(err)
        }
    }

    // Check to see if we're missing any card images, and if so, try and get them
    imagesUpdated, imagesUpdateDuration, err = UpdateImages(cardDB, pricesAndStatsDB)
    if err != nil {
        log.Print(err)
    }

    // Add the update run stats to the db
    log.Printf("Adding update run stats to db...\n")
    err = carddb.AddSingleUpdateStatsToDb(pricesAndStatsDB,
        cardsUpdated,
        pricesUpdated,
        imagesUpdated,
        cardsUpdateDuration,
        pricesUpdateDuration,
        imagesUpdateDuration)
    if err != nil {
        log.Print(err)
    }
}

func UpdateCards(cardDB *sql.DB, priceAndStatsDB influx.Client) (time.Duration, error) {
    updateDuration := time.Duration(0)
    updateStart := time.Now()

    log.Printf("Cards are out of date, updating...\n")

    log.Printf("Downloading new cards...\n")
    allSets, err := mtgcards.DownloadAllPrintings(false)
    if err != nil {
        return updateDuration, err
    }

    log.Printf("Updating cards in db...\n")
    stats, err := carddb.ImportSetsToDB(cardDB, allSets)
    if err != nil {
        return updateDuration, err
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

    log.Printf("Adding card update stats to the db...\n")
    err = stats.AddToDb(priceAndStatsDB)
    if err != nil {
        return updateDuration, err
    }
    updateDuration = time.Since(updateStart)

    return updateDuration, nil
}

func UpdatePrices(priceAndStatsDB influx.Client, pricesDate time.Time) (time.Duration, error) {
    updateDuration := time.Duration(0)
    updateStart := time.Now()
    log.Printf("Prices are out of date, updating...\n")

    log.Printf("Downloading new prices...\n")
    allPrices, err := mtgcards.DownloadAllPrices(false)
    if err != nil {
        return updateDuration, err
    }

    log.Printf("Updating prices in db...\n")
    stats, err := carddb.ImportPricesToDb(priceAndStatsDB, pricesDate, allPrices)
    if err != nil {
        return updateDuration, err
    }
    log.Printf("New price records added for %d cards\n", stats.TotalCardRecords())
    log.Printf("Total new price records added: %d\n", stats.TotalPriceRecords())
    log.Printf("New MTGO price records added: %d\n", stats.MTGOPriceRecords())
    log.Printf("New MTGO foil price records added: %d\n", stats.MTGOFoilPriceRecords())
    log.Printf("New paper price records added: %d\n", stats.PaperPriceRecords())
    log.Printf("New paper foil price records added: %d\n", stats.PaperFoilPriceRecords())

    log.Printf("Adding price update stats to the db...\n")
    err = stats.AddToDb(priceAndStatsDB)
    if err != nil {
        return updateDuration, err
    }

    updateDuration = time.Since(updateStart)

    return updateDuration, nil
}

func UpdateImages(cardDB *sql.DB, pricesAndStatsDB influx.Client) (bool, time.Duration, error) {
    log.Printf("Checking for and downloading any missing card images...\n")
    updateStartTime := time.Now()
    updateStats, err := carddb.UpdateCardImages(cardDB)
    updateDuration := time.Since(updateStartTime)
    log.Printf("Cards needing images: %d\n", updateStats.CardsNeedingImages())
    log.Printf("Tokens needing images: %d\n", updateStats.TokensNeedingImages())
    log.Printf("Images downloaded: %d\n", updateStats.ImagesDownloaded())
    log.Printf("Images that failed to download: %d\n", updateStats.ImagesFailedToDownload())
    if err != nil {
        return false, time.Duration(0), err
    } else {
        imagesUpdated := (updateStats.CardsNeedingImages() > 0) ||
            (updateStats.TokensNeedingImages() > 0)

        if imagesUpdated {
            err = updateStats.AddToDb(pricesAndStatsDB)
            if err != nil {
                return false, time.Duration(0), err
            }
            return true, updateDuration, nil
        } else {
            return false, time.Duration(0), nil
        }
    }
}
