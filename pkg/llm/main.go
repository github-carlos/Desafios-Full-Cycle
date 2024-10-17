package llm
import (
	"context"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
  "fmt"
)


type LLMGenerator interface {
  Complete(prompt string) (string, error)
}

func NewGeminiGenerator() *GeminiGenerator {
  return &GeminiGenerator{}
}

type GeminiGenerator struct { }

func (g *GeminiGenerator) Complete(prompt string) (string, error) {
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

  resp, err := model.GenerateContent(ctx, genai.Text(prompt))
  if err != nil || len(resp.Candidates) == 0 {
    fmt.Println("Error", err)
    return "", err
  }
  
  return getResponse(resp), nil
}

func getResponse(resp *genai.GenerateContentResponse) string {
	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
        return fmt.Sprintf("%s", part)
			}
		}
	}
  return ""
}
