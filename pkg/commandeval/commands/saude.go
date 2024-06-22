package commands

import (
	"fmt"
	"trevas-bot/pkg/commandextractor"
	"trevas-bot/pkg/platform"
)

type SaudeCommand struct {
	key      string
	platform platform.WhatsAppIntegration
}

func (p SaudeCommand) Handler(input commandextractor.CommandInput) {
	fmt.Println("Running Saude Command")
	go p.platform.SendReply("Alto em \"Saúde Adicionada\" 👍 ", &input.EventMessage)
}

func (c SaudeCommand) GetKey() string {
	return c.key
}

func NewSaudeCommand() *SaudeCommand {
	return &SaudeCommand{key: "saude"}
}
