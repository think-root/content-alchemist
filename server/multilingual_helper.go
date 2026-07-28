package server

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
	"unicode/utf8"
)

// IsMultilingualText checks if text follows the multilingual format
// Format: ===(lang_code)content===(lang_code)content===
func IsMultilingualText(text string) bool {
	if text == "" {
		return false
	}

	// Check if text contains the multilingual pattern
	pattern := `===\([a-z]{2,3}\).*?===`
	matched, err := regexp.MatchString(pattern, text)
	if err != nil {
		return false
	}

	return matched && strings.HasPrefix(text, "===(") && strings.HasSuffix(text, "===")
}

// ParseMultilingualText extracts language-content pairs from multilingual text
// Returns a map where keys are language codes and values are content
func ParseMultilingualText(text string) map[string]string {
	result := make(map[string]string)

	if !IsMultilingualText(text) {
		return result
	}

	// Split by === and process each segment
	segments := strings.Split(text, "===")

	for _, segment := range segments {
		segment = strings.TrimSpace(segment)
		if segment == "" {
			continue
		}

		// Check if segment starts with (lang_code)
		if strings.HasPrefix(segment, "(") {
			// Find the closing parenthesis
			closeIdx := strings.Index(segment, ")")
			if closeIdx > 1 {
				langCode := segment[1:closeIdx]
				content := strings.TrimSpace(segment[closeIdx+1:])

				// Validate language code format
				if len(langCode) >= 2 && len(langCode) <= 3 && content != "" {
					result[langCode] = content
				}
			}
		}
	}

	return result
}

// BuildMultilingualText reconstructs multilingual text from language-content map
// Languages are sorted alphabetically for consistency
func BuildMultilingualText(langMap map[string]string) string {
	if len(langMap) == 0 {
		return ""
	}

	// Sort languages alphabetically for consistent output
	var languages []string
	for lang := range langMap {
		languages = append(languages, lang)
	}
	sort.Strings(languages)

	var parts []string
	for _, lang := range languages {
		content := langMap[lang]
		if content != "" {
			parts = append(parts, fmt.Sprintf("===(%s)%s", lang, content))
		}
	}

	return strings.Join(parts, "") + "==="
}

// UpdateLanguageInText updates specific language in multilingual text
// If text is not multilingual, it handles conversion based on target language
func UpdateLanguageInText(existingText, newText, targetLang string) (string, error) {
	// Validate target language
	if err := ValidateLanguageCodes([]string{targetLang}); err != nil {
		return "", fmt.Errorf("invalid target language: %w", err)
	}

	// Validate new text
	newText = strings.TrimSpace(newText)
	if newText == "" {
		return "", fmt.Errorf("new text cannot be empty")
	}

	if utf8.RuneCountInString(newText) > 1000 {
		return "", fmt.Errorf("text must not exceed 1000 characters")
	}

	if !utf8.ValidString(newText) {
		return "", fmt.Errorf("text must be valid UTF-8")
	}

	if IsMultilingualText(existingText) {
		// Parse existing multilingual content
		langMap := ParseMultilingualText(existingText)

		// Update the specific language
		langMap[targetLang] = newText

		// Rebuild multilingual text
		return BuildMultilingualText(langMap), nil
	} else {
		// Handle old format text
		existingText = strings.TrimSpace(existingText)

		if targetLang == "uk" {
			// Replace entire text for Ukrainian (backward compatibility)
			return newText, nil
		} else {
			// Convert to multilingual format
			langMap := map[string]string{
				"uk":       existingText,
				targetLang: newText,
			}
			return BuildMultilingualText(langMap), nil
		}
	}
}

// ExtractLanguageFromText extracts content for a specific language from multilingual text
// Returns the content and a boolean indicating if the language was found
func ExtractLanguageFromText(text, targetLang string) (string, bool) {
	if !IsMultilingualText(text) {
		// For old format, assume it's Ukrainian
		if targetLang == "uk" {
			return strings.TrimSpace(text), true
		}
		return "", false
	}

	langMap := ParseMultilingualText(text)
	content, exists := langMap[targetLang]
	return content, exists
}

// GetAvailableLanguages returns a list of available languages in multilingual text
func GetAvailableLanguages(text string) []string {
	if !IsMultilingualText(text) {
		return []string{"uk"} // Assume old format is Ukrainian
	}

	langMap := ParseMultilingualText(text)
	var languages []string
	for lang := range langMap {
		languages = append(languages, lang)
	}

	sort.Strings(languages)
	return languages
}

