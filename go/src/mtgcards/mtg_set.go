package mtgcards

import "encoding/binary"
import "encoding/hex"
import "hash"
import "hash/fnv"
import "sort"

func hashToHexString(hashVal hash.Hash) string {
    hashBytes := make([]byte, 0, hashVal.Size())
    hashBytes = hashVal.Sum(hashBytes)
    return hex.EncodeToString(hashBytes)
}

type MTGSet struct {
	BaseSetSize int `json:"baseSetSize"`
	Block string `json:"block"`
	Cards []MTGCard `json:"cards"`
	Code string `json:"code"`
	IsForeignOnly bool `json:"isForeignOnly"`
	IsFoilOnly bool `json:"isFoilOnly"`
	IsOnlineOnly bool `json:"isOnlineOnly"`
	IsPartialPreview bool `json:"isPartialPreview"`
	KeyruneCode string `json:"keyruneCode"`
	MCMName string `json:"mcmName"`
	MCMId int `json:"mcmId"`
	MTGOCode string `json:"mtgoCode"`
	Name string `json:"name"`
	ParentCode string `json:"parentCode"`
	ReleaseDate string `json:"releaseDate"`
	TCGPlayerGroupId int `json:"tcgplayerGroupId"`
	TotalSetSize int `json:"totalSetSize"`
    Tokens []MTGToken `json:"tokens"`
	Translations map[string]string `json:"translations"`
	Type string `json:"type"`

	hash string
	hashValid bool
}

func (set *MTGSet) Hash() string {
	if !set.hashValid {
        hash := fnv.New128a()
		binary.Write(hash, binary.BigEndian, set.BaseSetSize)
		hash.Write([]byte(set.Block))

        // Cards
		for idx := range set.Cards {
            card := &set.Cards[idx]
            hash.Write([]byte(card.Hash()))
		}

        // Tokens
        for idx := range set.Tokens {
            token := &set.Tokens[idx]
            hash.Write([]byte(token.Hash()))
        }

		hash.Write([]byte(set.Code))
		binary.Write(hash, binary.BigEndian, set.IsForeignOnly)
		binary.Write(hash, binary.BigEndian, set.IsFoilOnly)
		binary.Write(hash, binary.BigEndian, set.IsOnlineOnly)
		binary.Write(hash, binary.BigEndian, set.IsPartialPreview)
		hash.Write([]byte(set.KeyruneCode))
		hash.Write([]byte(set.MCMName))
		binary.Write(hash, binary.BigEndian, set.MCMId)
		hash.Write([]byte(set.MTGOCode))
		hash.Write([]byte(set.Name))
		hash.Write([]byte(set.ParentCode))
		hash.Write([]byte(set.ReleaseDate))
		binary.Write(hash, binary.BigEndian, set.TCGPlayerGroupId)
		binary.Write(hash, binary.BigEndian, set.TotalSetSize)
		// Since go maps don't have a defined iteration order,
		// Ensure a repeatable hash by sorting the keyset, and using
		// that to define the iteration order
		translationLangs := make([]string, 0, len(set.Translations))
		for lang, _ := range set.Translations {
			translationLangs = append(translationLangs, lang)
		}
		sort.Strings(translationLangs)
		for _, lang := range translationLangs {
			hash.Write([]byte(lang))
			hash.Write([]byte(set.Translations[lang]))
		}
		hash.Write([]byte(set.Type))

        set.hash = hashToHexString(hash)

		set.hashValid = true
	}

	return set.hash
}

func (set *MTGSet) Canonicalize() {
    // Cards
	sort.Sort(CardByUUID(set.Cards))
	for idx := range set.Cards {
        // Need to access by index here so we're updating the cards
        // themselves, not copies
		set.Cards[idx].Canonicalize()
	}

    // Tokens
    sort.Sort(TokenByUUID(set.Tokens))
    for idx := range set.Tokens {
        // Same as above
        set.Tokens[idx].Canonicalize()
    }
}

type CardByUUID []MTGCard

func (cards CardByUUID) Len() int {
	return len(cards)
}

func (cards CardByUUID) Less(i, j int) bool {
	return cards[i].UUID < cards[j].UUID
}

func (cards CardByUUID) Swap(i, j int) {
	cards[i], cards[j] = cards[j], cards[i]
}
