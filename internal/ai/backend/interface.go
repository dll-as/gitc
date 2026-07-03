package backend

import "context"

type Provider interface {
	Name() string
	Generate(ctx context.Context, req *GenerateRequest) (*GenerateResponse, error)
}

type GenerateRequest struct {
	Prompt string

	Model string

	MaxTokens   int
	Temperature float64
	TopP        float64
}

type GenerateResponse struct {
	Message string

	Model string

	FinishReason string

	Usage *Usage
}

type Usage struct {
	PromptTokens     int
	CompletionTokens int
	TotalTokens      int
}
