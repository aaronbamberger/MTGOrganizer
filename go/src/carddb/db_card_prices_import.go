package carddb

import influx "github.com/influxdata/influxdb1-client/v2"
import "fmt"
import "mtgcards"
import "sort"
import "time"

const POINTS_PER_WRITE = 1000

func ImportPricesToDb(
        influxClient influx.Client,
        lastImportTime time.Time,
        prices map[string]mtgcards.MTGCardPrices) (PricesUpdateStats, error) {
    importStats := PricesUpdateStats{}

    // Create the points to be pushed to the db
    bpConfig := influx.BatchPointsConfig{Database: "mtg_cards"}
    bp, err := influx.NewBatchPoints(bpConfig)
    if err != nil {
        return importStats, err
    }

    totalRecords := len(prices)
    currentRecord := 0

    for card, cardPrices := range prices {
        fmt.Printf("Processing price record %d of %d\r", currentRecord, totalRecords)

        // Batch up records to write to influx so that we don't overflow a single
        // request, but also don't take too long
        if len(bp.Points()) > POINTS_PER_WRITE {
            err = influxClient.Write(bp)
            if err != nil {
                return importStats, err
            }

            importStats.AddToTotalPriceRecords(len(bp.Points()))

            bp, err = influx.NewBatchPoints(bpConfig)
            if err != nil {
                return importStats, err
            }
        }

        newPriceRecordsAdded := false

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
        if newPriceRecords > 0 {
            newPriceRecordsAdded = true
        }
        importStats.AddToMTGOPriceRecords(newPriceRecords)

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
        if newPriceRecords > 0 {
            newPriceRecordsAdded = true
        }
        importStats.AddToMTGOFoilPriceRecords(newPriceRecords)

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
        if newPriceRecords > 0 {
            newPriceRecordsAdded = true
        }
        importStats.AddToPaperPriceRecords(newPriceRecords)

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
        if newPriceRecords > 0 {
            newPriceRecordsAdded = true
        }
        importStats.AddToPaperFoilPriceRecords(newPriceRecords)

        if newPriceRecordsAdded {
            importStats.AddToTotalCardRecords(1)
        }

        currentRecord += 1
    }

    // Flush any remaining points to the db
    if len(bp.Points()) > 0 {
        err = influxClient.Write(bp)
        if err != nil {
            return importStats, err
        }

        importStats.AddToTotalPriceRecords(len(bp.Points()))
    }

    fmt.Printf("\n")

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
