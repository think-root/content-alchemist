package routers

import (
	"content-alchemist/database"
	"content-alchemist/server"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"
)

type getRepositoryRequestBody struct {
	Limit        int    `json:"limit"`
	Posted       *bool  `json:"posted"`
	SortBy       string `json:"sort_by"`
	SortOrder    string `json:"sort_order"`
	Page         int    `json:"page"`
	PageSize     int    `json:"page_size"`
	TextLanguage string `json:"text_language"`
}

type getRepositoryItem struct {
	ID         int64      `json:"id"`
	Posted     bool       `json:"posted"`
	URL        string     `json:"url"`
	Text       string     `json:"text"`
	DateAdded  *time.Time `json:"date_added"`
	DatePosted *time.Time `json:"date_posted"`
}

type getRepositoryResponse struct {
	All        int                 `json:"all"`
	Posted     int                 `json:"posted"`
	Unposted   int                 `json:"unposted"`
	Items      []getRepositoryItem `json:"items"`
	Page       int                 `json:"page"`
	PageSize   int                 `json:"page_size"`
	TotalPages int                 `json:"total_pages"`
	TotalItems int                 `json:"total_items"`
}

var trueVal = true
var falseVal = false

func countRepositories() (all, posted, unposted int, err error) {
	if all, err = database.CountRepositories(nil); err != nil {
		return
	}
	if posted, err = database.CountRepositories(&trueVal); err != nil {
		return
	}
	unposted, err = database.CountRepositories(&falseVal)
	return
}

// ParseMultilingualText parses multilingual text and extracts content for the specified language
func ParseMultilingualText(text, languageCode string) (string, error) {
	if text == "" {
		return "", nil
	}

	// Check if text is in old format (plain text without language markers)
	if !strings.Contains(text, "(") || !strings.Contains(text, ")") {
		// Old format - plain text
		if languageCode == "" {
			return text, nil // Return as is if no language specified
		}
		return "", fmt.Errorf("no text available for language: %s", languageCode)
	}

	// Check if text is in single language format: (code)text
	singleLangPattern := `^\(([a-z]{2})\)(.*)$`
	if matched, _ := regexp.MatchString(singleLangPattern, text); matched {
		re := regexp.MustCompile(singleLangPattern)
		matches := re.FindStringSubmatch(text)
		if len(matches) == 3 {
			textLang := matches[1]
			textContent := matches[2]
			
			if languageCode == "" {
				return textContent, nil // Return content if no specific language requested
			}
			if textLang == languageCode {
				return textContent, nil
			}
			return "", fmt.Errorf("no text available for language: %s", languageCode)
		}
	}

	// Check if text is in multilingual format: ===(code1)text1===(code2)text2===
	if strings.HasPrefix(text, "===") && strings.HasSuffix(text, "===") {
		// Remove leading and trailing ===
		content := strings.TrimPrefix(text, "===")
		content = strings.TrimSuffix(content, "===")
		
		// Split by === to get language sections
		sections := strings.Split(content, "===")
		
		languageTexts := make(map[string]string)
		for _, section := range sections {
			if section == "" {
				continue
			}
			
			// Parse each section: (code)text
			langPattern := `^\(([a-z]{2})\)(.*)$`
			re := regexp.MustCompile(langPattern)
			matches := re.FindStringSubmatch(section)
			if len(matches) == 3 {
				lang := matches[1]
				content := matches[2]
				languageTexts[lang] = content
			}
		}
		
		if languageCode == "" {
			// Return Ukrainian by default if available, otherwise first available
			if ukText, exists := languageTexts["uk"]; exists {
				return ukText, nil
			}
			for _, content := range languageTexts {
				return content, nil
			}
			return text, nil // Fallback to original text
		}
		
		if content, exists := languageTexts[languageCode]; exists {
			return content, nil
		}
		return "", fmt.Errorf("no text available for language: %s", languageCode)
	}

	// If we can't parse the format, treat as old format
	if languageCode == "" {
		return text, nil
	}
	return "", fmt.Errorf("no text available for language: %s", languageCode)
}

