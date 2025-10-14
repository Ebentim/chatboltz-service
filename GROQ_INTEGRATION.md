# Groq AI Provider Integration

This document describes the complete Groq AI provider implementation for the Helix AI agent system.

## Overview

Groq provides extremely fast inference for various AI models, including large language models and audio processing models. This integration allows agents to use Groq's high-performance infrastructure for:

- **Text Generation**: Using models like Llama, Mixtral, and Gemma
- **Audio Transcription**: Using Whisper models
- **Streaming Responses**: Real-time text generation

## Supported Models

### Text Generation Models

- `llama-3.3-70b-versatile`: Latest Llama model with high versatility
- `llama-3.1-8b-instant`: Fast 8B parameter model for quick responses
- `llama-3.1-70b-versatile`: High-quality 70B parameter model
- `gemma2-9b-it`: Google's instruction-tuned Gemma model
- `mixtral-8x7b-32768`: Mistral's mixture of experts model with large context

### Audio Models

- `whisper-large-v3`: OpenAI's Whisper for audio transcription
- `whisper-large-v3-turbo`: Faster variant of Whisper
- `distil-whisper-large-v3-en`: English-optimized distilled Whisper

## Implementation Details

### Provider Structure

```go
type GroqAIProvider struct {
    client *openai.Client  // Uses OpenAI-compatible client
}
```

### Key Features

1. **OpenAI-Compatible API**: Uses the OpenAI Go client with Groq's base URL
2. **Complete LLMProvider Interface**: Implements all required methods
3. **Streaming Support**: Real-time response generation
4. **Multimodal Handling**: Converts image inputs to text descriptions
5. **Error Handling**: Comprehensive error wrapping and reporting
6. **Configuration Support**: Supports temperature, max_tokens, top_p parameters

### Factory Integration

The provider is integrated into the factory system:

```go
// In entity/model.go
const (
    OpenAI LLMProvider = iota
    Anthropic
    Meta
    Google
    Groq  // Added
    StabilityAI
    HuggingFace
)

// In factory.go
entity.Groq: func(apiKey string) (LLMProvider, error) {
    return NewGroqAIClient(apiKey, "https://api.groq.com/openai/v1")
},
```

### Agent Configuration

Agents can use Groq by setting:

```go
agent := entity.Agent{
    AiProvider: "groq",
    AiModel:    "llama-3.3-70b-versatile",
    // ... other fields
}
```

## Usage Examples

### Basic Chat Completion

```go
provider, err := NewGroqAIClient(apiKey, "https://api.groq.com/openai/v1")
if err != nil {
    log.Fatal(err)
}

conversation := Conversation{
    Messages: []Message{
        {Role: RoleSystem, Content: "You are a helpful assistant."},
        {Role: RoleUser, Content: "Hello!"},
    },
}

config := map[string]interface{}{
    "model":       "llama-3.3-70b-versatile",
    "temperature": 0.7,
    "max_tokens":  1000,
}

response, err := provider.CompleteConversation(conversation, config)
```

### Streaming Response

```go
err := provider.CompleteConversationStream(conversation, config, func(chunk string, done bool) error {
    if done {
        fmt.Println("\n[Stream complete]")
        return nil
    }
    fmt.Print(chunk)
    return nil
})
```

### Through Agent System

```go
factory := NewProviderFactory()
agent := entity.Agent{
    AiProvider: "groq",
    AiModel:    "llama-3.1-8b-instant",
}

provider, err := factory.GetProviderFromAgent(agent, apiKey)
response, err := provider.CompleteConversation(conversation, config)
```

## API Configuration

### Environment Variables

- `GROQ_API_KEY`: Your Groq API key from [console.groq.com](https://console.groq.com)

### Base URL

- Production: `https://api.groq.com/openai/v1`
- The provider uses OpenAI-compatible endpoints

### Rate Limits

Groq has generous rate limits but refer to their documentation for current limits:

- Free tier: High requests per minute
- Paid tier: Even higher limits with priority access

## Performance Characteristics

### Speed

- **Extremely Fast**: Groq is optimized for speed with custom hardware
- **Low Latency**: Sub-second response times for most models
- **High Throughput**: Can handle many concurrent requests

### Model Selection

- **llama-3.1-8b-instant**: Best for speed-critical applications
- **llama-3.3-70b-versatile**: Best for quality and complex reasoning
- **mixtral-8x7b-32768**: Best for large context applications

## Error Handling

The provider includes comprehensive error handling:

- **API Errors**: Wrapped with "Groq API error:" prefix
- **Network Errors**: Properly propagated
- **Streaming Errors**: Handled in callback mechanism
- **Authentication**: Clear error messages for invalid API keys

## Limitations

1. **Vision**: Currently not supported (images converted to text descriptions)
2. **Text-to-Speech**: Not supported by Groq (use TTS fallback services)
3. **Fine-tuning**: Limited fine-tuning options compared to other providers

## Integration Status

✅ **Complete Features:**

- [x] LLMProvider interface implementation
- [x] Factory integration
- [x] Agent system integration
- [x] Chat completion
- [x] Streaming responses
- [x] Multimodal message handling
- [x] Model capability mapping
- [x] Error handling
- [x] Configuration parameters

🔄 **Future Enhancements:**

- [ ] Native audio transcription API integration
- [ ] Batch processing support
- [ ] Fine-tuning integration when available
- [ ] Advanced function calling support

## Testing

The integration includes comprehensive tests covering:

- Provider creation through factory
- Agent-based provider instantiation
- Interface compliance
- Error handling
- Capability reporting
- Model configuration

## Troubleshooting

### Common Issues

1. **401 Unauthorized**: Check your `GROQ_API_KEY`
2. **Model not found**: Verify model name in configuration
3. **Rate limiting**: Implement retry logic with backoff
4. **Network timeouts**: Configure appropriate timeouts

### Debug Mode

Enable debug logging to see API requests:

```go
// Add logging to see request details
fmt.Printf("Making request to Groq with model: %s\n", model)
```

The Groq AI provider integration is now complete and production-ready!
