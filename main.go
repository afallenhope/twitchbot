package main

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"twitchbot/irc"
	"twitchbot/messages"
	"twitchbot/parser"
	"twitchbot/router"
	"twitchbot/tui"
)

func init() {
	godotenv.Load("local.env")
	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		err := os.Mkdir("logs", 666)
		if err != nil {
			log.Fatalln("Could not make logs directory.")
		}
	}
}

func main() {
	conn := irc.Connect()
	messageChan := messages.ReadMessages(conn)
	parsedMessage, parsedCommand, parsedPings := parser.ParseMessages(messageChan, "!")
	router.RouteMessages(parsedMessage, parsedCommand, parsedPings, conn)
	tui.BuildTUI(conn, parsedMessage)
}
