package mtgcards

import "encoding/binary"
import "encoding/hex"
import "hash"
import "hash/fnv"
import "reflect"
import "fmt"
import "sort"
import "unicode"
import "unicode/utf8"

func objectHash(object interface{}) string {
    objectType := reflect.TypeOf(object)
    objectValue := reflect.ValueOf(object)

    hash := fnv.New128a()

    switch objectType.Kind() {
    case reflect.Bool:
        fallthrough
    case reflect.Int:
        fallthrough
    case reflect.Int8:
        fallthrough
    case reflect.Int16:
        fallthrough
    case reflect.Int32:
        fallthrough
    case reflect.Int64:
        fallthrough
    case reflect.Uint:
        fallthrough
    case reflect.Uint8:
        fallthrough
    case reflect.Uint16:
        fallthrough
    case reflect.Uint32:
        fallthrough
    case reflect.Uint64:
        fallthrough
    case reflect.Float32:
        fallthrough
    case reflect.Float64:
        binary.Write(hash, binary.BigEndian, objectValue.Interface())

    case reflect.String:
        hash.Write([]byte(objectValue.String()))

    case reflect.Struct:
        // First, obtain a list of all of the object fields, sorted by name
        // We always want to hash the fields in the same order, so we do it by
        // field name order
        fieldNames := make([]string, 0, objectType.NumField())
        for i := 0; i < objectType.NumField(); i++ {
            fieldNames = append(fieldNames, objectType.Field(i).Name)
        }
        sort.Strings(fieldNames)

        for _, fieldName := range fieldNames {
            // Skip unexported fields
            firstRune, _ := utf8.DecodeRuneInString(fieldName)
            if unicode.IsLower(firstRune) {
                continue
            }

            field := objectValue.FieldByName(fieldName).Interface()
            hash.Write([]byte(objectHash(field)))
        }

    case reflect.Slice:
        // Recursively add the hash of each element to the hash
        // We also assume that the slice has already been sorted,
        // else this hash won't be deterministic
        for i := 0; i < objectValue.Len(); i++ {
            hash.Write([]byte(objectHash(objectValue.Index(i).Interface())))
        }

    case reflect.Map:
        // We only handle string-indexed maps right now, so make sure
        // that's the case here
        keys := objectValue.MapKeys()
        if len(keys) > 0 {
            keyKind := keys[0].Kind()
            if keyKind != reflect.String {
                panic(fmt.Errorf("Unexpected non-string map key %v for object %s",
                    keyKind, objectType))
            }
            // Since go maps don't have a defined iteration order,
            // ensure a repeatable hash by sorting the keyset, and using
            // that to define the iteration order
            sort.Sort(ByString(keys))

            // Add each key, and the hash of each value, to the hash
            for _, key := range keys {
                hash.Write([]byte(key.String()))
                hash.Write([]byte(objectHash(objectValue.MapIndex(key).Interface())))
            }
        }

    default:
        panic(fmt.Errorf("Unexpected object kind: %v for object %s",
            objectType.Kind(), objectType))
    }

    return hashToHexString(hash)
}

func hashToHexString(hashVal hash.Hash) string {
    hashBytes := make([]byte, 0, hashVal.Size())
    hashBytes = hashVal.Sum(hashBytes)
    return hex.EncodeToString(hashBytes)
}

type ByString []reflect.Value

func (strings ByString) Len() int {
    return len(strings)
}

func (strings ByString) Less(i, j int) bool {
    return strings[i].String() < strings[j].String()
}

func (strings ByString) Swap(i, j int) {
    strings[i], strings[j] = strings[j], strings[i]
}
