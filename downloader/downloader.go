package downloader

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func Download(url string) error {

	videoFile, err := os.Create("baofengxuehao.mp4")
	if err != nil {
		log.Fatal("failed to create files")
		return err
	}

	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer resp.Body.Close()

	fmt.Println("Response status:", resp.Status)
	io.Copy(videoFile, resp.Body)

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}
