package backend

import (
	"context"
	"encoding/json"
	"fmt"
	"maps"
	"strings"

	"github.com/dll-as/gitc/internal/ai/transport"
	"github.com/dll-as/gitc/internal/config"
)

// OpenAICompatibleProvider implements the Provider interface for OpenAI-compatible APIs
type OpenAICompatibleProvider struct {
	name   string
	config *config.Config
	client *transport.Client
}

const systemPrompt = "You are an AI assistant that generates concise and meaningful Git commit messages."

// NewOpenAICompatibleProvider creates a new OpenAI-compatible provider
func NewOpenAICompatibleProvider(cfg *config.Config) (Provider, error) {
	client, err := transport.New(transport.Config{
		Timeout:    cfg.Backend.Timeout,
		MaxRetries: cfg.Backend.MaxRetries,
		Proxy:      cfg.Backend.Proxy,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP client: %w", err)
	}

	return &OpenAICompatibleProvider{
		name:   "openai-compatible",
		config: cfg,
		client: client,
	}, nil
}

// Name returns the provider name
func (p *OpenAICompatibleProvider) Name() string {
	return p.name
}

// Generate generates a commit message
func (p *OpenAICompatibleProvider) Generate(
	ctx context.Context,
	req *GenerateRequest,
) (*GenerateResponse, error) {

	payload := map[string]any{
		"model": p.config.Backend.Model,
		"messages": []Message{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: req.Prompt},
		},
		"max_tokens":  req.MaxTokens,
		"temperature": req.Temperature,
		"top_p":       req.TopP,
		"stream":      false,
	}

	maps.Copy(payload, p.config.Backend.ExtraBody)

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	headers := map[string]string{
		"Content-Type": "application/json",
	}

	if p.config.Backend.APIKey != "" {
		headers["Authorization"] = "Bearer " + p.config.Backend.APIKey
	}

	for k, v := range p.config.Backend.ExtraHeaders {
		headers[k] = v
	}

	respBody, err := p.client.DoWithContext(ctx, &transport.Request{
		Method:  "POST",
		URL:     strings.TrimRight(p.config.Backend.BaseURL, "/") + "/chat/completions",
		Headers: headers,
		Body:    body,
	})
	if err != nil {
		return nil, fmt.Errorf("API request failed: %w", err)
	}

	var resp OpenAIResponse
	if err = json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	if resp.Error.Message != "" {
		return nil, fmt.Errorf("API error: %s", resp.Error.Message)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("provider returned no choices")
	}

	msg := strings.TrimSpace(resp.Choices[0].Message.Content)
	if msg == "" {
		return nil, fmt.Errorf("provider returned empty message")
	}

	return &GenerateResponse{
		Message: msg,
		Model:   resp.Model,
		Usage: &Usage{
			PromptTokens:     resp.Usage.PromptTokens,
			CompletionTokens: resp.Usage.CompletionTokens,
			TotalTokens:      resp.Usage.TotalTokens,
		},
	}, nil
}
