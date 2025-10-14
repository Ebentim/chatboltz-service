package main

/*
import (
	"fmt"
	"log"

	"github.com/alpinesboltltd/boltz-ai/internal/config"
	"github.com/alpinesboltltd/boltz-ai/internal/entity"
	"github.com/alpinesboltltd/boltz-ai/internal/usecase"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

// Example implementation of AgentRepository
type MockAgentRepo struct{}

func (r *MockAgentRepo) GetAgentByID(id string) (*entity.Agent, error) {
	return &entity.Agent{
		Id:         id,
		AiProvider: "google",
		AiModel:    "gemini-2.5-flash",
	}, nil
}

func (r *MockAgentRepo) GetAgentBehavior(agentID string) (*entity.AgentBehavior, error) {
	return &entity.AgentBehavior{
		SystemInstructionId: "sys-123",
		Temperature:         0.7,
		MaxTokens:           1000,
	}, nil
}

func (r *MockAgentRepo) GetSystemInstruction(id string) (*entity.SystemInstruction, error) {
	return &entity.SystemInstruction{
		Id:      id,
		Content: "You are a helpful customer support agent for our e-commerce platform.",
	}, nil
}

func (r *MockAgentRepo) GetPromptTemplate(id string) (*entity.PromptTemplate, error) {
	return &entity.PromptTemplate{
		Id:      id,
		Content: "Please respond professionally and helpfully.",
	}, nil
}

func main() {
	// Initialize the chat service with repository
	repo := &MockAgentRepo{}
	chatService := usecase.NewChatService(repo)

	godotenv.Load(".env")
	var cfg config.Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal(err)
	}
	// Process a message - agent config is automatically cached
	response, err := chatService.ProcessMessage(
		"agent-123",
		"Hello, I need help with my order",
		cfg.GEMINI_API_KEY,
	)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Response:", response)
}
*/
/*

package aiprovider

import (
	"context"
	"time"

	"github.com/alpinesboltltd/boltz-ai/internal/entity"
)

type StreamCallback func(chunk string, done bool) error

type StreamingProvider interface {
	CompleteConversationStream(conversation Conversation, config map[string]interface{}, callback StreamCallback) error
	CompleteMultimodalConversationStream(messages []MultimodalMessage, config map[string]interface{}, callback StreamCallback) error
}

// FastResponse provides immediate acknowledgment while processing continues
type FastResponse struct {
	Immediate string `json:"immediate"`
	StreamID  string `json:"stream_id"`
	Status    string `json:"status"`
}

// ResponseCache for sub-500ms responses
type ResponseCache struct {
	cache map[string]CachedResponse
	ttl   time.Duration
}

type CachedResponse struct {
	Response  string
	Timestamp time.Time
}

func NewResponseCache(ttl time.Duration) *ResponseCache {
	return &ResponseCache{
		cache: make(map[string]CachedResponse),
		ttl:   ttl,
	}
}

func (c *ResponseCache) Get(key string) (string, bool) {
	resp, exists := c.cache[key]
	if !exists || time.Since(resp.Timestamp) > c.ttl {
		return "", false
	}
	return resp.Response, true
}

func (c *ResponseCache) Set(key string, response string) {
	c.cache[key] = CachedResponse{
		Response:  response,
		Timestamp: time.Now(),
	}
}

// FastLLMProvider wraps providers with speed optimizations
type FastLLMProvider struct {
	provider LLMProvider
	cache    *ResponseCache
	timeout  time.Duration
}

func NewFastLLMProvider(provider LLMProvider, cacheTimeout time.Duration) *FastLLMProvider {
	return &FastLLMProvider{
		provider: provider,
		cache:    NewResponseCache(cacheTimeout),
		timeout:  450 * time.Millisecond, // Leave 50ms buffer
	}
}

func (f *FastLLMProvider) CompleteConversation(conversation Conversation, config map[string]interface{}) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), f.timeout)
	defer cancel()

	// Check cache first
	cacheKey := f.generateCacheKey(conversation, config)
	if cached, found := f.cache.Get(cacheKey); found {
		return cached, nil
	}

	// Use channel for timeout handling
	resultChan := make(chan string, 1)
	errorChan := make(chan error, 1)

	go func() {
		result, err := f.provider.CompleteConversation(conversation, config)
		if err != nil {
			errorChan <- err
			return
		}
		f.cache.Set(cacheKey, result)
		resultChan <- result
	}()

	select {
	case result := <-resultChan:
		return result, nil
	case err := <-errorChan:
		return "", err
	case <-ctx.Done():
		return "Processing your request...", nil // Immediate response
	}
}

func (f *FastLLMProvider) generateCacheKey(conversation Conversation, config map[string]interface{}) string {
	// Simple hash of last message + model
	if len(conversation.Messages) == 0 {
		return ""
	}
	lastMsg := conversation.Messages[len(conversation.Messages)-1].Content
	model := ""
	if m, ok := config["model"].(string); ok {
		model = m
	}
	return lastMsg[:min(50, len(lastMsg))] + "_" + model
}

func (f *FastLLMProvider) GetCapabilities() entity.ModelCapabilities {
	return f.provider.GetCapabilities()
}

func (f *FastLLMProvider) CompleteMultimodalConversation(messages []MultimodalMessage, config map[string]interface{}) (string, error) {
	return f.provider.CompleteMultimodalConversation(messages, config)
}


*/

