package llm

import (
	"bytes"
	"content-alchemist/config"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

func OpenAI(requestBody map[string]any) (string, error) {
	apiKey := config.OPENAI_TOKEN
	if apiKey == "" {
		return "", fmt.Errorf("OPENAI_TOKEN is not set in configuration")
	}

	if _, ok := requestBody["messages"]; !ok {
		return "", fmt.Errorf("messages field is required in request body")
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: time.Second * 60}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("OpenAI API returned non-200 status code: %d, body: %s", resp.StatusCode, string(bodyBytes))
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

	return responseData.Choices[0].Message.Content, nil
}
