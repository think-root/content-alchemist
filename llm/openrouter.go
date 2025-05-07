package llm

import (
	"bytes"
	"content-alchemist/config"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func OpenRouter(requestBody map[string]any) (string, error) {
	apiKey := config.OPENROUTER_TOKEN
	if apiKey == "" {
		return "", fmt.Errorf("OPENROUTER_TOKEN is not set in configuration")
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequest("POST", "https://openrouter.ai/api/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("HTTP_REFERER", "https://github.com/think-root/content-alchemist")
	req.Header.Set("X-Title", "content-alchemist")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("OpenRouter API returned non-200 status code: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	var responseData map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&responseData); err != nil {
		return "", fmt.Errorf("failed to decode response body: %w", err)
	}

	choices, ok := responseData["choices"].([]any)
	if !ok || len(choices) == 0 {
		return "", fmt.Errorf("no choices in response or invalid format")
	}

	firstChoice, ok := choices[0].(map[string]any)
	if !ok {
		return "", fmt.Errorf("invalid choice format in response")
	}

	message, ok := firstChoice["message"].(map[string]any)
	if !ok {
		return "", fmt.Errorf("message not found in response choice")
	}

	content, ok := message["content"].(string)
	if !ok {
		return "", fmt.Errorf("content not found in message or not a string")
	}

	return content, nil
}
