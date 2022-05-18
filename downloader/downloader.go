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
	rowURLExtractor = extractor.New()

	downloadableURL, err := rowURLExtractor.ExtractRowURL(url)
	if err != nil {
		log.Fatal("can't extract row URL")
		return 0, nil
	}

	videoFile, err := os.Create("viedo.mp4")
	if err != nil {
		log.Fatal("failed to create files")
		return 0, err
	}

	resp, err := http.Get(downloadableURL[0])
	if err != nil {
		log.Fatal(err)
		return 0, err
	}
	defer resp.Body.Close()

	io.Copy(videoFile, resp.Body)

	return resp.StatusCode, nil
}
