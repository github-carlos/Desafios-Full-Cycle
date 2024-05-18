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

var top3Icons map[int]string = map[int]string{
	1: "🥇",
	2: "🥈",
	3: "🥉",
}

type Top5Command struct {
	key      string
	platform platform.WhatsAppIntegration
}

func (p Top5Command) Handler(input commandextractor.CommandInput) {
	fmt.Println("Running Top5 Command")
	participants, _ := p.platform.GetParticipantsOfGroup(&input.EventMessage)

	rand.Shuffle(len(participants), func(i, j int) {
		participants[i], participants[j] = participants[j], participants[i]
	})

  title := strings.ToUpper(input.Payload)

  if title == "" {
    p.platform.SendReply("É necessário um título para o Top 5", &input.EventMessage)
    return
  }

  p.platform.SendReply("⌛⌛ *...Calculando...* ⌛⌛", &input.EventMessage)
  time.Sleep(3 * time.Second)

	text := fmt.Sprintf("❗ ATENÇÃO ❗\n\n⚠️ RESULTADO *TOP 5 %s* ⚠️\n\n", title)

	defaultNumberOfChoosen := 5

	if len(participants) < defaultNumberOfChoosen {
		defaultNumberOfChoosen = len(participants)
	}

	for indice, user := range participants[:defaultNumberOfChoosen] {
		position := indice + 1
		medal := top3Icons[position]
		line := fmt.Sprintf("%d. %s @%s\n", position, medal, user)
		text += line
	}

  p.platform.SendText(platformTypes.SendTextInput{Text: text, Mentions: participants[:defaultNumberOfChoosen]}, &input.EventMessage)
}

func (c Top5Command) GetKey() string {
	return c.key
}

func NewTop5Command() *Top5Command {
	return &Top5Command{key: "top5"}
}
