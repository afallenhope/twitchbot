package tui

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"io"
	"log"
	"os"
	"gitlab.com/afallenhope/twitchbot/irc"
	"gitlab.com/afallenhope/twitchbot/parser"
)

var CHANNEL string
var NICK string

func init() {
	CHANNEL = os.Getenv("TWITCH_CHANNEL")
	NICK = os.Getenv("TWITCH_NICK")
}

var textView = tview.NewTextView()
var inputView = tview.NewInputField()

func PassMessage(message string, messagesChan chan parser.PrivMessage) {
	if len(message) > 0 {
		messagesChan <- parser.PrivMessage{Text: message, User: irc.NICK, Channel: irc.CHANNEL, IsCommand: false}
	}
}

func MessageTUI(message string) {
	_, err := fmt.Fprintf(textView, message)

	if err != nil {
		log.Fatalln("Error updating UI (textView)", err)
	}
}

func BuildTUI(conn io.Writer, messageChan chan parser.PrivMessage) {

	app := tview.NewApplication()

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyCtrlC {
			app.Stop()
			os.Exit(0)
		}
		return event
	})

	textView.SetDynamicColors(true).
		SetRegions(true).
		SetChangedFunc(func() {
			app.Draw()
		})

	inputView.SetLabel("Enter Message: ")
	inputView.SetDoneFunc(func(key tcell.Key) {
		textSend := inputView.GetText()
		fmt.Fprintf(conn, "PRIVMSG #%s :%s\r\n", CHANNEL, textSend)
		PassMessage(textSend, messageChan)
		inputView.SetText("")
	})

	flexView := tview.NewFlex().
		SetDirection(0).
		AddItem(textView, 0, 1, false).
		AddItem(inputView, 0, 1, true)

	if err := app.SetRoot(flexView, true).Run(); err != nil {
		log.Fatalln(err)
	}
}
