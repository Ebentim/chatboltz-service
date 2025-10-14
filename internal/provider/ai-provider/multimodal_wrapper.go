package aiprovider

import (
	"fmt"
	"net/http"

	"github.com/alpinesboltltd/boltz-ai/internal/entity"
)

// MultimodalWrapper wraps non-multimodal providers with external TTS/STT
type MultimodalWrapper struct {
	llmProvider LLMProvider
	ttsProvider TTSProvider
	sttProvider STTProvider
}

func NewMultimodalWrapper(llm LLMProvider, ttsType entity.TTSProvider, sttType entity.STTProvider, ttsKey, sttKey string) *MultimodalWrapper {
	var tts TTSProvider
	var stt STTProvider

	switch ttsType {
	case entity.ElevenLabs:
		tts = NewElevenLabsProvider(ttsKey)
	case entity.OpenAITTS:
		tts = NewOpenAITTSProvider(ttsKey)
	}

	switch sttType {
	case entity.Deepgram:
		stt = NewDeepgramProvider(sttKey)
	case entity.OpenAISTT:
		stt = NewOpenAISTTProvider(sttKey)
	}

	return &MultimodalWrapper{
		llmProvider: llm,
		ttsProvider: tts,
		sttProvider: stt,
	}
}

func (w *MultimodalWrapper) CompleteConversation(conversation Conversation, config map[string]interface{}) (string, error) {
	return w.llmProvider.CompleteConversation(conversation, config)
}

func (w *MultimodalWrapper) CompleteMultimodalConversation(messages []MultimodalMessage, config map[string]interface{}) (string, error) {
	// Process audio inputs with STT
	processedMessages := make([]MultimodalMessage, 0, len(messages))
	for _, msg := range messages {
		if msg.MediaType == "audio" && len(msg.MediaData) > 0 {
			text, err := w.sttProvider.SpeechToText(msg.MediaData, config)
			if err != nil {
				return "", fmt.Errorf("STT failed: %w", err)
			}
			msg.Content = text
			msg.MediaType = "text"
		}
		processedMessages = append(processedMessages, msg)
	}

	// Use underlying provider's multimodal capability if available
	if multimodal, ok := w.llmProvider.(interface {
		CompleteMultimodalConversation([]MultimodalMessage, map[string]interface{}) (string, error)
	}); ok {
		return multimodal.CompleteMultimodalConversation(processedMessages, config)
	}

	// Fallback to text-only conversation
	conv := Conversation{}
	for _, msg := range processedMessages {
		conv.Messages = append(conv.Messages, Message{
			Role:    Role(msg.Role),
			Content: msg.Content,
		})
	}
	return w.llmProvider.CompleteConversation(conv, config)
}

func (w *MultimodalWrapper) CompleteConversationStream(conversation Conversation, config map[string]interface{}, callback StreamCallback) error {
	return w.llmProvider.CompleteConversationStream(conversation, config, callback)
}

func (w *MultimodalWrapper) GetCapabilities() entity.ModelCapabilities {
	caps := w.llmProvider.GetCapabilities()
	// Add voice capability through external providers
	if w.ttsProvider != nil && w.sttProvider != nil {
		caps.Voice = true
	}
	return caps
}

// Placeholder implementations for external providers
type ElevenLabsProvider struct {
	apiKey string
}

func NewElevenLabsProvider(apiKey string) *ElevenLabsProvider {
	return &ElevenLabsProvider{apiKey: apiKey}
}

func (p *ElevenLabsProvider) TextToSpeech(text string, config map[string]interface{}) ([]byte, error) {
	// TODO: Implement ElevenLabs TTS API call
	return nil, fmt.Errorf("ElevenLabs TTS not implemented")
}

type DeepgramProvider struct {
	apiKey string
}

func NewDeepgramProvider(apiKey string) *DeepgramProvider {
	return &DeepgramProvider{apiKey: apiKey}
}

func (p *DeepgramProvider) SpeechToText(audio []byte, config map[string]interface{}) (string, error) {
	// TODO: Implement Deepgram STT API call
	return "", fmt.Errorf("Deepgram STT not implemented")
}
func (p *DeepgramProvider) TextToSpeech(audio []byte, config map[string]interface{}) (string, error) {
	// TODO: Implement Deepgram STT API call
	resp, err := http.Get("https://api.deepgram.com/v1/speal?model=aura-2-thalia-en")

	if err != nil {
		return "", fmt.Errorf("Deepgram returned error")
	}

	if resp.StatusCode == 200 {
		resp.Body.Close()
	}

	return "", fmt.Errorf("Deepgram TTS not implemented")
}

type OpenAITTSProvider struct {
	apiKey string
}

func NewOpenAITTSProvider(apiKey string) *OpenAITTSProvider {
	return &OpenAITTSProvider{apiKey: apiKey}
}

func (p *OpenAITTSProvider) TextToSpeech(text string, config map[string]interface{}) ([]byte, error) {
	// TODO: Implement OpenAI TTS API call
	return nil, fmt.Errorf("OpenAI TTS not implemented")
}

type OpenAISTTProvider struct {
	apiKey string
}

func NewOpenAISTTProvider(apiKey string) *OpenAISTTProvider {
	return &OpenAISTTProvider{apiKey: apiKey}
}

func (p *OpenAISTTProvider) SpeechToText(audio []byte, config map[string]interface{}) (string, error) {
	// TODO: Implement OpenAI STT API call
	return "", fmt.Errorf("OpenAI STT not implemented")
}
