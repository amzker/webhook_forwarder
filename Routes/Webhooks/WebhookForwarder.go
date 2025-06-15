package Webhooks

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"webhook/Database/RedisManager"

	"github.com/google/uuid"
)

type RedisPublishMessage struct {
	Method  string              `json:"method"`
	URL     string              `json:"url"`
	Headers map[string][]string `json:"headers"`
	Body    string              `json:"body,omitempty"`
}

type WebHookResponseStruct struct {
	Message    string `json:"message"`
	Success    bool   `json:"success"`
	Channel    string `json:"channel"`
	LisnerUrl  string `json:"lisner_url,nullable"`
	StatusCode int    `json:"status_code"`
}

func Webhook_forwarder(ResponseWriter http.ResponseWriter, request *http.Request) {
	channel_to_publish := request.PathValue("channel")
	ResponseWriter.Header().Set("Content-Type", "application/json")
	if uuid.Validate(channel_to_publish) != nil {
		// http.Error(ResponseWriter, "In valid UUID", http.StatusBadRequest)
		rspreturn := WebHookResponseStruct{
			Message:    "Invalid UUID",
			Success:    false,
			Channel:    "",
			LisnerUrl:  "",
			StatusCode: http.StatusBadRequest,
		}
		json.NewEncoder(ResponseWriter).Encode(rspreturn)
		log.Printf("Invalid UUID %s", channel_to_publish)
		return
	}
	bodybytes, berr := io.ReadAll(request.Body)
	if berr != nil {
		log.Printf("Failed to read body,%s", berr)
		// http.Error(ResponseWriter, "Failed to read body", http.StatusBadRequest)
		rspreturn := WebHookResponseStruct{
			Message:    "Failed to read body",
			Success:    false,
			Channel:    "",
			LisnerUrl:  "",
			StatusCode: http.StatusBadRequest,
		}
		ResponseWriter.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(ResponseWriter).Encode(rspreturn)
		return
	}
	messageObj := RedisPublishMessage{
		Method:  request.Method,
		URL:     request.RequestURI,
		Headers: request.Header,
		Body:    string(bodybytes),
	}
	request.Body.Close()
	jsonObj, jerr := json.Marshal(messageObj)
	if jerr != nil {
		// i feel like this is redundant as obj is struct
		log.Printf("json object was not able to be created %s", jerr)
		// http.Error(ResponseWriter, "Non Encodable Request Object", http.StatusBadRequest)
		rspreturn := WebHookResponseStruct{
			Message:    "Non Encodable Request Object",
			Success:    false,
			Channel:    "",
			LisnerUrl:  "",
			StatusCode: http.StatusBadRequest,
		}
		ResponseWriter.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(ResponseWriter).Encode(rspreturn)
		return
	}

	success := RedisManager.PushtoPubsubChannel(channel_to_publish, jsonObj)
	if success {
		log.Printf("Successfully Got requests from %s and pushed to redis", channel_to_publish)
		rspreturn := WebHookResponseStruct{
			Message:    "Successfully Got requests from " + channel_to_publish + " and Forwarded to Event Stream",
			Success:    true,
			Channel:    channel_to_publish,
			LisnerUrl:  "/events/" + channel_to_publish,
			StatusCode: http.StatusOK,
		}
		ResponseWriter.WriteHeader(http.StatusOK)
		json.NewEncoder(ResponseWriter).Encode(rspreturn)
		return
	} else {
		log.Printf("Failed to Create request forward for %s", channel_to_publish)
		// http.Error(ResponseWriter, "Failed to Create request forward", http.StatusInternalServerError)
		rspreturn := WebHookResponseStruct{
			Message:    "Failed to Create request forward",
			Success:    false,
			Channel:    "",
			LisnerUrl:  "",
			StatusCode: http.StatusInternalServerError,
		}
		ResponseWriter.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(ResponseWriter).Encode(rspreturn)
		return
	}
	// hehe , pretty sure this is not the way to do it.
	// http.Error(ResponseWriter, "Success", http.StatusOK)
	// return
}
