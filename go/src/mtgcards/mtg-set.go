package mtgcards

import "encoding/binary"
import "hash"
import "hash/fnv"
import "sort"

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
	Translations MTGSetNameTranslations `json:"translations"`
	Type string `json:"type"`
	hash hash.Hash
	hashValid bool
}

func (set *MTGSet) Hash() hash.Hash {
	if !set.hashValid {
		set.hash = fnv.New128a()
		binary.Write(set.hash, binary.BigEndian, set.BaseSetSize)
		set.hash.Write([]byte(set.Block))
		for _, card := range set.Cards {
			cardHash := card.Hash()
			cardHashBytes := make([]byte, 0, cardHash.Size())
			set.hash.Write(cardHash.Sum(cardHashBytes))
		}
		set.hash.Write([]byte(set.Code))
		binary.Write(set.hash, binary.BigEndian, set.IsForeignOnly)
		binary.Write(set.hash, binary.BigEndian, set.IsFoilOnly)
		binary.Write(set.hash, binary.BigEndian, set.IsOnlineOnly)
		binary.Write(set.hash, binary.BigEndian, set.IsPartialPreview)
		set.hash.Write([]byte(set.KeyruneCode))
		set.hash.Write([]byte(set.MCMName))
		binary.Write(set.hash, binary.BigEndian, set.MCMId)
		set.hash.Write([]byte(set.MTGOCode))
		set.hash.Write([]byte(set.Name))
		set.hash.Write([]byte(set.ParentCode))
		set.hash.Write([]byte(set.ReleaseDate))
		binary.Write(set.hash, binary.BigEndian, set.TCGPlayerGroupId)
		binary.Write(set.hash, binary.BigEndian, set.TotalSetSize)
		translationsHash := set.Translations.Hash()
		translationsHashBytes := make([]byte, 0, translationsHash.Size())
		set.hash.Write(translationsHash.Sum(translationsHashBytes))
		set.hash.Write([]byte(set.Type))
		set.hashValid = true
	}

	return set.hash
}

func (set *MTGSet) Canonicalize() {
	sort.Sort(ByUUID(set.Cards))
	for _, card := range set.Cards {
		card.Canonicalize()
	}
}
