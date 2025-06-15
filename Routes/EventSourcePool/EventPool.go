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
	channel_to_listen := request.PathValue("channel")
	log.Printf("Channel to listen: %s", channel_to_listen)
	if uuid.Validate(channel_to_listen) != nil {
		log.Printf("Invalid UUID %s", channel_to_listen)
		http.Error(writer, "Invalid UUID", http.StatusBadRequest) // in stream? , i will need to look into it.
		return
	}

	headers := map[string]string{
		"Content-Type":                "text/event-stream",
		"Cache-Control":               "no-cache",
		"Connection":                  "keep-alive",
		"Access-Control-Allow-Origin": "*",
	}
	for key, value := range headers {
		writer.Header().Set(key, value)
	}

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
			if msg != nil && msg.Payload != "" {
				fmt.Fprintf(writer, "_data_:%s\n\n", msg.Payload)
				flusher.Flush()
			}
		case <-time.After(time.Second * 600):
			log.Printf("Close connection as 600 Secounds has been passed")
			http.Error(writer, "Timeout: Close connection as 600 Secounds has been passed", http.StatusRequestTimeout)
			return
		case <-connection_close.Done():
			log.Print("Client disconnected from stream")
			return
		}
	}
}
