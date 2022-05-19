package downloader

import (
	"io"
	"log"
	"net/http"
	"os"

	"github.com/zzLinus/GoTUITODOList/extractor"
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
	if err != nil {
		panic("can't extract data from this given url")
	}

	videoFile, err := os.Create("viedo.mp4")
	if err != nil {
		log.Fatal("failed to create files")
		return 0, err
	}

	resp, err := http.Get(data.DownloadableURL)
	if err != nil {
		log.Fatal(err)
		return 0, err
	}
	defer resp.Body.Close()

	io.Copy(videoFile, resp.Body)

	return resp.StatusCode, nil
}
