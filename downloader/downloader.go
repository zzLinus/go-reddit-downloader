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
	prefixedurlBeg  = "https://v.redd.it/"
	prefixedurlEnd  = "/DASH_720.mp4"
)

func New() *Downloader {
	return &Downloader{}
}

func (*Downloader) Download(url string) (int, error) {
	rowURLExtractor = extractor.New()

	downloadableURL, err := rowURLExtractor.RowURLExtractor(url)
	if err != nil {
		log.Fatal("can't extract row URL")
		return 0, nil
	}

	downloadableURL[0] = prefixedurlBeg + downloadableURL[0] + prefixedurlEnd

	videoFile, err := os.Create("baofengxuehao.mp4")
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
