package routers

import (
	"chappie_server/database"
	"chappie_server/llm"
	"chappie_server/parser"
	"encoding/json"
	"log"
	"net/http"
)

type autoGenerateRequest struct {
	MaxRepos           int    `json:"max_repos"`
	Since              string `json:"since"`
	SpokenLanguageCode string `json:"spoken_language_code"`
}

type autoGenerateResponse struct {
	Status    string   `json:"status"`
	Added     []string `json:"added"`
	DontAdded []string `json:"dont_added"`
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

	if reqBody.MaxRepos <= 0 || reqBody.Since == "" || reqBody.SpokenLanguageCode == "" {
		http.Error(w, "All fields (max_repos, since, spoken_language_code) are required", http.StatusBadRequest)
		return
	}

	repos, err := parser.GetTrendingRepos(reqBody.MaxRepos, reqBody.Since, reqBody.SpokenLanguageCode)
	if err != nil {
		log.Printf("Error fetching trending repositories: %v", err)
		http.Error(w, "Failed to fetch trending repositories", http.StatusInternalServerError)
		return
	}

	response := autoGenerateResponse{
		Status:    "ok",
		Added:     []string{},
		DontAdded: []string{},
	}

	for _, repo := range repos {
		log.Printf("Processing repository: %s", repo.URL)

		exists, err := database.SearchPostInDB(repo.URL)
		if err != nil {
			log.Printf("Error searching repository in DB for URL %s: %v", repo.URL, err)
			response.Status = "error"
			response.DontAdded = append(response.DontAdded, repo.URL)
			continue
		}

		if exists {
			log.Printf("Repository already exists in DB: %s", repo.URL)
			response.DontAdded = append(response.DontAdded, repo.URL)
			continue
		}

		repoReadme, err := parser.GetRepoReadme(repo.URL)
		if err != nil {
			log.Printf("Error fetching repo readme for URL %s: %v", repo.URL, err)
			response.Status = "error"
			response.DontAdded = append(response.DontAdded, repo.URL)
			continue
		}

		processedText, err := llm.Mistral(repoReadme)
		if err != nil {
			log.Printf("Error processing text with LLM for URL %s: %v", repo.URL, err)
			response.Status = "error"
			response.DontAdded = append(response.DontAdded, repo.URL)
			continue
		}

		if err := database.AddRepositoryToDB(repo.URL, processedText); err != nil {
			log.Printf("Error adding repository to database for URL %s: %v", repo.URL, err)
			response.Status = "error"
			response.DontAdded = append(response.DontAdded, repo.URL)
			continue
		}

		response.Added = append(response.Added, repo.URL)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
