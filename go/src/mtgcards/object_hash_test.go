package mtgcards

import "fmt"
import "testing"
import "math/rand"
import "reflect"
import "unicode"
import "unicode/utf8"

func loadTestCardData() (map[string]MTGSet, error) {
    testCardDataGz, err := tryGz(true, true, "TestCardData", decodeSets)
    if err != nil {
        return nil, err
    }

    testCardData, ok := testCardDataGz.(map[string]MTGSet)
    if !ok {
        return nil, fmt.Errorf("Unable to parse test card data JSON to correct type")
    }

    return testCardData, nil
}

func shuffleSlices(object interface{}) {
    objectType := reflect.TypeOf(object)
    objectValue := reflect.ValueOf(object)

    switch objectType.Kind() {
    case reflect.Slice:
        // Randomly shuffle the slice
        rand.Shuffle(objectValue.Len(), func(i, j int) {
            // Need to convert back and forth from interface
            // or else the values act like pointers and the swap doesn't work
            elem1 := objectValue.Index(i).Interface()
            elem2 := objectValue.Index(j).Interface()
            objectValue.Index(j).Set(reflect.ValueOf(elem1))
            objectValue.Index(i).Set(reflect.ValueOf(elem2))
        })

        // Recursively try and shuffle each slice element
        for i := 0; i < objectValue.Len(); i++ {
            shuffleSlices(objectValue.Index(i).Interface())
        }

    case reflect.Struct:
        // Recursively try and shuffle each field
        for i := 0; i < objectType.NumField(); i++ {
            // Skip unexported fields
            field := objectType.Field(i)
            firstRune, _ := utf8.DecodeRuneInString(field.Name)
            if unicode.IsLower(firstRune) {
                continue
            }

            shuffleSlices(objectValue.Field(i).Interface())
        }

    case reflect.Map:
        // Recursively try and shuffle each map entry
        keys := objectValue.MapKeys()
        for _, key := range keys {
            shuffleSlices(objectValue.MapIndex(key).Interface())
        }
    }
}

func checkObjectIsCanonical(object interface{}, t *testing.T) {
    objectType := reflect.TypeOf(object)
    objectValue := reflect.ValueOf(object)

    // The only types of objects we care about are slices, structs, and maps,
    // everything else is a scalar so is already canonical
    switch objectType.Kind() {
    case reflect.Slice:
        if objectValue.Len() > 0 {
            if objectValue.Index(0).Kind() == reflect.String {
                for i := 0; i < objectValue.Len() - 1; i++ {
                    if objectValue.Index(i).String() > objectValue.Index(i + 1).String() {
                        t.Errorf("String slice %s isn't sorted", objectType.Name())
                    }
                }
            } else {
                // Recursively check if each object is canonical
                for i := 0; i < objectValue.Len(); i++ {
                    checkObjectIsCanonical(objectValue.Index(i).Interface(), t)
                }
            }
        }

    case reflect.Struct:
        // Recursively check if each field is canonical
        for i := 0; i < objectType.NumField(); i++ {
            // Skip unexported fields
            field := objectType.Field(i)
            firstRune, _ := utf8.DecodeRuneInString(field.Name)
            if unicode.IsLower(firstRune) {
                continue
            }

            checkObjectIsCanonical(objectValue.Field(i).Interface(), t)
        }

    case reflect.Map:
        // Recursively check if each map entry is canonical
        keys := objectValue.MapKeys()
        for _, key := range keys {
            checkObjectIsCanonical(objectValue.MapIndex(key).Interface(), t)
        }
    }

}

func TestSetCanonicalize(t *testing.T) {
    cardData, err := loadTestCardData()
    if err != nil {
        t.Fatal(err)
    }

    for setName := range cardData {
        set := cardData[setName]

        // First, randomize the set data
        shuffleSlices(set)

        // Now, canonicalize the set, and check to see if it's sorted
        set.Canonicalize()

        checkObjectIsCanonical(set, t)
    }
}

func TestObjectHashStability(t *testing.T) {
    cardData1, err := loadTestCardData()
    if err != nil {
        t.Fatal(err)
    }
    cardData2, err := loadTestCardData()
    if err != nil {
        t.Fatal(err)
    }

    if len(cardData1) != len(cardData2) {
        t.Errorf("Sizes of test card data (%d and %d) don't match",
            len(cardData1), len(cardData2))
    }

    for setName := range cardData1 {
        //setName := "GS1"
        set1 := cardData1[setName]
        set2 := cardData2[setName]

        // Randomly shuffle all the slices in the test data, so we make sure our
        // hash stability isn't just due to using data that happens to be in the
        // same order
        shuffleSlices(set1)
        shuffleSlices(set2)

        hashShuffle1 := objectHash(set1)
        hashShuffle2 := objectHash(set2)

        // Now, canonicalize the sets, and make sure the hashes match
        set1.Canonicalize()
        set2.Canonicalize()

        hashCorrect1 := objectHash(set1)
        hashCorrect2 := objectHash(set2)

        // Some sets only have 1 card, in that case, these checks won't make
        // sense because there's nothing to randomize
        if len(set1.Cards) > 1 {
            if hashShuffle1 == hashShuffle2 {
                t.Errorf(`Set %s has a consistent hash between runs (%s == %s)
                    "even though data was randomized`,
                    setName, hashShuffle1, hashShuffle2)
            }

            // First, make sure that the hash has changed after canonicalizing (it should)
            if hashShuffle1 == hashCorrect1 {
                t.Errorf("Set %s has consistent hash between shuffled and canonical (%s == %s)",
                    setName, hashShuffle1, hashCorrect1)
            }
            if hashShuffle2 == hashCorrect2 {
                t.Errorf("Set %s has consistent hash between shuffled and canonical (%s == %s)",
                    setName, hashShuffle2, hashCorrect2)
            }
        }

        if hashCorrect1 != hashCorrect2 {
            t.Errorf("Set %s doesn't have consistent hash between runs (%s != %s) after canonicalizing",
                setName, hashCorrect1, hashCorrect2)
            for i, card := range set1.Cards {
                otherCard := set2.Cards[i]
                if card.Hash() != otherCard.Hash() {
                    t.Logf("Card %s hash doesn't match card %s", card.Name, otherCard.Name)
                }
            }
            for i, token := range set1.Tokens {
                otherToken := set2.Tokens[i]
                if token.Hash() != otherToken.Hash() {
                    t.Logf("Token %s hash doesn't match token %s", token.Name, otherToken.Name)
                }
            }

        }
    }
}
