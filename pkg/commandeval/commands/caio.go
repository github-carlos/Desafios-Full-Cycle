package commands

import (
	"fmt"
	"strconv"

	"math/rand"
	"trevas-bot/pkg/commandextractor"
	"trevas-bot/pkg/platform"
)

var phrasesCaio = []string{
	"ACORDA ZÉ, TÃO LEVANDO SUA CRYPTON!",
}

type CaioCommand struct {
	key      string
	Platform platform.WhatsAppIntegration
}

func (b CaioCommand) Handler(commandInput commandextractor.CommandInput) {
	fmt.Println("Running Caio Command")

	b.sendRandomPhrase(&commandInput)
}

func (b CaioCommand) sendRandomPhrase(commandInput *commandextractor.CommandInput) {
  payload := commandInput.Payload;
  convertedIndex, _ := strconv.Atoi(payload)

  phraseIndex := 0;
  if convertedIndex <= len(phrasesCaio) && convertedIndex > 0 {
    phraseIndex = convertedIndex - 1;
  } else {
    phraseIndex = rand.Intn(len(phrasesCaio))
  }

	go b.Platform.SendReply(phrasesCaio[phraseIndex], &commandInput.EventMessage)
}

func (c CaioCommand) GetKey() string {
	return c.key
}

func NewCaioCommand() *CaioCommand {
	return &CaioCommand{key: "caio"}
}
