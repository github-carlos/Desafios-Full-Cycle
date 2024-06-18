package converter

import (
	"bytes"
	"errors"
	"fmt"
	"image/png"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/nfnt/resize"
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

func Webm2Mp4(webm []byte) ([]byte, error) {
	now := time.Now()
	inputPath := fmt.Sprintf("temp/%d-temp.webm", now.Unix())
	outputPath := fmt.Sprintf("temp/%d-temp.mp4", now.Unix())

	err := os.WriteFile(inputPath, webm, 0644)

	if err != nil {
		fmt.Println("Error saving input to use in converter", err)
		return nil, nil
	}

	// ffmpeg -i 272422a\ \[272422a\].webm  -strict experimental video.mp4
	cmd := exec.Command("ffmpeg", "-i", inputPath, outputPath)

	if err := cmd.Run(); err != nil {
		fmt.Println("Erro ao converter o arquivo:", err)
		_ = os.Remove(inputPath)
		_ = os.Remove(outputPath)
		return nil, errors.New("Ocorreu um erro ao processar seu video... por favor, tente novamente")
	}

	mp4Bytes, err := os.ReadFile(outputPath)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return nil, errors.New("Ocorreu um erro ao processar seu video... por favor, tente novamente")
	}

	_ = os.Remove(inputPath)
	_ = os.Remove(outputPath)

	return mp4Bytes, nil
}

type GenThumbVideoInput struct {
	Path  string
	Video []byte
}

func GenThumbVideo(input GenThumbVideoInput) ([]byte, error) {

	var err error
	video := input.Video
	videoPath := input.Path
	now := time.Now()
	thumbnailPath := fmt.Sprintf("temp/%d.png", now.Unix())

	if videoPath == "" {
		videoPath = fmt.Sprintf("temp/%d.mp4", now.Unix())
		err = os.WriteFile(videoPath, video, 0644)

	}

  fmt.Println("video thumb url reading", videoPath)

	if video == nil {
		video, err = os.ReadFile(videoPath)
		if err != nil {
			fmt.Println("Error generating Video thumb")
			return nil, err
		}

	}

	if err != nil {
		return nil, err
	}

	cmd := exec.Command("ffmpeg", "-i", videoPath, "-ss", "00:00:01.000", "-vframes", "1", "-q:v", "2", thumbnailPath)
	err = cmd.Run()
	if err != nil {
		return nil, err
	}

	thumbnail, err := os.ReadFile(thumbnailPath)
	if err != nil {
		return nil, err
	}

	err = os.Remove(videoPath)
	if err != nil {
		return nil, err
	}

	err = os.Remove(thumbnailPath)
	if err != nil {
		return nil, err
	}
	resized, err := ResizeImage(thumbnail, 72, 72)
	if err != nil {
		return nil, err
	}
	return resized, nil
}

func ResizeImage(image []byte, width uint, height uint) ([]byte, error) {
	imgz, err := png.Decode(strings.NewReader(string(image)))
	if err != nil {
		return nil, err
	}
	img := resize.Thumbnail(width, height, imgz, resize.Lanczos3)
	buf := new(bytes.Buffer)
	err = png.Encode(buf, img)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
