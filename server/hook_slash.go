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

// SlashCommandHandler creates a http.HandlerFunc for handling slash command webhooks
func SlashCommandHandler(bots map[string][]robots.Robot) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		d := schema.NewDecoder()
		command := new(robots.SlashCommand)
		err = d.Decode(command, r.PostForm)
		if err != nil {
			log.Println("Couldn't parse post request:", err)
		}
		if command.Command == "" || command.Token == "" {
			log.Printf("[DEBUG] Ignoring request from unidentified source: %s - %s", command.Token, r.Host)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		command.Robot = command.Command[1:]

		if token := os.Getenv(fmt.Sprintf("%s_SLACK_TOKEN", strings.ToUpper(command.Robot))); token != "" && token != command.Token {
			log.Printf("[DEBUG] Ignoring request from unidentified source: %s - %s", command.Token, r.Host)
			w.WriteHeader(http.StatusBadRequest)
		}
		robots := getRobots(bots, command.Robot)
		if len(robots) == 0 {
			plainResp(w, "No robot for that command yet :(")
			return
		}
		resp := ""
		for _, robot := range robots {
			resp += fmt.Sprintf("\n%s", robot.Run(&command.Payload))
		}
		plainResp(w, strings.TrimSpace(resp))
	}
}
