package extractor

import (
	"fmt"
	"testing"
)

func TestExtrac(t *testing.T) {
	testCases := []struct {
		rowURL string
	}{
		{
			rowURL: "https://www.reddit.com/r/space/comments/uj8sod/a_couple_of_days_ago_i_visited_this_place_an/",
		},
		{
			rowURL: "https://www.reddit.com/r/DotA2/comments/uq012r/til_how_useful_hurricane_bird_is/",
		},
	}
	for i, testCase := range testCases {
		fmt.Println("test:", i)
		_, err := New().RowURLExtractor(testCase.rowURL)
		if err != nil {
			t.Error(err)
		}
	}
}
