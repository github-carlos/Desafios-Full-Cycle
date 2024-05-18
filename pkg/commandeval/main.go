package commandeval

import (
	"fmt"
	"trevas-bot/pkg/commandeval/commands"
	"trevas-bot/pkg/commandextractor"
	"trevas-bot/pkg/platform"
)

type CommandEval struct {
	commandPrefix byte
	commands      map[string]commands.Commander
	platform      platform.WhatsAppIntegration
}

func NewCommandEval() *CommandEval {
	ping := commands.NewPingCommand()
	bola := commands.NewBolaCommand()
	sticker := commands.NewStickerCommand()
  download := commands.NewDownloadCommand()
  ze := commands.NewZeCommand()
  viadometro := commands.NewViadometroCommand()
  top5 := commands.NewTop5Command()
	const commandPrefix = '!'

	commands := make(map[string]commands.Commander)

	commands[ping.GetKey()] = ping
	commands[bola.GetKey()] = bola
	commands[sticker.GetKey()] = sticker
	commands[download.GetKey()] = download
  commands[ze.GetKey()] = ze
  commands["ze"] = ze
  commands[viadometro.GetKey()] = viadometro
  commands[top5.GetKey()] = top5

	whatsApp := platform.NewWhatsAppIntegration()

	return &CommandEval{commandPrefix, commands, *whatsApp}
}

func (c CommandEval) Handle(commandInput *commandextractor.CommandInput) error {

	command := c.commands[commandInput.Command]

	if command == nil {
		fmt.Println("Command not found")
		return nil
	}

	command.Handler(*commandInput)

	return nil
}
