package extractor

import (
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/zzLinus/GoTUITODOList/fakeheaders"
)

type Extractor struct {
}

func New() *Extractor {
	return &Extractor{}
}

func GetHTMLPage(rowURL string) (string, error) {
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

func (*Extractor) RowURLExtractor(rowURL string, path string) (string, error) {
	html, err := GetHTMLPage(rowURL)
	if err != nil {
		log.Fatal("Failed to get html page")
	}
	urls := MatchOneOf(html, `https://v.redd.it/.*/HLSPlaylist`)

	if urls == nil {
		fmt.Println("can't match anything")
	}

	for _, url := range urls {
		fmt.Println(url)
	}

	respHTML, err := os.Create(path)
	if err != nil {
		fmt.Println("can't create file")
		return "", err
	}

	io.WriteString(respHTML, html)

	return "", nil
}

func MatchOneOf(text string, patterns ...string) []string {
	var (
		re    *regexp.Regexp
		value []string
	)
	for _, pattern := range patterns {
		// (?flags): set flags within current group; non-capturing
		// s: let . match \n (default false)
		// https://github.com/google/re2/wiki/Syntax
		re = regexp.MustCompile(pattern)
		value = re.FindStringSubmatch(text)
		if len(value) > 0 {
			return value
		}
	}
	return nil
}
