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

	mediaBytes, error := p.platform.ExtractMediaBytes(&input.EventMessage)

	isReplying := p.platform.IsReplying(&input.EventMessage)

	if !isReplying && error != nil {
		p.platform.SendReply(error.Error(), &input.EventMessage)
		return
	}

	if isReplying && mediaBytes == nil {
		repliedJid := p.platform.GetJidReplied(&input.EventMessage)
		repliedText := p.platform.GetQuotedText(&input.EventMessage)
		profilePic, _ := p.platform.GetProfilePicture(repliedJid)
		fmt.Println(repliedJid, repliedText, profilePic)
	}

	if mediaBytes == nil {
		p.platform.SendReply("Envie uma imagem, video, gif ou marque uma conversa.", &input.EventMessage)
		return
	}

	go p.platform.SendReaction(&input.EventMessage, platform.BullReaction)
	p.platform.SendReply("Um momento meu amigo bovino", &input.EventMessage)

	isVideo := p.platform.HasVideo(&input.EventMessage)
	webpMedia, err := converter.Img2Webp(mediaBytes, isVideo)

	if err != nil {
		go p.platform.SendReaction(&input.EventMessage, platform.ErrorReaction)
		return
	}

	err = p.platform.SendSticker(webpMedia, true, &input.EventMessage, true)

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
