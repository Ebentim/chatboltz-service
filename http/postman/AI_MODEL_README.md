# AI Models API Documentation

## Overview
The AI Models API allows management of AI models used by agents. Only **SuperAdmin** users can create, update, or delete models. All authenticated users can view models.

## Authentication
All endpoints require JWT authentication via Bearer token in the Authorization header.

## Endpoints

### 1. Create AI Model
**POST** `/api/ai-models`

**Authorization:** SuperAdmin only

**Request Body:**
```json
{
  "name": "gpt-4o",
  "provider": "openai",
  "credits_per_1k": 10,
  "supports_text": true,
  "supports_vision": true,
  "supports_voice": true,
  "is_reasoning": false
}
```

**Response (201):**
```json
{
  "ai_model": {
    "id": "uuid",
    "name": "gpt-4o",
    "provider": "openai",
    "credits_per_1k": 10,
    "supports_text": true,
    "supports_vision": true,
    "supports_voice": true,
    "is_reasoning": false,
    "created_at": "2025-02-01T10:00:00Z",
    "updated_at": "2025-02-01T10:00:00Z"
  }
}
```

---

### 2. Get AI Model
**GET** `/api/ai-models/:modelId`

**Authorization:** Any authenticated user

**Response (200):**
```json
{
  "ai_model": {
    "id": "uuid",
    "name": "gpt-4o",
    "provider": "openai",
    "credits_per_1k": 10,
    "supports_text": true,
    "supports_vision": true,
    "supports_voice": true,
    "is_reasoning": false,
    "created_at": "2025-02-01T10:00:00Z",
    "updated_at": "2025-02-01T10:00:00Z"
  }
}
```

---

### 3. List All AI Models
**GET** `/api/ai-models`

**Authorization:** Any authenticated user

**Response (200):**
```json
{
  "ai_models": [
    {
      "id": "uuid-1",
      "name": "gpt-4o",
      "provider": "openai",
      "credits_per_1k": 10,
      "supports_text": true,
      "supports_vision": true,
      "supports_voice": true,
      "is_reasoning": false
    },
    {
      "id": "uuid-2",
      "name": "claude-3-5-sonnet-20241022",
      "provider": "anthropic",
      "credits_per_1k": 12,
      "supports_text": true,
      "supports_vision": true,
      "supports_voice": false,
      "is_reasoning": false
    }
  ]
}
```

---

### 4. List AI Models by Provider
**GET** `/api/ai-models?provider=openai`

**Authorization:** Any authenticated user

**Query Parameters:**
- `provider` (string): Filter by provider (openai, anthropic, google, groq, meta)

**Response (200):**
```json
{
  "ai_models": [
    {
      "id": "uuid-1",
      "name": "gpt-4o",
      "provider": "openai",
      "credits_per_1k": 10,
      "supports_text": true,
      "supports_vision": true,
      "supports_voice": true,
      "is_reasoning": false
    },
    {
      "id": "uuid-2",
      "name": "gpt-4o-mini",
      "provider": "openai",
      "credits_per_1k": 2,
      "supports_text": true,
      "supports_vision": true,
      "supports_voice": true,
      "is_reasoning": false
    }
  ]
}
```

---

### 5. Update AI Model
**PUT** `/api/ai-models/:modelId`

**Authorization:** SuperAdmin only

**Request Body (all fields optional):**
```json
{
  "name": "gpt-4o-updated",
  "provider": "openai",
  "credits_per_1k": 12,
  "supports_text": true,
  "supports_vision": true,
  "supports_voice": true,
  "is_reasoning": false
}
```

**Response (200):**
```json
{
  "ai_model": {
    "id": "uuid",
    "name": "gpt-4o-updated",
    "provider": "openai",
    "credits_per_1k": 12,
    "supports_text": true,
    "supports_vision": true,
    "supports_voice": true,
    "is_reasoning": false,
    "created_at": "2025-02-01T10:00:00Z",
    "updated_at": "2025-02-01T11:00:00Z"
  }
}
```

---

### 6. Delete AI Model
**DELETE** `/api/ai-models/:modelId`

**Authorization:** SuperAdmin only

**Response (200):**
```json
{
  "message": "AI model deleted successfully"
}
```

**Note:** Deletion will fail if any agents are currently using this model (foreign key constraint).

---

## Supported Providers
- `openai` - OpenAI (GPT models)
- `anthropic` - Anthropic (Claude models)
- `google` - Google (Gemini models)
- `groq` - Groq (Fast inference)
- `meta` - Meta (Llama models)

## Model Capabilities
- `supports_text`: Model can process text input/output
- `supports_vision`: Model can process images/video
- `supports_voice`: Model can process audio input/output
- `is_reasoning`: Model uses chain-of-thought reasoning (e.g., o1, o1-mini)

## Error Responses

**403 Forbidden:**
```json
{
  "error": "Only superadmin can create AI models"
}
```

**404 Not Found:**
```json
{
  "error": "AI model not found"
}
```

**400 Bad Request:**
```json
{
  "error": "Invalid request format"
}
```

## Pre-seeded Models
The system comes with 15+ pre-configured models including:
- OpenAI: gpt-4o, gpt-4o-mini, gpt-4-turbo, o1, o1-mini
- Anthropic: claude-3-5-sonnet, claude-3-5-haiku, claude-3-opus
- Google: gemini-2.0-flash-exp, gemini-1.5-pro, gemini-1.5-flash
- Groq: llama-3.3-70b-versatile, llama-3.1-8b-instant
