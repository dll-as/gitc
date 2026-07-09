package backend

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/dll-as/gitc/internal/ai/transport"
	"github.com/dll-as/gitc/internal/config"
)

type OllamaProvider struct {
	name   string
	config *config.Config
	client *transport.Client
}

func NewOllamaProvider(cfg *config.Config) (Provider, error) {
	client, err := transport.New(transport.Config{
		Timeout:    cfg.Backend.Timeout,
		MaxRetries: cfg.Backend.MaxRetries,
	})
	if err != nil {
		return nil, err
	}

	return &OllamaProvider{
		name:   "ollama",
		config: cfg,
		client: client,
	}, nil
}

func (p *OllamaProvider) Name() string {
	return p.name
}

func (p *OllamaProvider) Generate(ctx context.Context, req *GenerateRequest) (*GenerateResponse, error) {
	payload := OllamaRequest{
		Model: p.config.Backend.Model,
		Messages: []Message{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: req.Prompt},
		},
		Stream: false,
		Options: &OllamaOptions{
			Temperature: req.Temperature,
			TopP:        req.TopP,
			NumPredict:  req.MaxTokens,
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	respBody, err := p.client.DoWithContext(ctx, &transport.Request{
		Method: "POST",
		URL:    strings.TrimRight(p.config.Backend.BaseURL, "/") + "/chat/completions",
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: body,
	})
	if err != nil {
		return nil, err
	}

	if p.config.Prompt.Debug {
		fmt.Println("========== Ollama Response ==========")
		fmt.Println(string(respBody))
		fmt.Println("=====================================")
	}

	var resp OllamaResponse
	if err = json.Unmarshal(respBody, &resp); err != nil {
		return nil, err
	}

	msg := strings.TrimSpace(resp.Message.Content)
	if msg == "" {
		return nil, fmt.Errorf("ollama returned empty message")
	}

	return &GenerateResponse{
		Message: msg,
		Model:   resp.Model,
	}, nil
}
