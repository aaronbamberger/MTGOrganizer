package main

import "backend"
import "net/http"
import "log"

func main() {
    http.HandleFunc("/api", backend.HandleApi)
    log.Printf("Starting listener...\n")
    http.ListenAndServe("192.168.50.185:8085", nil)
}

