package ai

import (
	"fmt"

	"github.com/dll-as/gitc/internal/ai/backend"
	"github.com/dll-as/gitc/internal/config"
)

func New(cfg *config.Config) (backend.Provider, error) {
	switch cfg.Backend.Backend {

	case config.BackendOpenAI:
		return backend.NewOpenAICompatibleProvider(cfg)

	case config.BackendOllama:
		return backend.NewOllamaProvider(cfg)

	case config.BackendAnthropic:
		return backend.NewOpenAICompatibleProvider(cfg)

	default:
		return nil, fmt.Errorf("unsupported backend %q", cfg.Backend.Backend)
	}
}
