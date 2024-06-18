package commands

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
	"trevas-bot/pkg/commandextractor"
	"trevas-bot/pkg/platform"
	platformTypes "trevas-bot/pkg/platform/types"
)

type SorteioCommand struct {
	key      string
	platform platform.WhatsAppIntegration
}

func (p SorteioCommand) Handler(input commandextractor.CommandInput) {
	fmt.Println("Running Sorteio Command")
	participants, _ := p.platform.GetParticipantsOfGroup(&input.EventMessage)

	rand.Shuffle(len(participants), func(i, j int) {
		participants[i], participants[j] = participants[j], participants[i]
	})

  title := strings.ToUpper(input.Payload)

  if title == "" {
    p.platform.SendReply("É necessário um tema para o sorteio.", &input.EventMessage)
    return
  }

  p.platform.SendReply("🎲🎲 *...Sorteando...* 🎲🎲", &input.EventMessage)
  time.Sleep(3 * time.Second)

  text := "\t\t❗ ATENÇÃO ❗\n\n 🎉 RESULTADO DO SORTEIO 🎉\n\n"

	defaultNumberOfChoosen := 1

	if len(participants) < defaultNumberOfChoosen {
		defaultNumberOfChoosen = len(participants)
	}

  var mentions []string

  mentionedNumber := extractUserNumber(title)

  if mentionedNumber != "" {
    mentions = append(mentions, mentionedNumber)
  }

	for _, user := range participants[:defaultNumberOfChoosen] {
		line := fmt.Sprintf("*E O GRANDE GANHADOR É....*\n\n  🎊🎊🥳🍾🎊 @%s 🎊🍾🥳🎊🎊\n\n", user)
		text += line
    text += fmt.Sprintf("*PARABÉNS, @%s*!!!! Você acaba de ganhar um(a) *%s*❗\n", user, title)
    mentions = append(mentions, user)
	}

  p.platform.SendText(platformTypes.SendTextInput{Text: text, Mentions: mentions}, &input.EventMessage)
}

func (c SorteioCommand) GetKey() string {
	return c.key
}

func NewSorteioCommand() *SorteioCommand {
	return &SorteioCommand{key: "sorteio"}
}
