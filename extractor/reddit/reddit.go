package reddit

import (
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/zzLinus/GoRedditDownloader/extractor"
	"github.com/zzLinus/GoRedditDownloader/fakeheaders"

	"fmt"
	"log"

	"github.com/zzLinus/GoRedditDownloader/utils"
)

const (
	redditMP4API = "https://v.redd.it/"
	audioURLPart = "/DASH_audio.mp4"
	res720       = "/DASH_720.mp4"
	res480       = "/DASH_480.mp4"
	res360       = "/DASH_360.mp4"
	res280       = "/DASH_280.mp4"
)

type redditExtractor struct{}

func init() {
	extractor.Register("reddit", New())
}

func New() extractor.Extractor {
	return &redditExtractor{}
}

func (*redditExtractor) ExtractRowURL(rowURL string) (*extractor.Data, error) {
	html, err := getHTMLPage(rowURL)
	if err != nil {
		log.Fatal("Failed to get html page")
		return &extractor.Data{}, err
	}

	return getData(html), nil
}

func getData(html string) *extractor.Data {
	var fileType = ""
	videoName := utils.MatchOneOf(html, `<title>.*<\/title>`)[0]
	if utils.MatchOneOf(html, `meta property="og:video" content=.*HLSPlaylist`) != nil {
		fileType = "mp4"
	} else if utils.MatchOneOf(html, `https:\/\/preview\.redd\.it\/.*gif`) != nil {
		fileType = "gif"
	}

	if fileType == "mp4" {
		url := utils.MatchOneOf(html, `https://v.redd.it/.*/HLSPlaylist`)[0]
		if url == "" {
			panic("can't match anything")
		}

		for i := len(url) - 1; i >= 0; i-- {
			if url[i] == '/' {
				url = url[:i]
				break
			}
		}
		for i := len(url) - 1; i >= 0; i-- {
			if url[i] == '/' {
				url = url[i+1:]
				break
			}
		}

		for i := len(videoName) - 1; i >= 0; i-- {
			if videoName[i] == '<' {
				videoName = videoName[:i]
				break
			}
		}

		for i := len(videoName) - 1; i >= 0; i-- {
			if videoName[i] == '>' {
				videoName = videoName[i+1:]
				break
			}
		}

		videoURL := fmt.Sprintf("%s%s%s", redditMP4API, url, res720)
		audioURL := fmt.Sprintf("%s%s%s", redditMP4API, url, audioURLPart)

		return &extractor.Data{FileType: fileType,
			VideoName:       videoName,
			DownloadableURL: videoURL,
			AudioURL:        audioURL,
		}
	} else if fileType == "gif" {
		url, urlU, urlL := "", "", ""
		urls := utils.MatchOneOf(html, `https:\/\/preview\.redd\.it\/.*?\.gif\?format=mp4.*?"`)
		if urls != nil {
			url = urls[0]
		}
		if url == "" {
			panic("can't match anything")
		}

		for i := len(url) - 1; i >= 0; i-- {
			if url[i] == '&' {
				urlU = url[:i+1]
				break
			}
		}

		for i := len(url) - 1; i >= 0; i-- {
			if url[i] == ';' {
				urlL = url[i+1 : len(url)-1]
				break
			}
		}

		url = urlU + urlL

		for i := len(videoName) - 1; i >= 0; i-- {
			if videoName[i] == '<' {
				videoName = videoName[:i]
				break
			}
		}

		for i := len(videoName) - 1; i >= 0; i-- {
			if videoName[i] == '>' {
				videoName = videoName[i+1:]
				break
			}
		}

		url = fmt.Sprintf("%s", url)

		// warning:i don't know why the .gif can't open after downloaded,but after rename it as .mp4 it dose play
		return &extractor.Data{FileType: "mp4", VideoName: videoName, DownloadableURL: url}
	}
	return &extractor.Data{}
}

func getHTMLPage(rowURL string) (string, error) {
	var (
		reTrytimes = 10
		resp       = &http.Response{}
	)
	transport := &http.Transport{
		Proxy:               http.ProxyFromEnvironment,
		DisableCompression:  true,
		TLSHandshakeTimeout: 10 * time.Second,
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{
		Transport: transport,
		Timeout:   15 * time.Minute,
	}

	req, err := http.NewRequest(http.MethodGet, rowURL, nil)
	if err != nil {
		log.Fatal("failed at create request")
		return "", err
	}

	for k, v := range fakeheaders.FakeHeaders {
		req.Header.Set(k, v)
	}

	req.Header.Set("Referer", "https://www.reddit.com/")
	req.Header.Set("Origin", "https://www.reddit.com")

	//retry after 1 second if request failed and reTrytimes > 0
	for ; reTrytimes > 0; reTrytimes-- {
		resp, err = client.Do(req)
		if (err != nil || resp.StatusCode > 400) && reTrytimes > 0 {
			time.Sleep(1 * time.Second)
		} else {
			break
		}
		if reTrytimes == 0 {
			return "", err
		}
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("failed to read response body part")
		return "", err
	}

	return string(body), nil
}
