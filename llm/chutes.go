package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"content-alchemist/config"
)

func Chutes(requestBody map[string]any) (string, error) {
	apiKey := config.CHUTES_API_TOKEN
	if apiKey == "" {
		return "", fmt.Errorf("CHUTES_API_TOKEN is not set in configuration")
	}

	log.Printf("Chutes: Starting request")

	if _, ok := requestBody["messages"]; !ok {
		return "", fmt.Errorf("messages field is required in request body")
	}

	// Ensure model is set, default to moonshotai/Kimi-K2-Instruct-0905 if not specified
	if _, ok := requestBody["model"]; !ok {
		requestBody["model"] = "moonshotai/Kimi-K2-Instruct-0905"
	}

	// Set default parameters for Chutes API
	if _, ok := requestBody["max_tokens"]; !ok {
		requestBody["max_tokens"] = 1024
	}
	if _, ok := requestBody["temperature"]; !ok {
		requestBody["temperature"] = 0.7
	}
	if _, ok := requestBody["stream"]; !ok {
		requestBody["stream"] = false
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequest("POST", "https://llm.chutes.ai/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	log.Printf("Chutes: Sending request to %s", req.URL)

	client := &http.Client{Timeout: time.Second * 60}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Chutes: HTTP request failed: %v", err)
		return "", fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer resp.Body.Close()

	log.Printf("Chutes: Response status: %d", resp.StatusCode)

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("Chutes API returned non-200 status code: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	var responseData struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&responseData); err != nil {
		return "", fmt.Errorf("failed to decode response body: %w", err)
	}

	if len(responseData.Choices) == 0 {
		return "", fmt.Errorf("no choices in response")
	}

	result := responseData.Choices[0].Message.Content
	log.Printf("Chutes: Success! Response length: %d characters", len(result))
	return result, nil
}
