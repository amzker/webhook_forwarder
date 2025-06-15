package Webhooks

import (
	"fmt"
	"log"
	"net/http"
	"webhook/Database/RedisManager"

	"github.com/google/uuid"
)

func Webhook_forwarder(ResponseWriter http.ResponseWriter, request *http.Request) {
	log.Printf("webhook request came")
	channel_to_publish := request.PathValue("webhook_path")
	if uuid.Validate(channel_to_publish) != nil {
		fmt.Fprint(ResponseWriter, "In valid UUID")
		return
	}
	success := RedisManager.PushtoPubsubChannel(channel_to_publish, request)
	if success {
		log.Printf("Successfully Got requests from %s and pushed to redis", channel_to_publish)
	} else {
		log.Panicf("Failed to Create request forward for %s", channel_to_publish)
	}
	return
}
