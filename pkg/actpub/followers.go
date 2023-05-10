package actpub

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/zjkmxy/fake-pub/pkg/config"
)

type AsOrderedCollection struct {
	AsObject

	TotalItems   int      `json:"totalItems"`
	OrderedItems []string `json:"orderedItems,omitempty"`
}

func FollowHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("Follower/ee request received: ", r.Method, " ", r.URL.String(), " ", r.Header)

	urlPaths := strings.Split(r.URL.Path, "/")
	if len(urlPaths) < 2 {
		log.Fatal("Wrong URL redirected", r.URL)
	}
	acctName := urlPaths[len(urlPaths)-1]
	collectionName := urlPaths[len(urlPaths)-2]
	collectionUrl := fmt.Sprintf(`%s/%s/%s`, config.UrlPrefix, collectionName, acctName)

	collectionData := &AsOrderedCollection{
		AsObject: AsObject{
			Type: "OrderedCollection",
			Id:   collectionUrl,
		},
		TotalItems:   0,
		OrderedItems: []string{},
	}
	for _, context := range PersonContext {
		collectionData.AtContext = append(collectionData.AtContext, context)
	}

	data, err := json.Marshal(collectionData)
	if err != nil {
		log.Print("ERROR: ", "Unable to marshal Collection: ", acctName)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}
