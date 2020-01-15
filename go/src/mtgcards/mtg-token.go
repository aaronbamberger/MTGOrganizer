package mtgcards

type MTGToken struct {
	MTGCardCommon
	ReverseRelated []string `json:"reverseRelated"`
}
