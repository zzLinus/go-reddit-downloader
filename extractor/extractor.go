package extractor

import (
	"log"
)

var extractorMap = make(map[string]Extractor)

type Extractor interface {
	ExtractRowURL(rowURL string, c chan SubscriptMsg) (*Data, error)
}

type SubscriptMsg struct {
	Msg string
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

func ExtractData(url string, c chan SubscriptMsg) (*Data, error) {
	extractor := extractorMap["reddit"]
	if extractor == nil {
		panic("can't get extractor")
	}
	data, err := extractor.ExtractRowURL(url, c)
	if err != nil {
		log.Fatal("some problem while extracting data")
	}
	return data, err
}
