package server

import "github.com/trinchan/slackbot/robots"

func getRobots(bots map[string][]robots.Robot, command string) []robots.Robot {
	if r, ok := bots[command]; ok {
		return r
	}
	return nil
}
