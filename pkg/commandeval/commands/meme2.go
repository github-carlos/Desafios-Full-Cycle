package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"

	// "math/rand"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
	"trevas-bot/pkg/commandextractor"
	"trevas-bot/pkg/converter"
	"trevas-bot/pkg/platform"
	"trevas-bot/pkg/platform/types"

	"github.com/PuerkitoBio/goquery"
	"go.mau.fi/whatsmeow/types/events"
)

type Category struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type Media struct {
	// html, image
	Type    string `json:"type"`
	Content string `json:"content"`
}

type PostItem struct {
	Category Category `json:"category"`
	Media    Media    `json:"media"`
	Title    string   `json:"title"`
	Type     string   `json:"type"`
}

type PostResponse struct {
	Posts []PostItem `json:"posts"`
}

type PostCommand struct {
	key      string
	platform platform.WhatsAppIntegration
}

func (p PostCommand) Handler(input commandextractor.CommandInput) {
	fmt.Println("Running Post Command")

	// Define the endpoint
	uri := "https://www.naointendo.com.br/api/posts"

	page := 7985 //rand.Intn(8000)

	url := fmt.Sprintf("%s?page=%d", uri, page)
	fmt.Println("URl", url)

	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		fmt.Println("Error creating request", err.Error())
		return
	}

	req.Header.Set("x-requested-with", "XMLHttpRequest")

	resp, err := client.Do(req)

	if err != nil {
		go p.platform.SendReply("Não foi possível buscar um post...", &input.EventMessage)
		return
	}

	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	// Parse the JSON response
	var postResponse PostResponse
	err = json.Unmarshal(body, &postResponse)

	if err != nil {
		return
	}

  err = p.handleResponse(postResponse.Posts, &input.EventMessage)

  if err != nil {
    p.Handler(input)
  }
}

func (p PostCommand) GetKey() string {
	return p.key
}

func NewPostCommand() *PostCommand {
	return &PostCommand{key: "post"}
}

func (p PostCommand) handleResponse(posts []PostItem, e *events.Message) error {

  post := getRandomPost(posts)
	fmt.Println("Post", post)

	switch post.Type {
	case "html":
    return p.handleHtmlPost(post, e)
	case "image":

		imageBytes, err := p.getImagePost(post.Media.Content)

		if err != nil {
			fmt.Println(err.Error())
			return err
		}

		p.platform.SendImg(types.SendImageInput{Image: imageBytes, Caption: post.Title}, e)
		break
	case "video":
		return p.sendVideoPost(post, e)
	default:
		fmt.Println("Post type not supported")
	}
  return nil
}

