package commands

import (
	"fmt"
	"trevas-bot/pkg/commandextractor"
	"trevas-bot/pkg/converter"
	"trevas-bot/pkg/platform"
)

type StickerCommand struct {
	key      string
	platform platform.WhatsAppIntegration
}

func (p StickerCommand) Handler(input commandextractor.CommandInput) {
	fmt.Println("Running Sticker Command")

  mediaBytes, err := p.platform.ExtractMediaBytes(&input.EventMessage)

  if err != nil {
    p.platform.SendReply(err.Error(), &input.EventMessage)
    return
  }

  p.platform.SendReaction(&input.EventMessage, platform.BullReaction)
	p.platform.SendReply("Um momento meu amigo bovino", &input.EventMessage)

  webpMedia, err := converter.Img2Webp(mediaBytes)

  if err != nil {
    p.platform.SendReaction(&input.EventMessage, platform.ErrorReaction)
    return
  }

  err = p.platform.SendSticker(webpMedia, true, &input.EventMessage)

  if err != nil {
    p.platform.SendReaction(&input.EventMessage, platform.ErrorReaction)
    p.platform.SendReply(err.Error(), &input.EventMessage)
    return
  }
}

func (c StickerCommand) GetKey() string {
	return c.key
}

func NewStickerCommand() *StickerCommand {
	return &StickerCommand{key: "sticker"}
}
