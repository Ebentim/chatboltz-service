package aiprovider

import (
	"github.com/alpinesboltltd/boltz-ai/internal/entity"
)

type LLMManager struct {
	factory *ProviderFactory
	configs map[string]entity.ProviderConfig // keyed by agent ID or user ID
}

func NewLLMManager() *LLMManager {
	return &LLMManager{
		factory: NewProviderFactory(),
		configs: make(map[string]entity.ProviderConfig),
	}
}

// GetProviderForAgent returns the appropriate LLM provider for an agent
func (m *LLMManager) GetProviderForAgent(agent entity.Agent, apiKey string) (LLMProvider, error) {
	return m.factory.GetProviderFromAgent(agent, apiKey)
}

// GetProviderForChat returns provider based on chat context
func (m *LLMManager) GetProviderForChat(agentId string, agent entity.Agent, apiKey string) (LLMProvider, error) {
	// You can add logic here to select provider based on:
	// - Agent configuration
	// - User preferences
	// - Load balancing
	// - Cost optimization

	return m.GetProviderForAgent(agent, apiKey)
}

// GetMultimodalProvider returns provider with multimodal capabilities
func (m *LLMManager) GetMultimodalProvider(agent entity.Agent, apiKey, ttsKey, sttKey string) (LLMProvider, error) {
	return m.factory.GetMultimodalProvider(agent, apiKey, ttsKey, sttKey)
}

// ProcessMultimodalMessage handles different input types based on agent capabilities
func (m *LLMManager) ProcessMultimodalMessage(agent entity.Agent, messages []MultimodalMessage, apiKey, ttsKey, sttKey string) (string, error) {
	provider, err := m.GetMultimodalProvider(agent, apiKey, ttsKey, sttKey)
	if err != nil {
		return "", err
	}

	config := m.BuildConfig(entity.AgentBehavior{}, agent.AiModel.Name)
	return provider.CompleteMultimodalConversation(messages, config)
}

// BuildConfig creates conversation config from agent behavior
func (m *LLMManager) BuildConfig(behavior entity.AgentBehavior, model string) map[string]interface{} {
	config := make(map[string]interface{})

	if model != "" {
		config["model"] = model
	}

	if behavior.Temperature > 0 {
		config["temperature"] = behavior.Temperature
	}

	if behavior.MaxTokens > 0 {
		config["max_tokens"] = behavior.MaxTokens
	}

	return config
}
