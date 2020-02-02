package mtgcards

import "encoding/binary"
import "hash/fnv"

type MTGCardAlternateLanguageInfo struct {
	FlavorText string `json:"flavorText"`
	Language string `json:"language"`
	MultiverseId int `json:"multiverseId"`
	Name string `json:"name"`
	Text string `json:"text"`
	Type string `json:"type"`

	hash string
	hashValid bool
}

func (langInfo *MTGCardAlternateLanguageInfo) Hash() string {
	if !langInfo.hashValid {
        hash := fnv.New128a()

		hash.Write([]byte(langInfo.FlavorText))
		hash.Write([]byte(langInfo.Language))
		binary.Write(hash, binary.BigEndian, langInfo.MultiverseId)
		hash.Write([]byte(langInfo.Name))
		hash.Write([]byte(langInfo.Text))
		hash.Write([]byte(langInfo.Type))

        langInfo.hash = hashToHexString(hash)
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
