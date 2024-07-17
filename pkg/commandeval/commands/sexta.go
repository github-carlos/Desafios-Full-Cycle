package commands

import (
	"fmt"
	"os"
	"trevas-bot/pkg/commandextractor"
	"trevas-bot/pkg/converter"
	"trevas-bot/pkg/platform"
)

type SextaCommand struct {
	key      string
	platform platform.WhatsAppIntegration
}

func (p SextaCommand) Handler(input commandextractor.CommandInput) {
	fmt.Println("Running Sexta Command")

	numberOfStickers := 5

	for i := 0; i < numberOfStickers; i++ {
		filePath := fmt.Sprintf("assets/group/sexta/%d.jpeg", i)
		img, err := os.ReadFile(filePath)

		if err != nil {
			fmt.Println("Failing getting Saude image")
			return
		}

		webp, err := converter.Img2Webp(img)

		if err != nil {
			fmt.Println("Failing converting img to webp")
			return
		}

		err = p.platform.SendSticker(webp, false, &input.EventMessage, false)

		if err != nil {
			fmt.Println("Failing sending message")
			return
		}
	}

}

func (c SextaCommand) GetKey() string {
	return c.key
}

func NewSextaCommand() *SextaCommand {
	return &SextaCommand{key: "sexta"}
}
