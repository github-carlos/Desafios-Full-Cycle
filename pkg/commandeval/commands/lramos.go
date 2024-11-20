package commands

import (
	"fmt"
	"os"
	"trevas-bot/pkg/commandextractor"
	"trevas-bot/pkg/converter"
	"trevas-bot/pkg/platform"
)

type LRamosCommand struct {
	key      string
	platform platform.WhatsAppIntegration
}

func (p LRamosCommand) Handler(input commandextractor.CommandInput) {
	fmt.Println("Running LRamos Command")
	filePath := "assets/group/ramos/ramos.jpeg"
	img, err := os.ReadFile(filePath)

	if err != nil {
		fmt.Println("Failing getting Ramos image")
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

func (c LRamosCommand) GetKey() string {
	return c.key
}

func NewLRamosCommand() *LRamosCommand {
	return &LRamosCommand{key: "lramos"}
}
