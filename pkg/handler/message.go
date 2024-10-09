package handler

import (
	"fmt"
	"trevas-bot/pkg/commandeval"
	"trevas-bot/pkg/commandextractor"
	"trevas-bot/pkg/store"

	"go.mau.fi/whatsmeow/types/events"
)

var CommandEval = commandeval.NewCommandEval()
var appDatabase, _ = store.NewAppDatabase()

func MessageHandler(eventMessage *events.Message) {
  channelJID := eventMessage.Info.Chat.String()
  jid := eventMessage.Info.Sender.String()
  name := eventMessage.Info.PushName
  isGroup := eventMessage.Info.IsGroup
  messageType := eventMessage.Info.Type
  commandInput, err := commandextractor.Extract(eventMessage)

	commandInput, err = commandextractor.Extract(eventMessage)

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

  input := store.SaveMessageInput{
    JID:          jid,
    Name:         name,
    ChannelJID:  channelJID,
    IsGroup: isGroup,
    Message:      commandInput.Text,
    MessageType:  messageType,
    Command:      commandInput.Command, 
    Timestamp:    eventMessage.Info.Timestamp.String(),
  }

  appDatabase.SaveMessage(input)

	if commandInput.Command == "" {
		return
	}

	fmt.Println("Message:", commandInput.Command, "Phone:", commandInput.EventMessage.Info.Sender)

	if err := CommandEval.Handle(commandInput); err != nil {
		fmt.Println(err)
	}
}
