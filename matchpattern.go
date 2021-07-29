package matchpattern

import (
	"net/url"
)

const (
	allURLs = "<all_urls>"
)

type matcherType int
type matcherFunc func(url *url.URL) (bool, error)

const (
	allURL matcherType = iota
)

func GetDefaultProtocols() []string {
	return []string{"http", "https"}
}

type matcher struct {
	mType   matcherType
	raw     string
	matcher matcherFunc
}

type MatchPattern struct {
	matchers []matcher
}

func NewMatchPattern(patterns []string, schemes []string) (*MatchPattern, error) {
	matchers, err := getMatchers(patterns, schemes)
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

func getMatchers(patterns, protocols []string) ([]matcher, error) {
	result := make([]matcher, len(patterns))
	var err error
	for i := 0; i < len(patterns); i++ {
		if patterns[i] == allURLs {
			result[i], err = getAllMatcher(patterns[i], protocols)
			if err != nil {
				return result, err
			}
		}
	}
	return result, nil
}

func getAllMatcher(pattern string, schemes []string) (matcher, error) {
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
