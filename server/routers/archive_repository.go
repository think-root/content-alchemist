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
)

const maxArchiveBatchSize = 100

type archiveRepositoryRequest struct {
	IDs  []int64  `json:"ids"`
	URLs []string `json:"urls"`
}

type archivedRepositoryItem struct {
	ArchiveID    int64      `json:"archive_id"`
	ID           *int64     `json:"id"`
	URL          string     `json:"url"`
	DateAdded    *time.Time `json:"date_added"`
	DatePosted   *time.Time `json:"date_posted"`
	DateArchived *time.Time `json:"date_archived"`
}

// toArchivedItem maps a stored archive row to its API representation. A zero
// DateArchived (dry-run preview, nothing was written) is reported as null.
func toArchivedItem(repo database.ArchivedRepository) archivedRepositoryItem {
	item := archivedRepositoryItem{
		ArchiveID:  repo.ID,
		ID:         repo.OriginalID,
		URL:        repo.URL,
		DateAdded:  repo.DateAdded,
		DatePosted: repo.DatePosted,
	}
	if !repo.DateArchived.IsZero() {
		dateArchived := repo.DateArchived
		item.DateArchived = &dateArchived
	}
	return item
}

type archiveRepositoryResponse struct {
	Archived []archivedRepositoryItem  `json:"archived"`
	Failed   []database.ArchiveFailure `json:"failed"`
}

// ArchiveRepository moves one or more published repositories into the archive.
// Repositories are removed from the main table, which frees their URL so they
// can be collected and published again later.
func ArchiveRepository(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		server.RespondJSON(w, http.StatusMethodNotAllowed, "error", "Only POST method is allowed", nil)
		return
	}

	var reqBody archiveRepositoryRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		server.RespondJSON(w, http.StatusBadRequest, "error", "Invalid request body", nil)
		return
	}

	if len(reqBody.IDs) == 0 && len(reqBody.URLs) == 0 {
		server.RespondJSON(w, http.StatusBadRequest, "error", "Either ids or urls must be provided", nil)
		return
	}
	if len(reqBody.IDs) > 0 && len(reqBody.URLs) > 0 {
		server.RespondJSON(w, http.StatusBadRequest, "error", "Provide either ids or urls, not both", nil)
		return
	}

	var identifiers []string
	var isID bool

	if len(reqBody.IDs) > 0 {
		if len(reqBody.IDs) > maxArchiveBatchSize {
			server.RespondJSON(w, http.StatusBadRequest, "error", fmt.Sprintf("A maximum of %d ids can be archived per request", maxArchiveBatchSize), nil)
			return
		}
		isID = true
		identifiers = make([]string, 0, len(reqBody.IDs))
		for _, id := range reqBody.IDs {
			if id <= 0 {
				server.RespondJSON(w, http.StatusBadRequest, "error", "Every id must be a positive integer", nil)
				return
			}
			identifiers = append(identifiers, strconv.FormatInt(id, 10))
		}
	} else {
		if len(reqBody.URLs) > maxArchiveBatchSize {
			server.RespondJSON(w, http.StatusBadRequest, "error", fmt.Sprintf("A maximum of %d urls can be archived per request", maxArchiveBatchSize), nil)
			return
		}
		identifiers = make([]string, 0, len(reqBody.URLs))
		for _, url := range reqBody.URLs {
			trimmed := strings.TrimSpace(url)
			if trimmed == "" {
				server.RespondJSON(w, http.StatusBadRequest, "error", "URL cannot be empty", nil)
				return
			}
			identifiers = append(identifiers, trimmed)
		}
	}

	archived, failed, err := database.ArchiveRepositories(identifiers, isID)
	if err != nil {
		log.Printf("Error archiving repositories: %v", err)
		server.RespondJSON(w, http.StatusInternalServerError, "error", "Failed to archive repositories", nil)
		return
	}

	items := make([]archivedRepositoryItem, len(archived))
	for i, repo := range archived {
		items[i] = toArchivedItem(repo)
	}

	payload := &archiveRepositoryResponse{
		Archived: items,
		Failed:   failed,
	}
	message := fmt.Sprintf("Archived %d of %d repositories", len(items), len(identifiers))
	server.RespondJSON(w, http.StatusOK, "ok", message, payload)
}
