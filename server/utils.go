package server

import (
	"encoding/json"
	"net/http"
)

func RespondJSON(w http.ResponseWriter, status int, statusMsg, message string, data any) {
	response := map[string]any{}

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
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
}
