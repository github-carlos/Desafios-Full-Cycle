package platform

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"image/jpeg"
	"net/http"
	"os"
	"time"
	"trevas-bot/pkg/platform/types"

	"github.com/nfnt/resize"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/proto/waE2E"
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
	BolaReacton       string = "⚽"
	BullReaction      string = "🐂"
)

type WhatsAppIntegration struct{}

func NewWhatsAppIntegration() *WhatsAppIntegration {
	return &WhatsAppIntegration{}
}

var Client *whatsmeow.Client

func SetWhatsAppClient(client *whatsmeow.Client) {
	Client = client
}

func (w WhatsAppIntegration) SendReply(text string, eventMessage *events.Message) {
	fmt.Println("Sending Reply", eventMessage)

	var msg = &waProto.Message{
		ExtendedTextMessage: &waProto.ExtendedTextMessage{
			Text: proto.String(text),
			ContextInfo: &waProto.ContextInfo{
				StanzaID:      proto.String(eventMessage.Info.ID),
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
}

func (w WhatsAppIntegration) SendText(input types.SendTextInput, eventMessage *events.Message) {
	fmt.Println("Sending Text", eventMessage)

	fmt.Println(eventMessage.Info.Sender.ToNonAD().String())

	var mentions []string

	for _, mention := range input.Mentions {
		mentions = append(mentions, fmt.Sprintf("%s@s.whatsapp.net", mention))
	}

	var msg = &waProto.Message{
		ExtendedTextMessage: &waProto.ExtendedTextMessage{
			Text: proto.String(input.Text),
			ContextInfo: &waProto.ContextInfo{
				MentionedJID: mentions,
			},
		},
	}

	_, err := Client.SendMessage(context.Background(), eventMessage.Info.Chat, msg)

	if err != nil {
		fmt.Printf("Error sending message: %v", err)
	}

	fmt.Println("Message sent:", input)
}

func (w WhatsAppIntegration) SendSticker(stickerBytes []byte, animated bool, eventMessage *events.Message, reply bool) error {

	if len(stickerBytes) > 1024*1024 {
		return errors.New("O arquivo enviado é muito grande.")
	}

	uploadedSticker, err := Client.Upload(context.Background(), stickerBytes, whatsmeow.MediaImage)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return errors.New("Ocorreu um erro ao fazer o upload da imagem... por favor, tente novamente.")
	}

	var contextInfo waProto.ContextInfo

	if reply {
		contextInfo = waProto.ContextInfo{
			StanzaID:      proto.String(eventMessage.Info.ID),
			Participant:   proto.String(eventMessage.Info.Sender.ToNonAD().String()),
			QuotedMessage: eventMessage.Message,
		}
	}

	msgToSend := &waProto.Message{
		StickerMessage: &waProto.StickerMessage{
			URL:           proto.String(uploadedSticker.URL),
			DirectPath:    proto.String(uploadedSticker.DirectPath),
			MediaKey:      uploadedSticker.MediaKey,
			IsAnimated:    proto.Bool(animated),
			IsAvatar:      proto.Bool(false),
			Mimetype:      proto.String("image/webp"),
			FileEncSHA256: uploadedSticker.FileEncSHA256,
			FileSHA256:    uploadedSticker.FileSHA256,
			FileLength:    proto.Uint64(uploadedSticker.FileLength),
			StickerSentTS: proto.Int64(time.Now().Unix()),
			ContextInfo:   &contextInfo,
		},
	}

	_, err = Client.SendMessage(context.Background(), eventMessage.Info.Chat, msgToSend)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return errors.New("Ocorreu um erro ao enviar a figurinha... por favor, tente novamente.")
	}

	return nil
}

func (w WhatsAppIntegration) SendImg(input types.SendImageInput, eventMessage *events.Message) error {
	uploadedImg, err := Client.Upload(context.Background(), input.Image, whatsmeow.MediaImage)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return errors.New("Ocorreu um erro ao fazer o upload do imagem... por favor, tente novamente.")
	}

	// decode jpeg into image.Image
	decodedImg, err := jpeg.Decode(bytes.NewReader(input.Image))

	if err != nil {
		fmt.Println(err)
	}

	m := resize.Thumbnail(72, 72, decodedImg, resize.Lanczos3)
	outPath := fmt.Sprintf("temp/images/%d.jpg", time.Now().Unix())
	out, err := os.Create(outPath)
	if err != nil {
		fmt.Println(err)
	}

	defer func() {
		out.Close()
		os.Remove(outPath)
	}()

	// write new image to file
	jpeg.Encode(out, m, nil)

	thumbnailBytes, err := os.ReadFile(outPath)

	msgToSend := &waProto.Message{
		ImageMessage: &waProto.ImageMessage{
			URL:           proto.String(uploadedImg.URL),
			DirectPath:    proto.String(uploadedImg.DirectPath),
			MediaKey:      uploadedImg.MediaKey,
			Mimetype:      proto.String(http.DetectContentType(input.Image)),
			FileEncSHA256: uploadedImg.FileEncSHA256,
			FileSHA256:    uploadedImg.FileSHA256,
			FileLength:    proto.Uint64(uploadedImg.FileLength),
			Caption:       &input.Caption,
			JPEGThumbnail: thumbnailBytes,
			ContextInfo: &waProto.ContextInfo{
				StanzaID:      proto.String(eventMessage.Info.ID),
				Participant:   proto.String(eventMessage.Info.Sender.ToNonAD().String()),
				QuotedMessage: eventMessage.Message,
			},
		},
	}

	_, err = Client.SendMessage(context.Background(), eventMessage.Info.Chat, msgToSend)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return errors.New("Ocorreu um erro ao enviar a imagem... por favor, tente novamente.")
	}
	fmt.Println("Imagem enviada com sucesso!")
	return nil
}

