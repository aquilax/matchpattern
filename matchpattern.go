// Package matchpattern implements Match patterns validation primarily used in browser extensions
// as defined in https://developer.mozilla.org/en-US/docs/Mozilla/Add-ons/WebExtensions/Match_patterns
package matchpattern

import (
	"fmt"
	"net/url"
	"strings"
)

const (
	// AllURLPattern is a magic pattern, that matches all urls with allowed schemes
	AllURLPattern   = "<all_urls>"
	allPattern      = "*://*/*"
	matchAll        = "*"
	matchSubdomains = "*."
	matchAllPath    = "/*"

	pathSeparator  = "/"
	querySeparator = "?"
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
	matchers []*matcher
}

// New creates new pattern matcher given list of patterns and match set
func New(patterns []string, ms MatchSet) (*MatchPattern, error) {
	matchers, err := getMatchers(patterns, ms)
	if err != nil {
		return nil, err
	}
	return &MatchPattern{matchers}, nil
}

// Matches validates the address against the provided patterns and returns true and no error if it matches
func (mp MatchPattern) Matches(address string) (bool, error) {
	addressUrl, err := url.Parse(address)
	if err != nil {
		return false, err
	}
	return mp.MatchesUrl(addressUrl)
}

// MatchesUrl validates the url against the provided patterns and returns true and no error if it matches
func (mp MatchPattern) MatchesUrl(address *url.URL) (bool, error) {
	for i := 0; i < len(mp.matchers); i++ {
		matches, err := mp.matchers[i].matcher(address)
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

func getMatchers(patterns []string, ms MatchSet) ([]*matcher, error) {
	result := make([]*matcher, len(patterns))
	var err error
	for i := 0; i < len(patterns); i++ {
		if patterns[i] == AllURLPattern {
			if result[i], err = getSchemeMatcher(patterns[i], ms.allUrlSchemes); err != nil {
				return result, err
			}
		} else if patterns[i] == allPattern {
			// TODO: may not need a special case
			if result[i], err = getSchemeMatcher(patterns[i], ms.allowedSchemes); err != nil {
				return result, err
			}
		} else {
			if result[i], err = getFullMatcher(patterns[i], ms); err != nil {
				return result, err
			}
		}
	}
	return result, nil
}

func getSchemeMatcher(pattern string, schemes []string) (*matcher, error) {
	return &matcher{
		allURL,
		pattern,
		func(url *url.URL) (bool, error) {
			return contains(url.Scheme, schemes), nil
		},
	}, nil
}

func getFullMatcher(pattern string, ms MatchSet) (*matcher, error) {
	segments := strings.Split(pattern, "://")
	if len(segments) != 2 {
		return nil, fmt.Errorf("invalid pattern %s", pattern)
	}
	if segments[0] != matchAll && !contains(segments[0], ms.allUrlSchemes) {
		return nil, fmt.Errorf("unsupported scheme %s", segments[0])
	}

	return &matcher{
		allURL,
		pattern,
		func(url *url.URL) (bool, error) {
			// validate scheme
			if !isValidSchemePattern(url, segments[0], ms.allowedSchemes) {
				return false, nil
			}
			// validate the rest
			return isValidHostPath(url, segments[1])
		},
	}, nil
}

func isValidSchemePattern(url *url.URL, scheme string, schemes []string) bool {
	if scheme == matchAll {
		return contains(url.Scheme, schemes)
	}
	return url.Scheme == scheme
}

func isValidHostPath(url *url.URL, hostPathPattern string) (bool, error) {
	var host, path string
	firstSlashIndex := strings.Index(hostPathPattern, pathSeparator)
	if firstSlashIndex == -1 {
		host = hostPathPattern
		return isValidHost(url, host)
	}
	host = hostPathPattern[:firstSlashIndex]
	path = hostPathPattern[firstSlashIndex:]
	if ok, err := isValidHost(url, host); !ok || err != nil {
		return ok, err
	}

	return isValidPath(url, path)
}

func contains(s string, list []string) bool {
	for i := 0; i < len(list); i++ {
		if s == list[i] {
			return true
		}
	}
	return false
}

func isValidHost(url *url.URL, hostPattern string) (bool, error) {
	if hostPattern == matchAll {
		return true, nil
	}
	if strings.HasPrefix(hostPattern, matchSubdomains) {
		// *. matches all subdomains
		hostSuffix := hostPattern[2:]
		return strings.HasSuffix(url.Host, hostSuffix), nil
	}
	return hostPattern == url.Host, nil
}

func isValidPath(url *url.URL, pathPattern string) (bool, error) {
	if pathPattern == matchAllPath {
		return true, nil
	}
	path := url.Path
	if path == "" {
		// always prefix path with path separator
		path = pathSeparator
	}
	if url.RawQuery != "" {
		// add the query string
		path = path + querySeparator + url.RawQuery
	}
	return isValidPathPattern(pathPattern, path)
}

func isValidPathPattern(pathPattern, path string) (bool, error) {
	// TODO: Refactor me
	segments := strings.Split(pathPattern, matchAll)
	rem := path
	if segments[len(segments)-1] != "" {
		// consume the capping match
		if !strings.HasSuffix(path, segments[len(segments)-1]) {
			return false, nil
		}
		index := len(path) - len(segments[len(segments)-1])
		if index == -1 {
			return false, nil
		}
		rem = rem[:index]
		segments[len(segments)-1] = ""
	}
	for i, segment := range segments {
		if segment == "" && i == len(segments)-1 {
			// trailing matchall
			rem = ""
		}
		index := strings.Index(rem, segment)
		if index == -1 {
			return false, nil
		}
		rem = rem[index+len(segment):]
	}
	return rem == "", nil
}
