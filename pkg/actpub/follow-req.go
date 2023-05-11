package actpub

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-fed/httpsig"
	"github.com/zjkmxy/fake-pub/pkg/config"
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
			Id:   config.UrlPrefix + "/a6d7d528-78af-4322-9cd7-04a20a27ab33",
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
	r.Header.Set("Accept", "application/activity+json")
	r.Header.Set("User-Agent", "Fake-ActivityPub/unknown")

	err = signer.SignRequest(privkey, keyUrl, r, wire)
	if err != nil {
		log.Fatal("Failed to sign the Post request: ", err)
	}

	// Replace hs2019 with rsa-sha256 for Misskey Mei 11 compatibility
	sigField := r.Header.Get("Signature")
	sigField = strings.Replace(sigField, "algorithm=\"hs2019\"", "algorithm=\"rsa-sha256\"", 1)
	r.Header.Set("Signature", sigField)

	return r
}
