package main

import "compress/gzip"
import "encoding/json"
//import "encoding/hex"
import "fmt"
import "mtgcards"
import "net/http"

type CardAndCount struct {
	Card string
	Count int
}

func main() {
	resp, err := http.Get("https://www.mtgjson.com/files/AllPrintings.json.gz")
	if err != nil {
		fmt.Println("Error while downloading: %s\n", err)
	}
	defer resp.Body.Close()
	fmt.Printf("Response status: %s\n", resp.Status)
	fmt.Printf("Response proto: %s\n", resp.Proto)
	fmt.Printf("Response length: %d\n", resp.ContentLength)
	fmt.Printf("Response encodings: %v\n", resp.TransferEncoding)

	decompressor, err := gzip.NewReader(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	decoder := json.NewDecoder(decompressor)
	var allSets map[string]mtgcards.MTGSet
	if err := decoder.Decode(&allSets); err != nil {
		fmt.Println(err)
		return
	}

	/*
	numCards := 0
	fmt.Printf("Number of sets retrieved: %d\n", len(allSets))
	for code, set := range allSets {
		fmt.Printf("Set %s (%s):\n", code, set.Name)
		fmt.Printf("\tNumber of cards in set: %d\n", len(set.Cards))
		numCards += len(set.Cards)
	}
	*/

	/*
	atomicDataHashMap := make(map[string]CardAndCount)

	for _, set := range allSets {
		for _, card := range set.Cards {
			cardAtomicHash := card.AtomicPropertiesHash()
			var hashBytes []byte
			hashBytes = cardAtomicHash.Sum(hashBytes)
			hashString := hex.EncodeToString(hashBytes)
			existingCard, ok := atomicDataHashMap[hashString]
			if ok {
				if existingCard.Card != card.Name {
					fmt.Printf("Collision at hash %s between card %s and card %s\n",
						hashString,
						atomicDataHashMap[hashString].Card, card.Name)
				} else {
					atomicDataHashMap[hashString] =
						CardAndCount{Card: existingCard.Card, Count: existingCard.Count + 1}
				}
			} else {
					atomicDataHashMap[hashString] =
						CardAndCount{Card: card.Name, Count: 1}
			}
		}
	}

	for _, cardAndCount := range atomicDataHashMap {
		if cardAndCount.Count > 1 {
			fmt.Printf("Card %s, count %d\n", cardAndCount.Card, cardAndCount.Count)
		}
		//fmt.Printf("Card %s (hash %s) number of printings: %d\n",
		//	cardAndCount.Card, hash, cardAndCount.Count)
	}

	fmt.Printf("Total cards retrieved: %d\n", numCards)
	fmt.Printf("Total unique cards retrieved: %d\n", len(atomicDataHashMap))
	*/

	maxSetBlockNameLen := 0
	maxSetCodeLen := 0
	maxSetKeyruneCodeLen := 0
	maxSetMcmNameLen := 0
	maxSetMTGOCodeLen := 0
	maxSetNameLen := 0
	maxSetParentCodeLen := 0
	maxSetTranslatedNameLen := 0
	maxSetTypeLen := 0

	maxArtistLen := 0
	maxBorderColorLen := 0
	maxDuelDeckLen := 0
	maxFlavorTextLen := 0
	maxAltLangFlavorTextLen := 0
	maxAltLangLanguageLen := 0
	maxAltLangNameLen := 0
	maxAltLangTextLen := 0
	maxAltLangTypeLen := 0
	maxFrameEffectsLen := 0
	maxFrameVersionLen := 0
	maxHandLen := 0
	maxLayoutLen := 0
	maxLifeLen := 0
	maxLoyaltyLen := 0
	maxManaCost := 0
	maxNameLen := 0
	maxNumberLen := 0
	maxOriginalTextLen := 0
	maxOriginalTypeLen := 0
	maxPowerLen := 0
	maxPurchaseUrlLen := 0
	maxRarityLen := 0
	maxRulingTextLen := 0
	maxSideLen := 0
	maxSubtypesLen := 0
	maxSupertypesLen := 0
	maxTextLen := 0
	maxToughnessLen := 0
	maxTypeLen := 0
	maxWatermarkLen := 0

	for _, set := range allSets {
		if len(set.Block) > maxSetBlockNameLen {
			maxSetBlockNameLen = len(set.Block)
		}
		if len(set.Code) > maxSetCodeLen {
			maxSetCodeLen = len(set.Code)
		}
		if len(set.KeyruneCode) > maxSetKeyruneCodeLen {
			maxSetKeyruneCodeLen = len(set.KeyruneCode)
		}
		if len(set.MCMName) > maxSetMcmNameLen {
			maxSetMcmNameLen = len(set.MCMName)
		}
		if len(set.MTGOCode) > maxSetMTGOCodeLen {
			maxSetMTGOCodeLen = len(set.MTGOCode)
		}
		if len(set.Name) > maxSetNameLen {
			maxSetNameLen = len(set.Name)
		}
		if len(set.ParentCode) > maxSetParentCodeLen {
			maxSetParentCodeLen = len(set.ParentCode)
		}
		if len(set.Translations.ChineseSimplified) > maxSetTranslatedNameLen {
			maxSetTranslatedNameLen = len(set.Translations.ChineseSimplified)
		}
		if len(set.Translations.ChineseTraditional) > maxSetTranslatedNameLen {
			maxSetTranslatedNameLen = len(set.Translations.ChineseTraditional)
		}
		if len(set.Translations.French) > maxSetTranslatedNameLen {
			maxSetTranslatedNameLen = len(set.Translations.French)
		}
		if len(set.Translations.German) > maxSetTranslatedNameLen {
			maxSetTranslatedNameLen = len(set.Translations.German)
		}
		if len(set.Translations.Italian) > maxSetTranslatedNameLen {
			maxSetTranslatedNameLen = len(set.Translations.Italian)
		}
		if len(set.Translations.Japanese) > maxSetTranslatedNameLen {
			maxSetTranslatedNameLen = len(set.Translations.Japanese)
		}
		if len(set.Translations.Korean) > maxSetTranslatedNameLen {
			maxSetTranslatedNameLen = len(set.Translations.Korean)
		}
		if len(set.Translations.PortugeseBrazil) > maxSetTranslatedNameLen {
			maxSetTranslatedNameLen = len(set.Translations.PortugeseBrazil)
		}
		if len(set.Translations.Russian) > maxSetTranslatedNameLen {
			maxSetTranslatedNameLen = len(set.Translations.Russian)
		}
		if len(set.Translations.Spanish) > maxSetTranslatedNameLen {
			maxSetTranslatedNameLen = len(set.Translations.Spanish)
		}
		if len(set.Type) > maxSetTypeLen {
			maxSetTypeLen = len(set.Type)
		}

		for _, card := range set.Cards {
			if len(card.Artist) > maxArtistLen {
				maxArtistLen = len(card.Artist)
			}
			if len(card.BorderColor) > maxBorderColorLen {
				maxBorderColorLen = len(card.BorderColor)
			}
			if len(card.DuelDeck) > maxDuelDeckLen {
				maxDuelDeckLen = len(card.DuelDeck)
			}
			if len(card.FlavorText) > maxFlavorTextLen {
				maxFlavorTextLen = len(card.FlavorText)
			}
			for _, langInfo := range card.AlternateLanguageData {
				if len(langInfo.FlavorText) > maxAltLangFlavorTextLen {
					maxAltLangFlavorTextLen = len(langInfo.FlavorText)
				}
				if len(langInfo.Language) > maxAltLangLanguageLen {
					maxAltLangLanguageLen = len(langInfo.Language)
				}
				if len(langInfo.Name) > maxAltLangNameLen {
					maxAltLangNameLen = len(langInfo.Name)
				}
				if len(langInfo.Text) > maxAltLangTextLen {
					maxAltLangTextLen = len(langInfo.Text)
				}
				if len(langInfo.Type) > maxAltLangTypeLen {
					maxAltLangTypeLen = len(langInfo.Type)
				}
			}
			for _, frameEffect := range card.FrameEffects {
				if len(frameEffect) > maxFrameEffectsLen {
					maxFrameEffectsLen = len(frameEffect)
				}
			}
			if len(card.FrameVersion) > maxFrameVersionLen {
				maxFrameVersionLen = len(card.FrameVersion)
			}
			if len(card.Hand) > maxHandLen {
				maxHandLen = len(card.Hand)
			}
			if len(card.Layout) > maxLayoutLen {
				maxLayoutLen = len(card.Layout)
			}
			if len(card.Life) > maxLifeLen {
				maxLifeLen = len(card.Life)
			}
			if len(card.Loyalty) > maxLoyaltyLen {
				maxLoyaltyLen = len(card.Loyalty)
			}
			if len(card.ManaCost) > maxManaCost {
				maxManaCost = len(card.ManaCost)
			}
			if len(card.Name) > maxNameLen {
				maxNameLen = len(card.Name)
			}
			if len(card.Number) > maxNumberLen {
				maxNumberLen = len(card.Number)
			}
			if len(card.OriginalText) > maxOriginalTextLen {
				maxOriginalTextLen = len(card.OriginalText)
			}
			if len(card.OriginalType) > maxOriginalTypeLen {
				maxOriginalTypeLen = len(card.OriginalType)
			}
			if len(card.Power) > maxPowerLen {
				maxPowerLen = len(card.Power)
			}
			if len(card.PurchaseURLs.Cardmarket) >= maxPurchaseUrlLen {
				maxPurchaseUrlLen = len(card.PurchaseURLs.Cardmarket)
			}
			if len(card.PurchaseURLs.TCGPlayer) >= maxPurchaseUrlLen {
				maxPurchaseUrlLen = len(card.PurchaseURLs.TCGPlayer)
			}
			if len(card.PurchaseURLs.MTGStocks) >= maxPurchaseUrlLen {
				maxPurchaseUrlLen = len(card.PurchaseURLs.MTGStocks)
			}
			if len(card.Rarity) > maxRarityLen {
				maxRarityLen = len(card.Rarity)
			}
			for _, ruling := range card.Rulings {
				if len(ruling.Text) >= maxRulingTextLen {
					maxRulingTextLen = len(ruling.Text)
				}
			}
			if len(card.Side) > maxSideLen {
				maxSideLen = len(card.Side)
			}
			for _, subtype := range card.Subtypes {
				if len(subtype) > maxSubtypesLen {
					maxSubtypesLen = len(subtype)
				}
			}
			for _, supertype := range card.Supertypes {
				if len(supertype) > maxSupertypesLen {
					maxSupertypesLen = len(supertype)
				}
			}
			if len(card.Text) > maxTextLen {
				maxTextLen = len(card.Text)
			}
			if len(card.Toughness) > maxToughnessLen {
				maxToughnessLen = len(card.Toughness)
			}
			if len(card.Type) > maxTypeLen {
				maxTypeLen = len(card.Type)
			}
			if len(card.Watermark) > maxWatermarkLen {
				maxWatermarkLen = len(card.Watermark)
			}
		}
	}

	fmt.Printf("Max set block name: %d\n", maxSetBlockNameLen)
	fmt.Printf("Max set code len: %d\n", maxSetCodeLen)
	fmt.Printf("Max set keyrune code len: %d\n", maxSetKeyruneCodeLen)
	fmt.Printf("Max set mcm name len: %d\n", maxSetMcmNameLen)
	fmt.Printf("Max set mtgo code len: %d\n", maxSetMTGOCodeLen)
	fmt.Printf("Max set name len: %d\n", maxSetNameLen)
	fmt.Printf("Max set parent code len: %d\n", maxSetParentCodeLen)
	fmt.Printf("Max set translated name len: %d\n", maxSetTranslatedNameLen)
	fmt.Printf("Max set type len: %d\n", maxSetTypeLen)

	fmt.Printf("Max artist len: %d\n", maxArtistLen)
	fmt.Printf("Max border color len: %d\n", maxBorderColorLen)
	fmt.Printf("Max duel deck len: %d\n", maxDuelDeckLen)
	fmt.Printf("Max flavor text len: %d\n", maxFlavorTextLen)
	fmt.Printf("Max alt lang flavor text len: %d\n", maxAltLangFlavorTextLen)
	fmt.Printf("Max alt lang language len: %d\n", maxAltLangLanguageLen)
	fmt.Printf("Max alt lang name len: %d\n", maxAltLangNameLen)
	fmt.Printf("Max alt lang text len: %d\n", maxAltLangTextLen)
	fmt.Printf("Max alt lang type len: %d\n", maxAltLangTypeLen)
	fmt.Printf("Max frame effects len: %d\n", maxFrameEffectsLen)
	fmt.Printf("Max frame version len: %d\n", maxFrameVersionLen)
	fmt.Printf("Max hand len: %d\n", maxHandLen)
	fmt.Printf("Max layout len: %d\n", maxLayoutLen)
	fmt.Printf("Max life len: %d\n", maxLifeLen)
	fmt.Printf("Max loyalty len: %d\n", maxLoyaltyLen)
	fmt.Printf("Max mana cost len: %d\n", maxManaCost)
	fmt.Printf("Max name len: %d\n", maxNameLen)
	fmt.Printf("Max number len: %d\n", maxNumberLen)
	fmt.Printf("Max original text len: %d\n", maxOriginalTextLen)
	fmt.Printf("Max original type len: %d\n", maxOriginalTypeLen)
	fmt.Printf("Max power len: %d\n", maxPowerLen)
	fmt.Printf("Max purchase url len: %d\n", maxPurchaseUrlLen)
	fmt.Printf("Max rarity len: %d\n", maxRarityLen)
	fmt.Printf("Max ruling text len: %d\n", maxRulingTextLen)
	fmt.Printf("Max side length: %d\n", maxSideLen)
	fmt.Printf("Max subtypes len: %d\n", maxSubtypesLen)
	fmt.Printf("Max supertypes len: %d\n", maxSupertypesLen)
	fmt.Printf("Max text length: %d\n", maxTextLen)
	fmt.Printf("Max toughness len: %d\n", maxToughnessLen)
	fmt.Printf("Max type len: %d\n", maxTypeLen)
	fmt.Printf("Max watermark len: %d\n", maxWatermarkLen)

	/*
	colorIdentityMap := make(map[string]int)
	colorIndicatorMap := make(map[string]int)
	colorMap := make(map[string]int)

	for _, set := range allSets {
		for _, card := range set.Cards {
			for _, color := range card.ColorIdentity {
				if numHits, ok := colorIdentityMap[color]; ok {
					colorIdentityMap[color] = numHits + 1
				} else {
					colorIdentityMap[color] = 1
				}
			}
			for _, color := range card.ColorIndicator {
				if numHits, ok := colorIndicatorMap[color]; ok {
					colorIndicatorMap[color] = numHits + 1
				} else {
					colorIndicatorMap[color] = 1
				}
			}
			for _, color := range card.Colors {
				if numHits, ok := colorMap[color]; ok {
					colorMap[color] = numHits + 1
				} else {
					colorMap[color] = 1
				}
			}
		}
	}

	for color, count := range colorIdentityMap {
		fmt.Printf("Color Identity %s, Hits %d\n", color, count)
	}
	for color, count := range colorIndicatorMap {
		fmt.Printf("Color Indicator %s, Hits %d\n", color, count)
	}
	for color, count := range colorMap {
		fmt.Printf("Color %s, Hits %d\n", color, count)
	}
	*/
}

