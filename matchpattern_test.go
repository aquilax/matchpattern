package matchpattern

import "testing"

func TestMatchPattern_Matches(t *testing.T) {
	tests := []struct {
		name string
		p    []string
		url  string
		want bool
	}{
		{"all_urls", []string{"<all_urls>"}, "http://example.org/", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, _ := NewMatchPattern(tt.p, GetDefaultProtocols())
			if got, _ := m.Matches(tt.url); got != tt.want {
				t.Errorf("MatchPattern.Matches() = %v, want %v", got, tt.want)
			}
		})
	}
}
