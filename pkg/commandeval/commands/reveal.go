package commands

import (
	"fmt"
	"trevas-bot/pkg/commandextractor"
	"trevas-bot/pkg/platform"
	"trevas-bot/pkg/platform/types"
)

type RevealCommand struct {
	key      string
	platform platform.WhatsAppIntegration
}

func (p RevealCommand) Handler(input commandextractor.CommandInput) {
	fmt.Println("Running Reveal Command")

	mediaBytes, error := p.platform.ExtractMediaBytes(&input.EventMessage)

  if error != nil {
		p.platform.SendReply(error.Error(), &input.EventMessage)
		return
  }

	isReplying := p.platform.IsReplying(&input.EventMessage)

	if !isReplying {
		p.platform.SendReply("Marque uma mensagem de visualização única", &input.EventMessage)
		return
	}

	if mediaBytes == nil {
		p.platform.SendReply("Mídia indisponível", &input.EventMessage)
		return
	}

	isVideo := p.platform.HasVideo(&input.EventMessage)

  if isVideo {
    go p.platform.SendVideo(types.SendVideoInput{VideoBytes: mediaBytes}, &input.EventMessage)
  } else {
    go p.platform.SendImg(types.SendImageInput{Image: mediaBytes}, &input.EventMessage)
  }
}

func (c RevealCommand) GetKey() string {
	return c.key
}

func NewRevealCommand() *RevealCommand {
	return &RevealCommand{key: "reveal"}
}
