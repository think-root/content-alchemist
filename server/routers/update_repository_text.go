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
	ID           *int64  `json:"id"`
	URL          *string `json:"url"`
	Text         string  `json:"text"`
	TextLanguage *string `json:"text_language,omitempty"`
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

	trimmedText := strings.TrimSpace(reqBody.Text)
	if trimmedText == "" {
		server.RespondJSON(w, http.StatusBadRequest, "error", "Text field is required and cannot be empty", nil)
		return
	}

	if reqBody.ID == nil && reqBody.URL == nil {
		server.RespondJSON(w, http.StatusBadRequest, "error", "Either id or url must be provided", nil)
		return
	}
	
	if reqBody.ID != nil && reqBody.URL != nil {
		server.RespondJSON(w, http.StatusBadRequest, "error", "Provide either id or url, not both", nil)
		return
	}

	if utf8.RuneCountInString(trimmedText) > 1000 {
		server.RespondJSON(w, http.StatusBadRequest, "error", "Text must not exceed 1000 characters", nil)
		return
	}

	if !utf8.ValidString(trimmedText) {
		server.RespondJSON(w, http.StatusBadRequest, "error", "Text must be valid UTF-8", nil)
		return
	}

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

	var finalText string
	var updatedLanguage string

	if reqBody.TextLanguage == nil {
		finalText = server.CleanMultilingualText(trimmedText)
	} else {
		targetLang := strings.ToLower(strings.TrimSpace(*reqBody.TextLanguage))

		if err := server.ValidateLanguageCodes([]string{targetLang}); err != nil {
			server.RespondJSON(w, http.StatusBadRequest, "error", fmt.Sprintf("Invalid language code '%s': %v", targetLang, err), nil)
			return
		}

		if server.IsMultilingualText(existingRepo.Text) {
			_, exists := server.ExtractLanguageFromText(existingRepo.Text, targetLang)
			if !exists {
				server.RespondJSON(w, http.StatusUnprocessableEntity, "error", fmt.Sprintf("language '%s' not found in existing content", targetLang), nil)
				return
			}
			langMap := server.ParseMultilingualText(existingRepo.Text)
			langMap[targetLang] = trimmedText
			finalText = server.BuildMultilingualText(langMap)
		} else {
			finalText = server.BuildMultilingualText(map[string]string{targetLang: trimmedText})
		}

		updatedLanguage = targetLang
	}

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

	availableLanguages := server.GetAvailableLanguages(updatedRepo.Text)

	var response updateRepositoryTextResponse
	if reqBody.TextLanguage == nil {
		response = updateRepositoryTextResponse{
			ID:                 updatedRepo.ID,
			URL:                updatedRepo.URL,
			Text:               updatedRepo.Text,
			AvailableLanguages: availableLanguages,
			UpdatedAt:          time.Now(),
		}
	} else {
		response = updateRepositoryTextResponse{
			ID:                 updatedRepo.ID,
			URL:                updatedRepo.URL,
			Text:               updatedRepo.Text,
			UpdatedLanguage:    updatedLanguage,
			AvailableLanguages: availableLanguages,
			UpdatedAt:          time.Now(),
		}
	}

	server.RespondJSON(w, http.StatusOK, "ok", "Repository text updated successfully", response)
}