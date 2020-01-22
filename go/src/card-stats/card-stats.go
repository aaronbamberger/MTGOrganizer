package main

import "log"
import "mtgcards"

func main() {
	allSets, err := mtgcards.DownloadAllPrintings(true)
	if err != nil {
		log.Fatal(err)
	}

	mtgcards.DevelopmentStats(allSets)
}

