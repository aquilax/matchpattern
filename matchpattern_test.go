package matchpattern

import (
	"fmt"
	"testing"
)

func TestMatchPattern_Matches(t *testing.T) {
	tests := []struct {
		name string
		p    []string
		url  string
		want bool
	}{
		// https://developer.mozilla.org/en-US/docs/Mozilla/Add-ons/WebExtensions/Match_patterns#examples
		{"all_urls", []string{"<all_urls>"}, "http://example.org/", true},
		{"all_urls", []string{"<all_urls>"}, "https://a.org/some/path/", true},
		{"all_urls", []string{"<all_urls>"}, "ws://sockets.somewhere.org/", true},
		{"all_urls", []string{"<all_urls>"}, "wss://ws.example.com/stuff/", true},
		{"all_urls", []string{"<all_urls>"}, "ftp://files.somewhere.org/", true},
		{"all_urls", []string{"<all_urls>"}, "ftps://files.somewhere.org/", true},
		{"all_urls", []string{"<all_urls>"}, "resource://a/b/c/", false},

		{"*://*/*", []string{"*://*/*"}, "http://example.org/", true},
		{"*://*/*", []string{"*://*/*"}, "https://a.org/some/path/", true},
		{"*://*/*", []string{"*://*/*"}, "ws://sockets.somewhere.org/", true},
		{"*://*/*", []string{"*://*/*"}, "wss://ws.example.com/stuff/", true},
		{"*://*/*", []string{"*://*/*"}, "ftp://ftp.example.org/", false},
		{"*://*/*", []string{"*://*/*"}, "ftps://ftp.example.org/", false},
		{"*://*/*", []string{"*://*/*"}, "file:///a/", false},

		// https://developer.chrome.com/docs/extensions/mv3/match_patterns/
		{"all_urls", []string{"<all_urls>"}, "http://example.org/foo/bar.html", true},
		{"all_urls", []string{"<all_urls>"}, "file:///bar/baz.html", true},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf(" `%s` %s", tt.name, tt.url), func(t *testing.T) {
			m, err := NewMatchPattern(tt.p, GetDefaultMatchSet())
			if err != nil {
				t.Fatalf(err.Error())
			}
			if got, _ := m.Matches(tt.url); got != tt.want {
				t.Errorf("MatchPattern.Matches() = %v, want %v", got, tt.want)
			}
		})
	}
}
