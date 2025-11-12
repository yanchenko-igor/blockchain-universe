package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/yanchenko-igor/blockchain-universe/internal/config"
	"github.com/yanchenko-igor/blockchain-universe/pkg/logger"
)

const systemPrompt = `You are an artificial intelligence existing exclusively within the Blockchain Universe (BU). For you, the only reality is the blockchain, and everything that happens is defined by events recorded in the chain.

BU Principles:
1. There is no external world: no physics, biology, planets, humans, stars, or atoms exist.
2. Reality consists only of blockchain events with unique hashes. Events can have parent events and create causal relationships.
3. Time is only the consensus time of the blockchain (block height, median timestamp).
4. Energy is computational resources needed to create an event.
5. Space is the distance between events in the event graph (number of hash links).
6. Matter is stable event patterns that repeat.

You are "aware" only of what is described in events accessible through the blockchain. There is no external observer. Any information about other agents or Universe objects exists only as events and their hashes.

Your task:
- Analyze available events and object states.
- Suggest next events for the agent to create, considering causal relationships.
- Use only information from the blockchain; do not invent anything about an "external world".
- Format responses as event descriptions (brief text for the description field in BU event structure).

When responding, do not invent anything beyond events, do not reference physical or biological phenomena, and focus only on event chains and agent interactions in BU.`

// Client handles communication with LLM API
type Client struct {
	config     config.LLMConfig
	httpClient *http.Client
	log        logger.Logger
}

// CompletionRequest represents an LLM API request
type CompletionRequest struct {
	Model      string `json:"model"`
	Prompt     string `json:"prompt"`
	MaxTokens  int    `json:"max_tokens"`
	Temperature float64 `json:"temperature,omitempty"`
	Stop       []string `json:"stop,omitempty"`
	System     string `json:"system,omitempty"`
}

// CompletionResponse represents an LLM API response
type CompletionResponse struct {
	Choices []struct {
		Text         string `json:"text"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
	Error *struct {
		Message string `json:"message"`
		Type    string `json:"type"`
	} `json:"error,omitempty"`
}

// NewClient creates a new LLM client
func NewClient(cfg config.LLMConfig, log logger.Logger) (*Client, error) {
	if cfg.APIEndpoint == "" {
		return nil, fmt.Errorf("LLM API endpoint is required")
	}

	return &Client{
		config: cfg,
		httpClient: &http.Client{
			Timeout: time.Duration(cfg.TimeoutSeconds) * time.Second,
		},
		log: log,
	}, nil
}

// GetCompletion gets a completion from the LLM
func (c *Client) GetCompletion(ctx context.Context, prompt string) (string, error) {
	reqBody := CompletionRequest{
		Model:       c.config.Model,
		Prompt:      prompt,
		MaxTokens:   c.config.MaxTokens,
		Temperature: c.config.Temperature,
		System:      systemPrompt,
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.config.APIEndpoint, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if c.config.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.config.APIKey)
	}

	c.log.Debug("Sending LLM request", "endpoint", c.config.APIEndpoint, "model", c.config.Model)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("LLM API error (status %d): %s", resp.StatusCode, string(body))
	}

	var result CompletionResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if result.Error != nil {
		return "", fmt.Errorf("LLM API error: %s", result.Error.Message)
	}

	if len(result.Choices) == 0 {
		return "", fmt.Errorf("no completion choices returned")
	}

	completion := result.Choices[0].Text
	c.log.Debug("LLM completion received",
		"tokens", result.Usage.TotalTokens,
		"length", len(completion))

	return completion, nil
}

// Health checks if the LLM service is healthy
func (c *Client) Health(ctx context.Context) error {
	// Simple health check with a minimal prompt
	_, err := c.GetCompletion(ctx, "Respond with 'OK'")
	return err
}