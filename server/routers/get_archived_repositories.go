package routers

import (
	"content-alchemist/database"
	"content-alchemist/server"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

const archiveDateOnlyLayout = "2006-01-02"

type getArchivedRepositoriesRequest struct {
	Limit     int    `json:"limit"`
	Page      int    `json:"page"`
	PageSize  int    `json:"page_size"`
	SortBy    string `json:"sort_by"`
	SortOrder string `json:"sort_order"`

	URL  string `json:"url"`
	Text string `json:"text"`

	DateAddedFrom    string `json:"date_added_from"`
	DateAddedTo      string `json:"date_added_to"`
	DatePostedFrom   string `json:"date_posted_from"`
	DatePostedTo     string `json:"date_posted_to"`
	DateArchivedFrom string `json:"date_archived_from"`
	DateArchivedTo   string `json:"date_archived_to"`

	TextLanguage string `json:"text_language"`
}

type getArchivedRepositoryItem struct {
	ID           int64      `json:"id"`
	OriginalID   *int64     `json:"original_id"`
	URL          string     `json:"url"`
	Text         string     `json:"text"`
	DateAdded    *time.Time `json:"date_added"`
	DatePosted   *time.Time `json:"date_posted"`
	DateArchived time.Time  `json:"date_archived"`
}

type getArchivedRepositoriesResponse struct {
	All        int                         `json:"all"`
	Items      []getArchivedRepositoryItem `json:"items"`
	Page       int                         `json:"page"`
	PageSize   int                         `json:"page_size"`
	TotalPages int                         `json:"total_pages"`
	TotalItems int                         `json:"total_items"`
}

// parseArchiveDate accepts either a full RFC3339 timestamp or a bare
// YYYY-MM-DD date. Upper bounds are exclusive, so a date-only upper bound is
// shifted to the start of the next day to keep the whole day inside the range.
func parseArchiveDate(value string, isUpperBound bool) (*time.Time, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return nil, nil
	}

	if parsed, err := time.Parse(time.RFC3339, trimmed); err == nil {
		return &parsed, nil
	}

	parsed, err := time.Parse(archiveDateOnlyLayout, trimmed)
	if err != nil {
		return nil, fmt.Errorf("must be an RFC3339 timestamp or a YYYY-MM-DD date")
	}
	if isUpperBound {
		parsed = parsed.AddDate(0, 0, 1)
	}
	return &parsed, nil
}

