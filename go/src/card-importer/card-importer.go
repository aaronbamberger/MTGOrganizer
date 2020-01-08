package main

import "compress/gzip"
import "encoding/json"
import "fmt"
import "mtgcards"
import "net/http"

func main() {
	resp, err := http.Get("https://www.mtgjson.com/files/AllPrintings.json.gz")
	if err != nil {
		fmt.Println("Error while downloading: %s\n", err)
	}
	defer resp.Body.Close()
	fmt.Printf("Response status: %s\n", resp.Status)
	fmt.Printf("Response proto: %s\n", resp.Proto)
	fmt.Printf("Response length: %d\n", resp.ContentLength)
	fmt.Printf("Response encodings: %v\n", resp.TransferEncoding)

	decompressor, err := gzip.NewReader(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	decoder := json.NewDecoder(decompressor)
	//var data map[string]interface{}
	var allSets map[string]mtgcards.MTGSet
	if err := decoder.Decode(&allSets); err != nil {
		fmt.Println(err)
		return
	}

	numCards := 0
	fmt.Printf("Number of sets retrieved: %d\n", len(allSets))
	for code, set := range allSets {
		fmt.Printf("Set %s (%s):\n", code, set.Name)
		fmt.Printf("\tNumber of cards in set: %d\n", len(set.Cards))
		numCards += len(set.Cards)
	}
	fmt.Printf("Total cards retrieved: %d\n", numCards)
}

