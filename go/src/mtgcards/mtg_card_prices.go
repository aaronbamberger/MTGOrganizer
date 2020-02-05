package mtgcards

import "bytes"
import "fmt"
import "strconv"
import "time"

type MTGCardPricesTopLevelDummy struct {
    Prices MTGCardPrices `json:"prices"`
}

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

func (prices MTGCardPrices) String() string {
    return fmt.Sprintf("{MTGO: %v, MTGOFoil: %v, Paper: %v, PaperFoil: %v}",
        prices.MTGO, prices.MTGOFoil, prices.Paper, prices.PaperFoil)
}

type MTGCardPriceRecord struct {
    Date time.Time
    Price float64
}

func (priceRecords *MTGCardPriceRecords) UnmarshalJSON(data []byte) error {
    // First, truncate the starting and ending object delimiters
    data = data[1:len(data)-1]
    if len(data) > 0 {
        records := bytes.Split(data, []byte(","))
        for _, record := range records {
            dateAndPrice := bytes.Split(record, []byte(":"))

            // For some reason, the input data set can sometimes have "null" for a price
            // skip these entries
            if bytes.Contains(dateAndPrice[1], []byte("null")) {
                continue
            }

            // Trim whitespace and quotes from the date and price strings
            dateString := string(bytes.Trim(dateAndPrice[0], "\" "))
            priceString := string(bytes.Trim(dateAndPrice[1], "\" "))

            date, err := time.Parse("2006-01-02", dateString)
            if err != nil {
                return err
            }
            price, err := strconv.ParseFloat(priceString, 64)
            if err != nil {
                return err
            }
            *priceRecords = append(*priceRecords, MTGCardPriceRecord{Date: date, Price: price})
        }
    }

    return nil
}
