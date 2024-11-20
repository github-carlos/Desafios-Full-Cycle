package commands

import (
	"fmt"
	"os"
	"trevas-bot/pkg/commandextractor"
	"trevas-bot/pkg/converter"
	"trevas-bot/pkg/platform"
)

type ZéCommand struct {
	key      string
	platform platform.WhatsAppIntegration
}

func (p ZéCommand) Handler(input commandextractor.CommandInput) {
	fmt.Println("Running Zé Command")
	filePath := "assets/group/ze.jpeg"
	img, err := os.ReadFile(filePath)

	if err != nil {
		fmt.Println("Failing getting Ze image")
		return
	}

	webp, err := converter.Img2Webp(img, false)

	if err != nil {
		fmt.Println("Failing converting img to webp")
		return
	}

	err = p.platform.SendSticker(webp, false, &input.EventMessage, true)

	if err != nil {
		fmt.Println("Failing sending message")
		return
	}
}

func (c ZéCommand) GetKey() string {
	return c.key
}

func NewZéCommand() *ZéCommand {
	return &ZéCommand{key: "zé"}
}
