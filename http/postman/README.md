# Helix Backend API - Postman Collections

This directory contains Postman collections organized by modules for the Helix Backend API.

## Collections

### Core Modules
- **auth.postman_collection.json** - Authentication endpoints (signup, login, verify)
- **agent.postman_collection.json** - Agent management (CRUD, appearance, behavior, channels, integrations)
- **system.postman_collection.json** - System management (instructions, templates)
- **chat.postman_collection.json** - Chat and WebSocket endpoints
- **scraper.postman_collection.json** - Web scraping endpoints
- **otp.postman_collection.json** - OTP generation and verification
- **training.postman_collection.json** - Training data and RAG endpoints

### Environment
- **helix-be.postman_environment.json** - Environment variables for all collections

## Setup Instructions

1. **Import Collections**: Import all `.postman_collection.json` files into Postman
2. **Import Environment**: Import `helix-be.postman_environment.json`
3. **Set Environment**: Select "Helix BE Environment" in Postman
4. **Configure Variables**:
   - Set `baseUrl` to your server URL (default: `http://localhost:8080`)
   - After authentication, set `token`, `userId`, `agentId` etc. from responses

## Usage Flow

### 1. Authentication
Start with the **Auth Module**:
1. Use "Signup with Email" or "Login with Email"
2. Copy the returned JWT token to the `token` environment variable
3. Copy the user ID to the `userId` environment variable

### 2. Agent Management
Use the **Agent Module**:
1. Create an agent with "Create Agent"
2. Copy the agent ID to the `agentId` environment variable
3. Configure agent appearance, behavior, channels, and integrations

### 3. System Configuration
Use the **System Module**:
1. Create system instructions and prompt templates
2. Copy IDs to respective environment variables

### 4. Training & Chat
Use **Training Module** for knowledge base and **Chat Module** for interactions.

## Environment Variables

| Variable | Description | Example |
|----------|-------------|---------|
| `baseUrl` | API base URL | `http://localhost:8080` |
| `token` | JWT authentication token | `eyJhbGciOiJIUzI1NiIs...` |
| `userId` | User ID from auth response | `b3b0b951-fa31-49d9-a862-45bb48a8358b` |
| `agentId` | Agent ID from create agent | `9e0c8d37-5520-4a16-903a-49dd0a30fa30` |
| `instructionId` | System instruction ID | `sys_inst_001` |
| `templateId` | Prompt template ID | `template_001` |

## Notes

- All authenticated endpoints require the `Authorization: Bearer {{token}}` header
- WebSocket endpoints in the Chat module require a WebSocket client
- File uploads in Training module use `multipart/form-data`
- Some endpoints are restricted to super admin users only