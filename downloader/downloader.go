package downloader

import (
	"io"
	"log"
	"net/http"
	"os"
)

type Downloader struct {
}

func New() *Downloader {
	return &Downloader{}
}

func (*Downloader) Download(url string) (int, error) {

	videoFile, err := os.Create("baofengxuehao.mp4")
	if err != nil {
		log.Fatal("failed to create files")
		return 0, err
	}

	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
		return 0, err
	}
	defer resp.Body.Close()

	io.Copy(videoFile, resp.Body)

	return resp.StatusCode, nil
}
