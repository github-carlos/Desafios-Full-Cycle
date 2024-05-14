package commands

import (
	"fmt"
	"os"

	"math/rand"
	"trevas-bot/pkg/commandextractor"
	"trevas-bot/pkg/converter"
	"trevas-bot/pkg/platform"
)

var phrases = []string{
	"Boa, Milhaça!",
	"E a véia tá no tchutchero, né?!",
	"Chato pá buné!",
	"Ticaracatica!",
	"Fiu fiu!",
	"O cara ficou muito pistola véi!",
	"Puta cheiro de chalalá!",
	"Faaaala, fiote!",
	"Tá sabendo legal!",
	"Vai, ô troxa!",
	"HÁ HÁ HÁ HÁ",
	"Nossa, velho!",
	"Completamente lelé...",
	"Puta nego chato, meu!",
	"Puta nego tonto!",
	"Deixa ele falar, meu!",
	"Putsa lá miséria, meu!",
	"Ahhh vá, é memo?",
	"Putcha que me paroca!",
	"Que jumento, véio!",
	"Há há, show de bola!",
	"Eu já vi nego cagão, mas olha... Parabéns, viu meu!",
	"Parabéns, Zé! HÁ HÁ HÁ HÁ! Vai ser pai de novo!",
}

var types = []string{
	"img",
	"sticker",
	"text",
	"text",
	"text",
}

type BolaCommand struct {
	key      string
	Platform platform.WhatsAppIntegration
}

func (b BolaCommand) Handler(commandInput commandextractor.CommandInput) {
	fmt.Println("Running Bola Command")

	b.Platform.SendReaction(&commandInput.EventMessage, platform.BolaReacton)

	messageType := types[rand.Intn(len(types))]

	if messageType == "text" {
		b.sendRandomPhrase(&commandInput)
		return
	}

	b.sendRandomImageSticker(&commandInput, messageType == "img")
}

func (b BolaCommand) sendRandomPhrase(commandInput *commandextractor.CommandInput) {
	randomPhrase := rand.Intn(len(phrases))
	error := b.Platform.SendReply(phrases[randomPhrase], &commandInput.EventMessage)
	if error != nil {
		fmt.Println("Error sending Bola random message")
	}

}
func (b BolaCommand) sendRandomImageSticker(commandInput *commandextractor.CommandInput, isImg bool) {
	photos, _ := os.ReadDir("assets/bola")

	randomPhoto := rand.Intn(len(photos))
	filePath := fmt.Sprintf("assets/bola/%s", photos[randomPhoto].Name())

	fmt.Println("Bola image", filePath)
	img, err := os.ReadFile(filePath)

	if err != nil {
		fmt.Println("Failing getting Bola image")
		return
	}

	webp, err := converter.Img2Webp(img)

	if err != nil {
		fmt.Println("Failing converting img to webp")
		return
	}

	if !isImg {
		err = b.Platform.SendSticker(webp, false, &commandInput.EventMessage)

		if err != nil {
			fmt.Println("Failing sending message")
			return
		}
		return
	}

	err = b.Platform.SendImg(webp, false, &commandInput.EventMessage)

	if err != nil {
		fmt.Println("Failing sending message")
		return
	}
}

func (c BolaCommand) GetKey() string {
	return c.key
}

func NewBolaCommand() *BolaCommand {
	return &BolaCommand{key: "bola"}
}
