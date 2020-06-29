package mtgcards

import "fmt"
import "reflect"
import "sort"
import "strings"

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
}

func (set *MTGSet) Diff(other *MTGSet) string {
    var b strings.Builder

    setType := reflect.TypeOf(*set)

    selfValue := reflect.ValueOf(*set)
    otherValue := reflect.ValueOf(*set)

    for i := 0; i < setType.NumField(); i++ {
        setField := setType.Field(i)
        selfField := selfValue.FieldByName(setField.Name)
        otherField := otherValue.FieldByName(setField.Name)
        switch selfField.Kind() {
        case reflect.Bool:
            if selfField.Bool() != otherField.Bool() {
                fmt.Fprintf(&b, "%s (%t | %t)\n", setField.Name,
                    selfField.Bool(), otherField.Bool())
            }
        case reflect.Int:
            if selfField.Int() != otherField.Int() {
                fmt.Fprintf(&b, "%s (%d | %d)\n", setField.Name,
                    selfField.Int(), otherField.Int())
            }
        case reflect.String:
            if selfField.String() != otherField.String() {
                fmt.Fprintf(&b, "%s (%s | %s)\n", setField.Name,
                    selfField.String(), otherField.String())
            }

        case reflect.Slice:
            if selfField.Len() != otherField.Len() {
                fmt.Fprintf(&b, "%s length (%d | %d)\n", setField.Name,
                    selfField.Len(), otherField.Len())
            }
            // TODO: Diff values

        case reflect.Map:
            // TODO: Diff maps

        }
    }

    return b.String()
}

func (set MTGSet) String() string {
    var b strings.Builder

    for i, card := range set.Cards {
        fmt.Fprintf(&b, "Card %d:\n%s\n", i, card)
    }
    for i, token := range set.Tokens {
        fmt.Fprintf(&b, "Token %d:\n%s\n", i, token.Name)
    }

    return b.String()
}

func (set *MTGSet) Hash() string {
    return objectHash(*set)
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
    sort.Sort(TokenByUUIDThenSide(set.Tokens))
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
