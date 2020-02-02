package mtgcards

import "hash/fnv"
import "sort"

type MTGToken struct {
	MTGCardCommon
	ReverseRelated []string `json:"reverseRelated"`

    hash string
    hashValid bool
}

func (token *MTGToken) Hash() string {
    if !token.hashValid {
        hash := fnv.New128a()

        // Start with the hash of the common properties
        hash.Write([]byte(token.MTGCardCommon.Hash()))

        for _, reverseRelated := range token.ReverseRelated {
            hash.Write([]byte(reverseRelated))
        }

        token.hash = hashToHexString(hash)
        token.hashValid = true
    }

    return token.hash
}

func (token *MTGToken) Canonicalize() {
    // First, canonicalize the common properties
    token.MTGCardCommon.Canonicalize()

    sort.Strings(token.ReverseRelated)
}

type TokenByUUID []MTGToken

func (tokens TokenByUUID) Len() int {
	return len(tokens)
}

func (tokens TokenByUUID) Less(i, j int) bool {
	return tokens[i].UUID < tokens[j].UUID
}

func (tokens TokenByUUID) Swap(i, j int) {
	tokens[i], tokens[j] = tokens[j], tokens[i]
}
