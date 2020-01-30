package mtgcards

import "context"
import "database/sql"
import "fmt"
import "log"
import "strings"

const (
    AtomicPropRequestGetRef = iota
    AtomicPropRequestRemoveRef
    AtomicPropRequestGetHash
)

type atomicPropRequest struct {
    RequestType int
    ResponseChan chan atomicPropResponse
    AtomicPropertiesHash string
    AtomicPropertiesId int64
    ColorIdentity []string
    ColorIndicator []string
    Colors []string
    ConvertedManaCost float32
    EDHRecRank int
    FaceConvertedManaCost float32
    Hand string
    IsReserved bool
    Layout string
    Life string
    Loyalty string
    ManaCost string
    MTGStocksId int
    Name string
    Power string
    ScryfallOracleId string
    Side string
    Text string
    Toughness string
    Type string
}

type atomicPropResponse struct {
    AtomicPropertiesId int64
    AtomicPropertiesHash string
    NewRecordAdded bool
    Error error
}

type atomicPropertiesDbQueries struct {
	NumAtomicPropertiesQuery *sql.Stmt
	AtomicPropertiesIdQuery *sql.Stmt
	AtomicPropertiesHashQuery *sql.Stmt
	InsertAtomicPropertiesQuery *sql.Stmt
    GetRefCntQuery *sql.Stmt
	UpdateRefCntQuery *sql.Stmt
    DeleteAtomicPropertiesQuery *sql.Stmt
}

func atomicPropDbThread(
        db *sql.DB,
        requests chan atomicPropRequest,
        results chan error,
        quit chan interface{}) {
    // Prepare the queries we'll need
    log.Printf("Starting atomic prop db thread...")
    queries, err := prepareDbAtomicQueries(db)
    if err != nil {
        queries.cleanup()
        results <- err
        return
    }
    defer queries.cleanup()
    log.Printf("Atomic prop db thread queries prepared")

    // Open a connection to the database
    ctx := context.Background()
	conn, err := db.Conn(ctx)
	if err != nil {
        results <- err
		return
	}
	defer conn.Close()

    // Indicate to the caller that we've successfully initialized and we're
    // ready to start processing requests
    results <- nil

    // Now, spin and wait for requests to come in
    done := false
    for !done {
        select {
        case req := <-requests:
            var resp atomicPropResponse

            // Start with a transaction
            tx, err := conn.BeginTx(ctx, nil)
            if err != nil {
                resp.Error = err
                req.ResponseChan <- resp
                continue
            }

            // Perform the request
            switch req.RequestType {
            case AtomicPropRequestGetRef:
                numAtomicPropertiesQuery := tx.Stmt(queries.NumAtomicPropertiesQuery)
                atomicPropertiesIdQuery := tx.Stmt(queries.AtomicPropertiesIdQuery)
                insertAtomicPropertiesQuery := tx.Stmt(queries.InsertAtomicPropertiesQuery)
                updateRefCntQuery := tx.Stmt(queries.UpdateRefCntQuery)

                atomicPropId, refCnt, exists, err := getAtomicPropertiesId(
                    numAtomicPropertiesQuery,
                    atomicPropertiesIdQuery,
                    req.AtomicPropertiesHash,
                    req.ScryfallOracleId)
                if err != nil {
                    tx.Rollback()
                    resp.Error = err
                    req.ResponseChan <- resp
                    continue
                }

                if exists {
                    // If this record already exists, just increase the refcount
                    // and return the id
                    resp.AtomicPropertiesId = atomicPropId
                    _, err := updateRefCntQuery.Exec(refCnt + 1, atomicPropId)
                    if err != nil {
                        tx.Rollback()
                        resp.Error = err
                        req.ResponseChan <- resp
                        continue
                    }
                    resp.Error = tx.Commit()
                    req.ResponseChan <- resp
                } else {
                    // If this record doesn't exist, add it, and return the id
                    resp.NewRecordAdded = true
                    atomicPropId, err := insertAtomicPropertiesRecord(
                        insertAtomicPropertiesQuery,
                        &req)
                    if err != nil {
                        tx.Rollback()
                        resp.Error = err
                        req.ResponseChan <- resp
                        continue
                    }
                    resp.AtomicPropertiesId = atomicPropId
                    resp.Error = tx.Commit()
                    req.ResponseChan <- resp
                }

            case AtomicPropRequestRemoveRef:
                getRefCntQuery := tx.Stmt(queries.GetRefCntQuery)
                updateRefCntQuery := tx.Stmt(queries.UpdateRefCntQuery)
                deleteAtomicPropertiesQuery := tx.Stmt(queries.DeleteAtomicPropertiesQuery)

                // Get the refcount of this record.  If it's over 1, just decrease
                // it by 1.  If it's at 1, delete the entire record
                result := getRefCntQuery.QueryRow(req.AtomicPropertiesId)
                var refCnt int
                if err := result.Scan(&refCnt); err != nil {
                    tx.Rollback()
                    resp.Error = err
                    req.ResponseChan <- resp
                    continue
                }

                if refCnt > 1 {
                    _, err := updateRefCntQuery.Exec(refCnt - 1, req.AtomicPropertiesId)
                    if err != nil {
                        tx.Rollback()
                        resp.Error = err
                        req.ResponseChan <- resp
                        continue
                    }
                } else {
                    _, err := deleteAtomicPropertiesQuery.Exec(req.AtomicPropertiesId)
                    if err != nil {
                        tx.Rollback()
                        resp.Error = err
                        req.ResponseChan <- resp
                        continue
                    }
                }

                resp.Error = tx.Commit()
                req.ResponseChan <- resp

            case AtomicPropRequestGetHash:
                atomicPropertiesHashQuery := tx.Stmt(queries.AtomicPropertiesHashQuery)

                result := atomicPropertiesHashQuery.QueryRow(req.AtomicPropertiesId)
                var atomicPropertiesHash string
                if err := result.Scan(&atomicPropertiesHash); err != nil {
                    tx.Rollback()
                    resp.Error = err
                    req.ResponseChan <- resp
                    continue
                }

                resp.AtomicPropertiesHash = atomicPropertiesHash
                resp.Error = tx.Commit()
                req.ResponseChan <- resp
            }
        case <-quit:
            done = true
        }
    }

    results <- nil
}

