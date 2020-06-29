package mtgcards

import "fmt"
import "sort"
import "strings"

type MTGToken struct {
	MTGCardCommon
	ReverseRelated []string `json:"reverseRelated"`

    hash string
    hashValid bool
}

func (token *MTGToken) Hash() string {
    return objectHash(*token)
    /*
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
    */
}

func (token MTGToken) String() string {
    var b strings.Builder

    fmt.Fprintf(&b, "%s", token.MTGCardCommon)
    fmt.Fprintf(&b, "Reverse related: %v\n", token.ReverseRelated)

    return b.String()
}

func (token *MTGToken) Canonicalize() {
    // First, canonicalize the common properties
    token.MTGCardCommon.Canonicalize()

    sort.Strings(token.ReverseRelated)
}

type TokenByUUIDThenSide []MTGToken

func (tokens TokenByUUIDThenSide) Len() int {
	return len(tokens)
}

func (tokens TokenByUUIDThenSide) Less(i, j int) bool {
    if tokens[i].UUID != tokens[j].UUID {
        return tokens[i].UUID < tokens[j].UUID
    } else {
        return tokens[i].Side < tokens[j].Side
    }
}

func (tokens TokenByUUIDThenSide) Swap(i, j int) {
	tokens[i], tokens[j] = tokens[j], tokens[i]
}
