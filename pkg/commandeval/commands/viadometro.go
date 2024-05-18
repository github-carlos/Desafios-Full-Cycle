package commands

import (
	"fmt"
	"math/rand"
	"time"
	"trevas-bot/pkg/commandextractor"
	"trevas-bot/pkg/platform"
)

var phrasesWhenIs = []string{
  "Esse da o cu com força!",
  "Esse da o cu que chora.",
  "Em plantação de mandioca esse usa o cu de enxada.",
  "Pelo bem e pelo mal, esse cu já levou pau.",
  "Pelo circo e pelo bozo, esse cu já deu gostoso",
  "Pelo som do clarinete, esse cu levou cacete",
  "Pelo som que ele faz, esse cu já deu demais",
  "Pelo ronco do motor, esse cu já fez amor",
  "Pelo som e pelo ronco, seu cu já levou tronco",
  "Pelo toque desse sinto, esse cu levou pepino",
  "Pelo som do clarinete, esse cu levou cacete",
  "Pelo som do piano, seu rabo levou cano",
  "Pelo cheiro da cebola",
  "Pela barba do profeta, esse cu não tem mais prega",
  "Em chuva de piroca esse aí põe o cu na goteira",
  "Pelo rosnar do cachorro, esse cu pede socorro",
  "Pela Glória Maria, leva pica noite e dia",
  "Pelo cheiro de bacalhau, esse cu agasalha pau",
  "Pelo cheiro da cebola, esse cu já levou rola",
  "Pela som da filosofia, esse cu já deu até cria",
  "Confirmado!! Esse mama cacete com vontade",
  "Hummmmmmmmmmm boiola",
}

var phrasesWhenIsNot = []string {
  "Esse é conhecido como esfolador de buceta",
}

type ViadometroCommand struct {
	key      string
	platform platform.WhatsAppIntegration
}

func (p ViadometroCommand) Handler(input commandextractor.CommandInput) {
	fmt.Println("Running Viadometro Command")
  p.platform.SendReply("⌛⌛ Calculando ⌛⌛", &input.EventMessage)
  time.Sleep(3 * time.Second)

  // 90%
  is := is(9)

  fmt.Println("Payload", input.Payload)

  var phrase string
  if is {
    phrase = getRandomPhrase(phrasesWhenIs)
  } else {
    phrase = getRandomPhrase(phrasesWhenIsNot)
  }
  fmt.Println(phrase)

  fmt.Println(input.EventMessage.Info.Sender.ToNonAD().String())
  // p.platform.SendText("Meontioning " + input.Payload, &input.EventMessage)
}

func (c ViadometroCommand) GetKey() string {
	return c.key
}

func NewViadometroCommand() *ViadometroCommand {
	return &ViadometroCommand{key: "teste"}
}

func is(chanceToBe int) bool {
  randomInt := rand.Intn(10)
  return randomInt <= chanceToBe;
}

func getRandomPhrase(phrases []string) string {
  randomInt := rand.Intn(len(phrases))
  return phrases[randomInt]
}
