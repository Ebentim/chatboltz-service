# AI Model Refactoring Summary

## Overview
Refactored the Agent entity to extract AI model information into a separate `ai_models` table. This allows centralized management of AI models with their capabilities and pricing across all agents.

## Changes Made

### 1. New Entity: `AiModel` (`internal/entity/ai_model.go`)
- **Fields:**
  - `id`: Primary key
  - `name`: Unique model name (e.g., "gpt-4o", "claude-3-5-sonnet-20241022")
  - `provider`: Provider name (openai, anthropic, google, groq, meta)
  - `credits_per_1k`: Cost per 1000 tokens
  - `supports_text`: Boolean for text capability
  - `supports_vision`: Boolean for vision/image capability
  - `supports_voice`: Boolean for voice/audio capability
  - `is_reasoning`: Boolean to identify reasoning models (e.g., o1, o1-mini)
  - `created_at`, `updated_at`: Timestamps

### 2. Refactored Entity: `Agent` (`internal/entity/agent.go`)
- **Removed fields:**
  - `AiModel` (string)
  - `AiProvider` (string)
  - `CreditsPer1k` (int)

- **Added fields:**
  - `AiModelId` (string): Foreign key to ai_models table
  - `AiModel` (*AiModel): Relationship to load model details

- **Updated structs:**
  - `AgentUpdate`: Now uses `AiModelId` instead of model/provider/credits fields
  - `AgentResponse`: Includes both `AiModelId` and `AiModel` for API responses

### 3. Repository Layer
- **New:** `internal/repository/ai_model.go`
  - `CreateAiModel`: Create new AI model
  - `GetAiModel`: Get by ID
  - `GetAiModelByName`: Get by name
  - `ListAiModels`: List all models
  - `ListAiModelsByProvider`: Filter by provider
  - `UpdateAiModel`: Update model details
  - `DeleteAiModel`: Delete model (restricted if agents use it)

- **Updated:** `internal/repository/agent.go`
  - `CreateAgent`: Now takes `aiModelId` instead of model/provider/credits
  - `GetAgent`: Preloads `AiModel` relationship
  - `GetAgentsByUserId`: Preloads `AiModel` for list queries

- **Updated:** `internal/repository/interfaces.go`
  - Added `AiModelRepositoryInterface`
  - Updated `AgentRepositoryInterface.CreateAgent` signature

### 4. Usecase Layer
- **New:** `internal/usecase/ai_model.go`
  - Full CRUD operations for AI models
  - Validation logic

- **Updated:** `internal/usecase/agent.go`
  - `CreateNewAgent`: Updated signature to use `aiModelId`

- **Updated:** `internal/usecase/chat.go`
  - Changed references from `config.Agent.AiModel` to `config.Agent.AiModel.Name`

### 5. Handler Layer
- **New:** `internal/handler/ai_model.go`
  - `CreateAiModel`: SuperAdmin only
  - `GetAiModel`: Get single model
  - `ListAiModels`: List all or filter by provider
  - `UpdateAiModel`: SuperAdmin only
  - `DeleteAiModel`: SuperAdmin only

- **Updated:** `internal/handler/agent.go`
  - Updated to use `req.AiModelId` instead of separate fields
  - Response includes full `AiModel` object

### 6. Provider Layer
- **Updated:** `internal/provider/ai-provider/factory.go`
  - `GetProviderFromAgent`: Now reads from `agent.AiModel.Provider`
  - Builds capabilities from `agent.AiModel` fields instead of hardcoded lookup

### 7. Database Migrations
- **Migration:** `20250201_create_ai_models_table.sql`
  - Creates `ai_models` table
  - Migrates existing agent data to ai_models
  - Updates agents table with foreign key
  - Drops old columns

- **Rollback:** `20250201_create_ai_models_table_rollback.sql`
  - Restores original schema if needed

- **Seed Data:** `20250201_seed_ai_models.sql`
  - Pre-populates common models:
    - OpenAI: gpt-4o, gpt-4o-mini, gpt-4-turbo, gpt-4, gpt-3.5-turbo, o1, o1-mini
    - Anthropic: claude-3-5-sonnet, claude-3-5-haiku, claude-3-opus
    - Google: gemini-2.0-flash-exp, gemini-1.5-pro, gemini-1.5-flash
    - Groq: llama-3.3-70b-versatile, llama-3.1-8b-instant

## Benefits

1. **Centralized Model Management**: All AI models defined in one place
2. **Distinct Model Identification**: Each agent references a specific model by ID
3. **Capability Tracking**: Know which models support text/vision/voice/reasoning
4. **Pricing Management**: Update pricing centrally without touching agents
5. **Provider Flexibility**: Support multiple providers with same interface
6. **Scalability**: Easy to add new models and providers

## API Changes

### Agent Creation (Before)
```json
{
  "name": "My Agent",
  "ai_model": "gpt-4o",
  "ai_provider": "openai",
  "credits_per_1k": 10
}
```

### Agent Creation (After)
```json
{
  "name": "My Agent",
  "ai_model_id": "uuid-of-gpt-4o-model"
}
```

### Agent Response (After)
```json
{
  "id": "agent-uuid",
  "name": "My Agent",
  "ai_model_id": "model-uuid",
  "ai_model": {
    "id": "model-uuid",
    "name": "gpt-4o",
    "provider": "openai",
    "credits_per_1k": 10,
    "supports_text": true,
    "supports_vision": true,
    "supports_voice": true,
    "is_reasoning": false
  }
}
```

## New Endpoints

- `POST /api/ai-models` - Create AI model (SuperAdmin)
- `GET /api/ai-models` - List all models
- `GET /api/ai-models?provider=openai` - Filter by provider
- `GET /api/ai-models/:modelId` - Get single model
- `PUT /api/ai-models/:modelId` - Update model (SuperAdmin)
- `DELETE /api/ai-models/:modelId` - Delete model (SuperAdmin)

## Migration Steps

1. Run migration: `20250201_create_ai_models_table.sql`
2. Run seed data: `20250201_seed_ai_models.sql`
3. Update application code (already done)
4. Test agent creation with new `ai_model_id` field
5. Verify existing agents load correctly with preloaded `AiModel`

## Rollback Plan

If issues arise, run: `20250201_create_ai_models_table_rollback.sql`

This will restore the original schema with inline model fields.
