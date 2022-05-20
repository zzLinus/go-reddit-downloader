package downloader

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/zzLinus/GoTUITODOList/extractor"
	"github.com/zzLinus/GoTUITODOList/utils"
)

type Downloader struct {
}

var (
	rowURLExtractor *extractor.Extractor
)

func New() *Downloader {
	return &Downloader{}
}

func (*Downloader) Download(url string) (int, error) {
	data, err := extractor.ExtractData(url)
	var (
		resp = &http.Response{}
	)
	var audioFile *os.File
	filep := []string{}
	if err != nil {
		panic("can't extract data from this given url")
	}
	if data.AudioURL != "" {

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
	videoFile, err := os.OpenFile("video.mp4", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal("failed to create files")
		return 0, err
	}

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

	mergErr := utils.MergeAudioVideo(filep, data.VideoName+".mp4")
	if mergErr != nil {
		fmt.Println("can't merg audio with video")
	}

	return resp.StatusCode, nil
}
