package server

import "testing"

func TestCleanMultilingualText(t *testing.T) {
	const wantTwo = "===(en)A terminal spinner.===(uk)Термінальний спінер.==="

	tests := []struct {
		name string
		text string
		want string
	}{
		{
			name: "already canonical",
			text: "===(en)A terminal spinner.===(uk)Термінальний спінер.===",
			want: wantTwo,
		},
		{
			name: "spaces around markers",
			text: "=== (en) A terminal spinner. === (uk) Термінальний спінер. ===",
			want: wantTwo,
		},
		{
			name: "two-character separator",
			text: "==(en)A terminal spinner.==(uk)Термінальний спінер.",
			want: wantTwo,
		},
		{
			name: "missing trailing separator",
			text: "===(en)A terminal spinner.===(uk)Термінальний спінер.",
			want: wantTwo,
		},
		{
			name: "newlines between segments",
			text: "===(en)A terminal spinner.\n===(uk)\nТермінальний спінер.\n===",
			want: wantTwo,
		},
		{
			name: "leading whitespace before first marker",
			text: "\n  ===(en)A terminal spinner.===(uk)Термінальний спінер.===",
			want: wantTwo,
		},
		{
			name: "single language",
			text: "===(uk)Термінальний спінер.===",
			want: "===(uk)Термінальний спінер.===",
		},
		{
			name: "plain text is untouched",
			text: "A terminal spinner for showing task progress.",
			want: "A terminal spinner for showing task progress.",
		},
		{
			name: "text before the first marker is never dropped",
			text: "Ось відповідь: ===(en)A terminal spinner.===",
			want: "Ось відповідь: ===(en)A terminal spinner.===",
		},
		{
			name: "empty",
			text: "",
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CleanMultilingualText(tt.text)
			if got != tt.want {
				t.Errorf("CleanMultilingualText()\n got: %q\nwant: %q", got, tt.want)
			}
		})
	}
}

// The point of the cleanup is that the result is parseable, so that every
// language segment is validated on its own instead of as one blob.
func TestCleanMultilingualTextIsParseable(t *testing.T) {
	cleaned := CleanMultilingualText("== (en) A terminal spinner. == (uk) Термінальний спінер.")

	if !IsMultilingualText(cleaned) {
		t.Fatalf("cleaned text is not recognized as multilingual: %q", cleaned)
	}

	langMap := ParseMultilingualText(cleaned)
	if len(langMap) != 2 {
		t.Fatalf("expected 2 language segments, got %d: %#v", len(langMap), langMap)
	}
	if langMap["en"] != "A terminal spinner." {
		t.Errorf("en segment = %q", langMap["en"])
	}
	if langMap["uk"] != "Термінальний спінер." {
		t.Errorf("uk segment = %q", langMap["uk"])
	}
}
