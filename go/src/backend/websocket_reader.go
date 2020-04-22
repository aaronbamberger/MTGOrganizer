package backend

import "net/http"
import "io/ioutil"
import "encoding/json"
import "log"
import "database/sql"
import _ "github.com/go-sql-driver/mysql"
import "github.com/gorilla/websocket"

func checkOrigin(req *http.Request) bool {
    origin := req.Header["Origin"]
    log.Printf("Websocket connection from origin: %s\n", origin)

    return true
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
    cardDB, err := sql.Open("mysql",
        "app_user:app_db_password@tcp(172.18.0.5)/mtg_cards?parseTime=true")
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
                case LoginChallengeCheck:
                    go checkLoginChallenge(message.Value, doneChan, respChan)
                case CardSearchRequest:
                    go cardSearch(cardDB, message.Value, doneChan, respChan)
                case CardDetailRequest:
                    go cardDetail(cardDB, message.Value, doneChan, respChan)
                }

            default:
                log.Printf("Received an unexpected message type: %d", messageType)
            }
        }
    }
}
