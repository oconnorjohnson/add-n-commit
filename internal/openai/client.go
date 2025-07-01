package openai

import (
	"context"
	"fmt"

	"github.com/sashabaranov/go-openai"
)

// Client wraps the OpenAI client
type Client struct {
	client      *openai.Client
	model       string
	temperature float32
}

// NewClient creates a new OpenAI client
func NewClient(apiKey, model string, temperature float32) *Client {
	return &Client{
		client:      openai.NewClient(apiKey),
		model:       model,
		temperature: temperature,
	}
}

// GenerateCommitMessage generates a commit message based on the diff
func (c *Client) GenerateCommitMessage(systemPrompt, diff string) (string, error) {
	if diff == "" {
		return "", fmt.Errorf("empty diff provided")
	}

	resp, err := c.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: c.model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: systemPrompt,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: diff,
				},
			},
			Temperature: c.temperature,
		},
	)

	if err != nil {
		return "", fmt.Errorf("failed to generate commit message: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no response from OpenAI")
	}

	return resp.Choices[0].Message.Content, nil
}

// GenerateCommitMessageWithContext generates a commit message with additional context
func (c *Client) GenerateCommitMessageWithContext(systemPrompt, diff, context string) (string, error) {
	if diff == "" {
		return "", fmt.Errorf("empty diff provided")
	}

	userMessage := diff
	if context != "" {
		userMessage = fmt.Sprintf("Context: %s\n\nDiff:\n%s", context, diff)
	}

	resp, err := c.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: c.model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: systemPrompt,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: userMessage,
				},
			},
			Temperature: c.temperature,
		},
	)

	if err != nil {
		return "", fmt.Errorf("failed to generate commit message: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no response from OpenAI")
	}

	return resp.Choices[0].Message.Content, nil
} 