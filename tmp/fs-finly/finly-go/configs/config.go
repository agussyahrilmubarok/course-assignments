package configs

import (
	"log"
	"os"
)

func GetAppPort() string {
	port := os.Getenv("APP_PORT")
	if port == "" {
		log.Fatal("APP_PORT is required but not set")
	}
	return port
}

func GetMongoURI() string {
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatal("MONGODB_URI is required but not set")
	}
	return uri
}

func GetMongoDBName() string {
	uri := os.Getenv("MONGODB_NAME")
	if uri == "" {
		log.Fatal("MONGODB_NAME is required but not set")
	}
	return uri
}

func GetSecretKey() string {
	key := os.Getenv("AUTH_SECRET")
	if key == "" {
		log.Fatal("AUTH_SECRET is required but not set")
	}
	return key
}
