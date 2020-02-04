package mtgcards

import "bytes"
import "encoding/json"
import "log"
import "time"

type MTGCardPrices struct {
    MTGO map[string]float32 `json:"mtgo"`
    MTGOFoil map[string]float32 `json:"mtgoFoil"`
    Paper map[string]float32 `json:"paper"`
    PaperFoil map[string]float32 `json:"paperFoil"`
}

type MTGCardPriceRecord struct {
    Date time.Time
    Price float64
}

type MTGCardPricesDummy struct {
    MTGO map[string]float32 `json:"mtgo"`
    MTGOFoil map[string]float32 `json:"mtgoFoil"`
    Paper map[string]float32 `json:"paper"`
    PaperFoil map[string]float32 `json:"paperFoil"`
}

func (prices *MTGCardPrices) UnmarshalJSON(data []byte) error {
    // First, trim off the leading and trailing object delimiters
    data = data[1:len(data)-1]
    var pricesDummy MTGCardPricesDummy
    keyAndValue := bytes.SplitN(data, []byte(":"), 2)
    if bytes.Contains(keyAndValue[0], []byte("prices")) {
        err := json.Unmarshal(keyAndValue[1], &pricesDummy)
        if err != nil {
            return err
        }
        prices.MTGO = pricesDummy.MTGO
        prices.MTGOFoil = pricesDummy.MTGOFoil
        prices.Paper = pricesDummy.Paper
        prices.PaperFoil = pricesDummy.PaperFoil
    } else {
        log.Print("Unexpected key and value %s", keyAndValue)
    }

    return nil
}
