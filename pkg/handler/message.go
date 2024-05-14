package handler

import (
	"fmt"
	"trevas-bot/pkg/commandeval"
	"trevas-bot/pkg/commandextractor"

	"go.mau.fi/whatsmeow/types/events"
)

var CommandEval = commandeval.NewCommandEval()

func MessageHandler(eventMessage *events.Message) {
	fmt.Println("Mensagem Recebida:", eventMessage)

	commandInput, err := commandextractor.Extract(eventMessage)

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	if commandInput.Command == "" {
		return
	}

	fmt.Println("Message:", commandInput.Command, "Phone:", commandInput.EventMessage.Info.Sender)

	if err := CommandEval.Handle(commandInput); err != nil {
		fmt.Println(err)
	}
}
