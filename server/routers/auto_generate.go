package routers

import (
	"content-alchemist/database"
	"content-alchemist/llm"
	"content-alchemist/parser"
	"content-alchemist/server"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type autoGenerateRequest struct {
	MaxRepos           int            `json:"max_repos"`
	Since              string         `json:"since"`
	SpokenLanguageCode string         `json:"spoken_language_code"`
	Resource           string         `json:"resource"`
	Period             string         `json:"period"`
	Language           string         `json:"language"`
	LLMProvider        string         `json:"llm_provider,omitempty"`
	LLMConfig          map[string]any `json:"llm_config,omitempty"`
	UseDirectURL       bool           `json:"use_direct_url,omitempty"`
	LLMOutputLanguage  string         `json:"llm_output_language,omitempty"`
}

type autoGenerateResponse struct {
	Status       string   `json:"status"`
	Added        []string `json:"added"`
	DontAdded    []string `json:"dont_added"`
	ErrorMessage string   `json:"error_message,omitempty"`
}

func AutoGenerate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	var reqBody autoGenerateRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if reqBody.MaxRepos <= 0 {
		http.Error(w, "Fields max_repos are required", http.StatusBadRequest)
		return
	}

	var repos []parser.Repository
	var err error

	if reqBody.Resource == "" {
		reqBody.Resource = "github"
	}

	// Validate Resource
	validResources := map[string]bool{
		"github":     true,
		"ossinsight": true,
	}
	if !validResources[reqBody.Resource] {
		http.Error(w, fmt.Sprintf("Invalid resource: %s. Allowed values: github, ossinsight", reqBody.Resource), http.StatusBadRequest)
		return
	}

	switch reqBody.Resource {
	case "ossinsight":
		// Validate Period
		validPeriods := map[string]bool{
			"past_24_hours": true,
			"past_week":     true,
			"past_month":    true,
			"past_3_months": true,
		}
		if reqBody.Period != "" && !validPeriods[reqBody.Period] {
			http.Error(w, fmt.Sprintf("Invalid period: %s. Allowed values: past_24_hours, past_week, past_month, past_3_months", reqBody.Period), http.StatusBadRequest)
			return
		}

		repos, err = parser.GetTrendingReposFromOssInsight(reqBody.MaxRepos, reqBody.Period, reqBody.Language)
	case "github":
		// Validate Since
		validSince := map[string]bool{
			"daily":   true,
			"weekly":  true,
			"monthly": true,
		}
		if reqBody.Since != "" && !validSince[reqBody.Since] {
			http.Error(w, fmt.Sprintf("Invalid since: %s. Allowed values: daily, weekly, monthly", reqBody.Since), http.StatusBadRequest)
			return
		}

		repos, err = parser.GetTrendingRepos(reqBody.MaxRepos, reqBody.Since, reqBody.SpokenLanguageCode)
	}

	if err != nil {
		log.Printf("Error fetching trending repositories (source: %s): %v", reqBody.Resource, err)
		http.Error(w, fmt.Sprintf("Failed to fetch trending repositories from %s", reqBody.Resource), http.StatusInternalServerError)
		return
	}

	response := autoGenerateResponse{
		Status:    "ok",
		Added:     []string{},
		DontAdded: []string{},
	}

	// Parse and validate language codes
	languageCodes := server.ParseLanguageCodes(reqBody.LLMOutputLanguage)
	if err := server.ValidateLanguageCodes(languageCodes); err != nil {
		log.Printf("Invalid language codes: %v", err)
		http.Error(w, fmt.Sprintf("Invalid language codes: %v", err), http.StatusBadRequest)
		return
	}

	for _, repo := range repos {
		log.Printf("Processing repository: %s", repo.URL)

		repoReadme, err := parser.GetRepoReadme(repo.URL)
		if err != nil {
			log.Printf("Error fetching repo readme for URL %s: %v", repo.URL, err)
			response.Status = "error"
			response.DontAdded = append(response.DontAdded, repo.URL)
			response.ErrorMessage = "Failed to fetch repository README"
			continue
		}

		var textToProcess string
		if reqBody.UseDirectURL {
			textToProcess = repo.URL
		} else {
			textToProcess = repoReadme
		}

		// Prepare LLM config with multilingual instructions
		llmConfig := reqBody.LLMConfig
		if llmConfig == nil {
			llmConfig = make(map[string]any)
		}

		// Add multilingual prompt instructions to the system message or create one
		multilingualPrompt := server.BuildMultilingualPrompt(languageCodes)

		// Handle messages in config
		if messages, exists := llmConfig["messages"]; exists {
			// Try different possible types for messages
			switch msgSlice := messages.(type) {
			case []map[string]any:
				// Look for existing system message
				systemMessageFound := false
				for i, msg := range msgSlice {
					if role, exists := msg["role"]; exists && role == "system" {
						if content, exists := msg["content"]; exists {
							msgSlice[i]["content"] = fmt.Sprintf("%s\n\n%s", content, multilingualPrompt)
						} else {
							msgSlice[i]["content"] = multilingualPrompt
						}
						systemMessageFound = true
						break
					}
				}

				// If no system message found, add one at the beginning
				if !systemMessageFound {
					systemMsg := map[string]any{
						"role":    "system",
						"content": multilingualPrompt,
					}
					msgSlice = append([]map[string]any{systemMsg}, msgSlice...)
					llmConfig["messages"] = msgSlice
				}
			case []any:
				// Handle []any type (common in JSON unmarshaling)
				var convertedMessages []map[string]any
				for _, msg := range msgSlice {
					if msgMap, ok := msg.(map[string]any); ok {
						convertedMessages = append(convertedMessages, msgMap)
					}
				}

				// Look for existing system message
				systemMessageFound := false
				for i, msg := range convertedMessages {
					if role, exists := msg["role"]; exists && role == "system" {
						if content, exists := msg["content"]; exists {
							convertedMessages[i]["content"] = fmt.Sprintf("%s\n\n%s", content, multilingualPrompt)
						} else {
							convertedMessages[i]["content"] = multilingualPrompt
						}
						systemMessageFound = true
						break
					}
				}

				// If no system message found, add one at the beginning
				if !systemMessageFound {
					systemMsg := map[string]any{
						"role":    "system",
						"content": multilingualPrompt,
					}
					convertedMessages = append([]map[string]any{systemMsg}, convertedMessages...)
				}
				llmConfig["messages"] = convertedMessages
			default:
				log.Printf("Warning: messages field has unexpected type: %T", messages)
			}
		} else {
			// No messages exist, create system message
			llmConfig["messages"] = []map[string]any{
				{
					"role":    "system",
					"content": multilingualPrompt,
				},
			}
		}

		processedText, err := llm.ProcessWithProvider(textToProcess, reqBody.LLMProvider, llmConfig)
		if err != nil {
			log.Printf("Error processing text with LLM for URL %s: %v", repo.URL, err)
			response.Status = "error"
			response.DontAdded = append(response.DontAdded, repo.URL)
			response.ErrorMessage = "Failed to process text with language model"
			continue
		}

		cleanedText := server.CleanMultilingualText(processedText)
		if err := database.AddRepositoryToDB(repo.URL, cleanedText); err != nil {
			log.Printf("Error adding repository to database for URL %s: %v", repo.URL, err)
			response.Status = "error"
			response.DontAdded = append(response.DontAdded, repo.URL)
			response.ErrorMessage = "Failed to add repository to database"
			continue
		}

		response.Added = append(response.Added, repo.URL)
	}

	w.Header().Set("Content-Type", "application/json")
	if response.Status == "error" {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
