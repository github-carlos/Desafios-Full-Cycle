package commands

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
	"trevas-bot/pkg/commandextractor"
	"trevas-bot/pkg/platform"

	"github.com/gographics/imagick/imagick"
)

var fonts = []string{
  "Chalkboard",
  "Trattatello",
  "Noteworthy",
  "Maker Felt",
  "Luminari",
  "Impact",
  "Herculanum",
  "Chalkduster",
}

type FigCommand struct {
	key      string
	platform platform.WhatsAppIntegration
}

func (p FigCommand) Handler(input commandextractor.CommandInput) {
	// Inicialize a biblioteca
	imagick.Initialize()
	defer imagick.Terminate()

  text := formatText(input.Payload);

	// Configure o tamanho do GIF
	width := getLengthBiggestLine(text) * 50
	height := len(strings.Split(text, "\n")) * 100

  if width < 300 {
    width = 300
  }

  if height < 300 {
    height = 300
  }

	frames := 10
	var delay uint = 10 // 10 centésimos de segundo

	// Criação do GIF
	anim := imagick.NewMagickWand()
	defer anim.Destroy()

  font := fmt.Sprintf("%s-Bold", fonts[rand.Intn(len(fonts))])

	for i := 0; i < frames; i++ {
		// Criação de um frame individual
		frame := imagick.NewMagickWand()
		defer frame.Destroy()

		// Plano de fundo
		bg := imagick.NewPixelWand()
		bg.SetColor("transparent")
		frame.NewImage(uint(width), uint(height), bg)

		// Configuração do texto
		dw := imagick.NewDrawingWand()
		pixel := imagick.NewPixelWand()
		defer dw.Destroy()
		defer pixel.Destroy()

		r, g, b := randomBrightColor()
		pixel.SetColor(fmt.Sprintf("rgb(%d,%d,%d)", r, g, b))
		dw.SetFillColor(pixel)

    // fontSize := float64(getLengthBiggestLine(text))
		dw.SetFontSize(80)
    dw.SetGravity(imagick.GRAVITY_CENTER)


		dw.SetFont(font) // Defina a fonte que tenha a variação negrito (por exemplo, Arial-Bold)
		dw.Annotation(0, 0, text)

		// Desenhar o texto no frame
		frame.DrawImage(dw)

		// Configure o tempo de exibição do frame
		frame.SetImageDelay(delay)

		// Adicione o frame ao GIF animado
		anim.AddImage(frame)
	}

	anim.SetImageFormat("webp")

  outputFileName := fmt.Sprintf("temp/%d.webp", time.Now().Unix())
  anim.WriteImages(outputFileName, true)

  sticker, _ := os.ReadFile(outputFileName)

  defer os.Remove(outputFileName)

  p.platform.SendSticker(sticker, false, &input.EventMessage, true)
}

func (c FigCommand) GetKey() string {
	return c.key
}

func NewFigCommand() *FigCommand {
	return &FigCommand{key: "fig"}
}

func formatText(text string) string {
  splitedText := strings.Split(text, " ")

  newText := ""

  for index := 0; index < len(splitedText); index++ {
    word := splitedText[index]
    nextWord := ""

    if index < len(splitedText) - 1 {
      nextWord = splitedText[index + 1]
    }

    if len(word) < 5 && len(nextWord) < 5 {
      newText += word + " " + nextWord
      index++;
    } else {
      newText += word
    }

    if nextWord != "" {
      newText += "\n"
    }
  }

  return newText
}

func randomBrightColor() (int, int, int) {
	rand.Seed(time.Now().UnixNano()) // Garante cores diferentes a cada execução
	min := 128                       // Define o brilho mínimo
	return min + rand.Intn(128), min + rand.Intn(128), min + rand.Intn(128)
}

func getLengthBiggestLine(text string) int {

  splittedText := strings.Split(text, "\n")

  size := len(splittedText[0])

  for index := 1; index < len(splittedText); index++ {
    if len(splittedText[index]) > size {
      size = len(splittedText[index])
    }
  }

  return size
}
