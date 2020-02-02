package mtgcards

import "hash/fnv"

type MTGCardRuling struct {
	Date string `json:"date"`
	Text string `json:"text"`

	hash string
	hashValid bool
}

func (ruling *MTGCardRuling) Hash() string {
	if !ruling.hashValid {
        hash := fnv.New128a()

		hash.Write([]byte(ruling.Date))
		hash.Write([]byte(ruling.Text))

        ruling.hash = hashToHexString(hash)
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
