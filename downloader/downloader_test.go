package downloader

import (
	"testing"

	"github.com/zzLinus/GoRedditDownloader/extractor"
)

func TestDownloa(t *testing.T) {
	testCases := []struct {
		rowURL string
	}{
		{
			rowURL: "https://www.reddit.com/r/DotA2/comments/uq012r/til_how_useful_hurricane_bird_is/",
		},
	}
	c := make(chan extractor.SubscriptMsg, 10)
	for _, testCase := range testCases {
		_, err := New().Download(testCase.rowURL, c)
		if err != nil {
			t.Error(err)
		}
	}
}
