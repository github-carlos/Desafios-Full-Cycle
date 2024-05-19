package commands

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"trevas-bot/pkg/commandextractor"
	"trevas-bot/pkg/platform"
	platformTypes "trevas-bot/pkg/platform/types"
)
type MemeResponse struct {
	PostLink  string `json:"postLink"`
	Subreddit string `json:"subreddit"`
	Title     string `json:"title"`
	URL       string `json:"url"`
	Nsfw      bool   `json:"nsfw"`
	Spoiler   bool   `json:"spoiler"`
	Author    string `json:"author"`
	Ups       int    `json:"ups"`
}

type MemeCommand struct {
	key      string
	platform platform.WhatsAppIntegration
}

func (p MemeCommand) Handler(input commandextractor.CommandInput) {
	fmt.Println("Running Meme Command")

	// Define the endpoint
	url := "https://meme-api.com/gimme/MemesBR"

	// Make the HTTP GET request
	resp, err := http.Get(url)

  if err != nil {
    go p.platform.SendReply("Não foi possível buscar um meme...", &input.EventMessage)
    return
  }

	defer resp.Body.Close()

	// Read the response body
  body, err := io.ReadAll(resp.Body)
	if err != nil {
    return
	}

	// Parse the JSON response
	var memeResponse MemeResponse
	err = json.Unmarshal(body, &memeResponse)
	if err != nil {
    return
	}

	// Print the data
	fmt.Printf("Title: %s\n", memeResponse.Title)
	fmt.Printf("URL: %s\n", memeResponse.URL)
	fmt.Printf("Author: %s\n", memeResponse.Author)
	fmt.Printf("Upvotes: %d\n", memeResponse.Ups)

  p.platform.SendText(platformTypes.SendTextInput{Text: memeResponse.Title}, &input.EventMessage)



}

func (c MemeCommand) GetKey() string {
	return c.key
}

func NewMemeCommand() *MemeCommand {
	return &MemeCommand{key: "meme"}
}