// GetArchivedRepositories returns a paginated, filtered view of the archive.
func GetArchivedRepositories(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		server.RespondJSON(w, http.StatusMethodNotAllowed, "error", "Only POST method is allowed", nil)
		return
	}

	var reqBody getArchivedRepositoriesRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		server.RespondJSON(w, http.StatusBadRequest, "error", "Invalid request body", nil)
		return
	}

	if reqBody.TextLanguage != "" {
		if strings.Contains(reqBody.TextLanguage, ",") {
			server.RespondJSON(w, http.StatusBadRequest, "error", "text_language parameter must contain only one language code", nil)
			return
		}
		if err := server.ValidateLanguageCodes([]string{reqBody.TextLanguage}); err != nil {
			server.RespondJSON(w, http.StatusBadRequest, "error", fmt.Sprintf("Invalid language code: %v", err), nil)
			return
		}
	}

	// A negative limit would reach the storage layer as "LIMIT -1", which SQLite
	// reads as "no limit" and would dump the whole archive while the response
	// still advertises the requested page_size.
	if reqBody.Limit < 0 || reqBody.Page < 0 || reqBody.PageSize < 0 {
		server.RespondJSON(w, http.StatusBadRequest, "error", "limit, page and page_size must not be negative", nil)
		return
	}

	if !database.IsValidArchivedSortBy(reqBody.SortBy) {
		server.RespondJSON(w, http.StatusBadRequest, "error", "sort_by must be one of: date_archived, date_posted, date_added, id", nil)
		return
	}
	if reqBody.SortOrder != "" && !strings.EqualFold(reqBody.SortOrder, "asc") && !strings.EqualFold(reqBody.SortOrder, "desc") {
		server.RespondJSON(w, http.StatusBadRequest, "error", "sort_order must be either asc or desc", nil)
		return
	}

	filter := database.ArchivedFilter{
		URL:       strings.TrimSpace(reqBody.URL),
		Text:      strings.TrimSpace(reqBody.Text),
		SortBy:    reqBody.SortBy,
		SortOrder: reqBody.SortOrder,
	}

	dateParams := []struct {
		name         string
		raw          string
		isUpperBound bool
		target       **time.Time
	}{
		{"date_added_from", reqBody.DateAddedFrom, false, &filter.DateAddedFrom},
		{"date_added_to", reqBody.DateAddedTo, true, &filter.DateAddedTo},
		{"date_posted_from", reqBody.DatePostedFrom, false, &filter.DatePostedFrom},
		{"date_posted_to", reqBody.DatePostedTo, true, &filter.DatePostedTo},
		{"date_archived_from", reqBody.DateArchivedFrom, false, &filter.DateArchivedFrom},
		{"date_archived_to", reqBody.DateArchivedTo, true, &filter.DateArchivedTo},
	}
	for _, param := range dateParams {
		parsed, err := parseArchiveDate(param.raw, param.isUpperBound)
		if err != nil {
			server.RespondJSON(w, http.StatusBadRequest, "error", fmt.Sprintf("Invalid %s: %v", param.name, err), nil)
			return
		}
		*param.target = parsed
	}

	paginationRequested := reqBody.Page > 0 || reqBody.PageSize > 0
	pageSize := reqBody.PageSize

	if paginationRequested || reqBody.Limit > 0 {
		if reqBody.Page < 1 {
			reqBody.Page = 1
		}
		if reqBody.PageSize < 1 {
			reqBody.PageSize = 10
		}

		filter.Limit = reqBody.Limit
		if filter.Limit == 0 {
			filter.Limit = reqBody.PageSize
		}
		// A page holds as many rows as the effective limit, so the offset and
		// total_pages must both be derived from it. Deriving the offset from
		// page_size while returning `limit` rows makes consecutive pages
		// overlap whenever the two differ.
		pageSize = filter.Limit
		filter.Offset = (reqBody.Page - 1) * filter.Limit
	}

	all, err := database.CountArchivedRepositories()
	if err != nil {
		log.Printf("Error counting archived repositories: %v", err)
		server.RespondJSON(w, http.StatusInternalServerError, "error", "Failed to count archived repositories", nil)
		return
	}

	repositories, totalItems, err := database.GetArchivedRepositories(filter)
	if err != nil {
		log.Printf("Error fetching archived repositories: %v", err)
		server.RespondJSON(w, http.StatusInternalServerError, "error", "Failed to fetch archived repositories", nil)
		return
	}

	items := make([]getArchivedRepositoryItem, len(repositories))
	for i, repo := range repositories {
		processedText := repo.Text
		if reqBody.TextLanguage != "" {
			processedText, err = processTextForLanguage(repo.Text, reqBody.TextLanguage)
			if err != nil {
				log.Printf("Error processing text for archived repository %d: %v", repo.ID, err)
				server.RespondJSON(w, http.StatusBadRequest, "error", err.Error(), nil)
				return
			}
		}

		items[i] = getArchivedRepositoryItem{
			ID:           repo.ID,
			OriginalID:   repo.OriginalID,
			URL:          repo.URL,
			Text:         processedText,
			DateAdded:    repo.DateAdded,
			DatePosted:   repo.DatePosted,
			DateArchived: repo.DateArchived,
		}
	}

	totalPages := 1
	if pageSize > 0 {
		totalPages = (totalItems + pageSize - 1) / pageSize
	}

	payload := &getArchivedRepositoriesResponse{
		All:        all,
		Items:      items,
		Page:       reqBody.Page,
		PageSize:   pageSize,
		TotalPages: totalPages,
		TotalItems: totalItems,
	}
	server.RespondJSON(w, http.StatusOK, "ok", "Archived repositories fetched successfully", payload)
}
