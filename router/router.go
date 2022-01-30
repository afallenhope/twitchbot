package router

import (
	"fmt"
	"io"
	"strings"
	"twitchbot/commands"
	"twitchbot/parser"
	"twitchbot/store"
	"twitchbot/tui"
)

func formatMessage(message parser.PrivMessage) string {
	response := fmt.Sprintf("[%s] (%s): %s\n", message.Channel, message.User, message.Text)
	return response
}

func formatCommand(command parser.CommandMessage) string {
	response := fmt.Sprintf("[%s] (%s): %s [%s]\n", command.Channel, command.User, command.Command, strings.Join(command.Parameters, ", "))
	return response
}

func RouteMessages(messagesChan <-chan parser.PrivMessage, commandsChan <-chan parser.CommandMessage, pingsChan <-chan parser.PongMessage, conn io.Writer) {
	logFile := store.CreateLogFile()
	loadedCommands := commands.GetCommands()

	go func() {
		for {
			select {
			case message := <-messagesChan:
				formattedMessage := formatMessage(message)
				tui.MessageTUI(formattedMessage)
				store.LogMessages(logFile, formattedMessage)
			case command := <-commandsChan:
				foundCommand, ok := loadedCommands[strings.TrimSpace(strings.ToLower(command.Command))]
				formattedCommand := formatCommand(command)

				if ok {
					fmt.Fprintf(conn, "PRIVMSG %s :%s\n", command.Channel, foundCommand.Response)
				} else {
					foundCommand = commands.Command{}
				}

				tui.MessageTUI(fmt.Sprintf("[%s] (%s): %s\n", command.Channel, command.User, foundCommand.Response))
				store.LogMessages(logFile, formattedCommand)
			case pong := <-pingsChan:
				response := "PING received\n"
				_ = fmt.Sprintf("PONG :%s", pong.Challenge)
				fmt.Fprintf(conn, response)
				store.LogMessages(logFile, response)
			}
		}
	}()
}
