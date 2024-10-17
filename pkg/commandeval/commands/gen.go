package commands

import (
	"fmt"
	"trevas-bot/pkg/commandextractor"
	"trevas-bot/pkg/llm"
	"trevas-bot/pkg/platform"
	platformTypes "trevas-bot/pkg/platform/types"
)

type GenCommand struct {
	key      string
	platform platform.WhatsAppIntegration
  llmGenerator llm.LLMGenerator
}

func (p GenCommand) Handler(input commandextractor.CommandInput) {


  if input.Payload == "" {
    p.platform.SendReply("Envie um texto junto ao comando", &input.EventMessage)
    return
  }

  prompt := fmt.Sprintf("Gere uma resposta de no máximo 300 caracteres\n %s", input.Payload)

  text, err := p.llmGenerator.Complete(prompt)

  if err != nil {
    fmt.Println("Error", err)
    p.platform.SendReply("Tente novamente mais tarde", &input.EventMessage)
    return;
  }

  p.platform.SendText(platformTypes.SendTextInput{Text: text}, &input.EventMessage)
}

func (c GenCommand) GetKey() string {
	return c.key
}

func NewGenCommand(llmGenerator llm.LLMGenerator) *GenCommand {
  return &GenCommand{key: "gen", llmGenerator: llmGenerator}
}

