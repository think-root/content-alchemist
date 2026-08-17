package server

import "testing"

func TestValidateLanguageCodes(t *testing.T) {
	tests := []struct {
		name    string
		codes   []string
		wantErr bool
	}{
		{name: "empty slice", codes: nil, wantErr: false},
		{name: "blank code is skipped", codes: []string{"  "}, wantErr: false},
		{name: "single valid", codes: []string{"en"}, wantErr: false},
		{name: "multiple valid", codes: []string{"uk", "en", "de"}, wantErr: false},
		{name: "uppercase is normalised", codes: []string{"EN"}, wantErr: false},
		{name: "padded is trimmed", codes: []string{" uk "}, wantErr: false},
		{name: "unknown code", codes: []string{"xx"}, wantErr: true},
		{name: "one unknown among valid", codes: []string{"en", "zz"}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateLanguageCodes(tt.codes)
			if tt.wantErr && err == nil {
				t.Fatalf("ValidateLanguageCodes(%v) = nil, want error", tt.codes)
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("ValidateLanguageCodes(%v) = %v, want nil", tt.codes, err)
			}
		})
	}
}

// The list must resolve from the embedded file alone: validation used to depend
// on a remote gist, and a 429 from that host broke every language-aware request.
func TestGetValidLanguageCodesUsesEmbeddedList(t *testing.T) {
	codes, err := getValidLanguageCodes()
	if err != nil {
		t.Fatalf("getValidLanguageCodes() error = %v", err)
	}

	if len(codes) < 180 {
		t.Fatalf("got %d language codes, want the full ISO 639-1 list", len(codes))
	}

	for _, want := range []string{"en", "uk"} {
		if _, ok := codes[want]; !ok {
			t.Errorf("language code %q missing from embedded list", want)
		}
	}

	if got := codes["en"].Name; got != "English" {
		t.Errorf("codes[\"en\"].Name = %q, want %q", got, "English")
	}
}
