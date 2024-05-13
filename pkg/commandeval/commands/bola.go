package commands

import "fmt"

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
}

func (BolaCommand) Handler(text string) {
  fmt.Println("Running Bola Command")
  fmt.Println(phrases[0])
}

func (c BolaCommand) GetKey() string {
  return c.key
}

func NewBolaCommand() *BolaCommand {
  return &BolaCommand{key: "bola"}
}
