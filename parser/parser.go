package parser

import (
	"log"
	"regexp"
	"strings"
)

type PrivMessage struct {
	Text      string
	Channel   string
	User      string
	IsCommand bool
}
type PongMessage struct {
	Challenge string
}

type CommandMessage struct {
	Command    string
	Channel    string
	User       string
	Parameters []string
}

func getMessage(message string, commandPrefix string) PrivMessage {
	pattern := `:([^ +]+)!([^ ]+)\s([^ ]+)\s([^ ]+)\s?:?(.*)`
	re, err := regexp.Compile(pattern)

	if err != nil {
		log.Fatalln("Failed to compile pattern")
	}

	result := re.FindStringSubmatch(message)

	command := isCommand(result[5], commandPrefix)
	return PrivMessage{Text: result[5], Channel: result[4], User: result[1], IsCommand: command}
}

func getPing(message string) PongMessage {
	parts := strings.Split(message, ":")
	return PongMessage{Challenge: parts[1]}
}

func getCommand(message string, commandPrefix string) CommandMessage {
	parsedMessage := getMessage(message, commandPrefix)
	parameters := strings.Split(parsedMessage.Text, " ")
	return CommandMessage{Command: parameters[0], Channel: parsedMessage.Channel, User: parsedMessage.User, Parameters: parameters[1:]}
}

func isMessage(message string) bool {
	return strings.Contains(message, "PRIVMSG")
}

func isCommand(message string, prefix string) bool {
	return string(message[0]) == prefix
}

func isPing(message string) bool {
	return string(message[0:4]) == "PING"
}

func ParseMessages(messagesChan <-chan string, commandPrefix string) (chan PrivMessage, <-chan CommandMessage, <-chan PongMessage) {

	messages := make(chan PrivMessage)
	commands := make(chan CommandMessage)
	pings := make(chan PongMessage)

	go func() {
		for message := range messagesChan {
			switch {
			case isPing(message):
				pongMessage := getPing(message)
				pings <- pongMessage
			case isMessage(message):
				parsedMessage := getMessage(message, commandPrefix)

				if parsedMessage.IsCommand {
					parsedCommand := getCommand(message, commandPrefix)
					commands <- parsedCommand
				} else {
					messages <- parsedMessage
				}
			}
		}
	}()
	return messages, commands, pings
}
