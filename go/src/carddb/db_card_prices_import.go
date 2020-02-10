package carddb

import influx "github.com/influxdata/influxdb1-client/v2"
import "fmt"
import "mtgcards"
import "sort"
import "time"

type PricesImportStats struct {
    CardRecordsAdded int
    MTGOPriceRecordsAdded int
    MTGOFoilPriceRecordsAdded int
    PaperPriceRecordsAdded int
    PaperFoilPriceRecordsAdded int
}

const POINTS_PER_WRITE = 1000

func ImportPricesToDb(
        influxClient influx.Client,
        lastImportTime time.Time,
        prices map[string]mtgcards.MTGCardPrices) (PricesImportStats, error) {
    importStats := PricesImportStats{}

    // Create the points to be pushed to the db
    bpConfig := influx.BatchPointsConfig{Database: "mtg_cards"}
    bp, err := influx.NewBatchPoints(bpConfig)
    if err != nil {
        return importStats, err
    }

    // Keep some stats
    cardRecordsAdded := 0
    mtgoPriceRecordsAdded := 0
    mtgoFoilPriceRecordsAdded := 0
    paperPriceRecordsAdded := 0
    paperFoilPriceRecordsAdded := 0

    totalRecords := len(prices)
    currentRecord := 0
    totalPointsImported := 0

    for card, cardPrices := range prices {
        fmt.Printf("Processing price record %d of %d\r", currentRecord, totalRecords)

        // Batch up records to write to influx so that we don't overflow a single
        // request, but also don't take too long
        if len(bp.Points()) > POINTS_PER_WRITE {
            err = influxClient.Write(bp)
            if err != nil {
                return importStats, err
            }

            totalPointsImported += len(bp.Points())

            bp, err = influx.NewBatchPoints(bpConfig)
            if err != nil {
                return importStats, err
            }
        }

        // MTGO
        newPriceRecords, err := maybeAddPoints(
            bp,
            card,
            cardPrices.MTGO,
            "mtgo",
            lastImportTime)
        if err != nil {
            return importStats, err
        }
        mtgoPriceRecordsAdded += newPriceRecords

        // MTGO Foil
        newPriceRecords, err = maybeAddPoints(
            bp,
            card,
            cardPrices.MTGOFoil,
            "mtgo_foil",
            lastImportTime)
        if err != nil {
            return importStats, err
        }
        mtgoFoilPriceRecordsAdded += newPriceRecords

        // Paper
        newPriceRecords, err = maybeAddPoints(
            bp,
            card,
            cardPrices.Paper,
            "paper",
            lastImportTime)
        if err != nil {
            return importStats, err
        }
        paperPriceRecordsAdded += newPriceRecords

        // Paper Foil
        newPriceRecords, err = maybeAddPoints(
            bp,
            card,
            cardPrices.PaperFoil,
            "paper_foil",
            lastImportTime)
        if err != nil {
            return importStats, err
        }
        paperFoilPriceRecordsAdded += newPriceRecords

        if mtgoPriceRecordsAdded > 0 ||
                mtgoFoilPriceRecordsAdded > 0 ||
                paperPriceRecordsAdded > 0 ||
                paperFoilPriceRecordsAdded > 0 {
            cardRecordsAdded += 1
        }

        currentRecord += 1
    }

    // Flush any remaining points to the db
    if len(bp.Points()) > 0 {
        err = influxClient.Write(bp)
        if err != nil {
            return importStats, err
        }

        totalPointsImported += len(bp.Points())
    }

    fmt.Printf("\n")
    fmt.Printf("Total points imported: %d\n", totalPointsImported)

    importStats.CardRecordsAdded = cardRecordsAdded
    importStats.MTGOPriceRecordsAdded = mtgoPriceRecordsAdded
    importStats.MTGOFoilPriceRecordsAdded = mtgoFoilPriceRecordsAdded
    importStats.PaperPriceRecordsAdded = paperPriceRecordsAdded
    importStats.PaperFoilPriceRecordsAdded = paperFoilPriceRecordsAdded

    return importStats, nil
}

func maybeAddPoints(
        bp influx.BatchPoints,
        card string,
        priceRecords mtgcards.MTGCardPriceRecords,
        measurementName string,
        lastImportDate time.Time) (int, error) {
    // First, sort the price records
    sort.Sort(priceRecords)

    pointsAdded := 0
    for _, priceRecord := range priceRecords {
        // Only import points that are after the last time we
        // imported data
        if priceRecord.Date.Before(lastImportDate) {
            continue
        }

        point, err := influx.NewPoint(
            measurementName,
            map[string]string{"card": card},
            map[string]interface{}{"price": priceRecord.Price},
            priceRecord.Date)
        if err != nil {
            return 0, err
        }
        bp.AddPoint(point)
        pointsAdded += 1
    }

    return pointsAdded, nil
}
