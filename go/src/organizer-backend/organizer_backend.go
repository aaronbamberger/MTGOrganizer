package main

import "carddb"
import "net/http"
import "log"
import "database/sql"
import _ "github.com/go-sql-driver/mysql"
import "github.com/gorilla/websocket"

type RequestType int

const (
    RequestCardSearch RequestType = iota
)

type TestMessage struct {
    Type RequestType `json:"type"`
    Value string `json:"value"`
}

func main() {
    http.HandleFunc("/api/test", handleTestApi)
    http.HandleFunc("/api", handleApi)
    log.Printf("Starting listener...\n")
    http.ListenAndServe("192.168.50.185:8085", nil)
}

func checkOrigin(req *http.Request) bool {
    origin := req.Header["Origin"]
    log.Printf("Websocket connection from origin: %s\n", origin)

    return true
}

func handleApi(resp http.ResponseWriter, req *http.Request) {
    log.Printf("Accepted connection...\n")

	// Connect to the mariadb database
	cardDB, err := sql.Open("mysql",
        "app_user:app_db_password@tcp(172.18.0.8)/mtg_cards?parseTime=true")
	if err != nil {
		log.Print(err)
        return
	}
	defer cardDB.Close()
	cardDB.SetMaxIdleConns(10)

    dbContext, err := carddb.NewDBFetchContext(cardDB)
    if err != nil {
        log.Print(err)
        return
    }

    upgrader := websocket.Upgrader{CheckOrigin: checkOrigin}
    conn, err := upgrader.Upgrade(resp, req, nil)
    if err != nil {
        log.Print(err)
        return
    }
    defer conn.Close()

    for {
        var message TestMessage
        err = conn.ReadJSON(&message)
        if err != nil {
            log.Print(err)
        } else {
            switch message.Type {
            case RequestCardSearch:
                cards, err := dbContext.SearchCardsByName(message.Value)
                if err != nil {
                    log.Print(err)
                    continue
                }
                err = conn.WriteJSON(cards)
                if err != nil {
                    log.Print(err)
                }

            default:
                log.Printf("Unknown request type: %d\n", message.Type)
            }
        }

    }
}

func handleTestApi(resp http.ResponseWriter, req *http.Request) {
    /*
    messageTypes := map[int]string{
        websocket.TextMessage: "text",
        websocket.BinaryMessage: "binary",
        websocket.CloseMessage: "close",
        websocket.PingMessage: "ping",
        websocket.PongMessage: "pong"}
    */

    log.Printf("Accepted connection...\n")
    upgrader := websocket.Upgrader{CheckOrigin: checkOrigin}
    conn, err := upgrader.Upgrade(resp, req, nil)
    if err != nil {
        log.Fatal(err)
    }

    for {
        var message TestMessage
        err = conn.ReadJSON(&message)
        if err != nil {
            log.Print(err)
        } else {
            log.Printf("Message received with type %d: %s\n", message.Type, message.Value)
        }

        /*
        if messageType, reader, err := conn.NextReader(); err != nil {
            log.Print(err)
            conn.Close()
            break
        } else {
            log.Printf("Received message with type %s\n", messageTypes[messageType])
            message, err := ioutil.ReadAll(reader)
            if err != nil {
                log.Print(err)
            } else {
                log.Printf("Message: %s\n", string(message))
            }
        }
        */
    }
}
