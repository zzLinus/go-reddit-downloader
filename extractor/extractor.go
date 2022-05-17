package extractor

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/zzLinus/GoTUITODOList/fakeheaders"
)

type Extractor struct {
}

func New() *Extractor {
	return &Extractor{}
}

func (*Extractor) RowURLExtractor(rowURL string) (string, error) {
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

	fmt.Println(req.Header)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("failed to recive a response")
		return "", err
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
		return "", err
	}

	return "", nil
}
