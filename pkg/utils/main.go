package utils

import (
	"strings"

	"go.mau.fi/whatsmeow/types/events"
)

var BOT_OWNER string = "556292147541"

func IsBotOwner(message *events.Message) bool {
  jid := message.Info.Sender.String()

  return strings.Contains(jid, BOT_OWNER)
}
