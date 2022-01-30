package commands

import (
	"encoding/json"
	"log"
	"os"
)

// Commands structure.
//
// {
//	"command" : "response",
//	"parameters": 0,
//  "prefix": "!"
// }

type Command struct {
	Name       string   `json:"name"`
	Response   string   `json:"response"`
	Parameters []string `json:"parameters"`
}

func GetCommands() map[string]Command {
	//content, err := ioutil.ReadFile("commands.json")
	f, err := os.Open("commands/commands.json")
	if err != nil {
		log.Fatalln("Failed to read commands file.", err)
	}

	loadedCommands := []Command{}
	json.NewDecoder(f).Decode(&loadedCommands)

	commandsMap := make(map[string]Command)
	for _, v := range loadedCommands {
		commandsMap[v.Name] = v
	}
	return commandsMap

}
