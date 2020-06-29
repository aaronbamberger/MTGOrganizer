package mtgcards

import "fmt"
import "strings"

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
    return objectHash(*langInfo)
    /*
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
    */
}

func (langInfo MTGCardAlternateLanguageInfo) String() string {
    var b strings.Builder

    fmt.Fprintf(&b, "Flavor text: %s\n", langInfo.FlavorText)
    fmt.Fprintf(&b, "Language: %s\n", langInfo.FlavorText)
    fmt.Fprintf(&b, "MultiverseId: %d\n", langInfo.MultiverseId)
    fmt.Fprintf(&b, "Name: %s\n", langInfo.Name)
    fmt.Fprintf(&b, "Text: %s\n", langInfo.Text)
    fmt.Fprintf(&b, "Type: %s\n", langInfo.Type)

    return b.String()
}

type ByLanguageThenName []MTGCardAlternateLanguageInfo

func (langInfo ByLanguageThenName) Len() int {
	return len(langInfo)
}

func (langInfo ByLanguageThenName) Less(i, j int) bool {
    if langInfo[i].Language != langInfo[j].Language {
        return langInfo[i].Language < langInfo[j].Language
    } else {
        return langInfo[i].Name < langInfo[j].Name
    }
}

func (langInfo ByLanguageThenName) Swap(i, j int) {
    langInfo[i], langInfo[j] = langInfo[j], langInfo[i]
}
