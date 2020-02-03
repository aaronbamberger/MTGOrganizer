package mtgcards

import "time"

type MTGJSONVersion struct {
	BuildDate time.Time `json:"date"`
	PricesDate time.Time `json:"pricesDate"`
	Version string `json:"version"`
}
