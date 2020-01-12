package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func serveDefaultGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello"))

	sendTelegramMessage("hello")
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", serveDefaultGet).Methods(http.MethodGet)
	r.HandleFunc("/papertrail", paperTrailPost).Methods(http.MethodPost)

	go startTelegramBot()
	log.Println("Application is started")
	log.Fatal(http.ListenAndServe(":8080", r))
}
