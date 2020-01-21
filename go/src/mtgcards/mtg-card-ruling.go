package mtgcards

import "hash"
import "hash/fnv"

type MTGCardRuling struct {
	Date string `json:"date"`
	Text string `json:"text"`

	hash hash.Hash
	hashValid bool
}

func (ruling *MTGCardRuling) Hash() hash.Hash {
	if !ruling.hashValid {
		ruling.hash = fnv.New128a()

		ruling.hash.Write([]byte(ruling.Date))
		ruling.hash.Write([]byte(ruling.Text))

		ruling.hashValid = true
	}

	return ruling.hash
}

type ByDate []MTGCardRuling

func (rulings ByDate) Len() int {
	return len(rulings)
}

func (rulings ByDate) Less(i, j int) bool {
	return rulings[i].Date < rulings[j].Date
}

func (rulings ByDate) Swap(i, j int) {
	rulings[i], rulings[j] = rulings[j], rulings[i]
}
