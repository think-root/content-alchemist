package routers

import (
	"content-alchemist/database"
	"content-alchemist/server"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type archiveOldRepositoriesRequest struct {
	Days   int  `json:"days"`
	DryRun bool `json:"dry_run"`
}

type archiveOldRepositoriesResponse struct {
	ArchivedCount int                      `json:"archived_count"`
	DryRun        bool                     `json:"dry_run"`
	Archived      []archivedRepositoryItem `json:"archived"`
}

// ArchiveOldRepositories archives every published repository whose publication
// date is older than the requested number of days. The age is measured by
// date_posted, not by date_added.
func ArchiveOldRepositories(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		server.RespondJSON(w, http.StatusMethodNotAllowed, "error", "Only POST method is allowed", nil)
		return
	}

	var reqBody archiveOldRepositoriesRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		server.RespondJSON(w, http.StatusBadRequest, "error", "Invalid request body", nil)
		return
	}

	if reqBody.Days < 1 {
		server.RespondJSON(w, http.StatusBadRequest, "error", "days must be a positive integer", nil)
		return
	}

	archived, err := database.ArchiveRepositoriesOlderThan(reqBody.Days, reqBody.DryRun)
	if err != nil {
		log.Printf("Error archiving old repositories: %v", err)
		server.RespondJSON(w, http.StatusInternalServerError, "error", "Failed to archive repositories", nil)
		return
	}

	items := make([]archivedRepositoryItem, len(archived))
	for i, repo := range archived {
		items[i] = toArchivedItem(repo)
	}

	payload := &archiveOldRepositoriesResponse{
		ArchivedCount: len(items),
		DryRun:        reqBody.DryRun,
		Archived:      items,
	}

	message := fmt.Sprintf("Archived %d repositories published more than %d days ago", len(items), reqBody.Days)
	if reqBody.DryRun {
		message = fmt.Sprintf("%d repositories published more than %d days ago would be archived", len(items), reqBody.Days)
	}
	server.RespondJSON(w, http.StatusOK, "ok", message, payload)
}
