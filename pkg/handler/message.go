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

  fmt.Println(commandInput.Text)
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
  appDatabase.SaveUser(eventMessage)

	if commandInput.Command == "" {
		return
	}

  commandInput.Store = appDatabase

	fmt.Println("Message:", commandInput.Command, "Phone:", commandInput.EventMessage.Info.Sender)

  if (appDatabase.CheckIfUserIsBlocked(eventMessage)) {
    fmt.Println("Ignoring Command")
    return
  }

	if err := CommandEval.Handle(commandInput); err != nil {
		fmt.Println(err)
	}
}
