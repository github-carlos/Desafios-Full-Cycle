package commandeval

import (
	"errors"
	"fmt"
	"strings"
	"trevas-bot/pkg/commandeval/commands"
)

type CommandEval struct {
	commandPrefix byte
	commands      map[string]commands.Command
}

func NewCommandEval() *CommandEval {
	ping := commands.NewPingCommand()

	const commandPrefix = '!'

	commands := make(map[string]commands.Command)

	commands[ping.GetKey()] = ping

	return &CommandEval{commandPrefix, commands}
}

func (c CommandEval) Handle(message string) error {
	clearedMessage := strings.Trim(message, " ")
	fmt.Println("Cleared Message:", clearedMessage)

	if clearedMessage[0] != c.commandPrefix {
    return errors.New("Comando Inválido!")
	}

  return nil
}
