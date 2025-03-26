package routers

import (
	"content-alchemist/database"
	"content-alchemist/server"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type getRepositoryRequestBody struct {
	Limit     int    `json:"limit"`
	Posted    *bool  `json:"posted"`
	SortBy    string `json:"sort_by"`
	SortOrder string `json:"sort_order"`
	Page      int    `json:"page"`
	PageSize  int    `json:"page_size"`
}

type getRepositoryItem struct {
	ID         int        `json:"id"`
	Posted     bool       `json:"posted"`
	URL        string     `json:"url"`
	Text       string     `json:"text"`
	DateAdded  *time.Time `json:"date_added"`
	DatePosted *time.Time `json:"date_posted"`
}

type getRepositoryResponse struct {
	All        int                 `json:"all"`
	Posted     int                 `json:"posted"`
	Unposted   int                 `json:"unposted"`
	Items      []getRepositoryItem `json:"items"`
	Page       int                 `json:"page"`
	PageSize   int                 `json:"page_size"`
	TotalPages int                 `json:"total_pages"`
	TotalItems int                 `json:"total_items"`
}

var trueVal = true
var falseVal = false

func countRepositories() (all, posted, unposted int, err error) {
	if all, err = database.CountRepositories(nil); err != nil {
		return
	}
	if posted, err = database.CountRepositories(&trueVal); err != nil {
		return
	}
	unposted, err = database.CountRepositories(&falseVal)
	return
}

func GetRepository(w http.ResponseWriter, r *http.Request) {
	var reqBody getRepositoryRequestBody
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		server.RespondJSON(w, http.StatusBadRequest, "error", "Invalid request body", nil)
		return
	}

	paginationRequested := reqBody.Page > 0 || reqBody.PageSize > 0

	if paginationRequested || reqBody.Limit > 0 {
		if reqBody.Page < 1 {
			reqBody.Page = 1
		}
		if reqBody.PageSize < 1 {
			reqBody.PageSize = 10
		}
	}

	all, posted, unposted, err := countRepositories()
	if err != nil {
		log.Printf("Error counting repositories: %v", err)
		server.RespondJSON(w, http.StatusInternalServerError, "error", "Failed to count repositories", nil)
		return
	}

	sortBy := reqBody.SortBy
	if sortBy == "" {
		if reqBody.Posted != nil && *reqBody.Posted {
			sortBy = "date_posted"
		} else {
			sortBy = "date_added"
		}
	}

	var offset int
	limit := reqBody.Limit

	if paginationRequested || limit > 0 {
		offset = (reqBody.Page - 1) * reqBody.PageSize
		if limit == 0 {
			limit = reqBody.PageSize
		}
	}

	repositories, totalItems, err := database.GetRepository(limit, offset, reqBody.Posted, sortBy, reqBody.SortOrder)
	if err != nil {
		log.Printf("Error fetching repositories: %v", err)
		server.RespondJSON(w, http.StatusInternalServerError, "error", "Failed to fetch repositories", nil)
		return
	}

	items := make([]getRepositoryItem, len(repositories))
	for i, repo := range repositories {
		items[i] = getRepositoryItem{
			ID:         repo.ID,
			Posted:     repo.Posted,
			URL:        repo.URL,
			Text:       repo.Text,
			DateAdded:  repo.DateAdded,
			DatePosted: repo.DatePosted,
		}
	}

	var totalPages int
	if paginationRequested || limit > 0 {
		totalPages = (totalItems + reqBody.PageSize - 1) / reqBody.PageSize
	} else {
		totalPages = 1
	}

	payload := &getRepositoryResponse{
		All:        all,
		Posted:     posted,
		Unposted:   unposted,
		Items:      items,
		Page:       reqBody.Page,
		PageSize:   reqBody.PageSize,
		TotalPages: totalPages,
		TotalItems: totalItems,
	}
	server.RespondJSON(w, http.StatusOK, "ok", "Repositories fetched successfully", payload)
}
