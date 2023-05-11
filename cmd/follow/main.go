package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/zjkmxy/fake-pub/pkg/actpub"
)

func main() {
	if len(os.Args) < 4 {
		log.Fatal("Usage: ", os.Args[0], " <inbox-url>", " <follower-url>", " <followee-url>")
	}

	r := actpub.MakeFollowReq(os.Args[1], os.Args[2], os.Args[3])
	r.Write(os.Stdout)
	fmt.Println()

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
