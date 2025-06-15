package RedisManager

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"webhook/Config"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func init() {
	redis_db, err := strconv.Atoi(Config.REDIS_DB)
	if err != nil {
		log.Printf("Redis Db is not int Using Default DB 0 %s", err)
	}
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprint(Config.REDIS_HOST, ":", Config.REDIS_PORT),
		Password: Config.REDIS_PASSWORD,
		DB:       redis_db,
		Protocol: 2,
		PoolSize: 500,
	})
}

func GetRedisClient() *redis.Client {
	return RedisClient
}

type RedisPublishMessage struct {
	Method  string              `json:"method"`
	URL     string              `json:"url"`
	Headers map[string][]string `json:"headers"`
	Body    string              `json:"body,omitempty"`
}

func PushtoPubsubChannel(channel string, request *http.Request) bool {
	log.Printf("Got Requests to push to %s", channel)
	ctx := context.Background()
	rerr := RedisClient.Ping(ctx).Err()
	if rerr != nil {
		log.Printf("redis connection error %s", rerr)
		return false
	}

	bodybytes, berr := io.ReadAll(request.Body)
	if berr != nil {
		log.Printf("Failed to read body,%s", berr)
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
		log.Printf("json object was not able to be created %s", jerr)
	}

	perr := RedisClient.Publish(ctx, channel, jsonObj).Err()
	if perr != nil {
		log.Printf("Error during Publishing to redis %s", perr)
		return false
	}
	return true
}

func SubscribedToPubsubChannel(channel string) *redis.PubSub {
	ctx := context.Background()
	subscribed := RedisClient.Subscribe(ctx, channel)
	return subscribed
}
