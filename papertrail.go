package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"
)

type webHookData struct {
	Events []eventData `json:"events"`
}

type eventData struct {
	ID                int64     `json:"id"`
	SourceIP          string    `json:"source_ip"`
	Program           string    `json:"program"`
	Message           string    `json:"message"`
	ReceivedAt        time.Time `json:"received_at"`
	DisplayReceivedAt string    `json:"display_received_at"`
	SourceName        string    `json:"source_name"`
	HostName          string    `json:"hostname"`
	Severity          string    `json:"severity"`
	Facility          string    `json:"facility"`
}

// https://help.papertrailapp.com/kb/how-it-works/web-hooks#parsing
// it parses incoming message, generate notification and sends it
func paperTrailPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	query := r.URL.Query()
	message := query.Get("message")
	payload := r.FormValue("payload")
	payload, err := url.QueryUnescape(payload)
	if err != nil {
		log.Printf("Failed to unesacape url: %v ", err)
		return
	}

	data := []byte(payload)

	var webhookData webHookData
	err = json.Unmarshal(data, &webhookData)
	if err != nil {
		log.Printf("Failed to parse json: %v", err)
		return
	}

	var result string
	if message != "" {
		result = message
	} else {
		result = ""
		for _, entry := range webhookData.Events {
			result += fmt.Sprintf("```\nHost:%s, Program:%s\n%s \n```\n", entry.HostName, entry.Program, entry.Message)
		}
	}

	sendTelegramMessage(result)

}
