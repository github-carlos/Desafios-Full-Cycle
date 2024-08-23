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
  // porn := commands.NewPornCommand()
  ramos := commands.NewRamosCommand()
  meme := commands.NewMemeCommand()
  gen := commands.NewGenCommand()
  sorteio := commands.NewSorteioCommand()
  video := commands.NewVideoCommand()
  saude := commands.NewSaudeCommand()
  img := commands.NewImgCommand()
  sexta := commands.NewSextaCommand()
  lramos := commands.NewLRamosCommand()
  zeImg := commands.NewZéCommand()
	const commandPrefix = '!'

	commands := make(map[string]commands.Commander)

	commands[ping.GetKey()] = ping
	commands[bola.GetKey()] = bola
	commands[sticker.GetKey()] = sticker
	commands[download.GetKey()] = download
  commands[ze.GetKey()] = ze
  commands[viadometro.GetKey()] = viadometro
  commands[top5.GetKey()] = top5
  // commands[porn.GetKey()] = porn
  commands[ramos.GetKey()] = ramos
  commands[meme.GetKey()] = meme
  commands[gen.GetKey()] = gen
  commands[sorteio.GetKey()] = sorteio
  commands[video.GetKey()] = video
  commands[saude.GetKey()] = saude
  commands[img.GetKey()] = img
  commands[sexta.GetKey()] = sexta
  commands[lramos.GetKey()] = lramos
  commands[zeImg.GetKey()] = zeImg
  commands[ze.GetKey()] = ze

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
