package main

import (
	"net/http"

	_ "github.com/trinchan/slackbot/importer"
	"github.com/trinchan/slackbot/robots"
	"github.com/trinchan/slackbot/server"
)

func main() {
	http.HandleFunc("/slack", server.SlashCommandHandler(robots.Robots))
	http.HandleFunc("/slack_hook", server.HookHandler(robots.Robots))
	server.Start()
}
