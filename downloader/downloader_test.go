package downloader

import "testing"

func TestDownloa(t *testing.T) {
	testCases := []struct {
		rowURL string
	}{
		{
			rowURL: "https://www.reddit.com/r/DotA2/comments/uq012r/til_how_useful_hurricane_bird_is/",
		},
	}
	for _, testCase := range testCases {
		_, err := New().Download(testCase.rowURL)
		if err != nil {
			t.Error(err)
		}
	}
}
