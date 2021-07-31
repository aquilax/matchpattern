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

		{"*://*.mozilla.org/*", []string{"*://*.mozilla.org/*"}, "http://mozilla.org/", true},
		{"*://*.mozilla.org/*", []string{"*://*.mozilla.org/*"}, "https://mozilla.org/", true},
		{"*://*.mozilla.org/*", []string{"*://*.mozilla.org/*"}, "http://a.mozilla.org/", true},
		{"*://*.mozilla.org/*", []string{"*://*.mozilla.org/*"}, "http://a.b.mozilla.org/", true},
		{"*://*.mozilla.org/*", []string{"*://*.mozilla.org/*"}, "https://b.mozilla.org/path/", true},
		{"*://*.mozilla.org/*", []string{"*://*.mozilla.org/*"}, "ws://ws.mozilla.org/", true},
		{"*://*.mozilla.org/*", []string{"*://*.mozilla.org/*"}, "wss://secure.mozilla.org/something", true},

		{"*://*.mozilla.org/*", []string{"*://*.mozilla.org/*"}, "ftp://mozilla.org/", false},
		{"*://*.mozilla.org/*", []string{"*://*.mozilla.org/*"}, "http://mozilla.com/", false},
		{"*://*.mozilla.org/*", []string{"*://*.mozilla.org/*"}, "http://firefox.org/", false},

		{"*://mozilla.org/", []string{"*://mozilla.org/"}, "http://mozilla.org/", true},
		{"*://mozilla.org/", []string{"*://mozilla.org/"}, "https://mozilla.org/", true},
		{"*://mozilla.org/", []string{"*://mozilla.org/"}, "ws://mozilla.org/", true},
		{"*://mozilla.org/", []string{"*://mozilla.org/"}, "wss://mozilla.org/", true},
		{"*://mozilla.org/", []string{"*://mozilla.org/"}, "ftp://mozilla.org/", false},
		{"*://mozilla.org/", []string{"*://mozilla.org/"}, "http://a.mozilla.org/", false},
		{"*://mozilla.org/", []string{"*://mozilla.org/"}, "http://mozilla.org/a", false},

		{"ftp://mozilla.org/", []string{"ftp://mozilla.org/"}, "ftp://mozilla.org", true},
		{"ftp://mozilla.org/", []string{"ftp://mozilla.org/"}, "http://mozilla.org/", false},
		{"ftp://mozilla.org/", []string{"ftp://mozilla.org/"}, "ftp://sub.mozilla.org/", false},
		{"ftp://mozilla.org/", []string{"ftp://mozilla.org/"}, "ftp://mozilla.org/path", false},

		{"https://*/path", []string{"https://*/path"}, "https://mozilla.org/path", true},
		{"https://*/path", []string{"https://*/path"}, "https://a.mozilla.org/path", true},
		{"https://*/path", []string{"https://*/path"}, "https://something.com/path", true},
		{"https://*/path", []string{"https://*/path"}, "http://mozilla.org/path", false},
		{"https://*/path", []string{"https://*/path"}, "https://mozilla.org/path/", false},
		{"https://*/path", []string{"https://*/path"}, "https://mozilla.org/a", false},
		{"https://*/path", []string{"https://*/path"}, "https://mozilla.org/", false},
		{"https://*/path", []string{"https://*/path"}, "https://mozilla.org/path?foo=1", false},

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
