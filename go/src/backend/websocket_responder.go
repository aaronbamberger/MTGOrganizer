package backend

import "log"
import "github.com/gorilla/websocket"

func websocketResponder(conn *websocket.Conn,
        done <-chan interface{},
        responses <-chan ResponseMessage) {
    processing := true

    for processing {
        select {
        case <-done:
            processing = false
        case response := <-responses:
            err := conn.WriteJSON(response)
            if err != nil {
                log.Print(err)
            }
        }
    }
}
