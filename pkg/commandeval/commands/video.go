// yt-dlp -f "bestvideo[filesize<20M]" "ytsearch:filosofo piton" 
package commands

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
	"trevas-bot/pkg/commandextractor"
	"trevas-bot/pkg/converter"
	"trevas-bot/pkg/platform"
)

type VideoCommand struct {
	key      string
	platform platform.WhatsAppIntegration
}

func (p VideoCommand) Handler(input commandextractor.CommandInput) {
	fmt.Println("Running Video Command")

  // p.platform.SendReply("Comando desativado por enquanto.", &input.EventMessage)
  // return;

  go p.platform.SendReaction(&input.EventMessage, platform.LoadingReaction)

  downloadPath := "temp/downloads/"
  now := time.Now()

  prefixFileName := fmt.Sprintf("%d", now.Unix())
  fileName := downloadPath + prefixFileName  + ".%(ext)s"

  searchVideo := fmt.Sprintf("ytsearch:%s", input.Payload)

  cmd := exec.Command("yt-dlp", "-f", "bestvideo[filesize<20M][height<=?480]+bestaudio/best", searchVideo, "--output", fileName, "--no-playlist")


  err := cmd.Run()

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
      videoBytes, err := os.ReadFile(downloadedFilePath)

      if err != nil {
        go p.platform.SendReaction(&input.EventMessage, platform.ErrorReaction)
        fmt.Println("Error trying to download media...", err.Error())
        p.platform.SendReply("Não foi possível baixar o video.", &input.EventMessage)
        return
      }

        videoBytes, err = converter.Webm2Mp4(videoBytes)
        
        if err != nil {
          go p.platform.SendReaction(&input.EventMessage, platform.ErrorReaction)
          fmt.Println("Error trying to converting media to mp4...", err.Error())
          p.platform.SendReply("Não foi possível baixar o video.", &input.EventMessage)
          return
        }

      thumbVideo, _ := converter.GenThumbVideo(converter.GenThumbVideoInput{Video: videoBytes})
      fmt.Println("THumbVideo", thumbVideo)

      err = p.platform.SendVideo(platform.SendVideoInput{VideoBytes: videoBytes, Thumbnail: thumbVideo}, &input.EventMessage)

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

func (c VideoCommand) GetKey() string {
	return c.key
}

func NewVideoCommand() *VideoCommand {
	return &VideoCommand{key: "video"}
}
