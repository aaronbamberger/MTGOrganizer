package carddb

import influx "github.com/influxdata/influxdb1-client/v2"
import "mtgcards"

func ImportPricesToDb(prices map[string]mtgcards.MTGCardPrices) error {
    clientConfig := influx.HTTPConfig{
        Addr: "http://172.18.0.3:8086",
        Username: "app_user",
        Password: "app_db_password"}

    client, err := influx.NewHTTPClient(clientConfig)
    if err != nil {
        return err
    }
    defer client.Close()

    // Create the points to be pushed to the db
    bpConfig := influx.BatchPointsConfig{Database: "mtg_cards"}
    bp, err := influx.NewBatchPoints(bpConfig)
    if err != nil {
        return err
    }

    for card, cardPrices := range prices {
        for _, priceRecord := range cardPrices.MTGO {
            point, err := influx.NewPoint("mtgo",
                map[string]string{"card": card}, 
                map[string]interface{}{"price": priceRecord.Price},
                priceRecord.Date)
            if err != nil {
                return err
            }
            bp.AddPoint(point)
        }

        for _, priceRecord := range cardPrices.MTGOFoil {
            point, err := influx.NewPoint(
                "mtgo_foil",
                map[string]string{"card": card},
                map[string]interface{}{"price": priceRecord.Price},
                priceRecord.Date)
            if err != nil {
                return err
            }
            bp.AddPoint(point)
        }

        for _, priceRecord := range cardPrices.Paper {
            point, err := influx.NewPoint(
                "paper",
                map[string]string{"card": card},
                map[string]interface{}{"price": priceRecord.Price},
                priceRecord.Date)
            if err != nil {
                return err
            }
            bp.AddPoint(point)
        }

        for _, priceRecord := range cardPrices.PaperFoil {
            point, err := influx.NewPoint(
                "paper_foil",
                map[string]string{"card": card},
                map[string]interface{}{"price": priceRecord.Price},
                priceRecord.Date)
            if err != nil {
                return err
            }
            bp.AddPoint(point)
        }
    }

    err = client.Write(bp)
    if err != nil {
        return err
    }

    return nil
}
