package entity

type LLMProvider int
type TTSProvider int
type STTProvider int
type Capability string

const (
	OpenAI LLMProvider = iota
	Anthropic
	Meta
	Google
	Groq
	StabilityAI
	HuggingFace
)

const (
	OpenAITTS TTSProvider = iota
	ElevenLabs
	StabilityAITTS
)

const (
	OpenAISTT STTProvider = iota
	Deepgram
)

const (
	CapabilityText   Capability = "text"
	CapabilityVoice  Capability = "voice"
	CapabilityVision Capability = "vision"
)

type ModelCapabilities struct {
	Text   bool `gorm:"type:boolean;default:false"`
	Voice  bool `gorm:"type:boolean;default:false"`
	Vision bool `gorm:"type:boolean;default:false"`
}

type ProviderConfig struct {
	Provider     LLMProvider
	APIKey       string
	TTSProvider  TTSProvider
	STTProvider  STTProvider
	TTSAPIKey    string
	STTAPIKey    string
	Capabilities ModelCapabilities
}
