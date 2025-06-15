package EventSourcePool

import (
	"fmt"
	"log"
	"net/http"
	"time"
	"webhook/Database/RedisManager"

	"github.com/google/uuid"
)

func EventStreamProvider(writer http.ResponseWriter, request *http.Request) {
	log.Print("Event Steam Request arrived")
	channel_to_listen := request.PathValue("channel")
	if uuid.Validate(channel_to_listen) != nil {
		log.Printf("Invalid UUID %s", channel_to_listen)
		return
	}
	writer.Header().Set("Content-Type", "text/event-stream")
	writer.Header().Set("Cache-Control", "no-cache")
	writer.Header().Set("Connection", "keep-alive")
	writer.Header().Set("Access-Control-Allow-Origin", "*")

	flusher, success := writer.(http.Flusher)
	if !success {
		http.Error(writer, "Streaming not supported", http.StatusInternalServerError)
		return
	}
	redis_event_stream := RedisManager.SubscribedToPubsubChannel(channel_to_listen)
	defer redis_event_stream.Close()

	connection_close := request.Context()
	events_channel := redis_event_stream.Channel()

	fmt.Fprint(writer, "ping\n\n")
	flusher.Flush()

	for {
		select {
		case msg := <-events_channel:
			if msg != nil {
				fmt.Fprintf(writer, "data: %s\n\n", msg.Payload)
				flusher.Flush()
			}
		case <-time.After(time.Second * 30):
			log.Printf("Close connection as 30 Secounds has been passed")
			return
		case <-connection_close.Done():
			log.Print("Client disconnected from stream")
			return
		}
	}
}
