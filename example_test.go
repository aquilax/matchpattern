package matchpattern_test

import (
	"fmt"
	"strings"

	"github.com/aquilax/matchpattern"
)

func ExampleMatchPattern_Matches() {
	cases := []struct {
		patterns []string
		urls     []string
	}{
		{
			[]string{matchpattern.AllURLPattern},
			[]string{"https://www.example.com"},
		},
		{
			[]string{"*://*example.com/"},
			[]string{"https://www.example.com", "http://www.example.com"},
		},
	}
	for _, example := range cases {
		mp, err := matchpattern.New(example.patterns, matchpattern.GetDefaultMatchSet())
		if err != nil {
			panic(err)
		}
		for _, url := range example.urls {
			if ok, err := mp.Matches(url); err != nil {
				panic(err)
			} else {
				if ok {
					fmt.Printf("%s is a match for %s\n", url, strings.Join(example.patterns, ", "))
				} else {
					fmt.Printf("%s is not a match for %s\n", url, strings.Join(example.patterns, ", "))
				}
			}
		}
	}

	// Output:
	// https://www.example.com is a match for <all_urls>
	// https://www.example.com is not a match for *://*example.com/
	// http://www.example.com is not a match for *://*example.com/
}
