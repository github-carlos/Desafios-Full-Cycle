package platform

import (
	"context"
	"fmt"

	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
)

type WhatsAppIntegration struct { }

func NewWhatsAppIntegration() *WhatsAppIntegration {
  return &WhatsAppIntegration{}
}

var Client *whatsmeow.Client
func SetWhatsAppClient(client *whatsmeow.Client) {
  Client = client
}

func (w WhatsAppIntegration) SendText(text string, eventMessage *events.Message) error {
  fmt.Println("Sending Message")
  return nil
}


func (w WhatsAppIntegration) SendReply(text string, eventMessage *events.Message) error {
  fmt.Println("Sending Reply", eventMessage)

	var msg = &waProto.Message{
		ExtendedTextMessage: &waProto.ExtendedTextMessage{
			Text: proto.String(text),
			ContextInfo: &waProto.ContextInfo{
				StanzaId:      proto.String(eventMessage.Info.ID),
				Participant:   proto.String(eventMessage.Info.Sender.ToNonAD().String()),
				QuotedMessage: eventMessage.Message,
			},
		},
	}

	_, err := Client.SendMessage(context.Background(), eventMessage.Info.Chat, msg)

	if err != nil {
		fmt.Printf("Error sending message: %v", err)
	}

	fmt.Println("Message sent:", text)
  return nil
}
