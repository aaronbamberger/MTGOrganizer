package backend

import "context"
import "fmt"
import "log"
import "database/sql"
import "strings"

type CardSearchResult struct {
    Name string `json:"name"`
    SetName string `json:"setName"`
    SetKeyruneCode string `json:"setKeyruneCode"`
}

func cardSearch(db *sql.DB,
        request string,
        done chan interface{},
        respChan chan ResponseMessage) {

    dbConn, err := db.Conn(context.Background())
    if err != nil {
        log.Print(err)
        return
    }
    defer dbConn.Close()

    processedName := strings.ReplaceAll(request, "s", "'?s")
    res, err := dbConn.QueryContext(
        context.Background(),
        `SELECT DISTINCT
        all_cards.name, sets.name, sets.keyrune_code
        FROM
        all_cards INNER JOIN sets ON all_cards.set_id = sets.set_id
        WHERE all_cards.name RLIKE ?
        ORDER BY all_cards.name ASC`,
        fmt.Sprintf("\\b%s", processedName))
    if err != nil {
        log.Print(err)
        return
    }
    defer res.Close()

    cards := make([]CardSearchResult, 0)

    for res.Next() {
        var cardName string
        var setName string
        var setKeyruneCode string
        err = res.Scan(&cardName, &setName, &setKeyruneCode)
        if err != nil {
            log.Print(err)
            return
        }
        cards = append(cards,
            CardSearchResult{
                Name: cardName,
                SetName: setName,
                SetKeyruneCode: setKeyruneCode})
    }
    if err = res.Err(); err != nil {
        log.Print(err)
        return
    }

    respChan <- ResponseMessage{Type:CardSearchResponse, Value: cards}
}
