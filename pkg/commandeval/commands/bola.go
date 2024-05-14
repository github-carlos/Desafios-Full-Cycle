package commands

import (
	"fmt"
	"math/rand"
	"trevas-bot/pkg/commandextractor"
	"trevas-bot/pkg/platform"
)

var phrases = []string{
    "Boa, Milhaça!",
    "E a véia tá no tchutchero, né?!",
    "Chato pá buné!",
    "Ticaracatica!",
    "Fiu fiu!",
    "O cara ficou muito pistola véi!",
    "Puta cheiro de chalalá!",
    "Faaaala, fiote!",
    "Tá sabendo legal!",
    "Vai, ô troxa!",
    "HÁ HÁ HÁ HÁ",
    "Nossa, velho!",
    "Completamente lelé...",
    "Puta nego chato, meu!",
    "Puta nego tonto!",
    "Deixa ele falar, meu!",
    "Putsa lá miséria, meu!",
    "Ahhh vá, é memo?",
    "Putcha que me paroca!",
    "Que jumento, véio!",
    "Há há, show de bola!",
    "Eu já vi nego cagão, mas olha... Parabéns, viu meu!",
    "Parabéns, Zé! HÁ HÁ HÁ HÁ! Vai ser pai de novo!",
}

type BolaCommand struct {
  key string
  Platform platform.WhatsAppIntegration
}

func (b BolaCommand) Handler(commandInput commandextractor.CommandInput) {
  fmt.Println("Running Bola Command")
  randomPhrase := rand.Intn(len(phrases))
  error := b.Platform.SendReply(phrases[randomPhrase], &commandInput.EventMessage)

  if error != nil {
    fmt.Println("Error sending message")
  }
}

func (c BolaCommand) GetKey() string {
  return c.key
}

func NewBolaCommand() *BolaCommand {
  return &BolaCommand{key: "bola"}
}
