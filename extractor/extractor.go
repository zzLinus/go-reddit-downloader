package extractor

import (
	"bufio"
	"fmt"
	"log"
	"net/http"

	"github.com/zzLinus/GoTUITODOList/fakeheaders"
)

func Extractor(rowURL string) string {
	client := &http.Client{}

	req, err := http.NewRequest(http.MethodGet, rowURL, nil)
	if err != nil {
		log.Fatal("failed at create request")
	}

	for k, v := range fakeheaders.FakeHeaders {
		req.Header.Set(k, v)
	}

	req.Header.Set("Referer", "https://www.reddit.com/")
	req.Header.Set("Origin", "https://www.reddit.com")
	// fmt.Println(req.Header)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("failed to recive a response")
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return ""
}
