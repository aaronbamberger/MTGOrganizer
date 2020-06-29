package backend

import "net/http"
import "io/ioutil"
import "encoding/json"
import "fmt"
import "log"
import "database/sql"
import _ "github.com/go-sql-driver/mysql"
import "github.com/gorilla/websocket"

const (
    DB_HOST = "card_db:3306"
    CARD_DB = "mtg_cards"
    USER_DB = "users"
    APP_DB_USER = "app_user"
    APP_DB_PW = "app_db_password"
    LOGIN_DB_USER = "login_user"
    LOGIN_DB_PW = "login_user_password"
)

func checkOrigin(req *http.Request) bool {
    origin := req.Header["Origin"]
    log.Printf("Websocket connection from origin: %s\n", origin)

    return true
}

func dbConnStr(user string, pw string, db string) string {
    return fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true",
            user, pw, DB_HOST, db)
}

func HandleApi(resp http.ResponseWriter, req *http.Request) {
    log.Printf("Accepted connection from %s...\n", req.RemoteAddr)

    upgrader := websocket.Upgrader{CheckOrigin: checkOrigin}
    conn, err := upgrader.Upgrade(resp, req, nil)
    if err != nil {
        log.Print(err)
        return
    }
    defer conn.Close()

    log.Printf("Successfully upgraded accepted connection from %s to websocket\n",
        req.RemoteAddr)

    // Connect to the mariadb database
    cardDB, err := sql.Open("mysql", dbConnStr(APP_DB_USER, APP_DB_PW, CARD_DB))
	if err != nil {
		log.Print(err)
        return
	}

	defer cardDB.Close()
	cardDB.SetMaxIdleConns(10)

    doneChan := make(chan interface{})

    // Start a goroutine to handle writing responses back to the websocket
    respChan := make(chan ResponseMessage)
    go websocketResponder(conn, doneChan, respChan)

    done := false
    for !done {
        if messageType, reader, err := conn.NextReader(); err != nil {
            log.Print(err)
            close(doneChan)
            conn.Close()
            done = true
        } else {
            switch messageType {
            case websocket.CloseMessage:
                close(doneChan)
                done = true
            case websocket.TextMessage:
                rawMessage, err := ioutil.ReadAll(reader)
                if err != nil {
                    log.Print(err)
                    continue
                }
                var message RequestMessage
                err = json.Unmarshal([]byte(rawMessage), &message)
                if err != nil {
                    log.Print(err)
                    continue
                }
                switch message.Type {
                case ApiTypesRequest:
                    go apiTypes(doneChan, respChan)
                case CardSearchRequest:
                    var searchName string
                    err = json.Unmarshal([]byte(message.Value), &searchName)
                    if err != nil {
                        log.Print(err)
                        continue
                    }
                    go cardSearch(cardDB, searchName, doneChan, respChan)
                case CardDetailRequest:
                    var cardUUID string
                    err = json.Unmarshal([]byte(message.Value), &cardUUID)
                    if err != nil {
                        log.Print(err)
                        continue
                    }
                    go cardDetail(cardDB, cardUUID, doneChan, respChan)
                }

            default:
                log.Printf("Received an unexpected message type: %d", messageType)
            }
        }
    }
}
