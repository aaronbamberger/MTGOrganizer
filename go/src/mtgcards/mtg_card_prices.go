package mtgcards

import "fmt"
import "encoding/json"
import "time"

type MTGCardPriceRecords []MTGCardPriceRecord

func (priceRecords MTGCardPriceRecords) Len() int {
	return len(priceRecords)
}

func (priceRecords MTGCardPriceRecords) Less(i, j int) bool {
	return priceRecords[i].Date.Before(priceRecords[j].Date)
}

func (priceRecords MTGCardPriceRecords) Swap(i, j int) {
	priceRecords[i], priceRecords[j] = priceRecords[j], priceRecords[i]
}

type MTGCardPrices struct {
    MTGO MTGCardPriceRecords `json:"mtgo"`
    MTGOFoil MTGCardPriceRecords `json:"mtgoFoil"`
    Paper MTGCardPriceRecords `json:"paper"`
    PaperFoil MTGCardPriceRecords `json:"paperFoil"`
}

type MTGCardPricesDummy MTGCardPrices

func (prices MTGCardPrices) String() string {
    return fmt.Sprintf("{MTGO: %v, MTGOFoil: %v, Paper: %v, PaperFoil: %v}",
        prices.MTGO, prices.MTGOFoil, prices.Paper, prices.PaperFoil)
}

type MTGCardPriceRecord struct {
    Date time.Time
    Price float64
}

type mtgCardPricesTopLevelDummy struct {
    Prices MTGCardPricesDummy `json:"prices"`
}

func (cardPrices *MTGCardPrices) UnmarshalJSON(data []byte) error {
    // First, unmarshal into a dummy object to deal with the top-level
    // "prices" key
    var pricesDummy mtgCardPricesTopLevelDummy
    err := json.Unmarshal(data, &pricesDummy)
    if err != nil {
        return err
    }

    // The top-level dummy object contains another dummy object for the actual
    // prices which is just a type alias to the actual prices object.  This is just
    // so that it's a different type, because else when we try to unmarshal into
    // the top-level object, we'll end up recursively calling this custom
    // unmarshal function, when what we actually want is the standard
    // unmarshalling.  Since it's a different type, we need to cast here,
    // but we know it's compatible since it's just a type alias
    *cardPrices = MTGCardPrices(pricesDummy.Prices)

    return nil
}

func (priceRecords *MTGCardPriceRecords) UnmarshalJSON(data []byte) error {
    // First, unpack the data into a map of date to price, since that's its
    // native format
    priceMap := make(map[string]float64)
    err := json.Unmarshal(data, &priceMap)
    if err != nil {
        return err
    }

    // Now, convert the map entries into price record objects
    for dateString, price := range priceMap {
        date, err := time.Parse("2006-01-02", dateString)
        if err != nil {
            return err
        }
        newRecord := MTGCardPriceRecord{Date: date, Price: price}
        *priceRecords = append(*priceRecords, newRecord)
    }

    return nil
}