func (p PostCommand) handleHtmlPost(post PostItem, e *events.Message) error {
	// it can be video or images inside HTML
	fmt.Println("Handling Post HTML", post.Media.Content)
	reader := strings.NewReader(post.Media.Content)

	doc, err := goquery.NewDocumentFromReader(reader)

	if err != nil {
		fmt.Println("Error creating goquery instance")
	}

	iframe := doc.Find("iframe").First()
	fmt.Println("iframe", iframe)

	if iframe.Length() > 0 {
		srcVideo, exists := iframe.Attr("src")

		if !exists {
			return nil
		}

		downloadPath := "./temp/downloads/"
		now := time.Now()

		prefixFileName := fmt.Sprintf("%d", now.Unix())
		fileName := downloadPath + prefixFileName + ".%(ext)s"

		fmt.Println("Src video", srcVideo)

		cmd := exec.Command("yt-dlp", "-vU", "--output", fileName, srcVideo)

		err := cmd.Run()

		if err != nil && !strings.Contains(err.Error(), "100") {
			fmt.Println("Error trying to download media...", err.Error())
			return err
		}

		defer func() {
			os.Remove(fileName)
		}()

		fmt.Println("DownloadPath", downloadPath)
		downloadsFiles, err := os.ReadDir(downloadPath)

		fmt.Println("Downloaded Files", downloadsFiles)

		if err != nil {
			fmt.Println("Error trying to download media...", err.Error())
			return err
		}

		var downloadedFilePath string
		for _, file := range downloadsFiles {
			fmt.Println("FileName", file.Name())
			if strings.HasPrefix(file.Name(), prefixFileName) {
				fmt.Println("Sending Video...")
				downloadedFilePath = downloadPath + file.Name()
				videoBytes, err := os.ReadFile(downloadedFilePath)

				if err != nil {
					fmt.Println("Error trying to download media...", err.Error())
					return err
				}

				videoBytes, err = converter.Webm2Mp4(videoBytes)

				if err != nil {
					fmt.Println("Error trying to converting media to mp4...", err.Error())
					return err
				}

				thumbVideo, _ := converter.GenThumbVideo(converter.GenThumbVideoInput{Video: videoBytes})
				err = p.platform.SendVideo(types.SendVideoInput{VideoBytes: videoBytes, Thumbnail: thumbVideo, Caption: post.Title}, e)

				if err != nil {
					fmt.Println("Error trying to Send media...", err.Error())
					return err
				}

				_ = os.Remove(downloadedFilePath)
				break
			}
		}
	}

	imgs := doc.Find("img")

	imgs.Each(func(i int, s *goquery.Selection) {
		imgLink, exists := s.Attr("src")

		if !exists {
			return
		}

		data, err := http.Get(imgLink)

		if err != nil || data.StatusCode > 300 {
			return
		}

		defer data.Body.Close()

		imageData, err := io.ReadAll(data.Body)

		if err != nil {
			return
		}

		p.platform.SendImg(types.SendImageInput{Image: imageData}, e)
	})
  return nil
}

func (p PostCommand) getImagePost(imageUrl string) ([]byte, error) {
	data, err := http.Get(imageUrl)

	if err != nil {
		return nil, errors.New("Error downloading image")
	}

	defer data.Body.Close()

	imageData, err := io.ReadAll(data.Body)

	if err != nil {
		return nil, errors.New("Error getting image buffer")
	}

	return imageData, nil
}

func (p PostCommand) sendVideoPost(post PostItem, e *events.Message) error {

	if post.Media.Type != "youtube" {
		return errors.New("Video Media Type Not Supported")
	}

	downloadPath := "./temp/downloads/"
	now := time.Now()

	prefixFileName := fmt.Sprintf("%d", now.Unix())
	fileName := downloadPath + prefixFileName + ".%(ext)s"

  videoSrc := fmt.Sprintf("https://www.youtube.com/watch?v=%s", post.Media.Content)

  fmt.Println("video", videoSrc)
	cmd := exec.Command("yt-dlp", "-vU", videoSrc, "--output", fileName, "--max-filesize", "50M", "--no-playlist")

	err := cmd.Run()

	if err != nil && !strings.Contains(err.Error(), "100") {
    fmt.Println(err.Error())
		return errors.New("Error trying to download media...")
	}

	downloadsFiles, err := os.ReadDir(downloadPath)

	if err != nil {
    return err
	}

	var downloadedFilePath string
	for _, file := range downloadsFiles {
		if strings.HasPrefix(file.Name(), prefixFileName) {
			downloadedFilePath = downloadPath + file.Name()
			videoBytes, err := os.ReadFile(downloadedFilePath)

			if err != nil {
				return err
			}

			if strings.HasSuffix(downloadedFilePath, "webm") {
				videoBytes, err = converter.Webm2Mp4(videoBytes)

				if err != nil {
					return err
				}
			}

			thumbVideo, _ := converter.GenThumbVideo(converter.GenThumbVideoInput{Video: videoBytes})

      err = p.platform.SendVideo(types.SendVideoInput{VideoBytes: videoBytes, Thumbnail: thumbVideo, Caption: post.Title}, e)

			if err != nil {
				return err
			}

			_ = os.Remove(downloadedFilePath)
			break
		}
	}
  return nil
}

func getRandomPost(posts []PostItem) PostItem {
  for true {
    postNumber := 0//rand.Intn(len(posts))
    post := posts[postNumber]

    if post.Media.Type != "gif" {
      fmt.Println("Post Number", postNumber)
      return post
    }
  }
  return PostItem{}
}
