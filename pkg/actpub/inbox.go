package actpub

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/zjkmxy/fake-pub/pkg/config"
)

func InboxHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("Inbox received: ", r.Method, " ", r.URL.String(), " ", r.Header)

	urlPaths := strings.Split(r.URL.Path, "/")
	if len(urlPaths) < 2 {
		log.Fatal("Wrong URL redirected", r.URL)
	}
	acctName := urlPaths[len(urlPaths)-1]

	objId := time.Now().Format(TimeFmtStr)
	objUrl := fmt.Sprintf(`%s/inbox/%s/%s`, config.UrlPrefix, acctName, objId)

	log.Print("CREATED: ", objUrl, "\n")
	body, _ := io.ReadAll(r.Body)
	fmt.Print(string(body), "\n")

	w.Header().Add("Location", objUrl)
	w.WriteHeader(http.StatusCreated)
}
