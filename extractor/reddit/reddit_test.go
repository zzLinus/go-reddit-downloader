package reddit

import (
	"fmt"
	"testing"

	"github.com/zzLinus/GoRedditDownloader/extractor"
)

func TestExtrac(t *testing.T) {
	testCases := []struct {
		rowURL string
	}{
		{
			rowURL: "https://www.reddit.com/r/ProgrammerHumor/comments/uqovco/my_code_works/",
		},
		{
			rowURL: "https://www.reddit.com/r/AnimatedPixelArt/comments/uomu32/animation_for_astral_ascent/",
		},
		{
			rowURL: "https://www.reddit.com/r/space/comments/uj8sod/a_couple_of_days_ago_i_visited_this_place_an/?utm_medium=android_app&utm_source=share",
		},
		{
			rowURL: "https://www.reddit.com/r/DotA2/comments/uq012r/til_how_useful_hurricane_bird_is/?utm_medium=android_app&utm_source=share",
		},
	}
	c := make(chan extractor.SubscriptMsg, 10)
	for i, testCase := range testCases {
		fmt.Println("test:", i)
		data, err := New().ExtractRowURL(testCase.rowURL, c)
		if err != nil {
			t.Error(err)
		}
		fmt.Println(data)
	}
}
