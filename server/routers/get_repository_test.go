package routers

import "testing"

func TestParseMultilingualTextFallsBackWhenRequestedLanguageMissing(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		language string
		want     string
	}{
		{
			name:     "plain text is returned for any requested language",
			text:     "Plain repository description",
			language: "uk",
			want:     "Plain repository description",
		},
		{
			name:     "single language text falls back to its available content",
			text:     "(en)English repository description",
			language: "uk",
			want:     "English repository description",
		},
		{
			name:     "multilingual text prefers requested language when available",
			text:     "===(en)English repository description===(uk)Український опис===",
			language: "en",
			want:     "English repository description",
		},
		{
			name:     "multilingual text falls back to Ukrainian when requested language is missing",
			text:     "===(en)English repository description===(uk)Український опис===",
			language: "pl",
			want:     "Український опис",
		},
		{
			name:     "multilingual text falls back to first available language without Ukrainian",
			text:     "===(en)English repository description===(de)Deutsche Beschreibung===",
			language: "uk",
			want:     "English repository description",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseMultilingualText(tt.text, tt.language)
			if err != nil {
				t.Fatalf("expected fallback text without error, got %v", err)
			}
			if got != tt.want {
				t.Fatalf("unexpected text: got %q, want %q", got, tt.want)
			}
		})
	}
}
