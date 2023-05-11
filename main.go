package main

import (
	"io"
	"log"
	"net/http"

	"github.com/zjkmxy/fake-pub/pkg/actpub"
)

func FallbackHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("Unknown request received: ", r.Method, " ", r.URL.String(), " ", r.Header)
	body, _ := io.ReadAll(r.Body)
	log.Print(string(body))
	w.WriteHeader(http.StatusNotFound)
}

func main() {
	http.HandleFunc("/.well-known/webfinger", actpub.WebfingerHandler)
	http.HandleFunc("/.well-known/nodeinfo", actpub.WellknownNodeinfoHandler)
	http.HandleFunc("/nodeinfo/2.0", actpub.NodeinfoHandler)
	http.HandleFunc("/actor/", actpub.ActorHandler)
	http.HandleFunc("/site-actor", actpub.SiteActorHandler)
	http.HandleFunc("/followers/", actpub.FollowHandler)
	http.HandleFunc("/following/", actpub.FollowHandler)
	http.HandleFunc("/inbox/", actpub.InboxHandler)
	http.HandleFunc("/", FallbackHandler)

	log.Fatal(http.ListenAndServe(":31000", nil))
}