// processTextForLanguage processes repository text based on the requested language
func processTextForLanguage(text, languageCode string) (string, error) {
	return ParseMultilingualText(text, languageCode)
}

func GetRepository(w http.ResponseWriter, r *http.Request) {
	var reqBody getRepositoryRequestBody
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		server.RespondJSON(w, http.StatusBadRequest, "error", "Invalid request body", nil)
		return
	}

	// Validate text_language parameter if provided
	if reqBody.TextLanguage != "" {
		// Check if multiple language codes are provided (comma-separated)
		if strings.Contains(reqBody.TextLanguage, ",") {
			server.RespondJSON(w, http.StatusBadRequest, "error", "text_language parameter must contain only one language code", nil)
			return
		}
		
		// Validate the language code
		if err := server.ValidateLanguageCodes([]string{reqBody.TextLanguage}); err != nil {
			server.RespondJSON(w, http.StatusBadRequest, "error", fmt.Sprintf("Invalid language code: %v", err), nil)
			return
		}
	}

	paginationRequested := reqBody.Page > 0 || reqBody.PageSize > 0

	if paginationRequested || reqBody.Limit > 0 {
		if reqBody.Page < 1 {
			reqBody.Page = 1
		}
		if reqBody.PageSize < 1 {
			reqBody.PageSize = 10
		}
	}

	all, posted, unposted, err := countRepositories()
	if err != nil {
		log.Printf("Error counting repositories: %v", err)
		server.RespondJSON(w, http.StatusInternalServerError, "error", "Failed to count repositories", nil)
		return
	}

	sortBy := reqBody.SortBy
	if sortBy == "" {
		if reqBody.Posted != nil && *reqBody.Posted {
			sortBy = "date_posted"
		} else {
			sortBy = "date_added"
		}
	}

	var offset int
	limit := reqBody.Limit

	if paginationRequested || limit > 0 {
		offset = (reqBody.Page - 1) * reqBody.PageSize
		if limit == 0 {
			limit = reqBody.PageSize
		}
	}

	repositories, totalItems, err := database.GetRepository(limit, offset, reqBody.Posted, sortBy, reqBody.SortOrder)
	if err != nil {
		log.Printf("Error fetching repositories: %v", err)
		server.RespondJSON(w, http.StatusInternalServerError, "error", "Failed to fetch repositories", nil)
		return
	}

	items := make([]getRepositoryItem, len(repositories))
	for i, repo := range repositories {
		var processedText string
		if reqBody.TextLanguage == "" {
			processedText = repo.Text
		} else {
			var err error
			processedText, err = processTextForLanguage(repo.Text, reqBody.TextLanguage)
			if err != nil {
				log.Printf("Error processing text for repository %d: %v", repo.ID, err)
				server.RespondJSON(w, http.StatusBadRequest, "error", err.Error(), nil)
				return
			}
		}

		items[i] = getRepositoryItem{
			ID:         repo.ID,
			Posted:     repo.Posted == 1,
			URL:        repo.URL,
			Text:       processedText,
			DateAdded:  repo.DateAdded,
			DatePosted: repo.DatePosted,
		}
	}

	var totalPages int
	if paginationRequested || limit > 0 {
		totalPages = (totalItems + reqBody.PageSize - 1) / reqBody.PageSize
	} else {
		totalPages = 1
	}

	payload := &getRepositoryResponse{
		All:        all,
		Posted:     posted,
		Unposted:   unposted,
		Items:      items,
		Page:       reqBody.Page,
		PageSize:   reqBody.PageSize,
		TotalPages: totalPages,
		TotalItems: totalItems,
	}
	server.RespondJSON(w, http.StatusOK, "ok", "Repositories fetched successfully", payload)
}
