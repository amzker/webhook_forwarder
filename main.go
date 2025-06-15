package main

import (
	"log"
	"net/http"
	"webhook/Routes/EventSourcePool"
	"webhook/Routes/Webhooks"
)

func main() {
	http.HandleFunc("/webhook/{channel}", Webhooks.Webhook_forwarder)
	http.HandleFunc("/events/{channel}", EventSourcePool.EventStreamProvider)

	log.Println("Server listening on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
