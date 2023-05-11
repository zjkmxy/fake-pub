package actpub

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/zjkmxy/fake-pub/pkg/config"
)

type WellknownNodeinfo struct {
	Links []WebfingerLink `json:"links"`
}

type NodeInfoMeta struct {
	NodeName        string            `json:"nodeName"`
	NodeDescription string            `json:"nodeDescription"`
	Maintainer      map[string]string `json:"maintainer"`
}

type NodeInfoV2 struct {
	Version           string              `json:"version"`
	Software          map[string]string   `json:"software"`
	Protocols         []string            `json:"protocols"`
	Services          map[string][]string `json:"services"`
	Usage             map[string]any      `json:"usage"`
	OpenRegistrations bool                `json:"openRegistrations"`
	Metadata          NodeInfoMeta        `json:"metadata"`
}

const JsonNodeInfo = `
{"version":"2.0",
"software":{"name":"fake-activitypub","version":"unknown"},
"protocols":["activitypub"],"services":{"outbound":[],"inbound":[]},
"usage":{},"openRegistrations":false,"metadata":{
"nodeName":"kinu-fake-pub.duckdns.org"}}`

func WellknownNodeinfoHandler(w http.ResponseWriter, r *http.Request) {
	wfData := &WellknownNodeinfo{
		Links: []WebfingerLink{
			{
				Rel:  "http://nodeinfo.diaspora.software/ns/schema/2.0",
				Href: fmt.Sprintf("%s/nodeinfo/2.0", config.UrlPrefix),
			},
		},
	}
	data, err := json.Marshal(wfData)
	if err != nil {
		log.Print("ERROR: ", "Unable to marshal nodeinfo: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

func NodeinfoHandler(w http.ResponseWriter, r *http.Request) {
	nodeInfo := &NodeInfoV2{
		Version: "2.0",
		Software: map[string]string{
			"name":    "fake-activitypub",
			"version": "unknown",
		},
		Protocols: []string{"activitypub"},
		Services: map[string][]string{
			"outbound": {},
			"inbound":  {},
		},
		Usage: map[string]any{
			"users": map[string]string{},
		},
		OpenRegistrations: false,
		Metadata: NodeInfoMeta{
			NodeName: config.HostDomain,
			Maintainer: map[string]string{
				"name":    "kinu",
				"misskey": "@kinu@misskey.dev",
			},
			NodeDescription: "A dummy server set to try out the ActivityPub protocol. " +
				"All users are robots on this server. All keys & notes are fake. " +
				"If any malicious/suspecious behavior is found, contact the maintainer *immediately*.",
		},
	}
	data, err := json.Marshal(nodeInfo)
	if err != nil {
		log.Print("ERROR: ", "Unable to marshal nodeinfo: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}
