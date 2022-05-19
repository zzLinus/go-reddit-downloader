package reddit

import (
	"fmt"
	"testing"
)

func TestExtrac(t *testing.T) {
	testCases := []struct {
		rowURL string
	}{
		{
			rowURL: "https://www.reddit.com/r/space/comments/uj8sod/a_couple_of_days_ago_i_visited_this_place_an/?utm_medium=android_app&utm_source=share",
		},
		{
			rowURL: "https://www.reddit.com/r/DotA2/comments/uq012r/til_how_useful_hurricane_bird_is/?utm_medium=android_app&utm_source=share",
		},
	}
	for i, testCase := range testCases {
		fmt.Println("test:", i)
		data, err := New().ExtractRowURL(testCase.rowURL)
		if err != nil {
			t.Error(err)
		}
		fmt.Println(data)
	}
}
