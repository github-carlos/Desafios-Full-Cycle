package commands

import (
	"fmt"
	"os"
	"trevas-bot/pkg/commandextractor"
	"trevas-bot/pkg/converter"
	"trevas-bot/pkg/platform"
)

type SaudeCommand struct {
	key      string
	platform platform.WhatsAppIntegration
}

func (p SaudeCommand) Handler(input commandextractor.CommandInput) {
	fmt.Println("Running Saude Command")
	go p.platform.SendReply("Alto em \"Saúde Adicionada\" 👍 ", &input.EventMessage)

  filePath := "assets/group/saude.jpeg"
	img, err := os.ReadFile(filePath)

	if err != nil {
		fmt.Println("Failing getting Saude image")
		return
	}

	webp, err := converter.Img2Webp(img)

	if err != nil {
		fmt.Println("Failing converting img to webp")
		return
	}

  err = p.platform.SendSticker(webp, false, &input.EventMessage)

  if err != nil {
    fmt.Println("Failing sending message")
    return
  }
}

func (c SaudeCommand) GetKey() string {
	return c.key
}

func NewSaudeCommand() *SaudeCommand {
	return &SaudeCommand{key: "saude"}
}
