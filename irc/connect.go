package irc

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net"
	"os"
)

var OAUTH_TOKEN string
var NICK string
var CHANNEL string

func init() {
	err := godotenv.Load("local.env")

	if err != nil {
		log.Fatalln("Error loading local.env file. Err: %s", err)
	}

	// read from our ENV variables.
	OAUTH_TOKEN = os.Getenv("TWITCH_OAUTH_TOKEN")
	NICK = os.Getenv("TWITCH_NICK")
	CHANNEL = os.Getenv("TWITCH_CHANNEL")

	if len(OAUTH_TOKEN) == 0 {
		log.Fatalln("TWITCH_OAUTH_TOKEN not defined or missing.")
	}

	if len(NICK) == 0 {
		log.Fatalln("TWITCH_USERNAME not defined or missing")
	}

	if len(CHANNEL) == 0 {
		log.Fatalln("TWITCH_CHANNEL not defined or missing")
	}
}

func Connect() net.Conn {
	conn, err := net.Dial("tcp", "irc.chat.twitch.tv:6667")
	if err != nil {
		log.Fatalln("Failed to dial to IRC", err)
	}

	fmt.Fprintf(conn, "PASS %s\r\n", OAUTH_TOKEN)
	fmt.Fprintf(conn, "NICK %s\r\n", NICK)
	fmt.Fprintf(conn, "JOIN #%s\r\n", CHANNEL)
	fmt.Fprintf(conn, "CAP REQ :twitch.tv/membership\r\n")
	fmt.Fprintf(conn, "CAP REQ :twitch.tv/tags\r\n")
	fmt.Fprintf(conn, "CAP REQ :twitch.tv/commands\r\n")

	return conn
}