/*
package usecase

import (
	"context"
	"fmt"
	"strings"
	"time"

	aiprovider "github.com/alpinesboltltd/boltz-ai/internal/provider/ai-provider"
)

// FastChatService optimized for sub-500ms responses
type FastChatService struct {
	*ChatService
	intentClassifier *IntentClassifier
	quickResponses   map[string]string
}

type IntentClassifier struct {
	patterns map[string][]string
}

func NewFastChatService(repo AgentRepository) *FastChatService {
	return &FastChatService{
		ChatService:      NewChatService(repo),
		intentClassifier: NewIntentClassifier(),
		quickResponses:   getQuickResponses(),
	}
}

func NewIntentClassifier() *IntentClassifier {
	return &IntentClassifier{
		patterns: map[string][]string{
			"greeting":     {"hello", "hi", "hey", "good morning", "good afternoon"},
			"help":         {"help", "support", "assist", "problem", "issue"},
			"order":        {"order", "purchase", "buy", "payment", "checkout"},
			"status":       {"status", "track", "where is", "delivery", "shipping"},
			"cancel":       {"cancel", "refund", "return", "stop"},
			"information":  {"what", "how", "when", "where", "why", "info"},
		},
	}
}

func getQuickResponses() map[string]string {
	return map[string]string{
		"greeting":     "Hello! How can I help you today?",
		"help":         "I'm here to help! What specific issue can I assist you with?",
		"status":       "I can help you check your order status. Could you provide your order number?",
		"information":  "I'd be happy to provide information. What would you like to know about?",
	}
}

// ProcessMessageFast provides sub-500ms responses using multiple optimization strategies
func (s *FastChatService) ProcessMessageFast(agentID, userMessage, apiKey string) (string, error) {
	start := time.Now()

	// Strategy 1: Quick pattern matching (< 10ms)
	if intent := s.intentClassifier.ClassifyIntent(userMessage); intent != "" {
		if quickResponse, exists := s.quickResponses[intent]; exists {
			return quickResponse, nil
		}
	}

	// Strategy 2: Cache lookup (< 50ms)
	cacheKey := s.generateCacheKey(agentID, userMessage)
	if cached, found := s.responseCache.Get(cacheKey); found {
		return cached, nil
	}

	// Strategy 3: Concurrent processing with timeout (< 450ms)
	ctx, cancel := context.WithTimeout(context.Background(), 400*time.Millisecond)
	defer cancel()

	resultChan := make(chan string, 1)
	errorChan := make(chan error, 1)

	go func() {
		result, err := s.ChatService.ProcessMessage(agentID, userMessage, apiKey)
		if err != nil {
			errorChan <- err
			return
		}
		s.responseCache.Set(cacheKey, result)
		resultChan <- result
	}()

	select {
	case result := <-resultChan:
		elapsed := time.Since(start)
		if elapsed > 500*time.Millisecond {
			// Log slow response for optimization
			fmt.Printf("Slow response: %v for message: %s\n", elapsed, userMessage[:min(50, len(userMessage))])
		}
		return result, nil
	case err := <-errorChan:
		return "", err
	case <-ctx.Done():
		// Strategy 4: Immediate acknowledgment with background processing
		go s.processInBackground(agentID, userMessage, apiKey, cacheKey)
		return s.generateImmediateResponse(userMessage), nil
	}
}

func (s *FastChatService) processInBackground(agentID, userMessage, apiKey, cacheKey string) {
	result, err := s.ChatService.ProcessMessage(agentID, userMessage, apiKey)
	if err == nil {
		s.responseCache.Set(cacheKey, result)
	}
}

func (s *FastChatService) generateImmediateResponse(userMessage string) string {
	intent := s.intentClassifier.ClassifyIntent(userMessage)
	switch intent {
	case "help":
		return "I'm analyzing your request and will provide detailed help shortly..."
	case "order":
		return "Looking up your order information now..."
	case "status":
		return "Checking the current status for you..."
	default:
		return "Processing your request..."
	}
}

func (s *FastChatService) generateCacheKey(agentID, message string) string {
	// Create deterministic cache key
	msgHash := message
	if len(message) > 100 {
		msgHash = message[:100]
	}
	return fmt.Sprintf("%s_%s", agentID, strings.ReplaceAll(msgHash, " ", "_"))
}

func (ic *IntentClassifier) ClassifyIntent(message string) string {
	message = strings.ToLower(message)

	for intent, patterns := range ic.patterns {
		for _, pattern := range patterns {
			if strings.Contains(message, pattern) {
				return intent
			}
		}
	}
	return ""
}

// StreamResponse provides real-time streaming for longer responses
func (s *FastChatService) StreamResponse(agentID, userMessage, apiKey string, callback func(chunk string)) error {
	// Immediate acknowledgment
	callback("I'm processing your request")

	// Stream the actual response
	config, err := s.agentCache.GetAgentConfig(agentID)
	if err != nil {
		return err
	}

	provider, err := s.llmManager.GetProviderForAgent(config.Agent, apiKey)
	if err != nil {
		return err
	}

	// Check if provider supports streaming
	if streamProvider, ok := provider.(aiprovider.StreamingProvider); ok {
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
		return streamProvider.CompleteConversationStream(conversation, llmConfig, func(chunk string, done bool) error {
			callback(chunk)
			return nil
		})
	}

	// Fallback to regular processing
	result, err := s.ProcessMessage(agentID, userMessage, apiKey)
	if err != nil {
		return err
	}
	callback(result)
	return nil
}
*/
