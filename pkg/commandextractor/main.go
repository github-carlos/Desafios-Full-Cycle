package commandextractor

import (
	"strings"

	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
)

type CommandInput struct {
  Command string
  Payload interface{}
  Sender types.JID
}

const commandPrefix = '!'

func Extract(eventMessage *events.Message) (*CommandInput, error) {
	text := extractText(eventMessage)

	if text == "" {
		return &CommandInput{}, nil
	}

  if text[0] != commandPrefix {
    return &CommandInput{}, nil
  }

	return &CommandInput{
    Command: text,
    Payload: text,
    Sender: eventMessage.Info.Sender,
  }, nil
}

func extractText(eventMessage *events.Message) string {
	message := eventMessage.Message
	text := message.GetConversation()

	if text == "" {
		text = message.ImageMessage.GetCaption()
	}

	if text == "" {
		text = message.VideoMessage.GetCaption()
	}

	if text == "" {
		text = message.ExtendedTextMessage.GetText()
	}
	return strings.Trim(text, " ")
}

