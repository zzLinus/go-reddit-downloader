package main

// A simple program demonstrating the spinner component from the Bubbles
// component library.

import (
	"github.com/zzLinus/GoTUITODOList/downloader"
	"github.com/zzLinus/GoTUITODOList/tuiapp"
)

func main() {
	rowURL := "https://v.redd.it/8akffrc6fqx81/DASH_720.mp4"
	// url := extractor.Extractor(rowURL)
	err := downloader.Download(rowURL)
	if err != nil {
		panic(err)
	}

	p := tuiapp.New()
	p.Start()
}
