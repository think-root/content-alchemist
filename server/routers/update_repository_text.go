package routers

import (
	"content-alchemist/database"
	"content-alchemist/server"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

type updateRepositoryTextRequest struct {
	ID   *int64  `json:"id"`
	URL  *string `json:"url"`
	Text string  `json:"text"`
}

type updateRepositoryTextResponse struct {
	ID                 int64     `json:"id"`
	URL                string    `json:"url"`
	Text               string    `json:"text"`
	UpdatedLanguage    string    `json:"updated_language,omitempty"`
	AvailableLanguages []string  `json:"available_languages"`
	UpdatedAt          time.Time `json:"updated_at"`
}

func UpdateRepositoryText(w http.ResponseWriter, r *http.Request) {
	var reqBody updateRepositoryTextRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		server.RespondJSON(w, http.StatusBadRequest, "error", "Invalid request body", nil)
		return
	}

	// Extract target language from query parameter (default to "uk")
	targetLang := r.URL.Query().Get("lang")
	if targetLang == "" {
		targetLang = "uk"
	}

	// Validate target language
	if err := server.ValidateLanguageCodes([]string{targetLang}); err != nil {
		server.RespondJSON(w, http.StatusBadRequest, "error", fmt.Sprintf("Invalid language code '%s': %v", targetLang, err), nil)
		return
	}

	// Validate text field first (most important validation)
	trimmedText := strings.TrimSpace(reqBody.Text)
	if trimmedText == "" {
		server.RespondJSON(w, http.StatusBadRequest, "error", "Text field is required and cannot be empty", nil)
		return
	}

	// Validate that exactly one identifier is provided
	if reqBody.ID == nil && reqBody.URL == nil {
		server.RespondJSON(w, http.StatusBadRequest, "error", "Either id or url must be provided", nil)
		return
	}
	
	if reqBody.ID != nil && reqBody.URL != nil {
		server.RespondJSON(w, http.StatusBadRequest, "error", "Provide either id or url, not both", nil)
		return
	}

	// Check text length (1000 characters max)
	if utf8.RuneCountInString(trimmedText) > 1000 {
		server.RespondJSON(w, http.StatusBadRequest, "error", "Text must not exceed 1000 characters", nil)
		return
	}

	// Validate UTF-8 encoding
	if !utf8.ValidString(trimmedText) {
		server.RespondJSON(w, http.StatusBadRequest, "error", "Text must be valid UTF-8", nil)
		return
	}

	// Determine identifier and type
	var identifier string
	var isID bool

	if reqBody.ID != nil {
		if *reqBody.ID <= 0 {
			server.RespondJSON(w, http.StatusBadRequest, "error", "ID must be a positive integer", nil)
			return
		}
		identifier = strconv.FormatInt(*reqBody.ID, 10)
		isID = true
	} else {
		if strings.TrimSpace(*reqBody.URL) == "" {
			server.RespondJSON(w, http.StatusBadRequest, "error", "URL cannot be empty", nil)
			return
		}
		identifier = strings.TrimSpace(*reqBody.URL)
		isID = false
	}

	// Get existing repository to check current text format
	existingRepo, err := database.GetRepositoryByIDOrURL(identifier, isID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			server.RespondJSON(w, http.StatusNotFound, "error", err.Error(), nil)
			return
		}
		log.Printf("Error fetching existing repository: %v", err)
		server.RespondJSON(w, http.StatusInternalServerError, "error", "Failed to fetch existing repository", nil)
		return
	}

	// Apply intelligent multilingual update logic
	finalText, err := server.UpdateLanguageInText(existingRepo.Text, trimmedText, targetLang)
	if err != nil {
		log.Printf("Error updating language in text: %v", err)
		server.RespondJSON(w, http.StatusBadRequest, "error", fmt.Sprintf("Failed to update text: %v", err), nil)
		return
	}

	// Update repository text with the processed final text
	updatedRepo, err := database.UpdateRepositoryTextByIDOrURL(identifier, finalText, isID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			server.RespondJSON(w, http.StatusNotFound, "error", err.Error(), nil)
			return
		}
		log.Printf("Error updating repository text: %v", err)
		server.RespondJSON(w, http.StatusInternalServerError, "error", "Failed to update repository text", nil)
		return
	}

	// Get available languages for response
	availableLanguages := server.GetAvailableLanguages(updatedRepo.Text)

	// Prepare response
	response := updateRepositoryTextResponse{
		ID:                 updatedRepo.ID,
		URL:                updatedRepo.URL,
		Text:               updatedRepo.Text,
		UpdatedLanguage:    targetLang,
		AvailableLanguages: availableLanguages,
		UpdatedAt:          time.Now(),
	}

	server.RespondJSON(w, http.StatusOK, "ok", "Repository text updated successfully", response)
}