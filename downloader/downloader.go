package downloader

import (
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/zzLinus/GoRedditDownloader/extractor"
	"github.com/zzLinus/GoRedditDownloader/utils"
)

type Downloader struct {
}

var (
	rowURLExtractor *extractor.Extractor
)

func New() *Downloader {
	return &Downloader{}
}

func (*Downloader) Download(url string, c chan extractor.SubscriptMsg) (int, error) {
	now := time.Now()
	data, err := extractor.ExtractData(url, c)

	var (
		resp = &http.Response{}
	)
	var audioFile *os.File
	filep := []string{}
	if err != nil {
		panic("can't extract data from this given url")
	}
	if data.AudioURL != "" {
		c <- extractor.SubscriptMsg{Msg: "Downloading audio data", Duration: time.Now().Sub(now)}
		now = time.Now()
		audioFile, err = os.OpenFile("autio.mp4", os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			log.Fatal("failed to create files")
			return 0, err
		}

		for reTrytimes := 10; reTrytimes > 0; reTrytimes-- {
			resp, err = http.Get(data.AudioURL)
			if (err != nil || resp.StatusCode > 400) && reTrytimes > 0 {
				time.Sleep(1 * time.Second)
			} else {
				break
			}
			if reTrytimes == 0 {
				return 0, err
			}
		}
		defer resp.Body.Close()
		filep = append(filep, audioFile.Name())

		io.Copy(audioFile, resp.Body)
	}

	//if there is no such file create one and give it right
	videoFile, err := os.OpenFile("video."+data.FileType, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal("failed to create files")
		return 0, err
	}

	c <- extractor.SubscriptMsg{Msg: "Downloading video data", Duration: time.Now().Sub(now)}
	now = time.Now()
	for reTrytimes := 10; reTrytimes > 0; reTrytimes-- {
		resp, err = http.Get(data.DownloadableURL)
		if (err != nil || resp.StatusCode > 400) && reTrytimes > 0 {
			time.Sleep(1 * time.Second)
		} else {
			break
		}
		if reTrytimes == 0 {
			return 0, err
		}
	}

	defer resp.Body.Close()

	io.Copy(videoFile, resp.Body)

	filep = append(filep, videoFile.Name())

	if data.FileType != "jpg" {
		mergErr := utils.MergeAudioVideo(filep, data.VideoName+".mp4")
		c <- extractor.SubscriptMsg{Msg: "Merging video with audio", Duration: time.Now().Sub(now)}
		now = time.Now()
		if mergErr != nil {
			log.Fatal("can't merg audio with video")
		}
	}

	time.Sleep(3 * time.Second)
	return resp.StatusCode, nil
}
