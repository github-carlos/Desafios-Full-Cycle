package commandextractor

import (
	"strings"

	"go.mau.fi/whatsmeow/types/events"
)

type CommandInput struct {
	Command      string
	Payload     string 
	EventMessage events.Message
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

	command := extractCommand(text)
	payload := extractPayload(text)

	return &CommandInput{
		Command:      command,
		Payload:      payload,
		EventMessage: *eventMessage,
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

func extractCommand(text string) string {
	splitedText := strings.Split(text, " ")
	return splitedText[0][1:]
}

func extractPayload(text string) string {
	splitedText := strings.Split(text, " ")
	return strings.Join(splitedText[1:], " ")
}
