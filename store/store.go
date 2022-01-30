package store

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"
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

func CreateLogFile() *os.File {
	date := time.Now().Format("02-01-2006")
	filename := fmt.Sprintf("logs/chat-%s-%s.log", os.Getenv("TWITCH_CHANNEL"), date)

	f, fileError := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)

	if fileError != nil {
		log.Fatalln(fileError)
	}

	return f
}

func LogMessages(f *os.File, message string) {

	_, err := f.WriteString(message)
	if err != nil {
		log.Fatalln("Failed to write message to the logfile.", err)
	}
}
