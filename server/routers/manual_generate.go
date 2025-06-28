package routers

import (
	"content-alchemist/database"
	"content-alchemist/llm"
	"content-alchemist/parser"
	"encoding/json"
	"log"
	"net/http"
	"slices"
	"strings"
)

type manualGenerateRequest struct {
	URL          string         `json:"url"`
	LLMProvider  string         `json:"llm_provider,omitempty"`
	LLMConfig    map[string]any `json:"llm_config,omitempty"`
	UseDirectURL bool           `json:"use_direct_url,omitempty"`
}

type manualGenerateResponse struct {
	Status       string   `json:"status"`
	Added        []string `json:"added"`
	DontAdded    []string `json:"dont_added"`
	ErrorMessage string   `json:"error_message,omitempty"`
}

func ManualGenerate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	var reqBody manualGenerateRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if reqBody.URL == "" {
		http.Error(w, "URL is required", http.StatusBadRequest)
		return
	}

	urls := strings.Fields(reqBody.URL)
	response := manualGenerateResponse{
		Status:    "ok",
		Added:     []string{},
		DontAdded: []string{},
	}

	filteredURLs, err := parser.FilterExistingURLs(urls)
	if err != nil {
		log.Printf("Error filtering existing URLs: %v", err)
		http.Error(w, "Failed to filter existing repositories", http.StatusInternalServerError)
		return
	}

	for _, url := range urls {
		if !slices.Contains(filteredURLs, url) {
			log.Printf("Repository already exists in DB: %s", url)
			response.DontAdded = append(response.DontAdded, url)
		}
	}

	for _, url := range filteredURLs {
		log.Printf("Processing repository: %s", url)

		repoReadme, err := parser.GetRepoReadme(url)
		if err != nil {
			log.Printf("Error fetching repo readme for URL %s: %v", url, err)
			response.Status = "error"
			response.DontAdded = append(response.DontAdded, url)
			response.ErrorMessage = "Failed to fetch repository README"
			continue
		}

		var textToProcess string
		if reqBody.UseDirectURL {
			textToProcess = url
		} else {
			textToProcess = repoReadme
		}

		processedText, err := llm.ProcessWithProvider(textToProcess, reqBody.LLMProvider, reqBody.LLMConfig)
		if err != nil {
			log.Printf("Error processing text with LLM for URL %s: %v", url, err)
			response.Status = "error"
			response.DontAdded = append(response.DontAdded, url)
			response.ErrorMessage = "Failed to process text with language model"
			continue
		}

		if err := database.AddRepositoryToDB(url, processedText); err != nil {
			log.Printf("Error adding repository to database for URL %s: %v", url, err)
			response.Status = "error"
			response.DontAdded = append(response.DontAdded, url)
			response.ErrorMessage = "Failed to add repository to database"
			continue
		}

		response.Added = append(response.Added, url)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
