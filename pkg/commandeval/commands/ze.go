package commands

import (
	"fmt"
	"os"
	"strconv"

	"math/rand"
	"trevas-bot/pkg/commandextractor"
	"trevas-bot/pkg/converter"
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
  "ThGTech possui mais de 500gb de Café Pelé.",
  "Miranda frequenta boate gay no interior de SP.",
  "LN desvia diariamente 65 marmitas entre 11 e 13h.",
  "L Ramos mora sozinho e recebe visitas diárias de machos pauzudos.",
  "Caio frequenta cabarés todos os fins de semana na capital paulista.",
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
  payload := commandInput.Payload;
  convertedIndex, _ := strconv.Atoi(payload)

  if commandInput.Payload == "amimir" {

    filePath := "assets/group/amimir.jpeg"
    img, err := os.ReadFile(filePath)

    if err != nil {
      fmt.Println("Failing getting amimir image")
      return
    }

    webp, err := converter.Img2Webp(img)

    if err != nil {
      fmt.Println("Failing converting img to webp")
      return
    }

    err = b.Platform.SendSticker(webp, false, &commandInput.EventMessage, true)

    if err != nil {
      fmt.Println("Failing sending message")
      return
    }
    return;
  }

  phraseIndex := 0;
  if convertedIndex <= len(phrasesZe) && convertedIndex > 0 {
    phraseIndex = convertedIndex - 1;
  } else {
    phraseIndex = rand.Intn(len(phrasesZe))
  }

	go b.Platform.SendReply(phrasesZe[phraseIndex], &commandInput.EventMessage)
}

func (c ZeCommand) GetKey() string {
	return c.key
}

func NewZeCommand() *ZeCommand {
	return &ZeCommand{key: "zé"}
}