// ValidateMultilingualContent validates all languages and content in multilingual text
func ValidateMultilingualContent(text string) error {
	if text == "" {
		return fmt.Errorf("text cannot be empty")
	}

	if !IsMultilingualText(text) {
		// Validate as single text
		if utf8.RuneCountInString(text) > 1000 {
			return fmt.Errorf("text must not exceed 1000 characters")
		}
		if !utf8.ValidString(text) {
			return fmt.Errorf("text must be valid UTF-8")
		}
		return nil
	}

	langMap := ParseMultilingualText(text)
	if len(langMap) == 0 {
		return fmt.Errorf("no valid language content found in multilingual text")
	}

	// Validate each language code
	var languages []string
	for lang := range langMap {
		languages = append(languages, lang)
	}

	if err := ValidateLanguageCodes(languages); err != nil {
		return fmt.Errorf("invalid language codes in multilingual text: %w", err)
	}

	// Validate each content
	for lang, content := range langMap {
		if content == "" {
			return fmt.Errorf("content for language '%s' cannot be empty", lang)
		}
		if utf8.RuneCountInString(content) > 1000 {
			return fmt.Errorf("content for language '%s' must not exceed 1000 characters", lang)
		}
		if !utf8.ValidString(content) {
			return fmt.Errorf("content for language '%s' must be valid UTF-8", lang)
		}
	}

	return nil
}

// minDescriptionLength is the minimum acceptable length (in runes) of a single
// generated description segment. The prompt asks for 200-500 chars, so anything
// well below that is almost certainly a truncated answer or a refusal.
const minDescriptionLength = 120

// refusalMarkers are lowercase substrings that strongly indicate the LLM
// refused to produce a description (e.g. because the README had no real
// content) instead of generating one. Extend this list as new patterns appear.
var refusalMarkers = []string{
	"надайте вміст",
	"надайте текст",
	"надайте більше",
	"не дозволяє мені",
	"не можу проаналізувати",
	"не можу згенерувати",
	"не можу створити",
	"будь ласка, надайте",
	"я не маю доступу",
	"немає доступу",
	"недостатньо інформації",
	"недостатньо даних",
	"readme.md",
	"i cannot",
	"i can't",
	"unable to",
	"please provide",
}

// IsLikelyRefusal reports whether the given text looks like an LLM refusal
// rather than a genuine repository description.
func IsLikelyRefusal(text string) bool {
	lower := strings.ToLower(text)
	for _, marker := range refusalMarkers {
		if strings.Contains(lower, marker) {
			return true
		}
	}
	return false
}

// ValidateGeneratedDescription checks that an LLM-generated description (in
// single or multilingual format) is worth storing: non-empty, long enough, and
// not a refusal message. It returns a descriptive error when the text should be
// rejected.
func ValidateGeneratedDescription(text string) error {
	trimmed := strings.TrimSpace(text)
	if trimmed == "" {
		return fmt.Errorf("generated description is empty")
	}

	// Collect the per-language segments (or the whole text for the old format).
	var segments []string
	if IsMultilingualText(trimmed) {
		langMap := ParseMultilingualText(trimmed)
		if len(langMap) == 0 {
			return fmt.Errorf("no valid language content found in generated description")
		}
		for lang, content := range langMap {
			segments = append(segments, fmt.Sprintf("[%s] %s", lang, content))
		}
	} else {
		segments = []string{trimmed}
	}

	for _, segment := range segments {
		content := strings.TrimSpace(segment)
		if idx := strings.Index(content, "] "); strings.HasPrefix(content, "[") && idx != -1 {
			content = content[idx+2:]
		}

		if utf8.RuneCountInString(content) < minDescriptionLength {
			return fmt.Errorf("generated description segment too short (%d < %d chars)",
				utf8.RuneCountInString(content), minDescriptionLength)
		}
		if IsLikelyRefusal(content) {
			return fmt.Errorf("generated description looks like an LLM refusal")
		}
	}

	return nil
}

// CleanMultilingualText cleans and standardizes the multilingual text format.
func CleanMultilingualText(text string) string {
	// Don't process empty strings.
	if text == "" {
		return ""
	}

	// Trim leading/trailing whitespace from the whole string.
	processedText := strings.TrimSpace(text)

	// Check if the text is intended to be in the multilingual format.
	// A simple check for the initial pattern is enough.
	if !strings.Contains(processedText, "===(") {
		// If it's not a multilingual text, return it as is.
		return text
	}

	// Regex to remove spaces between '===' and '(en)'.
	re1 := regexp.MustCompile(`===\s*\n*\s*\(([a-zA-Z]{2,3})\)`)
	processedText = re1.ReplaceAllString(processedText, "===($1)")

	// Regex to remove spaces between '(en)' and '==='
	re2 := regexp.MustCompile(`\)\s*\n*\s*===`)
	processedText = re2.ReplaceAllString(processedText, ")===")

	// Ensure the text ends with "===", but don't add it if it's already there.
	if !strings.HasSuffix(processedText, "===") {
		processedText += "==="
	}

	return processedText
}
