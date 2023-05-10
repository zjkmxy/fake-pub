package actpub

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/zjkmxy/fake-pub/pkg/config"
)

type WebfingerLink struct {
	Rel      string `json:"rel"`
	Type     string `json:"type,omitempty"`
	Href     string `json:"href,omitempty"`
	Template string `json:"template,omitempty"`
}

type Webfinger struct {
	Subject string          `json:"subject"`
	Links   []WebfingerLink `json:"links"`
}

func WebfingerHandler(w http.ResponseWriter, r *http.Request) {
	acct, err := url.QueryUnescape(r.URL.Query().Get("resource"))
	if err != nil || !strings.HasPrefix(acct, "acct:") || !strings.HasSuffix(acct, "@"+config.HostDomain) {
		log.Print("WARNING: ", "Bad webfinger request: ", acct)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	acctName := strings.TrimSuffix(strings.TrimPrefix(acct, "acct:"), "@"+config.HostDomain)
	wfData := &Webfinger{
		Subject: acct,
		Links: []WebfingerLink{
			{
				Rel:  "self",
				Type: "application/activity+json",
				Href: fmt.Sprintf("%s/actor/%s", config.UrlPrefix, acctName),
			},
			{
				Rel:  "http://webfinger.net/rel/profile-page",
				Type: "text/html",
				Href: fmt.Sprintf("%s/profile-page/%s", config.UrlPrefix, acctName),
			},
		},
	}
	data, err := json.Marshal(wfData)
	if err != nil {
		log.Print("ERROR: ", "Unable to marshal webfinger: ", acct)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}
