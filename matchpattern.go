package matchpattern

import (
	"fmt"
	"net/url"
)

const (
	allURLPattern = "<all_urls>"
	allPattern    = "*://*/*"
)

type matcherType int
type matcherFunc func(url *url.URL) (bool, error)

const (
	allURL matcherType = iota
	pattern
)

type MatchSet struct {
	allUrlSchemes  []string
	allowedSchemes []string
}

type matcher struct {
	mType   matcherType
	raw     string
	matcher matcherFunc
}

type MatchPattern struct {
	matchers []matcher
}

func NewMatchPattern(patterns []string, ms MatchSet) (*MatchPattern, error) {
	matchers, err := getMatchers(patterns, ms)
	if err != nil {
		return nil, err
	}
	return &MatchPattern{matchers}, nil
}

func (m MatchPattern) Matches(address string) (bool, error) {
	addressUrl, err := url.Parse(address)
	if err != nil {
		return false, err
	}
	return m.MatchesUrl(addressUrl)
}

func (m MatchPattern) MatchesUrl(address *url.URL) (bool, error) {
	for i := 0; i < len(m.matchers); i++ {
		matches, err := m.matchers[i].matcher(address)
		if matches || err != nil {
			return matches, err
		}
	}
	return false, nil
}

// GetDefaultMatchSet returns the allowed the extension schemes
// as defined in https://developer.mozilla.org/en-US/docs/Mozilla/Add-ons/WebExtensions/Match_patterns
func GetDefaultMatchSet() MatchSet {
	return MatchSet{
		[]string{"http", "https", "ws", "wss", "ftp", "ftps", "data", "file"},
		[]string{"http", "https", "ws", "wss"},
	}
}

// GetChromeExtensionMatchSet returns the allowed the extension schemes
// as defined in https://developer.chrome.com/docs/extensions/mv3/match_patterns/
func GetChromeExtensionMatchSet() MatchSet {
	return MatchSet{
		[]string{"http", "https", "file", "ftp", "urn"},
		[]string{"http", "https"},
	}
}

func getMatchers(patterns []string, ms MatchSet) ([]matcher, error) {
	result := make([]matcher, len(patterns))
	var err error
	for i := 0; i < len(patterns); i++ {
		if patterns[i] == allURLPattern {
			if result[i], err = getSchemeMatcher(patterns[i], ms.allUrlSchemes); err != nil {
				return result, err
			}
		} else if patterns[i] == allPattern {
			// TODO: may not need a special case
			if result[i], err = getSchemeMatcher(patterns[i], ms.allowedSchemes); err != nil {
				return result, err
			}
		} else {
			return result, fmt.Errorf("unknown pattern %s", patterns[i])
		}
	}
	return result, nil
}

func getSchemeMatcher(pattern string, schemes []string) (matcher, error) {
	return matcher{
		allURL,
		pattern,
		func(url *url.URL) (bool, error) {
			return contains(url.Scheme, schemes), nil
		},
	}, nil
}

func contains(s string, list []string) bool {
	for i := 0; i < len(list); i++ {
		if s == list[i] {
			return true
		}
	}
	return false
}
