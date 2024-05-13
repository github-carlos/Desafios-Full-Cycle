package commandeval

import (
	"trevas-bot/pkg/commandeval/commands"
	"trevas-bot/pkg/commandextractor"
)

type CommandEval struct {
	commandPrefix byte
	commands      map[string]commands.Commander
}

func NewCommandEval() *CommandEval {
	ping := commands.NewPingCommand()

	const commandPrefix = '!'

	commands := make(map[string]commands.Commander)

	commands[ping.GetKey()] = ping

	return &CommandEval{commandPrefix, commands}
}

func (c CommandEval) Handle(commandInput commandextractor.CommandInput) error {
	// TODO: extractKey
	const key = "ping"
	return nil
}
