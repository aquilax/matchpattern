package matchpattern

import (
	"fmt"
	"testing"
)

func TestMatchPattern_MatchesMDN(t *testing.T) {
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

		{"https://mozilla.org/*", []string{"https://mozilla.org/*"}, "https://mozilla.org/", true},
		{"https://mozilla.org/*", []string{"https://mozilla.org/*"}, "https://mozilla.org/path", true},
		{"https://mozilla.org/*", []string{"https://mozilla.org/*"}, "https://mozilla.org/another", true},
		{"https://mozilla.org/*", []string{"https://mozilla.org/*"}, "https://mozilla.org/path/to/doc", true},
		{"https://mozilla.org/*", []string{"https://mozilla.org/*"}, "https://mozilla.org/path/to/doc?foo=1", true},
		{"https://mozilla.org/*", []string{"https://mozilla.org/*"}, "http://mozilla.org/path", false},
		{"https://mozilla.org/*", []string{"https://mozilla.org/*"}, "https://mozilla.com/path", false},

		{"https://mozilla.org/a/b/c/", []string{"https://mozilla.org/a/b/c/"}, "https://mozilla.org/a/b/c/", true},
		{"https://mozilla.org/a/b/c/", []string{"https://mozilla.org/a/b/c/"}, "https://mozilla.org/a/b/c/#section1", true},
		{"https://mozilla.org/a/b/c/", []string{"https://mozilla.org/a/b/c/"}, "https://mozilla.org/a/b/c/?1", false},

		{"https://mozilla.org/*/b/*/", []string{"https://mozilla.org/*/b/*/"}, "https://mozilla.org/a/b/c/", true},
		{"https://mozilla.org/*/b/*/", []string{"https://mozilla.org/*/b/*/"}, "https://mozilla.org/d/b/f/", true},
		{"https://mozilla.org/*/b/*/", []string{"https://mozilla.org/*/b/*/"}, "https://mozilla.org/a/b/c/d/", true},
		{"https://mozilla.org/*/b/*/", []string{"https://mozilla.org/*/b/*/"}, "https://mozilla.org/a/b/c/d/#section1", true},
		{"https://mozilla.org/*/b/*/", []string{"https://mozilla.org/*/b/*/"}, "https://mozilla.org/a/b/c/d/?foo=/", true},
		{"https://mozilla.org/*/b/*/", []string{"https://mozilla.org/*/b/*/"}, "https://mozilla.org/a?foo=21314&bar=/b/&extra=c/", true},
		{"https://mozilla.org/*/b/*/", []string{"https://mozilla.org/*/b/*/"}, "https://mozilla.org/b/*/", false},
		{"https://mozilla.org/*/b/*/", []string{"https://mozilla.org/*/b/*/"}, "https://mozilla.org/a/b/", false},
		{"https://mozilla.org/*/b/*/", []string{"https://mozilla.org/*/b/*/"}, "https://mozilla.org/a/b/c/d/?foo=bar", false},

		{"file:///blah/*", []string{"file:///blah/*"}, "file:///blah/", true},
		{"file:///blah/*", []string{"file:///blah/*"}, "file:///blah/bleh", true},
		{"file:///blah/*", []string{"file:///blah/*"}, "file:///bleh/", false},
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
func TestMatchPattern_MatchesChromeExtensions(t *testing.T) {
	tests := []struct {
		name string
		p    []string
		url  string
		want bool
	}{
		// https://developer.chrome.com/docs/extensions/mv3/match_patterns/
		{"all_urls", []string{"<all_urls>"}, "http://example.org/foo/bar.html", true},
		{"all_urls", []string{"<all_urls>"}, "file:///bar/baz.html", true},

		{"http://*/*", []string{"http://*/*"}, "http://www.google.com/", true},
		{"http://*/*", []string{"http://*/*"}, "http://example.org/foo/bar.html", true},

		{"http://*/foo*", []string{"http://*/foo*"}, "http://example.com/foo/bar.html", true},
		{"http://*/foo*", []string{"http://*/foo*"}, "http://www.google.com/foo", true},

		{"https://*.google.com/foo*bar", []string{"https://*.google.com/foo*bar"}, "https://www.google.com/foo/baz/bar", true},
		{"https://*.google.com/foo*bar", []string{"https://*.google.com/foo*bar"}, "https://docs.google.com/foobar", true},

		{"http://example.org/foo/bar.html", []string{"http://example.org/foo/bar.html"}, "http://example.org/foo/bar.html", true},

		{"file:///foo*", []string{"file:///foo*"}, "file:///foo/bar.html", true},
		{"file:///foo*", []string{"file:///foo*"}, "file:///foo", true},

		{"http://127.0.0.1/*", []string{"http://127.0.0.1/*"}, "http://127.0.0.1/", true},
		{"http://127.0.0.1/*", []string{"http://127.0.0.1/*"}, "http://127.0.0.1/foo/bar.html", true},

		{"*://mail.google.com/*", []string{"*://mail.google.com/*"}, "http://mail.google.com/foo/baz/bar", true},
		{"*://mail.google.com/*", []string{"*://mail.google.com/*"}, "https://mail.google.com/foobar", true},

		// Needs extra care
		// {"urn:*", []string{"urn:*"}, "urn:uuid:54723bea-c94e-480e-80c8-a69846c3f582", true},
		// {"urn:*", []string{"urn:*"}, "urn:uuid:cfa40aff-07df-45b2-9f95-e023bcf4a6da", true},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf(" `%s` %s", tt.name, tt.url), func(t *testing.T) {
			m, err := NewMatchPattern(tt.p, GetChromeExtensionMatchSet())
			if err != nil {
				t.Fatalf(err.Error())
			}
			if got, _ := m.Matches(tt.url); got != tt.want {
				t.Errorf("MatchPattern.Matches() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isValidPathPattern(t *testing.T) {
	tests := []struct {
		pathPattern string
		path        string
		want        bool
	}{
		{"/*/b/*/", "/a/b/c/d/", true},
		{"/", "/a", false},
		{"/test", "/test", true},
		{"/test/*", "/test/one", true},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf(": `%s` matches `%s`", tt.path, tt.pathPattern), func(t *testing.T) {
			got, err := isValidPathPattern(tt.pathPattern, tt.path)
			if err != nil {
				t.Fatalf("isValidPathPattern() error = %v", err)
			}
			if got != tt.want {
				t.Errorf("isValidPathPattern() = %v, want %v", got, tt.want)
			}
		})
	}
}
