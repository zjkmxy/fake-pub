package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	HostDomain string
	Proto      string
	UrlPrefix  string
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	HostDomain = os.Getenv("HOST_DOMAIN")
	Proto = os.Getenv("PROTO")
	UrlPrefix = Proto + "://" + HostDomain
}
