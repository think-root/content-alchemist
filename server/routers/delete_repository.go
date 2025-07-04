package routers

import (
	"content-alchemist/database"
	"content-alchemist/server"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type deleteRepositoryRequest struct {
	ID  *int64  `json:"id"`
	URL *string `json:"url"`
}

func DeleteRepository(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		server.RespondJSON(w, http.StatusMethodNotAllowed, "error", "Only DELETE method is allowed", nil)
		return
	}

	var reqBody deleteRepositoryRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		server.RespondJSON(w, http.StatusBadRequest, "error", "Invalid request body", nil)
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

	// Delete repository
	err := database.DeleteRepositoryByIDOrURL(identifier, isID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			server.RespondJSON(w, http.StatusNotFound, "error", err.Error(), nil)
			return
		}
		log.Printf("Error deleting repository: %v", err)
		server.RespondJSON(w, http.StatusInternalServerError, "error", "Failed to delete repository", nil)
		return
	}

	server.RespondJSON(w, http.StatusOK, "ok", "Repository deleted successfully", nil)
}