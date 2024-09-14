package commands

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
	"trevas-bot/pkg/commandextractor"
	"trevas-bot/pkg/converter"
	"trevas-bot/pkg/platform"
	"trevas-bot/pkg/platform/types"

	"github.com/PuerkitoBio/goquery"
)

type PornCommand struct {
	key      string
	platform platform.WhatsAppIntegration
}

func (p PornCommand) Handler(input commandextractor.CommandInput) {
	fmt.Println("Running Porn Command")

  searchTerm := strings.Trim(input.Payload, " ")
  searchTerm = strings.ReplaceAll(searchTerm, " ", "+")
  page := rand.Intn(5 - 1) + 1
  link := fmt.Sprintf("https://pt.pornhub.com/gifs/search?search=%s&page=%d", searchTerm, page)
  fmt.Println("link", link)
  res, err := http.Get(link)

  if err != nil {
    fmt.Println("Error curling")
    return
  }
  go p.platform.SendReaction(&input.EventMessage, platform.LoadingReaction)
  defer res.Body.Close()
  if res.StatusCode != 200 {
    go p.platform.SendReaction(&input.EventMessage, platform.ErrorReaction)
    p.platform.SendReply("Ocorreu um erro ao baixar mídia. Tente novamente.", &input.EventMessage)
    return
  }

  // Load the HTML document
  doc, err := goquery.NewDocumentFromReader(res.Body)
  if err != nil {
    fmt.Println(err)
    return
  }


  var lastFoundLink string
  var webmLink string
  // Find the review items

  sortedNumber := rand.Intn(30)

  doc.Find(".gifVideo").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the title
    tempLink, exists := s.Attr("data-webm")

    if exists {
      lastFoundLink = tempLink
    }

    if i == sortedNumber && exists {
      webmLink = tempLink
      return
    }
	})

  if webmLink == "" {
    webmLink = lastFoundLink
  }

  downloadPath := "temp/downloads/"
  now := time.Now()

  prefixFileName := fmt.Sprintf("%d", now.Unix())
  fileName := downloadPath + prefixFileName  + ".%(ext)s"

  // gifLink := fmt.Sprintf("https://el.phncdn.com/gif/%s.gif", extractGifId(webmLink))

  // cmd := exec.Command("yt-dlp", "-vU", gifLink, "--output", fileName)
  cmd := exec.Command("yt-dlp", "-vU", webmLink, "--output", fileName)

  err = cmd.Run()

  if err != nil {
    go p.platform.SendReaction(&input.EventMessage, platform.ErrorReaction)
    fmt.Println("Error trying to download media...", err.Error())
    p.platform.SendReply("Não foi possível fazer o download.", &input.EventMessage)
    return
  }
  
  downloadsFiles, err := os.ReadDir(downloadPath)

  if err != nil {
    go p.platform.SendReaction(&input.EventMessage, platform.ErrorReaction)
    fmt.Println("Error trying to download media...", err.Error())
    p.platform.SendReply("Não foi possível fazer o download.", &input.EventMessage)
    return
  }

  var downloadedFilePath string
  for _, file := range downloadsFiles {
    if strings.HasPrefix(file.Name(), prefixFileName) {
      fmt.Println("Sending Video...")
      downloadedFilePath = downloadPath + file.Name()
      webmBytes, err := os.ReadFile(downloadedFilePath)

      if err != nil {
        go p.platform.SendReaction(&input.EventMessage, platform.ErrorReaction)
        fmt.Println("Error trying to download media...", err.Error())
        p.platform.SendReply("Não foi possível baixar o video.", &input.EventMessage)
        return
      }

      videoBytes, err := converter.Webm2Mp4(webmBytes)

      if err != nil {
        go p.platform.SendReaction(&input.EventMessage, platform.ErrorReaction)
        fmt.Println("Error trying to download media...", err.Error())
        p.platform.SendReply("Não foi possível baixar o video.", &input.EventMessage)
        return
      }

      videoThumbnail, _ := converter.GenThumbVideo(converter.GenThumbVideoInput{Path: downloadedFilePath})

      err = p.platform.SendVideo(types.SendVideoInput{VideoBytes: videoBytes, Thumbnail: videoThumbnail}, &input.EventMessage)

      if err != nil {
        go p.platform.SendReaction(&input.EventMessage, platform.ErrorReaction)
        fmt.Println("Error trying to Send media...", err.Error())
        p.platform.SendReply(err.Error(), &input.EventMessage)
        return
      }

      _ = os.Remove(downloadedFilePath)
      go p.platform.SendReaction(&input.EventMessage, platform.SuccessReaction)
      break
    }
  }
}

func (c PornCommand) GetKey() string {
	return c.key
}

func NewPornCommand() *PornCommand {
	return &PornCommand{key: "porn"}
}

func extractGifId(url string) string {
  re := regexp.MustCompile(`(\d+)\D*\.webm$`)
	match := re.FindStringSubmatch(url)

	if len(match) > 1 {
    return match[1]
	}
  return ""
}
