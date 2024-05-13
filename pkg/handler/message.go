package handler

import (
	"fmt"
	"trevas-bot/pkg/commandeval"

	"go.mau.fi/whatsmeow/types/events"
)

var CommandEval = commandeval.NewCommandEval()

func MessageHandler(eventMessage *events.Message) {
  fmt.Println("Mensagem Recebida:", eventMessage)
  text := extractText(eventMessage)
  fmt.Println(text)

  if err := CommandEval.Handle(text); err != nil {
    fmt.Println(err)
  }
}

func extractText(eventMessage *events.Message) string {
  message := eventMessage.Message;
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

  return text
}
