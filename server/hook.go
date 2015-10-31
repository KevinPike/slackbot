package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/trinchan/slackbot/Godeps/_workspace/src/github.com/gorilla/schema"
	"github.com/trinchan/slackbot/robots"
)

// HookHandler creates a http.HandlerFunc to handle outgoing webhooks
func HookHandler(bots map[string][]robots.Robot) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		d := schema.NewDecoder()
		command := new(robots.OutgoingWebHook)
		err = d.Decode(command, r.PostForm)
		if err != nil {
			log.Println("Couldn't parse post request:", err)
		}
		if command.Text == "" || command.Token != os.Getenv(fmt.Sprintf("%s_OUT_TOKEN", strings.ToUpper(command.TeamDomain))) {
			log.Printf("[DEBUG] Ignoring request from unidentified source: %s - %s", command.Token, r.Host)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		com := strings.TrimPrefix(command.Text, command.TriggerWord+" ")
		c := strings.Split(com, " ")
		command.Robot = c[0]
		command.Text = strings.Join(c[1:], " ")

		robots := getRobots(bots, command.Robot)
		if len(robots) == 0 {
			jsonResp(w, "No robot for that command yet :(")
			return
		}
		resp := ""
		for _, robot := range robots {
			resp += fmt.Sprintf("\n%s", robot.Run(&command.Payload))
		}
		w.WriteHeader(http.StatusOK)
		jsonResp(w, strings.TrimSpace(resp))
	}
}
