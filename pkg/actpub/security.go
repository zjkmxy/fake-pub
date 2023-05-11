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

type AsPublicKey struct {
	KeyId        string `json:"id"`
	Type         string `json:"type,omitempty"`
	Owner        string `json:"owner"`
	PublicKeyPem string `json:"publicKeyPem"`
}

func SiteActorHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("Site-Actor request received: ", r.Method, " ", r.URL.String(), " ", r.Header)

	urlPaths := strings.Split(r.URL.Path, "/")
	if len(urlPaths) < 1 {
		log.Fatal("Wrong URL redirected", r.URL)
	}
	acctUrl := fmt.Sprintf("%s/site-actor", config.UrlPrefix)

	pubkeyPem, err := os.ReadFile("data/public.pem")
	if err != nil {
		log.Print("ERROR: ", "Unable to read data/public.pem: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	actorData := &ActorPerson{
		AsObject: AsObject{
			Type: "Application",
			Id:   acctUrl,
		},
		PublicKey: &AsPublicKey{
			KeyId:        acctUrl + "#main-key",
			Type:         "Key",
			Owner:        acctUrl,
			PublicKeyPem: string(pubkeyPem),
		},
		ManualApproval: func(v bool) *bool { return &v }(true),
		Username:       config.HostDomain,
		Endpoints: map[string]string{
			"sharedInbox": fmt.Sprintf("%s/inbox/shared", config.UrlPrefix),
		},
		Inbox:  fmt.Sprintf("%s/inbox/shared", config.UrlPrefix),
		Outbox: fmt.Sprintf("%s/outbox/shared", config.UrlPrefix),
	}
	for _, context := range PersonContext {
		actorData.AtContext = append(actorData.AtContext, context)
	}

	data, err := json.Marshal(actorData)
	if err != nil {
		log.Print("ERROR: ", "Unable to marshal Site-Actor:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Print("Responded with: ", string(data))

	w.Write(data)
}
