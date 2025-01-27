package routers

import (
	"chappie/database"
	"chappie/server"
	"encoding/json"
	"log"
	"net/http"
)

type updatePostedRequest struct {
	URL    string `json:"url"`
	Posted bool   `json:"posted"`
}

func UpdatePostedStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		server.RespondJSON(w, http.StatusMethodNotAllowed, "error", "Only POST method is allowed", nil)
		return
	}

	var reqBody updatePostedRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		server.RespondJSON(w, http.StatusBadRequest, "error", "Invalid request body", nil)
		return
	}

	if err := database.UpdatePostedStatusByURL(reqBody.URL, reqBody.Posted); err != nil {
		log.Printf("Error updating posted status: %v", err)
		server.RespondJSON(w, http.StatusInternalServerError, "error", "Failed to update posted status", nil)
		return
	}

	server.RespondJSON(w, http.StatusOK, "ok", "Posted status updated successfully", nil)
}
