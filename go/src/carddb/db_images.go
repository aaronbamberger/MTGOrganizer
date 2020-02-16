package carddb

import "database/sql"
import "fmt"
import "io"
import "log"
import "net/http"
import "os"
import "time"

const (
    ScryfallImageURLTemplate = "https://api.scryfall.com/cards/%s?format=image&version=png"
    LocalDownloadLocation = "/var/card-importer/card-images/"
)

func UpdateCardImages(cardDB *sql.DB) (ImageUpdateStats, error) {
    updateStats := ImageUpdateStats{}

    cardInfoQuery, err := cardDB.Prepare(`SELECT
        uuid,
        scryfall_id
        FROM all_cards
        WHERE image_cached_locally = False`)
    if err != nil {
        return updateStats, err
    }
    defer cardInfoQuery.Close()

    tokenInfoQuery, err := cardDB.Prepare(`SELECT
        uuid,
        scryfall_id
        FROM all_tokens
        WHERE image_cached_locally = False`)
    if err != nil {
        return updateStats, err
    }
    defer tokenInfoQuery.Close()

    cardUpdateQuery, err := cardDB.Prepare(`UPDATE all_cards
        SET image_cached_locally = True
        WHERE uuid = ?`)
    if err != nil {
        return updateStats, err
    }
    defer cardUpdateQuery.Close()

    tokenUpdateQuery, err := cardDB.Prepare(`UPDATE all_tokens
        SET image_cached_locally = True
        WHERE uuid = ?`)
    if err != nil {
        return updateStats, err
    }
    defer tokenUpdateQuery.Close()

    // Check to see if there are any cards missing images
    // If there are, try and download them
    cardsMissingImages, err := checkForNewImages(cardInfoQuery)
    if err != nil {
        return updateStats, err
    }
    if len(cardsMissingImages) > 0 {
        updateStats.AddToCardsNeedingImages(len(cardsMissingImages))
        downloadNewImages(cardsMissingImages, cardUpdateQuery, &updateStats)
    }

    // Check to see if there are any tokens missing images
    // If there are, try and download them
    tokensMissingImages, err := checkForNewImages(tokenInfoQuery)
    if err != nil {
        return updateStats, err
    }
    if len(tokensMissingImages) > 0 {
        updateStats.AddToTokensNeedingImages(len(tokensMissingImages))
        downloadNewImages(tokensMissingImages, tokenUpdateQuery, &updateStats)
    }

    return updateStats, nil
}

func checkForNewImages(infoQuery *sql.Stmt) (map[string]string, error) {
    // Check and see whether there are any images we don't have yet
    missingImages := make(map[string]string)
    queryResult, err := infoQuery.Query()
    if err != nil {
        return nil, err
    }
    defer queryResult.Close()

    for queryResult.Next() {
        var uuid string
        var scryfallId string
        err = queryResult.Scan(&uuid, &scryfallId)
        if err != nil {
            return nil, err
        }
        missingImages[uuid] = scryfallId
    }
    if err = queryResult.Err(); err != nil {
        return nil, err
    }

    return missingImages, nil
}

func downloadNewImages(missingImages map[string]string,
        updateQuery *sql.Stmt,
        updateStats *ImageUpdateStats) {

    totalRecords := len(missingImages)
    currentRecord := 0
    // Try and download all of the missing images we can get from scryfall
    for uuid, scryfallId := range missingImages {
        fmt.Printf("Processing missing image record %d of %d\r", currentRecord, totalRecords)
        currentRecord += 1

        startTime := time.Now()
        scryfallImageUrl := fmt.Sprintf(ScryfallImageURLTemplate, scryfallId)
        resp, err := http.Get(scryfallImageUrl)
        if err != nil {
            log.Print(err)
            updateStats.AddToImagesFailedToDownload(1)
            continue
        }

        if resp.StatusCode != http.StatusOK {
            log.Print(resp.Status)
            resp.Body.Close()
            updateStats.AddToImagesFailedToDownload(1)
            continue
        }

        localFileName := fmt.Sprintf("%s%s.png", LocalDownloadLocation, uuid)
        localFile, err := os.Create(localFileName)
        if err != nil {
            log.Print(err)
            resp.Body.Close()
            updateStats.AddToImagesFailedToDownload(1)
            continue
        }

        // Copy the file from the http response to the local file
        _, err = io.Copy(localFile, resp.Body)
        if err != nil {
            log.Print(err)
            resp.Body.Close()
            localFile.Close()
            updateStats.AddToImagesFailedToDownload(1)
            continue
        }

        resp.Body.Close()
        localFile.Close()

        _, err = updateQuery.Exec(uuid)
        if err != nil {
            log.Print(err)
            continue
            updateStats.AddToImagesFailedToDownload(1)
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

        updateStats.AddToImagesDownloaded(1)
    }
}
