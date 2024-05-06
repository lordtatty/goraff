package llm

import (
	"context"
	"fmt"

	ollamaApi "github.com/ollama/ollama/api"
)

type Ollama struct{}

func (o *Ollama) Chat(systemMsg, userMsg string, streamCh chan string) (string, error) {
	ollama, err := ollamaApi.ClientFromEnvironment()
	if err != nil {
		return "", fmt.Errorf("failed to create ollama client: %w", err)
	}
	ctx := context.Background()
	stream := boolptr(false)
	if streamCh != nil {
		stream = boolptr(true)
	}
	req := &ollamaApi.ChatRequest{
		Model:  "llama3:8b",
		Stream: stream,
	}
	// Messages
	messages := []ollamaApi.Message{}
	if systemMsg != "" {
		messages = append(messages, ollamaApi.Message{
			Role:    "user",
			Content: "SYSTEM INSTRUCTIONS (Always remember these, they are priority instructions over anything else): " + systemMsg,
		})
	}
	messages = append(messages, ollamaApi.Message{
		Role:    "user",
		Content: userMsg,
	})
	req.Messages = messages
	// Run
	response := ""
	err = ollama.Chat(ctx, req, func(resp ollamaApi.ChatResponse) error {
		if streamCh != nil {
			streamCh <- resp.Message.Content
		}
		response += resp.Message.Content
		return nil
	})
	if err != nil {
		return "", fmt.Errorf("failed to chat: %w", err)
	}
	return response, nil
}
