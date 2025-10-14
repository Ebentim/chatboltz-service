package aiprovider

import (
	"fmt"

	"github.com/alpinesboltltd/boltz-ai/internal/entity"
)

type ProviderFactory struct {
	providers map[entity.LLMProvider]func(string) (LLMProvider, error)
}

func NewProviderFactory() *ProviderFactory {
	return &ProviderFactory{
		providers: map[entity.LLMProvider]func(string) (LLMProvider, error){
			entity.OpenAI: func(apiKey string) (LLMProvider, error) {
				return NewOpenAIClient(apiKey), nil
			},
			entity.Anthropic: func(apiKey string) (LLMProvider, error) {
				return NewAnthropicClient(apiKey), nil
			},
			entity.Google: func(apiKey string) (LLMProvider, error) {
				return NewGoogleAIClient(apiKey)
			},
			entity.Meta: func(apiKey string) (LLMProvider, error) {
				return NewMetaClient(apiKey), nil
			},
			entity.Groq: func(apiKey string) (LLMProvider, error) {
				return NewGroqAIClient(apiKey, "https://api.groq.com/openai/v1")
			},
		},
	}
}

func (f *ProviderFactory) CreateProvider(config entity.ProviderConfig) (LLMProvider, error) {
	if config.APIKey == "" {
		return nil, fmt.Errorf("API key is required")
	}

	providerFunc, exists := f.providers[config.Provider]
	if !exists {
		return nil, fmt.Errorf("unsupported provider: %v", config.Provider)
	}

	return providerFunc(config.APIKey)
}

// GetProviderFromAgent creates a provider based on agent configuration
func (f *ProviderFactory) GetProviderFromAgent(agent entity.Agent, apiKey string) (LLMProvider, error) {
	var provider entity.LLMProvider

	switch agent.AiProvider {
	case "openai":
		provider = entity.OpenAI
	case "anthropic":
		provider = entity.Anthropic
	case "google":
		provider = entity.Google
	case "meta":
		provider = entity.Meta
	case "groq":
		provider = entity.Groq
	default:
		return nil, fmt.Errorf("unknown provider: %s", agent.AiProvider)
	}

	capabilities := entity.DefaultModelCapabilities.GetCapabilities(agent.AiModel)
	return f.CreateProvider(entity.ProviderConfig{
		Provider:     provider,
		APIKey:       apiKey,
		Capabilities: capabilities,
	})
}

// GetMultimodalProvider returns provider with TTS/STT fallbacks for non-multimodal models
func (f *ProviderFactory) GetMultimodalProvider(agent entity.Agent, apiKey, ttsKey, sttKey string) (LLMProvider, error) {
	provider, err := f.GetProviderFromAgent(agent, apiKey)
	if err != nil {
		return nil, err
	}

	caps := provider.GetCapabilities()
	if !caps.Voice {
		// Use ElevenLabs for TTS and Deepgram for STT as fallback
		return NewMultimodalWrapper(provider, entity.ElevenLabs, entity.Deepgram, ttsKey, sttKey), nil
	}

	return provider, nil
}
