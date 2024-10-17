package commands

import (
	"fmt"
	"strings"
	"trevas-bot/pkg/commandextractor"
	"trevas-bot/pkg/platform"
	"trevas-bot/pkg/store"
	"trevas-bot/pkg/utils"
)

type IgnoreCommand struct {
	key      string
	platform platform.WhatsAppIntegration
}

func (p IgnoreCommand) Handler(input commandextractor.CommandInput) {
	fmt.Println("Running Ignore Command")

	if !utils.IsBotOwner(&input.EventMessage) {
		go p.platform.SendReply("Só o meu dono pode usar esse comando :c", &input.EventMessage)
		return
	}

	commands := strings.Split(input.Payload, " ")

	if len(commands) < 2 {
		return
	}

	command := commands[0]
	target := commands[1]

	number := strings.Split(target, "@")[1]
	switch command {
	case "add":
		add(target, input.Store, number)
		go p.platform.SendReply("Bovino ignorado com sucesso.", &input.EventMessage)
		return
	case "remove":
		remove(target, input.Store, number)
		go p.platform.SendReply("Bovino removido da lista de ignorados.", &input.EventMessage)
		return
	}
}

func add(target string, db *store.AppDatabase, number string) {
	fmt.Println("Running command Ignore")
	db.BlockUserByNumber(number)
}

func remove(target string, db *store.AppDatabase, number string) {
	fmt.Println("Running remove Ignore")
	db.UnblockUserByNumber(number)
}

func (c IgnoreCommand) GetKey() string {
	return c.key
}

func NewIgnoreCommand() *IgnoreCommand {
	return &IgnoreCommand{key: "ignore"}
}
