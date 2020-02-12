package main

import "fmt"
import "log"
import "mtgcards"

func main() {
    oldSets, err := mtgcards.DebugParseAllPrintingsGz("AllPrintings.json.gz.old")
    if err != nil {
        log.Fatal(err)
    }

    newSets, err := mtgcards.DebugParseAllPrintingsGz("AllPrintings.json.gz")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Old sets: %d, new sets: %d\n", len(oldSets), len(newSets))

    expectedUpdatedSets := 0
    expectedNewCards := 0
    expectedUpdatedCards := 0
    expectedNewTokens := 0
    expectedUpdatedTokens := 0
    for setCode, _ := range oldSets {
        oldSet := oldSets[setCode]
        oldSet.Canonicalize()
        newSet := newSets[setCode]
        newSet.Canonicalize()

        if oldSet.Hash() != newSet.Hash() {
            expectedUpdatedSets += 1
            fmt.Printf("Hash mismatch for set %s (%s %s)\n", setCode, oldSet.Hash(), newSet.Hash())
            fmt.Print(oldSet.Diff(&newSet))

            expectedNewCards += len(newSet.Cards) - len(oldSet.Cards)
            expectedNewTokens += len(newSet.Tokens) - len(oldSet.Tokens)

            mismatchedCards := 0
            mismatchedTokens := 0
            uuidToOldCard := make(map[string]*mtgcards.MTGCard)
            uuidToNewCard := make(map[string]*mtgcards.MTGCard)
            uuidToOldToken := make(map[string]*mtgcards.MTGToken)
            uuidToNewToken := make(map[string]*mtgcards.MTGToken)

            for idx, _  := range oldSet.Cards {
                card := &oldSet.Cards[idx]
                uuidToOldCard[card.UUID] = card
            }
            for idx, _  := range newSet.Cards {
                card := &newSet.Cards[idx]
                uuidToNewCard[card.UUID] = card
            }
            for idx, _  := range oldSet.Tokens {
                token := &oldSet.Tokens[idx]
                uuidToOldToken[token.UUID] = token
            }
            for idx, _  := range newSet.Tokens {
                token := &newSet.Tokens[idx]
                uuidToNewToken[token.UUID] = token
            }

            for uuid, _ := range uuidToOldCard {
                oldCard := uuidToOldCard[uuid]
                newCard := uuidToNewCard[uuid]

                if oldCard.Hash() != newCard.Hash() {
                    expectedUpdatedCards += 1
                    fmt.Printf("Hash mismatch for card %s (%s)\n",
                        uuid, oldCard.Name)
                    fmt.Print(oldCard.Diff(newCard))
                    mismatchedCards += 1
                }
            }

            for uuid, _ := range uuidToOldToken {
                oldToken := uuidToOldToken[uuid]
                newToken := uuidToNewToken[uuid]

                if oldToken.Hash() != newToken.Hash() {
                    expectedUpdatedTokens += 1
                    fmt.Printf("Hash mismatch for token %s (%s)\n",
                        uuid, oldToken.Name)
                    mismatchedTokens += 1
                }
            }

            fmt.Printf("Total old cards: %d, total new cards: %d, mismatched cards: %d\n",
                len(uuidToOldCard), len(uuidToNewCard), mismatchedCards)
            fmt.Printf("Total old tokens: %d, total new tokens: %d, mismatched tokens: %d\n",
                len(uuidToOldToken), len(uuidToNewToken), mismatchedTokens)
        }
    }

    fmt.Printf("Expected updated sets: %d\n", expectedUpdatedSets)
    fmt.Printf("Expected new cards: %d\n", expectedNewCards)
    fmt.Printf("Expected new tokens: %d\n", expectedNewTokens)
    fmt.Printf("Expected updated cards: %d\n", expectedUpdatedCards)
    fmt.Printf("Expected updated tokens: %d\n", expectedUpdatedTokens)
}

