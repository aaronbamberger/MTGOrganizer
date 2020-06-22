package main

import "backend"
import "net/http"
import "log"

func main() {
    http.HandleFunc("/backend/api", backend.HandleApi)
    http.HandleFunc("/backend/login/challenge", backend.HandleLoginChallenge)
    http.HandleFunc("/backend/login/creds", backend.HandleLoginCredentials)
    http.HandleFunc("/backend/consent/challenge", backend.HandleConsentChallenge)
    log.Printf("Starting listener...\n")
    http.ListenAndServe(":8085", nil)
    log.Printf("Exiting backend...\n")
}

