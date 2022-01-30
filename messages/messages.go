package messages

import (
	"bufio"
	"io"
	"log"
)

func ReadMessages(conn io.Reader) <-chan string {
	messageReader := bufio.NewReader(conn)
	messagesChan := make(chan string)

	go func() {
		for {
			message, err := messageReader.ReadString('\n')

			if err != nil {
				log.Fatalln("Failed to read message", err)
			}

			messagesChan <- message
		}
	}()

	return messagesChan
}
