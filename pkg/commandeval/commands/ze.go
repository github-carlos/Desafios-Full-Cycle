package commands

import (
	"fmt"

	"math/rand"
	"trevas-bot/pkg/commandextractor"
	"trevas-bot/pkg/platform"
)

var phrasesZe = []string{
	"Mermão vá se fuder você e seu futebol",
	"L Ramos nosso Deus",
	"Quero aprender informática ☝️",
	"🖕🖕🖕🖕🖕🖕🖕🖕🖕🖕🖕🖕🖕🖕🖕🖕🖕🖕🖕 VA TOMA NO CU WELLINGTON DOMINGUES 🖕🖕🖕🖕🖕🖕🖕🖕🖕🖕🖕🖕🖕🖕🖕🖕🖕🖕",
	"Alto em \"Saúde Adicionada\" 👍🏼",
	"Coldplay é bosta",
	"Chorão do Charlie Brown Jr era tão foda e tão sábio que morreu cagado e vomitado de tanto usar droga.",
  "🖕🖕🖕🖕🖕🖕🖕🖕🖕🖕🖕🖕🖕🖕🖕🖕🖕🖕🖕🖕🖕🖕 VA TOMA NO CU L RAMOS BURRO DO CARALHO 🖕🖕🖕🖕🖕🖕🖕🖕🖕🖕🖕🖕🖕🖕🖕🖕🖕🖕🖕🖕🖕🖕🖕🖕🖕🖕🖕",
  "L Ramos frequenta boate gay no Maranhão.",
  "ThGTech possui mais de 500gb de porno infantil.",
  "Miranda frequenta boate gay no interior de SP.",
  "LN desvia diariamente 65 marmitas entre 11 e 13h.",
  "L Ramos mora sozinho e recebe visitas diárias de machos pauzudos.",
}

type ZeCommand struct {
	key      string
	Platform platform.WhatsAppIntegration
}

func (b ZeCommand) Handler(commandInput commandextractor.CommandInput) {
	fmt.Println("Running Ze Command")

	b.sendRandomPhrase(&commandInput)

}

func (b ZeCommand) sendRandomPhrase(commandInput *commandextractor.CommandInput) {
	randomPhrase := rand.Intn(len(phrasesZe))
	go b.Platform.SendReply(phrasesZe[randomPhrase], &commandInput.EventMessage)
}

func (c ZeCommand) GetKey() string {
	return c.key
}

func NewZeCommand() *ZeCommand {
	return &ZeCommand{key: "zé"}
}
