package jobs

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
	"trevas-bot/pkg/converter"
	"trevas-bot/pkg/platform"
	"trevas-bot/pkg/platform/types"

	"github.com/PuerkitoBio/goquery"
	"go.mau.fi/whatsmeow"
)

type MemeSenderJob struct {
	platform platform.WhatsAppIntegration
}

func NewMemeSenderJob() *MemeSenderJob {
	return &MemeSenderJob{}
}

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

func (m MemeSenderJob) Execute(client *whatsmeow.Client) {
	fmt.Println("Running Post Command")

  currentTime := time.Now()

  if currentTime.Hour() < 7 {
    return
  }

	// Define the endpoint
	uri := "https://www.naointendo.com.br/api/posts"

	page := rand.Intn(8000)

	url := fmt.Sprintf("%s?page=%d", uri, page)
	fmt.Println("URl", url)

	clientHttp := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		fmt.Println("Error creating request", err.Error())
		return
	}

	req.Header.Set("x-requested-with", "XMLHttpRequest")

	resp, err := clientHttp.Do(req)

	if err != nil {
		fmt.Println("Não foi possível buscar um post...")
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

	var messages []types.Message

	groups, _ := client.GetJoinedGroups()

	for _, group := range groups {
		fmt.Println(group.GroupName.Name)
		fmt.Println("JID", group.JID)

		message := types.Message{Chat: group.JID}
		messages = append(messages, message)
	}

	err = m.handleResponse(postResponse.Posts, messages)

	if err != nil {
		m.Execute(client)
	}
}

func (m MemeSenderJob) CronConfig() string {
	return "*/30 * * * *" // every 15 minute
}

func (m MemeSenderJob) handleResponse(posts []PostItem, messages []types.Message) error {

	post := getRandomPost(posts)
	fmt.Println("Post", post)

	switch post.Type {
	case "html":
		return m.handleHtmlPost(post, messages)
	case "image":

		imageBytes, err := m.getImagePost(post.Media.Content)

		if err != nil {
			fmt.Println(err.Error())
			return err
		}

		for _, message := range messages {
			m.platform.SendImg(types.SendImageInput{Image: imageBytes, Caption: post.Title, Message: message}, nil)
		}
		break
	case "video":
		return m.sendVideoPost(post, messages)
	default:
		fmt.Println("Post type not supported")
	}
	return nil
}

func (m MemeSenderJob) handleHtmlPost(post PostItem, messages []types.Message) error {
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

		cmd := exec.Command("yt-dlp", "-vU", "--output", fileName, srcVideo)

		err := cmd.Run()

		if err != nil && !strings.Contains(err.Error(), "100") {
			fmt.Println("Error trying to download media...", err.Error())
			return err
		}

		defer func() {
			os.Remove(fileName)
		}()

		downloadsFiles, err := os.ReadDir(downloadPath)

		if err != nil {
			fmt.Println("Error trying to download media...", err.Error())
			return err
		}

		var downloadedFilePath string
		for _, file := range downloadsFiles {
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

				defer os.Remove(downloadedFilePath)

				for _, message := range messages {
					err = m.platform.SendVideo(types.SendVideoInput{VideoBytes: videoBytes, Thumbnail: thumbVideo, Caption: post.Title, Message: message}, nil)

					if err != nil {
						fmt.Println("Error trying to Send media...", err.Error())
						return err
					}
				}

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

		if err != nil {
			return
		}

		defer data.Body.Close()

    if data.StatusCode > 300 {
      return
    }

		imageData, err := io.ReadAll(data.Body)

		if err != nil {
			return
		}

		for _, message := range messages {
      m.platform.SendImg(types.SendImageInput{Image: imageData, Message: message}, nil)
			if err != nil {
				fmt.Println("Error trying to Send media...", err.Error())
				return 
			}
		}
	})
	return nil
}

func (m MemeSenderJob) getImagePost(imageUrl string) ([]byte, error) {
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

func (m MemeSenderJob) sendVideoPost(post PostItem, messages []types.Message) error {

	if post.Media.Type != "youtube" {
		return errors.New("Video Media Type Not Supported")
	}

	downloadPath := "./temp/downloads/"
	now := time.Now()

	prefixFileName := fmt.Sprintf("%d", now.Unix())
	fileName := downloadPath + prefixFileName + ".%(ext)s"

	videoSrc := fmt.Sprintf("https://www.youtube.com/watch?v=%s", post.Media.Content)

	cmd := exec.Command("yt-dlp", "-vU", videoSrc, "--output", fileName, "--max-filesize", "50M", "--no-playlist")

	err := cmd.Run()

	if err != nil && !strings.Contains(err.Error(), "100") {
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

			defer os.Remove(downloadedFilePath)

			for _, message := range messages {
				err = m.platform.SendVideo(types.SendVideoInput{VideoBytes: videoBytes, Thumbnail: thumbVideo, Caption: post.Title, Message: message}, nil)

				if err != nil {
					return err
				}
			}

			break
		}
	}
	return nil
}

func getRandomPost(posts []PostItem) PostItem {
	for true {
    postNumber := rand.Intn(len(posts))
    post := posts[postNumber]

    if post.Media.Type != "gif" {
      return post
    }
	}
	return PostItem{}
}
