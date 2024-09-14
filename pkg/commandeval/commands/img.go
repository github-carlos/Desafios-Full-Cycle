package commands

import (
	"fmt"
	"trevas-bot/pkg/commandextractor"
	"trevas-bot/pkg/converter"
	"trevas-bot/pkg/platform"
	"trevas-bot/pkg/platform/types"
)

type ImgCommand struct {
	key      string
	platform platform.WhatsAppIntegration
}

func (p ImgCommand) Handler(input commandextractor.CommandInput) {
	fmt.Println("Running Img Command")

  go p.platform.SendReaction(&input.EventMessage, platform.LoadingReaction)

	mediaBytes, err := p.platform.ExtractStickerMediaBytes(&input.EventMessage)

	if err != nil {
		p.platform.SendReply("Não foi possível concluir o comando :c", &input.EventMessage)
		return
	}

  jpgMedia, err := converter.Webp2Img(mediaBytes)

  if err != nil {
    p.platform.SendReaction(&input.EventMessage, platform.ErrorReaction)
    return
  }

  err = p.platform.SendImg(types.SendImageInput{Image: jpgMedia}, &input.EventMessage)

  if err != nil {
    p.platform.SendReaction(&input.EventMessage, platform.ErrorReaction)
    p.platform.SendReply(err.Error(), &input.EventMessage)
    return
  }

  go p.platform.SendReaction(&input.EventMessage, platform.SuccessReaction)
}

func (c ImgCommand) GetKey() string {
	return c.key
}

func NewImgCommand() *ImgCommand {
	return &ImgCommand{key: "img"}
}
