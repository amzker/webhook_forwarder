package RedisManager

import (
	"context"
	"fmt"
	"log"
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

func PushtoPubsubChannel(channel string, message []byte) bool {
	log.Printf("Got Requests to push to %s", channel)
	ctx := context.Background()
	rerr := RedisClient.Ping(ctx).Err()
	if rerr != nil {
		log.Printf("redis connection error %s", rerr)
		return false
	}

	perr := RedisClient.Publish(ctx, channel, message).Err()
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
