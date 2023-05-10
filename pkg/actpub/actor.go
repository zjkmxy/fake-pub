package actpub

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/zjkmxy/fake-pub/pkg/config"
)

var PersonContext []string = []string{
	"https://www.w3.org/ns/activitystreams",
	"https://w3id.org/security/v1",
}

type ActorPerson struct {
	AsObject

	Inbox     string            `json:"inbox"`
	Outbox    string            `json:"outbox"`
	Followers string            `json:"followers,omitempty"`
	Following string            `json:"following,omitempty"`
	Username  string            `json:"preferredUsername,omitempty"`
	Endpoints map[string]string `json:"endpoints,omitempty"`

	ManualApproval *bool `json:"manuallyApprovesFollowers,omitempty"`

	PublicKey *AsPublicKey `json:"publicKey,omitempty"`
}

func ActorHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("Actor request received: ", r.Method, " ", r.URL.String(), " ", r.Header)

	urlPaths := strings.Split(r.URL.Path, "/")
	if len(urlPaths) < 1 {
		log.Fatal("Wrong URL redirected", r.URL)
	}
	acctName := urlPaths[len(urlPaths)-1]
	acctUrl := fmt.Sprintf(`%s/actor/%s`, config.UrlPrefix, acctName)

	pubkeyPem, err := os.ReadFile("data/public.pem")
	if err != nil {
		log.Print("ERROR: ", "Unable to read data/public.pem: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	actorData := &ActorPerson{
		AsObject: AsObject{
			Type:    "Person",
			Id:      acctUrl,
			Tag:     []any{},
			Image:   "",
			Icon:    "",
			Summary: "This is a FAKE account.",
			Name:    acctName,
			URL:     fmt.Sprintf("%s/profile-page/%s", config.UrlPrefix, acctName),
		},
		PublicKey: &AsPublicKey{
			KeyId:        acctUrl + "#main-key",
			Type:         "Key",
			Owner:        acctUrl,
			PublicKeyPem: string(pubkeyPem),
		},
		ManualApproval: func(v bool) *bool { return &v }(false),
		Username:       acctName,
		Endpoints: map[string]string{
			"sharedInbox": fmt.Sprintf("%s/inbox/shared", config.UrlPrefix),
		},
		Inbox:     fmt.Sprintf("%s/inbox/%s", config.UrlPrefix, acctName),
		Outbox:    fmt.Sprintf("%s/outbox/%s", config.UrlPrefix, acctName),
		Followers: fmt.Sprintf("%s/followers/%s", config.UrlPrefix, acctName),
		Following: fmt.Sprintf("%s/following/%s", config.UrlPrefix, acctName),
	}

	data, err := json.Marshal(actorData)
	if err != nil {
		log.Print("ERROR: ", "Unable to marshal Person actor:", acctName)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}
