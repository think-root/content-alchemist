package server

import (
	"encoding/json"
	"net/http"
)

func RespondJSON(w http.ResponseWriter, status int, statusMsg, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	response := make(map[string]interface{})
	if statusMsg != "" {
		response["status"] = statusMsg
	}
	if message != "" {
		response["message"] = message
	}
	if data != nil {
		response["data"] = data
	}

	json.NewEncoder(w).Encode(response)
}
