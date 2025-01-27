package server

import (
	"encoding/json"
	"net/http"
)

func RespondJSON(w http.ResponseWriter, status int, statusMsg, message string, data interface{}) {
	response := map[string]interface{}{}

	if statusMsg != "" {
		response["status"] = statusMsg
	}
	if message != "" {
		response["message"] = message
	}
	if data != nil {
		response["data"] = data
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
}