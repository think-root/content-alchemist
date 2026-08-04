package parser

import "testing"

func TestNormalizeRepoURL(t *testing.T) {
	const canonical = "https://github.com/Xyri1/mailsink"

	valid := []struct {
		name string
		raw  string
		want string
	}{
		{name: "canonical", raw: canonical, want: canonical},
		{name: "tracking query", raw: canonical + "?twclid=25d0t936o3n1ofqchrde54lvzk", want: canonical},
		{name: "fragment", raw: canonical + "#readme", want: canonical},
		{name: "trailing slash", raw: canonical + "/", want: canonical},
		{name: "git suffix", raw: canonical + ".git", want: canonical},
		{name: "http scheme", raw: "http://github.com/Xyri1/mailsink", want: canonical},
		{name: "www host", raw: "https://www.github.com/Xyri1/mailsink", want: canonical},
		{name: "no scheme", raw: "github.com/Xyri1/mailsink", want: canonical},
		{name: "deep path", raw: canonical + "/tree/main/src", want: canonical},
		{name: "deep path with query", raw: canonical + "/blob/main/README.md?plain=1", want: canonical},
		{name: "surrounding whitespace", raw: "  " + canonical + "  ", want: canonical},
	}

	for _, tt := range valid {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NormalizeRepoURL(tt.raw)
			if err != nil {
				t.Fatalf("NormalizeRepoURL(%q) returned error: %v", tt.raw, err)
			}
			if got != tt.want {
				t.Errorf("NormalizeRepoURL(%q) = %q, want %q", tt.raw, got, tt.want)
			}
		})
	}

	invalid := []struct {
		name string
		raw  string
	}{
		{name: "empty", raw: ""},
		{name: "owner only", raw: "https://github.com/Xyri1"},
		{name: "other host", raw: "https://gitlab.com/Xyri1/mailsink"},
		{name: "reserved path", raw: "https://github.com/orgs/think-root/repositories"},
		{name: "trending", raw: "https://github.com/trending/go"},
		{name: "not a url", raw: "just some text"},
		{name: "unsupported scheme", raw: "ftp://github.com/Xyri1/mailsink"},
	}

	for _, tt := range invalid {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NormalizeRepoURL(tt.raw)
			if err == nil {
				t.Errorf("NormalizeRepoURL(%q) = %q, want error", tt.raw, got)
			}
		})
	}
}
