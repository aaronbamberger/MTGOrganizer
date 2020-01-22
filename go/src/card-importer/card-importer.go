package main

//import "bufio"
import "database/sql"
import _ "github.com/go-sql-driver/mysql"
//import _ "net/http/pprof"
//import "net/http"
import "log"
import "mtgcards"
//import "os"
import "sync"

func main() {
	/*
	go func() {
		log.Println(http.ListenAndServe("192.168.50.185:8085", nil))
	}()
	*/
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

	err = mtgcards.CreateDbQueries(db)
	if err != nil {
		log.Print(err)
		return
	}
	defer mtgcards.CloseDbQueries()

	totalSets := len(allSets)
	currentSet := 1
	//set := allSets["7ED"]

	//reader := bufio.NewReader(os.Stdin)
	//_, _ = reader.ReadString('\n')
	var setWaitGroup sync.WaitGroup
	for _, set := range allSets {
		log.Printf("Processing set with code %s (%d of %d)\n", set.Code, currentSet, totalSets)
		currentSet += 1
		setWaitGroup.Add(1)
		go mtgcards.MaybeInsertSetToDb(db, &setWaitGroup, set)
	}
	//_, _ = reader.ReadString('\n')
	log.Printf("Waiting on all set goroutines to finish\n")
	setWaitGroup.Wait()
}
