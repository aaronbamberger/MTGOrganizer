package mtgcards

import "encoding/binary"
import "hash"
import "hash/fnv"

type MTGCardAlternateLanguageInfo struct {
	FlavorText string `json:"flavorText"`
	Language string `json:"language"`
	MultiverseId int `json:"multiverseId"`
	Name string `json:"name"`
	Text string `json:"text"`
	Type string `json:"type"`
}

func (langInfo MTGCardAlternateLanguageInfo) Hash() hash.Hash {
	hashRes := fnv.New128a()

	hashRes.Write([]byte(langInfo.FlavorText))
	hashRes.Write([]byte(langInfo.Language))
	binary.Write(hashRes, binary.BigEndian, langInfo.MultiverseId)
	hashRes.Write([]byte(langInfo.Name))
	hashRes.Write([]byte(langInfo.Text))
	hashRes.Write([]byte(langInfo.Type))

	return hashRes
}

type ByLanguage []MTGCardAlternateLanguageInfo

func (langInfo ByLanguage) Len() int {
	return len(langInfo)
}

func (langInfo ByLanguage) Less(i, j int) bool {
	return langInfo[i].Language < langInfo[j].Language
}

func (langInfo ByLanguage) Swap(i, j int) {
	langInfo[i], langInfo[j] = langInfo[j], langInfo[i]
}
