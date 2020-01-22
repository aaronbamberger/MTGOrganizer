package main

import "database/sql"
import _ "github.com/go-sql-driver/mysql"
import "log"
import "mtgcards"

func main() {
	allSets, err := mtgcards.DownloadAllPrintings(true)
	if err != nil {
		log.Fatal(err)
	}

	// Connect to the database
	db, err := sql.Open("mysql", "app_user:app_db_password@tcp(172.18.0.3)/mtg_cards")
	defer db.Close()
	if err != nil {
		log.Fatal(err)
	}

	db.SetMaxIdleConns(1000)

	err = mtgcards.ImportSetsToDb(db, allSets)
	if err != nil {
		log.Fatal(err)
	}
}
