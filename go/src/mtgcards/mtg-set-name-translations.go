package mtgcards

import "hash"
import "hash/fnv"

type MTGSetNameTranslations struct {
	ChineseSimplified string `json:"Chinese Simplified"`
	ChineseTraditional string `json:"Chinese Traditional"`
	French string `json:"French"`
	German string `json:"German"`
	Italian string `json:"Italian"`
	Japanese string `json:"Japanese"`
	Korean string `json:"Korean"`
	PortugeseBrazil string `json:"Portugese (Brazil)"`
	Russian string `json:"Russian"`
	Spanish string `json:"Spanish"`
}

func (translations MTGSetNameTranslations) Hash() hash.Hash {
	hashRes := fnv.New128a()

	hashRes.Write([]byte(translations.ChineseSimplified))
	hashRes.Write([]byte(translations.ChineseTraditional))
	hashRes.Write([]byte(translations.French))
	hashRes.Write([]byte(translations.German))
	hashRes.Write([]byte(translations.Italian))
	hashRes.Write([]byte(translations.Japanese))
	hashRes.Write([]byte(translations.Korean))
	hashRes.Write([]byte(translations.PortugeseBrazil))
	hashRes.Write([]byte(translations.Russian))
	hashRes.Write([]byte(translations.Spanish))

	return hashRes
}
