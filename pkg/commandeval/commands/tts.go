package commands

import (
	"context"
	"fmt"
	"log"
	"trevas-bot/pkg/commandextractor"
	"trevas-bot/pkg/platform"

	texttospeech "cloud.google.com/go/texttospeech/apiv1"
	"cloud.google.com/go/texttospeech/apiv1/texttospeechpb"
)

var client *texttospeech.Client

type TTSCommand struct {
	key      string
	platform platform.WhatsAppIntegration
}

func (p TTSCommand) Handler(input commandextractor.CommandInput) {
	fmt.Println("Running TTS Command")

  text := input.Payload

  if len(text) > 300 {
    p.platform.SendReply("Texto muito grande, amigo. Máximo de 300 caractéres 👍", &input.EventMessage)
    return
  }

	// Perform the text-to-speech request on the text input with the selected
	// voice parameters and audio file type.
	req := texttospeechpb.SynthesizeSpeechRequest{
		// Set the text input to be synthesized.
		Input: &texttospeechpb.SynthesisInput{
      InputSource: &texttospeechpb.SynthesisInput_Text{Text: text},
		},
		// Build the voice request, select the language code ("en-US") and the SSML
		// voice gender ("neutral").
		Voice: &texttospeechpb.VoiceSelectionParams{
			LanguageCode: "pt-BR",
			SsmlGender:   texttospeechpb.SsmlVoiceGender_MALE,
		},
		// Select the type of audio file you want returned.
		AudioConfig: &texttospeechpb.AudioConfig{
			AudioEncoding: texttospeechpb.AudioEncoding_OGG_OPUS,
      Pitch: -10,
		},
	}

	ctx := context.Background()
	resp, err := client.SynthesizeSpeech(ctx, &req)
	if err != nil {
		log.Fatal(err)
	}

  err = p.platform.SendAudio(platform.SendAudioInput{AudioBytes: resp.AudioContent}, &input.EventMessage)

  if err != nil {
    p.platform.SendReply(err.Error(), &input.EventMessage)
  }
}

func (c TTSCommand) GetKey() string {
	return c.key
}

func NewTTSCommand() *TTSCommand {
	ctx := context.Background()

  var err error
  client, err = texttospeech.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// defer client.Close()
	return &TTSCommand{key: "tts"}
}
