package platform

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
)

const (
	ErrorReaction     string = "❌"
	ForbiddenReaction string = "🚫"
	SuccessReaction   string = "✅"
	LoadingReaction   string = "⏳"
	ConfigReaction    string = "⚙️"
	PingReaction      string = "🏓"
	LoveReaction      string = "❤️"
	LikeReaction      string = "👍"
	DislikeReaction   string = "👎"
  BolaReacton string = "⚽"
)

type WhatsAppIntegration struct{}

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

func (w WhatsAppIntegration) SendSticker(stickerBytes []byte, animated bool, eventMessage *events.Message) error {

	uploadedSticker, err := Client.Upload(context.Background(), stickerBytes, whatsmeow.MediaImage)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return errors.New("Ocorreu um erro ao fazer o upload da imagem... por favor, tente novamente.")
	}

	msgToSend := &waProto.Message{
		StickerMessage: &waProto.StickerMessage{
			Url:           proto.String(uploadedSticker.URL),
			DirectPath:    proto.String(uploadedSticker.DirectPath),
			MediaKey:      uploadedSticker.MediaKey,
			IsAnimated:    proto.Bool(animated),
			IsAvatar:      proto.Bool(false),
			Mimetype:      proto.String("image/webp"),
			FileEncSha256: uploadedSticker.FileEncSHA256,
			FileSha256:    uploadedSticker.FileSHA256,
			FileLength:    proto.Uint64(uploadedSticker.FileLength),
			StickerSentTs: proto.Int64(time.Now().Unix()),
		},
	}

	_, err = Client.SendMessage(context.Background(), eventMessage.Info.Chat, msgToSend)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return errors.New("Ocorreu um erro ao enviar a figurinha... por favor, tente novamente.")
	}

  return nil
}

func (w WhatsAppIntegration) SendImg(imgBytes []byte, animated bool, eventMessage *events.Message) error {

	uploadedImg, err := Client.Upload(context.Background(), imgBytes, whatsmeow.MediaImage)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return errors.New("Ocorreu um erro ao fazer o upload do imagem... por favor, tente novamente.")
	}

	msgToSend := &waProto.Message{
		ImageMessage: &waProto.ImageMessage{
			Url:           proto.String(uploadedImg.URL),
			DirectPath:    proto.String(uploadedImg.DirectPath),
			MediaKey:      uploadedImg.MediaKey,
			Mimetype:      proto.String("image/webp"),
			FileEncSha256: uploadedImg.FileEncSHA256,
			FileSha256:    uploadedImg.FileSHA256,
			FileLength:    proto.Uint64(uploadedImg.FileLength),
		},
	}

	_, err = Client.SendMessage(context.Background(), eventMessage.Info.Chat, msgToSend)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return errors.New("Ocorreu um erro ao enviar a imagem... por favor, tente novamente.")
	}

  return nil
}

func (w WhatsAppIntegration) SendReaction(eventMessage *events.Message, reaction string) {
	r := Client.BuildReaction(eventMessage.Info.Chat, eventMessage.Info.Sender, eventMessage.Info.ID, reaction)
	_, err := Client.SendMessage(context.Background(), eventMessage.Info.Chat, r)
	if err != nil {
		fmt.Println("Error sending reaction:", err)
	}
}
