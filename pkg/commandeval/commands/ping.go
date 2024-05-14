package commands

import (
	"fmt"
	"trevas-bot/pkg/commandextractor"
	"trevas-bot/pkg/platform"
)

type PingCommand struct {
	key      string
	platform platform.WhatsAppIntegration
}

func (p PingCommand) Handler(input commandextractor.CommandInput) {
	fmt.Println("Running Ping Command")
	error := p.platform.SendReply("pong", &input.EventMessage)
	if error != nil {
		fmt.Println("Error sending Ping command")
	}
}

func (c PingCommand) GetKey() string {
	return c.key
}

func NewPingCommand() *PingCommand {
	return &PingCommand{key: "ping"}
}
