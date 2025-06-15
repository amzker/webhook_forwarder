package Config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var REDIS_HOST string
var REDIS_PORT string
var REDIS_PASSWORD string
var REDIS_DB string

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Failed to load .env file")
	}

	REDIS_HOST = os.Getenv("REDIS_HOST")
	REDIS_PORT = os.Getenv("REDIS_PORT")
	REDIS_PASSWORD = os.Getenv("REDIS_PASSWORD")
	REDIS_DB = os.Getenv("REDIS_DB")
}
