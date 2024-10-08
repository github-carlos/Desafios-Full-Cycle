package types

import (
	waTypes "go.mau.fi/whatsmeow/types"
	waProto "go.mau.fi/whatsmeow/proto/waE2E"
)

type Message struct {
	Chat          waTypes.JID
	StanzaID      string
	Participant   string
	QuotedMessage *waProto.Message
}

type SendTextInput struct {
	Text     string
	Mentions []string
}

type SendImageInput struct {
	Image   []byte
	Caption string
	Message Message
}

type SendVideoInput struct {
	VideoBytes []byte
	Thumbnail  []byte
	Caption    string
	Message    Message
}
