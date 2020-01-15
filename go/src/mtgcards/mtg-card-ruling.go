package mtgcards

import "hash"
import "hash/fnv"

type MTGCardRuling struct {
	Date string `json:"date"`
	Text string `json:"text"`
}

func (ruling MTGCardRuling) Hash() hash.Hash {
	hashRes := fnv.New128a()

	hashRes.Write([]byte(ruling.Date))
	hashRes.Write([]byte(ruling.Text))

	return hashRes
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
