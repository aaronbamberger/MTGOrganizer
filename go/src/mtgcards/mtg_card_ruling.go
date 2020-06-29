package mtgcards

import "fmt"
import "strings"

type MTGCardRuling struct {
	Date string `json:"date"`
	Text string `json:"text"`

	hash string
	hashValid bool
}

func (ruling *MTGCardRuling) Hash() string {
    return objectHash(*ruling)
    /*
	if !ruling.hashValid {
        hash := fnv.New128a()

		hash.Write([]byte(ruling.Date))
		hash.Write([]byte(ruling.Text))

        ruling.hash = hashToHexString(hash)
		ruling.hashValid = true
	}

	return ruling.hash
    */
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
