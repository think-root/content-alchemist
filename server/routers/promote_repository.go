package routers

import (
	"content-alchemist/database"
	"content-alchemist/server"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type promoteRepositoryRequest struct {
	ID  *int64  `json:"id"`
	URL *string `json:"url"`
}

type promoteRepositoryResponse struct {
	ID              int64      `json:"id"`
	Posted          bool       `json:"posted"`
	URL             string     `json:"url"`
	Text            string     `json:"text"`
	DateAdded       *time.Time `json:"date_added"`
	DatePosted      *time.Time `json:"date_posted"`
	PublishPriority *int64     `json:"publish_priority"`
}

func PromoteRepository(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		server.RespondJSON(w, http.StatusMethodNotAllowed, "error", "Only PATCH method is allowed", nil)
		return
	}

	var reqBody promoteRepositoryRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		server.RespondJSON(w, http.StatusBadRequest, "error", "Invalid request body", nil)
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

	repo, err := database.PromoteRepositoryToNextByIDOrURL(identifier, isID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			server.RespondJSON(w, http.StatusNotFound, "error", err.Error(), nil)
			return
		}
		if strings.Contains(err.Error(), "already posted") {
			server.RespondJSON(w, http.StatusConflict, "error", "Repository is already posted", nil)
			return
		}
		log.Printf("Error promoting repository: %v", err)
		server.RespondJSON(w, http.StatusInternalServerError, "error", "Failed to promote repository", nil)
		return
	}

	response := promoteRepositoryResponse{
		ID:              repo.ID,
		Posted:          repo.Posted == 1,
		URL:             repo.URL,
		Text:            repo.Text,
		DateAdded:       repo.DateAdded,
		DatePosted:      repo.DatePosted,
		PublishPriority: repo.PublishPriority,
	}
	server.RespondJSON(w, http.StatusOK, "ok", "Repository promoted to publish next", response)
}
