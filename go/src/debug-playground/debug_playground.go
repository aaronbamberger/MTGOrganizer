package main

import "fmt"
import "log"
import "mtgcards"
import "strings"

const (
    NONE = 0x00
    MTGO = 0x01
    MTGO_FOIL = 0x02
    PAPER = 0x04
    PAPER_FOIL = 0x08
)

type PriceTypes uint

func (priceType PriceTypes) String() string {
    var b strings.Builder

    if priceType == NONE {
        return "None"
    }

    if priceType & MTGO == MTGO {
        b.WriteString("MTGO")
    }

    if priceType & MTGO_FOIL == MTGO_FOIL {
        if b.Len() > 0 {
            b.WriteString(",")
        }
        b.WriteString("MTGO_FOIL")
    }

    if priceType & PAPER == PAPER {
        if b.Len() > 0 {
            b.WriteString(",")
        }
        b.WriteString("PAPER")
    }

    if priceType & PAPER_FOIL == PAPER_FOIL {
        if b.Len() > 0 {
            b.WriteString(",")
        }
        b.WriteString("PAPER_FOIL")
    }

    return b.String()
}

func main() {
    allPrices, err := mtgcards.DownloadAllPrices(true)
    if err != nil {
        log.Fatal(err)
    }

    numCardsWithPrices := 0
    mtgoPrices := 0
    mtgoFoilPrices := 0
    paperPrices := 0
    paperFoilPrices := 0
    totalPrices := 0

    cardsWithSpecificPriceTypes := make(map[PriceTypes]int)
    cardsWithSpecificPriceTypes[NONE] = 0
    cardsWithSpecificPriceTypes[MTGO] = 0
    cardsWithSpecificPriceTypes[MTGO_FOIL] = 0
    cardsWithSpecificPriceTypes[PAPER] = 0
    cardsWithSpecificPriceTypes[PAPER_FOIL] = 0

    for _, prices := range allPrices {
        priceType := PriceTypes(NONE)
        newPriceRecords := false

        if len(prices.MTGO) > 0 {
            mtgoPrices += len(prices.MTGO)
            totalPrices += len(prices.MTGO)
            newPriceRecords = true
            priceType |= MTGO
        }

        if len(prices.MTGOFoil) > 0 {
            mtgoFoilPrices += len(prices.MTGOFoil)
            totalPrices += len(prices.MTGOFoil)
            newPriceRecords = true
            priceType |= MTGO_FOIL
        }

        if len(prices.Paper) > 0 {
            paperPrices += len(prices.Paper)
            totalPrices += len(prices.Paper)
            newPriceRecords = true
            priceType |= PAPER
        }

        if len(prices.PaperFoil) > 0 {
            paperFoilPrices += len(prices.PaperFoil)
            totalPrices += len(prices.PaperFoil)
            newPriceRecords = true
            priceType |= PAPER_FOIL
        }

        if newPriceRecords {
            numCardsWithPrices += 1
        }

        cardsWithSpecificPriceTypes[priceType] += 1
    }

    fmt.Printf("Total cards with price records: %d\n", numCardsWithPrices)
    fmt.Printf("Total price records: %d\n", totalPrices)
    fmt.Printf("MTGO price records: %d\n", mtgoPrices)
    fmt.Printf("MTGO foil price records: %d\n", mtgoFoilPrices)
    fmt.Printf("Paper price records: %d\n", paperPrices)
    fmt.Printf("Paper foil price records: %d\n", paperFoilPrices)
    for priceType, count := range cardsWithSpecificPriceTypes {
        fmt.Printf("Price type %s count: %d\n", priceType, count)
    }
}

