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

func Img2Webp(image []byte, isVideo bool) ([]byte, error) {

	now := time.Now()
	inputPath := fmt.Sprintf("./temp/%d.gif", now.Unix())
	outputPath := fmt.Sprintf("./temp/%d.webp", now.Unix())

	err := os.WriteFile(inputPath, image, 0644)

	if err != nil {
		fmt.Println("Error saving input to use in converter", err)
		return nil, nil
	}

  var cmd *exec.Cmd

  fmt.Println("isVideo", isVideo)
  if (isVideo) {
    cmd = exec.Command("ffmpeg",
        "-i", inputPath,
        "-vf", "fps=10,scale=512:512:force_original_aspect_ratio=decrease",
        "-c:v", "libwebp",
        "-q:v", "60",
        "-lossless", "0",
        "-loop", "0",
        "-preset", "default",
        "-an",
        "-vsync", "0",
      "-t", "6",
      "-q:v", "60",
        outputPath)
  } else {
    cmd = exec.Command("ffmpeg", "-i", inputPath, outputPath);
  }

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

func Webp2Img(sticker []byte) ([]byte, error) {

	now := time.Now()
	inputPath := fmt.Sprintf("./temp/%d.webp", now.Unix())
	outputPath := fmt.Sprintf("./temp/%d.jpg", now.Unix())

	err := os.WriteFile(inputPath, sticker, 0644)

	if err != nil {
		fmt.Println("Error saving input to use in converter", err)
		return nil, nil
	}

	cmd := exec.Command("ffmpeg", "-i", inputPath, outputPath)
	if err := cmd.Run(); err != nil {
		fmt.Println("Erro ao converter o arquivo:", err)
		_ = os.Remove(inputPath)
		_ = os.Remove(outputPath)
		return nil, errors.New("Ocorreu um erro ao processar sua figurinha... por favor, tente novamente")
	}

	imageBytes, err := os.ReadFile(outputPath)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return nil, errors.New("Ocorreu um erro ao processar sua figurinha... por favor, tente novamente")
	}

	// TODO por metadatas no sticker

	_ = os.Remove(inputPath)
	_ = os.Remove(outputPath)

	return imageBytes, nil
}

func Webm2Mp4(webm []byte) ([]byte, error) {
	now := time.Now()
	inputPath := fmt.Sprintf("./temp/%d-temp.webm", now.Unix())
	outputPath := fmt.Sprintf("./temp/%d-temp.mp4", now.Unix())

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
	thumbnailPath := fmt.Sprintf("./temp/%d.png", now.Unix())

	if videoPath == "" {
		videoPath = fmt.Sprintf("./temp/%d.mp4", now.Unix())
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
