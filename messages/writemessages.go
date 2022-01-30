package messages

import (
	"fmt"
	"log"
	"os"
)

var CHANNEL string

func init() {
	CHANNEL = os.Getenv("TWITCH_CHANNEL")

	if len(CHANNEL) == 0 {
		log.Fatalln("TWITCH_CHANNEL not defined or missing.")
	}
}

func WriteMessages(message string, f *os.File) {
	fmt.Fprintln(f, message)
	err := f.Sync()
	if err != nil {
		log.Fatalln("Failed to update file. %s", err)
	}
}
