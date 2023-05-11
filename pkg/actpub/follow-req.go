package actpub

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-fed/httpsig"
	"golang.org/x/crypto/ssh"
)

type FollowActivity struct {
	AsObject

	Actor  string `json:"actor"`
	Object string `json:"object"`
}

type AcceptActivity struct {
	AsObject

	Actor  string         `json:"actor"`
	Object FollowActivity `json:"object"`
}

func MakeFollowReq(inboxUrl string, actorId string, objectId string) *http.Request {
	var err error
	signer, _, err := httpsig.NewSigner(
		[]httpsig.Algorithm{httpsig.RSA_SHA256},
		httpsig.DigestSha256,
		[]string{httpsig.RequestTarget, "date", "host", "digest"},
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

	act := FollowActivity{
		AsObject: AsObject{
			Type: "Follow",
		},
		Actor:  actorId,
		Object: objectId,
	}
	for _, context := range PersonContext {
		act.AtContext = append(act.AtContext, context)
	}
	wire, err := json.Marshal(act)
	if err != nil {
		log.Fatal("Failed to marshal Follow activity: ", err)
	}

	keyUrl := actorId + "#main-key"
	r, err := http.NewRequest("POST", inboxUrl, bytes.NewBuffer(wire))
	if err != nil {
		log.Fatal("Failed to create Post request: ", err)
	}
	r.Header.Set("Content-Type", "application/activity+json")
	r.Header.Set("Date", time.Now().UTC().Format("Mon, 02 Jan 2006 15:04:05")+" GMT")
	r.Header.Set("Host", r.Host)

	err = signer.SignRequest(privkey, keyUrl, r, wire)
	if err != nil {
		log.Fatal("Failed to sign the Post request: ", err)
	}

	return r
}
