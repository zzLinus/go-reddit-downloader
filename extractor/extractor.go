package extractor

import (
	"log"
)

var extractorMap = make(map[string]Extractor)

type Extractor interface {
	ExtractRowURL(rowURL string) (*Data, error)
}

func Register(domain string, e Extractor) {
	extractorMap[domain] = e
}

type Data struct {
	FileType        string
	VideoName       string
	DownloadableURL string
	AudioURL        string
}

func ExtractData(url string) (*Data, error) {
	extractor := extractorMap["reddit"]
	if extractor == nil {
		panic("can't get extractor")
	}
	data, err := extractor.ExtractRowURL(url)
	if err != nil {
		log.Fatal("some problem while extracting data")
	}
	return data, err
}
