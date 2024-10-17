package commands

import (
	"fmt"
	"strings"
	"trevas-bot/pkg/commandextractor"
	"trevas-bot/pkg/llm"
	"trevas-bot/pkg/platform"
)

type InfoCommand struct {
	key      string
	platform platform.WhatsAppIntegration
  llmGenerator llm.LLMGenerator
}

func (p InfoCommand) Handler(input commandextractor.CommandInput) {
	fmt.Println("Running Info Command")

  schema := `
  messages (
  name TEXT, // user name
  channel_jid TEXT,
  message TEXT,
  type TEXT, // could be text, sticker, video, image
  command TEXT,
  timestamp DATE,
  created_at DATETIME
  )
  `

  channelJID := input.EventMessage.Info.Chat.String()
  payload := strings.ReplaceAll(input.Payload, "@", "")
  

  prompt := fmt.Sprintf(`
    You will be given a table schema surrounded by <> and your mission is to return and a user desire. Your mission is to return a Select statement that helps the user to extract info from that table. Return only the select statement and nothing more. This is the user query: %s \n This is the table schema: %s

    Here are some information that might help:
    channel_jid: %s

    always use the filter channel_jid.

    for name filter, always use the like operator with wildcard surrounded
    if its a phone number, use jid column on filter with like operator wildcard surrounded

    return in the following format:

    SELECT name, message FROM messages;

    always use a limit of max 50 rows`, payload, schema, channelJID)

  statement, err := p.llmGenerator.Complete(prompt)
  if err != nil {
    go p.platform.SendReply("Tente novamente mais tarde", &input.EventMessage)
    return
  }

  extractedSQL := strings.TrimPrefix(statement, "```sql\n")
  extractedSQL = strings.TrimSuffix(extractedSQL, "\n```")

  fmt.Println("Statement", extractedSQL)
  result, err := input.Store.ExecSql(extractedSQL)

  if err != nil {
    fmt.Println("Error", err)
    go p.platform.SendReply("Tente novamente mais tarde", &input.EventMessage)
    return
  }

  fmt.Println("RESULT", result)
  if result == "" {
    go p.platform.SendReply("Não consegui extrair essas informações. Tente de uma outra forma.", &input.EventMessage)
    return
  }

  prompt = fmt.Sprintf(`
    Elabore uma mensagem para compartilhar as informaçoes para a seguinte query: %s

    Use emojis para chamar atençao. Escreva de forma simples e direta o mais breve possivel
    Essas são todas as informaçoes:
    %s

    Use no máximo 300 caracteres. Nao envie ids ou jids ou qualquer outra informacao tecnica
    `, payload, result)

  text, err := p.llmGenerator.Complete(prompt)

  if err != nil {
    fmt.Println("Error generating last result", err)
    go p.platform.SendReply("Tente novamente mais tarde", &input.EventMessage)
    return
  }

  go p.platform.SendReply(text, &input.EventMessage)
}

func (c InfoCommand) GetKey() string {
	return c.key
}

func NewInfoCommand(llmGenerator llm.LLMGenerator) *InfoCommand {
  return &InfoCommand{key: "info", llmGenerator: llmGenerator}
}
