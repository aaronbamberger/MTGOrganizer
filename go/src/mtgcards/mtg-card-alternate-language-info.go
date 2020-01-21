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

	hash hash.Hash
	hashValid bool
}

func (langInfo *MTGCardAlternateLanguageInfo) Hash() hash.Hash {
	if !langInfo.hashValid {
		langInfo.hash = fnv.New128a()

		langInfo.hash.Write([]byte(langInfo.FlavorText))
		langInfo.hash.Write([]byte(langInfo.Language))
		binary.Write(langInfo.hash, binary.BigEndian, langInfo.MultiverseId)
		langInfo.hash.Write([]byte(langInfo.Name))
		langInfo.hash.Write([]byte(langInfo.Text))
		langInfo.hash.Write([]byte(langInfo.Type))

		langInfo.hashValid = true
	}

	return langInfo.hash
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
