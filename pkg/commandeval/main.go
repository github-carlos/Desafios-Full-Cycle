package commandeval

import (
	"fmt"
	"trevas-bot/pkg/commandeval/commands"
	"trevas-bot/pkg/commandextractor"
	"trevas-bot/pkg/llm"
	"trevas-bot/pkg/platform"
)

type CommandEval struct {
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
  ramos := commands.NewRamosCommand()
  meme := commands.NewMemeCommand()
  gen := commands.NewGenCommand(llm.NewGeminiGenerator())
  sorteio := commands.NewSorteioCommand()
  video := commands.NewVideoCommand()
  saude := commands.NewSaudeCommand()
  img := commands.NewImgCommand()
  sexta := commands.NewSextaCommand()
  lramos := commands.NewLRamosCommand()
  zeImg := commands.NewZéCommand()
  caio := commands.NewCaioCommand()
  // tts := commands.NewTTSCommand()
  post := commands.NewPostCommand()
  info := commands.NewInfoCommand()
  leo := commands.NewLeoCommand(llm.NewGeminiGenerator())
  gdiesel := commands.NewGDieselCommand()
  ignore := commands.NewIgnoreCommand()
  reveal := commands.NewRevealCommand()
  fig := commands.NewFigCommand()

	commands := make(map[string]commands.Commander)

	commands[ping.GetKey()] = ping
	commands[bola.GetKey()] = bola
	commands[sticker.GetKey()] = sticker
	commands[download.GetKey()] = download
  commands[ze.GetKey()] = ze
  commands[viadometro.GetKey()] = viadometro
  commands[top5.GetKey()] = top5
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
  commands[caio.GetKey()] = caio
  // commands[tts.GetKey()] = tts
  commands[post.GetKey()] = post
  commands[info.GetKey()] = info
  commands[leo.GetKey()] = leo
  commands[gdiesel.GetKey()] = gdiesel
  commands[ignore.GetKey()] = ignore
  commands[reveal.GetKey()] = reveal
  commands[fig.GetKey()] = fig

	whatsApp := platform.NewWhatsAppIntegration()

	return &CommandEval{commands, *whatsApp}
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