func (w WhatsAppIntegration) SendVideo(input types.SendVideoInput, eventMessage *events.Message) error {

	if len(input.VideoBytes) > 50*1024*1024 {
		return errors.New("Video muito grande para fazer Upload.")
	}

	uploadedVideo, err := Client.Upload(context.Background(), input.VideoBytes, whatsmeow.MediaVideo)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return errors.New("Ocorreu um erro ao fazer o upload do video... por favor, tente novamente.")
	}

	msgToSend := &waProto.Message{
		VideoMessage: &waProto.VideoMessage{
			URL:           proto.String(uploadedVideo.URL),
			DirectPath:    proto.String(uploadedVideo.DirectPath),
			MediaKey:      uploadedVideo.MediaKey,
			Mimetype:      proto.String("video/mp4"),
			FileEncSHA256: uploadedVideo.FileEncSHA256,
			FileSHA256:    uploadedVideo.FileSHA256,
			FileLength:    proto.Uint64(uploadedVideo.FileLength),
			Caption:       &input.Caption,
			JPEGThumbnail: input.Thumbnail,
			ContextInfo: &waProto.ContextInfo{
				StanzaID:      proto.String(eventMessage.Info.ID),
				Participant:   proto.String(eventMessage.Info.Sender.ToNonAD().String()),
				QuotedMessage: eventMessage.Message,
			},
		},
	}

	_, err = Client.SendMessage(context.Background(), eventMessage.Info.Chat, msgToSend)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return errors.New("Ocorreu um erro ao enviar o video... por favor, tente novamente.")
	}
	fmt.Println("video enviado com sucesso!")
	return nil
}

type SendAudioInput struct {
	AudioBytes []byte
}

func (w WhatsAppIntegration) SendAudio(input SendAudioInput, eventMessage *events.Message) error {

	if len(input.AudioBytes) > 50*1024*1024 {
		return errors.New("Audio muito grande para fazer Upload.")
	}

	uploadedAudio, err := Client.Upload(context.Background(), input.AudioBytes, whatsmeow.MediaAudio)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return errors.New("Ocorreu um erro ao fazer o upload do audio... por favor, tente novamente.")
	}

	msgToSend := &waProto.Message{
		AudioMessage: &waProto.AudioMessage{
			URL:           proto.String(uploadedAudio.URL),
			DirectPath:    proto.String(uploadedAudio.DirectPath),
			MediaKey:      uploadedAudio.MediaKey,
			Mimetype:      proto.String("audio/ogg; codecs=opus"),
			FileEncSHA256: uploadedAudio.FileEncSHA256,
			FileSHA256:    uploadedAudio.FileSHA256,
			FileLength:    proto.Uint64(uploadedAudio.FileLength),
			ContextInfo: &waProto.ContextInfo{
				StanzaID:      proto.String(eventMessage.Info.ID),
				Participant:   proto.String(eventMessage.Info.Sender.ToNonAD().String()),
				QuotedMessage: eventMessage.Message,
			},
		},
	}

	_, err = Client.SendMessage(context.Background(), eventMessage.Info.Chat, msgToSend)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return errors.New("Ocorreu um erro ao enviar o audio... por favor, tente novamente.")
	}
	fmt.Println("audio enviado com sucesso!")
	return nil
}

