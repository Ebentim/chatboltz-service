package usecase

import (
	"fmt"
	"sync"
	"time"

	aiprovider "github.com/alpinesboltltd/boltz-ai/internal/provider/ai-provider"
)

type ChatService struct {
	llmManager *aiprovider.LLMManager
	agentCache *AgentCache
	cache      map[string]string
	cacheMutex sync.RWMutex
}

func NewChatService(repo AgentRepository) *ChatService {
	return &ChatService{
		llmManager: aiprovider.NewLLMManager(),
		agentCache: NewAgentCache(repo, 30*time.Minute),
		cache:      make(map[string]string),
	}
}

func (s *ChatService) ProcessMessage(agentID, userMessage, apiKey string) (string, error) {
	// Quick cache check
	cacheKey := agentID + "_" + userMessage[:min(30, len(userMessage))]
	s.cacheMutex.RLock()
	if cached, found := s.cache[cacheKey]; found {
		s.cacheMutex.RUnlock()
		return cached, nil
	}
	s.cacheMutex.RUnlock()

	config, err := s.agentCache.GetAgentConfig(agentID)
	if err != nil {
		return "", fmt.Errorf("failed to get agent config: %w", err)
	}

	provider, err := s.llmManager.GetProviderForAgent(config.Agent, apiKey)
	if err != nil {
		return "", fmt.Errorf("failed to get provider: %w", err)
	}

	conversation := aiprovider.Conversation{}
	systemContent := config.SystemInstruction.Content
	if systemContent == "" {
		systemContent = "You are a helpful assistant."
	}

	conversation.Messages = append(conversation.Messages, aiprovider.Message{
		Role:    aiprovider.RoleSystem,
		Content: systemContent,
	})
	conversation.Messages = append(conversation.Messages, aiprovider.Message{
		Role:    aiprovider.RoleUser,
		Content: userMessage,
	})

	llmConfig := s.llmManager.BuildConfig(config.Behavior, config.Agent.AiModel)
	result, err := provider.CompleteConversation(conversation, llmConfig)
	if err != nil {
		return "", err
	}

	// Cache result
	s.cacheMutex.Lock()
	s.cache[cacheKey] = result
	s.cacheMutex.Unlock()

	return result, nil
}

// ProcessMessageStream provides streaming responses for sub-500ms initial response
func (s *ChatService) ProcessMessageStream(agentID, userMessage, apiKey string, callback aiprovider.StreamCallback) error {
	config, err := s.agentCache.GetAgentConfig(agentID)
	if err != nil {
		return fmt.Errorf("failed to get agent config: %w", err)
	}

	provider, err := s.llmManager.GetProviderForAgent(config.Agent, apiKey)
	if err != nil {
		return fmt.Errorf("failed to get provider: %w", err)
	}

	conversation := aiprovider.Conversation{}
	systemContent := config.SystemInstruction.Content
	if systemContent == "" {
		systemContent = "You are a helpful assistant."
	}

	conversation.Messages = append(conversation.Messages, aiprovider.Message{
		Role:    aiprovider.RoleSystem,
		Content: systemContent,
	})
	conversation.Messages = append(conversation.Messages, aiprovider.Message{
		Role:    aiprovider.RoleUser,
		Content: userMessage,
	})

	llmConfig := s.llmManager.BuildConfig(config.Behavior, config.Agent.AiModel)
	return provider.CompleteConversationStream(conversation, llmConfig, callback)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (s *ChatService) ProcessConversation(agentID string, messages []aiprovider.Message, apiKey string) (string, error) {
	config, err := s.agentCache.GetAgentConfig(agentID)
	if err != nil {
		return "", fmt.Errorf("failed to get agent config: %w", err)
	}

	provider, err := s.llmManager.GetProviderForAgent(config.Agent, apiKey)
	if err != nil {
		return "", fmt.Errorf("failed to get provider: %w", err)
	}

	conversation := aiprovider.Conversation{Messages: messages}
	llmConfig := s.llmManager.BuildConfig(config.Behavior, config.Agent.AiModel)
	return provider.CompleteConversation(conversation, llmConfig)
}

// ProcessMultimodalMessage handles text, voice, and vision inputs
func (s *ChatService) ProcessMultimodalMessage(agentID string, messages []aiprovider.MultimodalMessage, apiKey, ttsKey, sttKey string) (string, error) {
	config, err := s.agentCache.GetAgentConfig(agentID)
	if err != nil {
		return "", fmt.Errorf("failed to get agent config: %w", err)
	}
	// FIXME. agent must support the request type
	// Check if agent supports required capabilities
	// requiredCaps := s.getRequiredCapabilities(messages)
	// if !s.agentSupportsCapabilities(config.Agent, requiredCaps) {
	// 	return "", fmt.Errorf("agent does not support required capabilities: %v", requiredCaps)
	// }

	return s.llmManager.ProcessMultimodalMessage(config.Agent, messages, apiKey, ttsKey, sttKey)
}

func (s *ChatService) GetRequiredCapabilities(messages []aiprovider.MultimodalMessage) []string {
	caps := make(map[string]bool)
	for _, msg := range messages {
		switch msg.MediaType {
		case "audio":
			caps["voice"] = true
		case "image", "video":
			caps["vision"] = true
		default:
			caps["text"] = true
		}
	}

	result := make([]string, 0, len(caps))
	for cap := range caps {
		result = append(result, cap)
	}
	return result
}

// FIXME: need to ensure that agent's type is properly mapped to request
// func (s *ChatService) agentSupportsCapabilities(agent entity.Agent, required []string) bool {
// 	for _, req := range required {
// 		found := false
// 		for _, cap := range agent.AgentType {
// 			if cap == req {
// 				found = true
// 				break
// 			}
// 		}
// 		if !found {
// 			return false
// 		}
// 	}
// 	return true
// }
