package commands

import (
	"fmt"
	"strconv"

	"math/rand"
	"trevas-bot/pkg/commandextractor"
	"trevas-bot/pkg/platform"
)

var phrasesGDiesel = []string{
	"oh, nojeira!",
}

type GDieselCommand struct {
	key      string
	Platform platform.WhatsAppIntegration
}

func (b GDieselCommand) Handler(commandInput commandextractor.CommandInput) {
	fmt.Println("Running GDiesel Command")

	b.sendRandomPhrase(&commandInput)
}

func (b GDieselCommand) sendRandomPhrase(commandInput *commandextractor.CommandInput) {
  payload := commandInput.Payload;
  convertedIndex, _ := strconv.Atoi(payload)

  phraseIndex := 0;
  if convertedIndex <= len(phrasesGDiesel) && convertedIndex > 0 {
    phraseIndex = convertedIndex - 1;
  } else {
    phraseIndex = rand.Intn(len(phrasesGDiesel))
  }

	go b.Platform.SendReply(phrasesGDiesel[phraseIndex], &commandInput.EventMessage)
}

func (c GDieselCommand) GetKey() string {
	return c.key
}

func NewGDieselCommand() *GDieselCommand {
	return &GDieselCommand{key: "gdiesel"}
}
