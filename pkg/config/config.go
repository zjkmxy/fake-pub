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
	MailTarget string
	MailFrom   string
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	HostDomain = os.Getenv("HOST_DOMAIN")
	Proto = os.Getenv("PROTO")
	UrlPrefix = Proto + "://" + HostDomain
	MailTarget = os.Getenv("MAIL_TARGET")
	MailFrom = os.Getenv("MAIL_FROM")
}
