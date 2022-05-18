package extractor

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/zzLinus/GoTUITODOList/fakeheaders"
)

const (
	redditAPI = "https://v.redd.it/"
	res720    = "/DASH_720.mp4"
)

type Extractor struct {
	contentType int8
}

func New() *Extractor {
	return &Extractor{}
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

	for ; reTrytimes > 0; reTrytimes-- {
		resp, err = client.Do(req)
		if err != nil && reTrytimes > 0 {
			log.Fatal("failed to recive a response")
			return "", err
		} else {
			break
		}
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	return string(body), nil
}

func (*Extractor) ExtractRowURL(rowURL string) ([]string, error) {
	html, err := getHTMLPage(rowURL)
	if err != nil {
		log.Fatal("Failed to get html page")
	}

	urls := matchOneOf(html, `https://v.redd.it/.*/HLSPlaylist`)

	if urls == nil {
		fmt.Println("can't match anything")
	}

	for i, url := range urls {
		urls[i] = fmt.Sprintf("%s%s%s", redditAPI, url[18:31], res720)
		fmt.Println(urls[i])
	}

	return urls, nil
}

func matchOneOf(text string, patterns ...string) []string {
	var (
		re    *regexp.Regexp
		value []string
	)

	for _, pattern := range patterns {
		re = regexp.MustCompile(pattern)
		value = re.FindStringSubmatch(text)
		if len(value) > 0 {
			return value
		}
	}
	return nil
}
