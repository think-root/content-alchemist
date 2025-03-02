package routers

import (
	"content-alchemist/database"
	"content-alchemist/server"
	"encoding/json"
	"log"
	"net/http"
)

type getRepositoryRequestBody struct {
	Limit  int  `json:"limit"`
	Posted bool `json:"posted"`
}

type getRepositoryItem struct {
	ID     int    `json:"id"`
	Posted bool   `json:"posted"`
	URL    string `json:"url"`
	Text   string `json:"text"`
}

type getRepositoryResponse struct {
	All      int                 `json:"all"`
	Posted   int                 `json:"posted"`
	Unposted int                 `json:"unposted"`
	Items    []getRepositoryItem `json:"items"`
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

	all, posted, unposted, err := countRepositories()
	if err != nil {
		log.Printf("Error counting repositories: %v", err)
		server.RespondJSON(w, http.StatusInternalServerError, "error", "Failed to count repositories", nil)
		return
	}

	repositories, err := database.GetRepository(reqBody.Limit, &reqBody.Posted)
	if err != nil {
		log.Printf("Error fetching repositories: %v", err)
		server.RespondJSON(w, http.StatusInternalServerError, "error", "Failed to fetch repositories", nil)
		return
	}

	items := make([]getRepositoryItem, len(repositories))
	for i, repo := range repositories {
		items[i] = getRepositoryItem{repo.ID, repo.Posted, repo.URL, repo.Text}
	}

	payload := &getRepositoryResponse{
		All:      all,
		Posted:   posted,
		Unposted: unposted,
		Items:    items,
	}
	server.RespondJSON(w, http.StatusOK, "ok", "Repositories fetched successfully", payload)
}
