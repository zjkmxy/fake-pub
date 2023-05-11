package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-fed/httpsig"
	"github.com/zjkmxy/fake-pub/pkg/config"
	"golang.org/x/crypto/ssh"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: ", os.Args[0], " <actor-url>")
	}

	var err error
	signer, _, err := httpsig.NewSigner(
		[]httpsig.Algorithm{httpsig.RSA_SHA256},
		httpsig.DigestSha256,
		[]string{httpsig.RequestTarget, "date", "host", "accept"},
		httpsig.Signature,
		0,
	)
	if err != nil {
		log.Fatal("Failed to create HTTP signer: ", err)
	}

	privkeyPem, err := os.ReadFile("data/private.pem")
	if err != nil {
		log.Fatal("Unable to read data/private.pem: ", err)
	}
	privkey, err := ssh.ParseRawPrivateKey(privkeyPem)
	if err != nil {
		log.Fatal("Unable to parse the pem key: ", err)
	}

	keyUrl := config.UrlPrefix + "/site-actor#main-key"
	r, err := http.NewRequest("GET", os.Args[1], nil)
	if err != nil {
		log.Fatal("Failed to create get request: ", err)
	}
	r.Header.Set("Accept", "application/activity+json")
	r.Header.Set("Date", time.Now().UTC().Format("Mon, 02 Jan 2006 15:04:05")+" GMT")
	r.Header.Set("Host", r.Host)

	err = signer.SignRequest(privkey, keyUrl, r, nil)
	if err != nil {
		log.Fatal("Failed to sign the Get request: ", err)
	}

	client := http.Client{Timeout: time.Duration(10) * time.Second}
	resp, err := client.Do(r)
	if err != nil {
		log.Fatal("Error doing HTTP request: ", err)
	}
	log.Print("Response: ", resp.Status, "\n", resp.Header, "\n")

	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("%s", body)
}
