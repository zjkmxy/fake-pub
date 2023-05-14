package actpub

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os/exec"
	"strings"
	"time"

	"github.com/zjkmxy/fake-pub/pkg/config"
)

func InboxHandler(w http.ResponseWriter, r *http.Request) {
	// log.Print("Inbox received: ", r.Method, " ", r.URL.String(), " ", r.Header)

	urlPaths := strings.Split(r.URL.Path, "/")
	if len(urlPaths) < 2 {
		log.Fatal("Wrong URL redirected", r.URL)
	}
	acctName := urlPaths[len(urlPaths)-1]

	objId := url.PathEscape(time.Now().Format(TimeFmtStr))
	objUrl := fmt.Sprintf(`%s/inbox/%s/%s`, config.UrlPrefix, acctName, objId)

	log.Print("INBOX CREATED: ", objUrl, "\n")
	body, _ := io.ReadAll(r.Body)
	// fmt.Print(string(body), "\n")

	// Respond with Created
	w.Header().Add("Location", objUrl)
	w.WriteHeader(http.StatusCreated)

	// Send mail if it is not DELETE
	obj := map[string]any{}
	err := json.Unmarshal(body, &obj)
	if err != nil {
		log.Print("ERROR: ", "Unable to unmarshal inbox object: ", body)
		return
	}

	if typ, ok := obj["type"]; !ok || typ == "Delete" {
		// Ignore Mastodon broadcast delete feeds
		return
	}
	pretty, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		log.Print("ERROR: ", "Unable to remarshal inbox object: ", body)
		return
	}

	// Send mail
	subjectLine := fmt.Sprintf("Subject: ActPub: %s/%s", acctName, objId)
	toLine := fmt.Sprintf("To: %s", config.MailTarget)
	fromLine := fmt.Sprintf("From: fake-pub <%s>", config.MailFrom)
	mail := toLine + "\n" + fromLine + "\n" + subjectLine + "\n\n" + string(pretty) + "\n"

	cmd := exec.Command("sendmail", config.MailTarget)
	reader, writer := io.Pipe()
	cmd.Stdin = reader
	go func() {
		defer writer.Close()
		// the writer is connected to the reader via the pipe
		// so all data written here is passed on to the commands
		// standard input
		writer.Write([]byte(mail))
	}()
	err = cmd.Run()
	if err != nil {
		log.Print("ERROR: ", "Unable to sendmail: ", body)
	}
}
