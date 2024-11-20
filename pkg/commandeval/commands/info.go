package commands

import (
	"fmt"
	"strconv"
	"strings"
	"trevas-bot/pkg/commandextractor"
	"trevas-bot/pkg/platform"
)

type InfoCommand struct {
	key      string
	platform platform.WhatsAppIntegration
}

func (p InfoCommand) Handler(input commandextractor.CommandInput) {
	fmt.Println("Running Info Command")

  channelJID := input.EventMessage.Info.Chat.String()
  payload := strings.ReplaceAll(input.Payload, "@", "")

  if payload == "" {
    payload = input.EventMessage.Info.Sender.String()
  }

  // nome
  // numero de mensagens no grupo
  // hora e data da ultima mensagem
  // ultima mensagem
  // comando mais usado

  // go p.platform.SendReply("teste", &input.EventMessage)
  userInfo, err := input.Store.GetUserChannelInfo(payload, channelJID)

  if err != nil {
    p.platform.SendReply("Erro ao buscar informações do Usuário :c", &input.EventMessage);
    return;
  }

  var mostUsedCommands string

for _, command := range userInfo.MostUsedCommands {
    mostUsedCommands += fmt.Sprintf("\n*%s*: %s", command.Name, strconv.Itoa(command.Quantity));
  }

  text := fmt.Sprintf(`
*Nome*: %s

*Número de mensagens enviadas no grupo*: %s

*Última mensagem*: %s

*Hora última mensagem enviada*: %s

*Comandos mais usados*:
    %s
    `, userInfo.Name, strconv.Itoa(userInfo.GroupMessagesCount), userInfo.LastMessage, userInfo.TimeLastMessage, mostUsedCommands)

  p.platform.SendReply(text, &input.EventMessage);
}

func (c InfoCommand) GetKey() string {
	return c.key
}

func NewInfoCommand() *InfoCommand {
  return &InfoCommand{key: "info"}
}
