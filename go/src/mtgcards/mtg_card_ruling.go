package mtgcards

import "fmt"
import "strings"

type MTGCardRuling struct {
	Date string `json:"date"`
	Text string `json:"text"`
}

func (ruling *MTGCardRuling) Hash() string {
    return objectHash(*ruling)
}

func (ruling MTGCardRuling) String() string {
    var b strings.Builder

    fmt.Fprintf(&b, "Date: %s\n", ruling.Date)
    fmt.Fprintf(&b, "Text: %s\n", ruling.Text)

    return b.String()
}

type ByDateThenText []MTGCardRuling

func (rulings ByDateThenText) Len() int {
	return len(rulings)
}

func (rulings ByDateThenText) Less(i, j int) bool {
    if rulings[i].Date != rulings[j].Date {
        return rulings[i].Date < rulings[j].Date
    } else {
        return rulings[i].Text < rulings[j].Text
    }
}

func (rulings ByDateThenText) Swap(i, j int) {
	rulings[i], rulings[j] = rulings[j], rulings[i]
}