func (w WhatsAppIntegration) SendReaction(eventMessage *events.Message, reaction string) {
	r := Client.BuildReaction(eventMessage.Info.Chat, eventMessage.Info.Sender, eventMessage.Info.ID, reaction)
	_, err := Client.SendMessage(context.Background(), eventMessage.Info.Chat, r)
	if err != nil {
		fmt.Println("Error sending reaction:", err)
	}
}

func (w WhatsAppIntegration) ExtractMediaBytes(eventMessage *events.Message) ([]byte, error) {

	imageMedia := GetImageMessage(eventMessage)
	var videoMedia *waProto.VideoMessage

	if imageMedia == nil {
		videoMedia = GetVideoMessage(eventMessage)
	}

	if imageMedia == nil && videoMedia == nil {
		return nil, errors.New("Você precisa mandar uma imagem ou um vídeo para fazer uma figurinha")
	}

	var downloadedMedia []byte

	if imageMedia != nil {
		downloadedMedia, _ = Client.Download(imageMedia)
	}

	if videoMedia != nil {
		downloadedMedia, _ = Client.Download(videoMedia)
	}

	if downloadedMedia == nil {
		return nil, errors.New("Erro ao tentar baixar mídia")
	}
	return downloadedMedia, nil
}

func (w WhatsAppIntegration) ExtractStickerMediaBytes(eventMessage *events.Message) ([]byte, error) {
	stickerMedia := GetStickerMessage(eventMessage)

	if stickerMedia == nil {
		return nil, errors.New("Sticker Message not found")
	}

	downloadMedia, err := Client.Download(stickerMedia)

	if err != nil {
		return nil, errors.New("Error trying to Download Sticker")
	}

	return downloadMedia, nil
}

func (w WhatsAppIntegration) GetParticipantsOfGroup(msg *events.Message) ([]string, error) {
	groupInfo, err := Client.GetGroupInfo(msg.Info.Chat)

	if err != nil {
		fmt.Println("Error", err)
		return nil, errors.New("Error getting participants list")
	}

	var participants []string

	for _, jidParticipant := range groupInfo.Participants {
		if jidParticipant.JID.User != Client.Store.ID.User {
			participants = append(participants, jidParticipant.JID.User)
		}
	}

	return participants, nil
}

func GetImageMessage(msg *events.Message) (imageMsg *waProto.ImageMessage) {
	msg.SourceWebMsg.GetMediaData()
	if msg.Message.ImageMessage != nil {
		imageMsg = msg.Message.ImageMessage
	} else if msg.Message.ExtendedTextMessage != nil &&
		msg.Message.ExtendedTextMessage.ContextInfo != nil &&
		msg.Message.ExtendedTextMessage.ContextInfo.QuotedMessage != nil &&
		msg.Message.ExtendedTextMessage.ContextInfo.QuotedMessage.ImageMessage != nil {
		imageMsg = msg.Message.ExtendedTextMessage.ContextInfo.QuotedMessage.ImageMessage
	} else {
		return nil
	}
	return imageMsg
}

func GetVideoMessage(msg *events.Message) (videoMsg *waProto.VideoMessage) {
	if msg.Message.VideoMessage != nil {
		videoMsg = msg.Message.VideoMessage
	} else if msg.Message.ExtendedTextMessage != nil &&
		msg.Message.ExtendedTextMessage.ContextInfo != nil &&
		msg.Message.ExtendedTextMessage.ContextInfo.QuotedMessage != nil &&
		msg.Message.ExtendedTextMessage.ContextInfo.QuotedMessage.VideoMessage != nil {
		videoMsg = msg.Message.ExtendedTextMessage.ContextInfo.QuotedMessage.VideoMessage
	} else {
		return nil
	}
	return videoMsg
}

func GetStickerMessage(msg *events.Message) (stickerMsg *waProto.StickerMessage) {
	if msg.Message.StickerMessage != nil {
		stickerMsg = msg.Message.StickerMessage
	} else if msg.Message.ExtendedTextMessage != nil &&
		msg.Message.ExtendedTextMessage.ContextInfo != nil &&
		msg.Message.ExtendedTextMessage.ContextInfo.QuotedMessage != nil &&
		msg.Message.ExtendedTextMessage.ContextInfo.QuotedMessage.StickerMessage != nil {
		stickerMsg = msg.Message.ExtendedTextMessage.ContextInfo.QuotedMessage.StickerMessage
	} else {
		return nil
	}
	return stickerMsg
}
