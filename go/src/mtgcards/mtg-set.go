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
}

func (set MTGSet) SetHash() hash.Hash {
	hashRes := fnv.New128a()
	binary.Write(hashRes, binary.BigEndian, set.BaseSetSize)
	hashRes.Write([]byte(set.Block))
	for _, card := range set.Cards {
		cardHash := card.CardHash()
		cardHashBytes := make([]byte, 0, cardHash.Size())
		hashRes.Write(cardHash.Sum(cardHashBytes))
	}
	hashRes.Write([]byte(set.Code))
	binary.Write(hashRes, binary.BigEndian, set.IsForeignOnly)
	binary.Write(hashRes, binary.BigEndian, set.IsFoilOnly)
	binary.Write(hashRes, binary.BigEndian, set.IsOnlineOnly)
	binary.Write(hashRes, binary.BigEndian, set.IsPartialPreview)
	hashRes.Write([]byte(set.KeyruneCode))
	hashRes.Write([]byte(set.MCMName))
	binary.Write(hashRes, binary.BigEndian, set.MCMId)
	hashRes.Write([]byte(set.MTGOCode))
	hashRes.Write([]byte(set.Name))
	hashRes.Write([]byte(set.ParentCode))
	hashRes.Write([]byte(set.ReleaseDate))
	binary.Write(hashRes, binary.BigEndian, set.TCGPlayerGroupId)
	binary.Write(hashRes, binary.BigEndian, set.TotalSetSize)
	translationsHash := set.Translations.Hash()
	translationsHashBytes := make([]byte, 0, translationsHash.Size())
	hashRes.Write(translationsHash.Sum(translationsHashBytes))
	hashRes.Write([]byte(set.Type))
	return hashRes
}

func (set MTGSet) Canonicalize() {
	sort.Sort(ByUUID(set.Cards))
	for _, card := range set.Cards {
		card.Canonicalize()
	}
}
