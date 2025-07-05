package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
)

// LanguageCode represents a language code structure from the API
type LanguageCode struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

var (
	languageCodesCache map[string]LanguageCode
	cacheExpiry        time.Time
	cacheMutex         sync.RWMutex
	cacheValidFor      = 24 * time.Hour // Cache for 24 hours
)

// ValidateLanguageCodes validates language codes against the external API
func ValidateLanguageCodes(languageCodes []string) error {
	if len(languageCodes) == 0 {
		return nil
	}

	validCodes, err := getValidLanguageCodes()
	if err != nil {
		return fmt.Errorf("failed to fetch valid language codes: %w", err)
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

// getValidLanguageCodes fetches and caches language codes from the external API
func getValidLanguageCodes() (map[string]LanguageCode, error) {
	cacheMutex.RLock()
	if languageCodesCache != nil && time.Now().Before(cacheExpiry) {
		defer cacheMutex.RUnlock()
		return languageCodesCache, nil
	}
	cacheMutex.RUnlock()

	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	// Double-check after acquiring write lock
	if languageCodesCache != nil && time.Now().Before(cacheExpiry) {
		return languageCodesCache, nil
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get("https://gist.githubusercontent.com/Josantonius/b455e315bc7f790d14b136d61d9ae469/raw/language-codes.json")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch language codes: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch language codes, status: %d", resp.StatusCode)
	}

	var languageCodesMap map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&languageCodesMap); err != nil {
		return nil, fmt.Errorf("failed to decode language codes: %w", err)
	}

	// Create a map for faster lookups
	codeMap := make(map[string]LanguageCode)
	for code, name := range languageCodesMap {
		codeMap[strings.ToLower(code)] = LanguageCode{
			Code: code,
			Name: name,
		}
	}

	languageCodesCache = codeMap
	cacheExpiry = time.Now().Add(cacheValidFor)

	return codeMap, nil
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