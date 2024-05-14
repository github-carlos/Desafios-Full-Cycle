package converter

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"time"
)

func Img2Webp(image []byte) ([]byte, error) {

	now := time.Now()
	inputPath := fmt.Sprintf("temp/%d.gif", now.Unix())
	outputPath := fmt.Sprintf("temp/%d.webp", now.Unix())

  err := os.WriteFile(inputPath, image, 0644)

  if err != nil {
    fmt.Println("Error saving input to use in converter", err)
    return nil, nil
  }

	cmd := exec.Command("ffmpeg", "-i", inputPath, "-vcodec", "libwebp", "-preset", "default", "-loop", "0", "-an", "-qscale:v", "20", "-t", "00:00:10", "-compression_level", "100", "-vf", "scale='min(320,iw)':min'(320,ih)':force_original_aspect_ratio=decrease,fps=20, pad=320:320:-1:-1:color=white@0.0, split [a][b]; [a] palettegen=reserve_transparent=on:transparency_color=ffffff [p]; [b][p] paletteuse", outputPath)
	if err := cmd.Run(); err != nil {
		fmt.Println("Erro ao converter o arquivo:", err)
		_ = os.Remove(inputPath)
		_ = os.Remove(outputPath)
		return nil, errors.New("Ocorreu um erro ao processar sua figurinha... por favor, tente novamente")
	}

  stickerBytes, err := os.ReadFile(outputPath)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return nil, errors.New("Ocorreu um erro ao processar sua figurinha... por favor, tente novamente")
	}

  // TODO por metadatas no sticker

  _ = os.Remove(inputPath)
  _ = os.Remove(outputPath)

	return stickerBytes, nil
}
