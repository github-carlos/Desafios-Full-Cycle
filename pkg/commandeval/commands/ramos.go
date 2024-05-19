package commands

import (
	"fmt"

	"math/rand"
	"trevas-bot/pkg/commandextractor"
	"trevas-bot/pkg/platform"
)

var phrasesRamos = []string{
  "Ok ateu corno",
  "Estou triste amigos...",
  "Deus abençoe a todos",
  "Belissimo dia 👍",
  "Tristeza amigos",
}

type RamosCommand struct {
	key      string
	Platform platform.WhatsAppIntegration
}

func (b RamosCommand) Handler(commandInput commandextractor.CommandInput) {
	fmt.Println("Running Ramos Command")

	b.sendRandomPhrase(&commandInput)

}

func (b RamosCommand) sendRandomPhrase(commandInput *commandextractor.CommandInput) {
	randomPhrase := rand.Intn(len(phrasesRamos))
	go b.Platform.SendReply(phrasesRamos[randomPhrase], &commandInput.EventMessage)
}

func (c RamosCommand) GetKey() string {
	return c.key
}

func NewRamosCommand() *RamosCommand {
	return &RamosCommand{key: "ramos"}
}
