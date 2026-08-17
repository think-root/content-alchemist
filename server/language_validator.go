package server

import (
	"embed"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
)

// LanguageCode represents a single ISO 639-1 language code and its English name
type LanguageCode struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

// The ISO 639-1 list is baked into the binary on purpose. Fetching it at runtime
// made every request that carried a text_language depend on an external host,
// so a single 429 from that host turned into 400s across get-repository,
// manual-generate and auto-generate.
//
//go:embed language_codes.json
var languageCodesFS embed.FS

var (
	languageCodesOnce sync.Once
	languageCodes     map[string]LanguageCode
	languageCodesErr  error
)

// ValidateLanguageCodes validates language codes against the embedded ISO 639-1 list
func ValidateLanguageCodes(languageCodes []string) error {
	if len(languageCodes) == 0 {
		return nil
	}

	validCodes, err := getValidLanguageCodes()
	if err != nil {
		return fmt.Errorf("failed to load valid language codes: %w", err)
	}

	var invalidCodes []string
	for _, code := range languageCodes {
		code = strings.TrimSpace(strings.ToLower(code))
		if code == "" {
			continue
		}

		if _, exists := validCodes[code]; !exists {
			invalidCodes = append(invalidCodes, code)
		}
	}

	if len(invalidCodes) > 0 {
		return fmt.Errorf("invalid language codes: %s", strings.Join(invalidCodes, ", "))
	}

	return nil
}

// getValidLanguageCodes parses the embedded language code list once and reuses it
func getValidLanguageCodes() (map[string]LanguageCode, error) {
	languageCodesOnce.Do(func() {
		raw, err := languageCodesFS.ReadFile("language_codes.json")
		if err != nil {
			languageCodesErr = fmt.Errorf("failed to read embedded language codes: %w", err)
			return
		}

		var languageCodesMap map[string]string
		if err := json.Unmarshal(raw, &languageCodesMap); err != nil {
			languageCodesErr = fmt.Errorf("failed to decode embedded language codes: %w", err)
			return
		}

		if len(languageCodesMap) == 0 {
			languageCodesErr = fmt.Errorf("embedded language codes list is empty")
			return
		}

		// Create a map for faster lookups
		codeMap := make(map[string]LanguageCode, len(languageCodesMap))
		for code, name := range languageCodesMap {
			codeMap[strings.ToLower(code)] = LanguageCode{
				Code: code,
				Name: name,
			}
		}

		languageCodes = codeMap
	})

	return languageCodes, languageCodesErr
}

// ParseLanguageCodes parses comma-separated language codes string
func ParseLanguageCodes(languageCodesStr string) []string {
	if languageCodesStr == "" {
		return []string{"uk"} // Default to Ukrainian
	}

	codes := strings.Split(languageCodesStr, ",")
	var parsedCodes []string

	for _, code := range codes {
		code = strings.TrimSpace(code)
		if code != "" {
			parsedCodes = append(parsedCodes, code)
		}
	}

	if len(parsedCodes) == 0 {
		return []string{"uk"} // Default to Ukrainian if no valid codes
	}

	return parsedCodes
}

// BuildMultilingualPrompt builds the prompt instruction for multilingual generation
func BuildMultilingualPrompt(languageCodes []string) string {
	if len(languageCodes) == 1 {
		return fmt.Sprintf("Generate your response in the following format: (%s)your_response_text", languageCodes[0])
	}

	var examples []string
	for _, code := range languageCodes {
		examples = append(examples, fmt.Sprintf("===(%s)text_in_%s", code, code))
	}

	return fmt.Sprintf("Generate your response in the following multilingual format: %s===", strings.Join(examples, ""))
}
