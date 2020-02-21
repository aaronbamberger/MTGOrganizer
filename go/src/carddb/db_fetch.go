package carddb

import "context"
import "fmt"
import "database/sql"

type CardSearchResult struct {
    Name string `json:"name"`
    SetKeyruneCode string `json:"set_keyrune_code"`
}

type DBFetchContext struct {
    dbConn *sql.Conn
    cardSearchQuery *sql.Stmt
}

func NewDBFetchContext(db *sql.DB) (*DBFetchContext, error) {
    var fetchContext DBFetchContext
    var err error

    fetchContext.dbConn, err = db.Conn(context.Background())
    if err != nil {
        return nil, err
    }

    fetchContext.cardSearchQuery, err = fetchContext.dbConn.PrepareContext(
        context.Background(),
        `SELECT DISTINCT
        all_cards.name, sets.keyrune_code
        FROM
        all_cards INNER JOIN sets ON all_cards.set_id = sets.set_id
        WHERE all_cards.name RLIKE ?
        ORDER BY all_cards.name ASC`)
    if err != nil {
        return nil, err
    }

    return &fetchContext, nil
}

func (fc *DBFetchContext) SearchCardsByName(partialName string) ([]CardSearchResult, error) {
    res, err := fc.cardSearchQuery.Query(fmt.Sprintf("\\b%s", partialName))
    if err != nil {
        return nil, err
    }
    defer res.Close()

    cards := make([]CardSearchResult, 0)

    for res.Next() {
        var cardName string
        var setKeyruneCode string
        err = res.Scan(&cardName, &setKeyruneCode)
        if err != nil {
            return nil, err
        }
        cards = append(cards, CardSearchResult{Name: cardName, SetKeyruneCode: setKeyruneCode})
    }
    if err = res.Err(); err != nil {
        return nil, err
    }

    return cards, nil
}
