package llm

import (
	"bytes"
	"chappie/config"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type MistralPostResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

type MistralRequestBody struct {
	AgentID  string    `json:"agent_id"`
	Messages []Message `json:"messages"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

var (
	mistralApiUrl = "https://api.mistral.ai/v1/agents/completions"
	token         = config.MISTRAL_TOKEN
	agent         = config.MISTRAL_AGENT
	httpClient    = &http.Client{Timeout: time.Second * 60}
)

func Mistral(text string) (string, error) {
	requestBody := MistralRequestBody{
		AgentID: agent,
		Messages: []Message{
			{
				Role:    "user",
				Content: text,
			},
		},
	}

	var data MistralPostResponse
	var err error

	for i := 0; i < 5; i++ {
		jsonData, err := json.Marshal(requestBody)
		if err != nil {
			log.Printf("Error marshaling request body: %v", err)
			time.Sleep(time.Duration(i) * time.Minute)
			continue
		}

		req, err := http.NewRequest("POST", mistralApiUrl, bytes.NewBuffer(jsonData))
		if err != nil {
			log.Printf("Error creating request: %v", err)
			time.Sleep(time.Duration(i) * time.Minute)
			continue
		}

		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")

		resp, err := httpClient.Do(req)
		if err != nil {
			log.Printf("Error sending request: %v", err)
			time.Sleep(time.Duration(i) * time.Minute)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Printf("Received non-OK status: %d", resp.StatusCode)
			time.Sleep(time.Duration(i) * time.Minute)
			continue
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Error reading response body: %v", err)
			time.Sleep(time.Duration(i) * time.Minute)
			continue
		}

		err = json.Unmarshal(body, &data)
		if err != nil {
			log.Printf("Error unmarshaling response body: %v", err)
			time.Sleep(time.Duration(i) * time.Minute)
			continue
		}

		if len(data.Choices) == 0 || len(data.Choices[0].Message.Content) == 0 {
			log.Println("data.Choices is empty")
			time.Sleep(time.Duration(i) * time.Minute)
			continue
		}
		return data.Choices[0].Message.Content, nil
	}
	return "", err
}
