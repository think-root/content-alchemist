package llm

import (
	"fmt"
)

type LLMProvider interface {
	Process(text string, config map[string]any) (string, error)
}

func prepareMessagesConfig(text string, config map[string]any) map[string]any {
	if config == nil {
		config = make(map[string]any)
	}

	if _, ok := config["messages"]; !ok {
		config["messages"] = []map[string]any{}
	}

	messagesAny, ok := config["messages"]
	if !ok {
		messagesAny = []map[string]any{}
	}

	var messages []map[string]any
	switch v := messagesAny.(type) {
	case []map[string]any:
		messages = v
	case []any:
		messages = make([]map[string]any, 0, len(v))
		for _, item := range v {
			if msg, ok := item.(map[string]any); ok {
				messages = append(messages, msg)
			}
		}
	default:
		messages = []map[string]any{}
	}

	hasUserMessage := false
	for i, msg := range messages {
		if role, exists := msg["role"]; exists && role == "user" {
			hasUserMessage = true
			messages[i]["content"] = text
			break
		}
	}

	if !hasUserMessage {
		userMessage := map[string]any{
			"role":    "user",
			"content": text,
		}
		messages = append(messages, userMessage)
	}

	config["messages"] = messages
	return config
}

func ProcessWithProvider(text string, provider string, config map[string]any) (string, error) {
	switch provider {
	case "mistral_api", "":
		return MistralAPI(prepareMessagesConfig(text, config))
	case "openai":
		return OpenAI(prepareMessagesConfig(text, config))
	case "openrouter":
		return OpenRouter(prepareMessagesConfig(text, config))
	case "chutes":
		return Chutes(prepareMessagesConfig(text, config))
	default:
		return "", fmt.Errorf("unsupported LLM provider: %s", provider)
	}
}
