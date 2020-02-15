package main

import "database/sql"
import _ "github.com/go-sql-driver/mysql"
import "fmt"
import "io"
import "log"
import "net/http"
import "os"
import "time"

type CardImageRecord struct {
    ScryfallId string
    ImageCachedLocally bool
    IsToken bool
}

func main() {
	// Connect to the mariadb database
	cardDB, err := sql.Open("mysql",
        "app_user:app_db_password@tcp(172.18.0.3)/mtg_cards?parseTime=true")
	if err != nil {
		log.Fatal(err)
        return
	}
	defer cardDB.Close()
	cardDB.SetMaxIdleConns(1000)

    scryfallImageRecords := make(map[string]CardImageRecord)

    cardInfoQuery, err := cardDB.Prepare(`SELECT
        uuid,
        scryfall_id,
        image_cached_locally
        FROM all_cards`)
    if err != nil {
        log.Fatal(err)
    }
    defer cardInfoQuery.Close()

    tokenInfoQuery, err := cardDB.Prepare(`SELECT
        uuid,
        scryfall_id,
        image_cached_locally
        FROM all_tokens`)
    if err != nil {
        log.Fatal(err)
    }
    defer tokenInfoQuery.Close()

    cardUpdateQuery, err := cardDB.Prepare(`UPDATE all_cards
        SET image_cached_locally = 1
        WHERE uuid = ?`)
    if err != nil {
        log.Fatal(err)
    }
    defer cardUpdateQuery.Close()

    tokenUpdateQuery, err := cardDB.Prepare(`UPDATE all_tokens
        SET image_cached_locally = 1
        WHERE uuid = ?`)
    if err != nil {
        log.Fatal(err)
    }
    defer tokenUpdateQuery.Close()

    // First, get all of the cards and tokens we're downloading images for
    allCards, err := cardInfoQuery.Query()
    if err != nil {
        log.Fatal(err)
    }

    for allCards.Next() {
        var uuid string
        var scryfallId string
        var imageCachedLocally bool
        err = allCards.Scan(&uuid, &scryfallId, &imageCachedLocally)
        if err != nil {
            log.Print(err)
            continue
        }
        scryfallImageRecords[uuid] = CardImageRecord{
            ScryfallId: scryfallId,
            ImageCachedLocally: imageCachedLocally,
            IsToken: false}
    }
    if err = allCards.Err(); err != nil {
        log.Print(err)
    }
    allCards.Close()

    allTokens, err := tokenInfoQuery.Query()
    if err != nil {
        log.Fatal(err)
    }

    for allTokens.Next() {
        var uuid string
        var scryfallId string
        var imageCachedLocally bool
        err = allTokens.Scan(&uuid, &scryfallId, &imageCachedLocally)
        if err != nil {
            log.Print(err)
            continue
        }
        scryfallImageRecords[uuid] = CardImageRecord{
            ScryfallId: scryfallId,
            ImageCachedLocally: imageCachedLocally,
            IsToken: true}
    }
    if err = allTokens.Err(); err != nil {
        log.Print(err)
    }
    allTokens.Close()

    totalRecords := len(scryfallImageRecords)
    log.Printf("Card/Token records received: %d\n", totalRecords)

    totalStart := time.Now()

    currentRecord := 0
    // Try and download all of the card images we can get from scryfall
    for uuid, cardImageRecord := range scryfallImageRecords {
        fmt.Printf("Processing record %d of %d\r", currentRecord, totalRecords)
        currentRecord += 1
        // Skip cards whose images we've already downloaded
        if cardImageRecord.ImageCachedLocally {
            continue
        }

        startTime := time.Now()
        scryfallImageUrl := fmt.Sprintf("https://api.scryfall.com/cards/%s?format=image&version=png",
            cardImageRecord.ScryfallId)
        resp, err := http.Get(scryfallImageUrl)
        if err != nil {
            log.Print(err)
            continue
        }

        if resp.StatusCode != http.StatusOK {
            log.Print(resp.Status)
            resp.Body.Close()
            continue
        }

        localFileName := fmt.Sprintf("../../../web_content/card_face_images/%s.png", uuid)
        localFile, err := os.Create(localFileName)
        if err != nil {
            log.Print(err)
            resp.Body.Close()
            continue
        }

        // Copy the file from the http response to the local file
        _, err = io.Copy(localFile, resp.Body)
        if err != nil {
            log.Print(err)
            resp.Body.Close()
            localFile.Close()
            continue
        }

        resp.Body.Close()
        localFile.Close()

        // Update the database
        var updateQuery *sql.Stmt
        if cardImageRecord.IsToken {
            updateQuery = tokenUpdateQuery
        } else {
            updateQuery = cardUpdateQuery
        }

        _, err = updateQuery.Exec(uuid)
        if err != nil {
            log.Print(err)
            continue
        }

        duration := time.Since(startTime)

        // We want to rate limit requests to the server to be ~10/sec, so if the
        // download operation took less than 100ms, sleep the remainder
        if duration.Milliseconds() < 100 {
            waitDuration, err := time.ParseDuration(fmt.Sprintf("%dms", 100 - duration.Milliseconds()))
            if err != nil {
                log.Print(err)
                continue
            }
            time.Sleep(waitDuration)
        }
    }

    totalDuration := time.Since(totalStart)
    fmt.Printf("It took %v to download %d images\n", totalDuration, currentRecord)
}
