package generic

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/bytedance/sonic"
	"github.com/dll-as/gitc/internal/ai"
	"github.com/dll-as/gitc/pkg/utils"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpproxy"
)

// Default URLs for supported providers
const (
	defaultOpenAIURL   = "https://api.openai.com/v1/chat/completions"
	defaultGrokURL     = "https://api.x.ai/v1/chat/completions"
	defaultDeepSeekURL = "https://api.deepseek.com/v1/chat/completions"
	defaultOllamaURL   = "http://localhost:11434//api/generate"
	systemPrompt       = "You are an AI assistant that generates concise and meaningful Git commit messages."
)

// GenericProvider implements the AIProvider interface for OpenAI-compatible APIs
type GenericProvider struct {
	apiKey   string
	client   *fasthttp.Client
	url      string
	provider string
}

// NewGenericProvider creates a new provider for OpenAI-compatible APIs
func NewGenericProvider(apiKey, proxy, url, provider string) (*GenericProvider, error) {
	if apiKey == "" && provider != "ollama" {
		return nil, errors.New("API key is required")
	}

	if url == "" {
		switch provider {
		case "openai":
			url = defaultOpenAIURL
		case "grok":
			url = defaultGrokURL
		case "deepseek":
			url = defaultDeepSeekURL
		case "ollama":
			url = defaultOllamaURL
		default:
			return nil, fmt.Errorf("no default URL for provider: %s", provider)
		}
	}

	client := &fasthttp.Client{
		MaxConnsPerHost: 10,
	}

	if proxy != "" {
		client.Dial = fasthttpproxy.FasthttpHTTPDialer(proxy)
	}

	return &GenericProvider{
		apiKey:   apiKey,
		client:   client,
		url:      url,
		provider: provider,
	}, nil
}

// Request structures for different providers
type OpenAIRequest struct {
	Model       string          `json:"model"`
	Messages    []OpenAIMessage `json:"messages"`
	MaxTokens   int             `json:"max_tokens,omitempty"`
	Temperature float64         `json:"temperature,omitempty"`
	Stream      bool            `json:"stream,omitempty"`
}

type OllamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
	// Options struct {
	// 	Temperature float64 `json:"temperature,omitempty"`
	// 	NumPredict  int     `json:"num_predict,omitempty"`
	// } `json:"options,omitempty"`
}

type OpenAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenAIResponse struct {
	Choices []struct {
		Message OpenAIMessage `json:"message"`
	} `json:"choices"`
	Error struct {
		Message string `json:"message"`
	} `json:"error"`
}

type OllamaResponse struct {
	Model    string `json:"model"`
	Response string `json:"response"`
	Error    string `json:"error,omitempty"`
}

// GenerateCommitMessage generates a commit message using the API
func (p *GenericProvider) GenerateCommitMessage(ctx context.Context, diff string, opts ai.MessageOptions) (string, error) {
	// Adjust prompt based on provider if needed
	prompt := utils.GetPromptForSingleCommit(diff, opts)

	var reqBody []byte
	var err error

	if p.provider == "ollama" {
		ollamaReq := OllamaRequest{
			Model:  opts.Model,
			Prompt: prompt,
			Stream: false,
		}
		reqBody, err = sonic.Marshal(ollamaReq)
	} else {
		openaiReq := OpenAIRequest{
			Model: opts.Model,
			Messages: []OpenAIMessage{
				{"system", systemPrompt},
				{"user", prompt},
			},
			MaxTokens:   max(512, opts.MaxLength), // More tokens for complete messages
			Temperature: opts.Temperature,         // Slightly creative but controlled
			Stream:      false,
		}
		reqBody, err = sonic.Marshal(openaiReq)
	}
	if err != nil {
		return "", fmt.Errorf("failed to encode JSON: %v", err)
	}

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	req.SetRequestURI(p.url)
	req.Header.SetMethod("POST")

	// Only set Authorization header for providers that need it
	if p.provider != "ollama" && p.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+p.apiKey)
	}

	req.Header.Set("Content-Type", "application/json")
	req.SetBody(reqBody)

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	if err = p.client.DoRedirects(req, resp, opts.MaxRedirects); err != nil {
		return "", fmt.Errorf("API request failed: %w", err)
	}

	var commitMessage string
	if p.provider == "ollama" {
		var ollamaRes OllamaResponse
		if err = sonic.Unmarshal(resp.Body(), &ollamaRes); err != nil {
			return "", fmt.Errorf("failed to parse Ollama response: %v", err)
		}

		if ollamaRes.Error != "" {
			return "", fmt.Errorf("Ollama API error: %s", ollamaRes.Error)
		}

		commitMessage = strings.TrimSpace(ollamaRes.Response)
	} else {
		var openaiRes OpenAIResponse
		if err = sonic.Unmarshal(resp.Body(), &openaiRes); err != nil {
			return "", fmt.Errorf("failed to parse response: %v", err)
		}

		if openaiRes.Error.Message != "" {
			return "", fmt.Errorf("API error from %s: %s", p.provider, openaiRes.Error.Message)
		} else if len(openaiRes.Choices) == 0 {
			return "", fmt.Errorf("no response from %s", p.provider)
		}

		commitMessage = strings.TrimSpace(openaiRes.Choices[0].Message.Content)
	}
	if commitMessage == "" {
		return "", fmt.Errorf("empty commit message generated by %s", p.provider)
	}

	return commitMessage, nil
}
