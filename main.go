package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"
)

func serveDefaultGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello"))
}

func main() {

	config, err := loadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration. %s", err)
	}

	// monitor OS signal for exit
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		termCh := make(chan os.Signal, 1)
		signal.Notify(termCh, os.Interrupt, syscall.SIGINT)
		<-termCh
		log.Println("Shutdown...")
		cancel()
		os.Exit(0)
	}()

	r := mux.NewRouter()
	r.HandleFunc("/", serveDefaultGet).Methods(http.MethodGet)
	r.HandleFunc("/papertrail", paperTrailPost).Methods(http.MethodPost)

	go startTelegramBot(config, ctx)
	go startHealthcheck(config, ctx)
	log.Println("Application is started")
	log.Fatal(http.ListenAndServe(":8080", r))
}