func prepareDbAtomicQueries(db *sql.DB) (*atomicPropertiesDbQueries, error) {
	var err error
    var queries atomicPropertiesDbQueries

	queries.NumAtomicPropertiesQuery, err = db.Prepare(`SELECT COUNT(scryfall_oracle_id)
		FROM atomic_card_data
		WHERE card_data_hash = ?`)
	if err != nil {
		return &queries, err
	}

	queries.AtomicPropertiesIdQuery, err = db.Prepare(`SELECT atomic_card_data_id,
		ref_cnt,
		scryfall_oracle_id
		FROM atomic_card_data
		WHERE card_data_hash = ?`)
	if err != nil {
		return &queries, err
	}

	queries.AtomicPropertiesHashQuery, err = db.Prepare(`SELECT card_data_hash
		FROM atomic_card_data
		WHERE atomic_card_data_id = ?`)
	if err != nil {
		return &queries, err
	}

	queries.InsertAtomicPropertiesQuery, err = db.Prepare(`INSERT INTO atomic_card_data
		(card_data_hash, color_identity, color_indicator, colors, converted_mana_cost,
		edhrec_rank, face_converted_mana_cost, hand, is_reserved, layout, life,
		loyalty, mana_cost, mtgstocks_id, name, card_power, scryfall_oracle_id,
		side, text, toughness, card_type)
		VALUES
		(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return &queries, err
	}

    queries.GetRefCntQuery, err = db.Prepare(`SELECT ref_cnt
        FROM atomic_card_data
        WHERE atomic_card_data_id = ?`)
    if err != nil {
        return &queries, err
    }

	queries.UpdateRefCntQuery, err = db.Prepare(`UPDATE atomic_card_data
		SET ref_cnt = ?
		WHERE atomic_card_data_id = ?`)
    if err != nil {
        return &queries, err
    }

    queries.DeleteAtomicPropertiesQuery, err = db.Prepare(`DELETE FROM
        atomic_card_data
        WHERE
        atomic_card_data_id = ?`)
    if err != nil {
        return &queries, err
    }

	return &queries, nil
}

func (queries *atomicPropertiesDbQueries) cleanup() {
	if queries.NumAtomicPropertiesQuery != nil {
		queries.NumAtomicPropertiesQuery.Close()
	}

	if queries.AtomicPropertiesIdQuery != nil {
		queries.AtomicPropertiesIdQuery.Close()
	}

	if queries.AtomicPropertiesHashQuery != nil {
		queries.AtomicPropertiesHashQuery.Close()
	}

	if queries.InsertAtomicPropertiesQuery != nil {
		queries.InsertAtomicPropertiesQuery.Close()
	}

    if queries.GetRefCntQuery != nil {
        queries.GetRefCntQuery.Close()
    }

	if queries.UpdateRefCntQuery != nil {
		queries.UpdateRefCntQuery.Close()
	}

    if queries.DeleteAtomicPropertiesQuery != nil {
        queries.DeleteAtomicPropertiesQuery.Close()
    }
}

func getAtomicPropertiesId(
        numAtomicPropertiesQuery *sql.Stmt,
        atomicPropertiesIdQuery *sql.Stmt,
        atomicPropertiesHash string,
        scryfallOracleId string) (int64, int64, bool, error) {
	// First, check how many entries are already in the db with this card hash
	// If it's 0, this atomic data isn't in the db, so we can return without getting the id
	// If it's 1, we can just return the retrieved ID
	// If it's more than 1, we have a hash collision, so we use the scryfall_oracle_id to disambiguate

	var count int
	countResult := numAtomicPropertiesQuery.QueryRow(atomicPropertiesHash)
	if err := countResult.Scan(&count); err != nil {
		return 0, 0, false, err
	}

	if count == 0 {
		return 0, 0, false, nil
	}

	// Since count is at least 1, we need to query the actual ID
	var atomicPropertiesId int64
    var refCnt int64
	var rowScryfallOracleId string
	if count == 1 {
		// Only need to query the Id
		idResult := atomicPropertiesIdQuery.QueryRow(atomicPropertiesHash)
		if err := idResult.Scan(&atomicPropertiesId, &refCnt, &rowScryfallOracleId); err != nil {
			return 0, 0, false, err
		}
		return atomicPropertiesId, refCnt, true, nil
	} else {
		// Hash collision, so need to iterate and check the scryfall_oracle_id
		results, err := atomicPropertiesIdQuery.Query(atomicPropertiesHash)
		if err != nil {
			return 0, 0, false, err
		}
		defer results.Close()
		for results.Next() {
			if err := results.Err(); err != nil {
				return 0, 0, false, err
			}
			if err := results.Scan(&atomicPropertiesId, &refCnt, &rowScryfallOracleId); err != nil {
				return 0, 0, false, err
			}
			if scryfallOracleId == rowScryfallOracleId {
				return atomicPropertiesId, refCnt, true, nil
			}
		}

		// We shouldn't get here, since it means there are multiple entries with the correct
		// hash, but none that match the scryfall_oracle_id, so return an error
		return 0, 0, false, fmt.Errorf("Multiple atomic data with proper hash, but no matches")
	}
}

func insertAtomicPropertiesRecord(
        insertAtomicPropertiesQuery *sql.Stmt,
        req *atomicPropRequest) (int64, error) {
    // Build the set values needed for color_identity, color_indicator, and colors
    var colorIdentity sql.NullString
    var colorIndicator sql.NullString
    var colors sql.NullString
    var edhrecRank sql.NullInt32
    var hand sql.NullString
    var life sql.NullString
    var loyalty sql.NullString
    var name sql.NullString
    var side sql.NullString

    if len(req.ColorIdentity) > 0 {
        colorIdentity.String = strings.Join(req.ColorIdentity, ",")
        colorIdentity.Valid = true
    }

    if len(req.ColorIndicator) > 0 {
        colorIndicator.String = strings.Join(req.ColorIndicator, ",")
        colorIndicator.Valid = true
    }

    if len(req.Colors) > 0 {
        colors.String = strings.Join(req.Colors, ",")
        colors.Valid = true
    }

    if req.EDHRecRank != 0 {
        edhrecRank.Int32 = int32(req.EDHRecRank)
        edhrecRank.Valid = true
    }

    if len(req.Hand) > 0 {
        hand.String = req.Hand
        hand.Valid = true
    }

    if len(req.Life) > 0 {
        life.String = req.Life
        life.Valid = true
    }

    if len(req.Loyalty) > 0 {
        loyalty.String = req.Loyalty
        loyalty.Valid = true
    }

    if len(req.Name) > 0 {
        name.String = req.Name
        name.Valid = true
    }

    if len(req.Side) > 0 {
        side.String = req.Side
        side.Valid = true
    }

    res, err := insertAtomicPropertiesQuery.Exec(req.AtomicPropertiesHash,
        colorIdentity,
        colorIndicator,
        colors,
        req.ConvertedManaCost,
        edhrecRank,
        req.FaceConvertedManaCost,
        hand,
        req.IsReserved,
        req.Layout,
        life,
        loyalty,
        req.ManaCost,
        req.MTGStocksId,
        name,
        req.Power,
        req.ScryfallOracleId,
        side,
        req.Text,
        req.Toughness,
        req.Type)

    if err != nil {
        return 0, err
    }

    atomicPropId, err := res.LastInsertId()
    if err != nil {
        return 0, err
    }

    return atomicPropId, nil
}
