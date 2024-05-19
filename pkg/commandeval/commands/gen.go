package commands

import (
	"context"
	"fmt"
	"trevas-bot/pkg/commandextractor"
	"trevas-bot/pkg/platform"
	platformTypes "trevas-bot/pkg/platform/types"

	"github.com/google/generative-ai-go/genai"
  "google.golang.org/api/option"
)

type GenCommand struct {
	key      string
	platform platform.WhatsAppIntegration
}

func (p GenCommand) Handler(input commandextractor.CommandInput) {


  if input.Payload == "" {
    p.platform.SendReply("Envie um texto junto ao comando", &input.EventMessage)
    return
  }

	ctx := context.Background()
	// Access your API key as an environment variable (see "Set up your API key" above)
	client, err := genai.NewClient(ctx, option.WithAPIKey("AIzaSyBztQVIsaoX96ZbDDkvo0_cc6UGVnJwGqs"))
	if err != nil {
		fmt.Println(err)
	}
	defer client.Close()

	// For text-only input, use the gemini-pro model
  model := client.GenerativeModel("gemini-1.0-pro")
  model.SetTemperature(1)
  // model.SetMaxOutputTokens(100)
  // model.SystemInstruction = &genai.Content{
  //   Parts: []genai.Part{genai.Text("Você deve gerar respostas de no máximo 300 caracteres")},
  // }

  prompt := fmt.Sprintf("Gere uma resposta de no máximo 300 caracteres\n %s", input.Payload)

  resp, err := model.GenerateContent(ctx, genai.Text(prompt))
  if err != nil || len(resp.Candidates) == 0 {
    fmt.Println("Error", err)
    p.platform.SendReply("Tente novamente mais tarde", &input.EventMessage)
    return;
  }
  
  text := getResponse(resp)

  p.platform.SendText(platformTypes.SendTextInput{Text: text}, &input.EventMessage)
}

func (c GenCommand) GetKey() string {
	return c.key
}

func NewGenCommand() *GenCommand {
	return &GenCommand{key: "gen"}
}

func getResponse(resp *genai.GenerateContentResponse) string {
	for _, cand := range resp.Candidates {
    fmt.Println("Candidates", cand)
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
        return fmt.Sprintf("%s", part)
			}
		}
	}
  return ""
}
