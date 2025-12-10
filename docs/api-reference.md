# ChatBoltz Service API Reference

Base URL: `http://localhost:8080` (or your deployed URL)

## Authentication

### Signup
`POST /api/v1/auth/signup`

Request:
```json
{
  "email": "user@example.com",
  "password": "password123",
  "name": "John Doe"
}
```

### Login
`POST /api/v1/auth/login`

Request:
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

### Verify Token
`POST /api/v1/auth/verify`

Request:
```json
{
  "id_token": "firebase_id_token"
}
```

## Agents

### Create Agent
`POST /api/v1/agent/create`
Headers: `Authorization: Bearer <token>`

Request:
```json
{
  "userId": "user-id",
  "name": "Agent Name",
  "description": "Description",
  "ai_model": "gpt-4",
  "ai_provider": "openai",
  "agent_type": 0,
  "credits_per_1k": 10,
  "status": "active"
}
```

### Update Agent
`PATCH /api/v1/agent/update/:agentId`
Headers: `Authorization: Bearer <token>`

### Get Agent
`GET /api/v1/agent/:agentId`

### Get User Agents
`GET /api/v1/agent/agents/:userId`

### Delete Agent
`DELETE /api/v1/agent/:agentId`

## Agent Configuration

### Appearance
- Create: `POST /api/v1/agent/create/appearance`
- Get: `GET /api/v1/agent/:agentId/appearance`
- Update: `PATCH /api/v1/agent/:agentId/appearance`
- Delete: `DELETE /api/v1/agent/:agentId/appearance`

### Behavior
- Create: `POST /api/v1/agent/create/behavior`
- Get: `GET /api/v1/agent/:agentId/behavior`
- Update: `PATCH /api/v1/agent/:agentId/behavior`
- Delete: `DELETE /api/v1/agent/:agentId/behavior`

### Channels
- Create: `POST /api/v1/agent/create/channel`
- Get: `GET /api/v1/agent/:agentId/channel`
- Update: `PATCH /api/v1/agent/:agentId/channel`
- Delete: `DELETE /api/v1/agent/:agentId/channel`

### Integrations
- Create: `POST /api/v1/agent/create/integration`
- Get: `GET /api/v1/agent/:agentId/integration`
- Update: `PATCH /api/v1/agent/:agentId/integration`
- Delete: `DELETE /api/v1/agent/:agentId/integration`

### Stats
- Get: `GET /api/v1/agent/:agentId/stats`
- Delete: `DELETE /api/v1/agent/:agentId/stats`

## System (Admin Only)

### Instructions
- Create: `POST /api/v1/system/instructions`
- Get: `GET /api/v1/system/instructions/:id`
- List: `GET /api/v1/system/instructions`
- Update: `PATCH /api/v1/system/instructions/:id`
- Delete: `DELETE /api/v1/system/instructions/:id`

### Templates
- Create: `POST /api/v1/system/templates`
- Get: `GET /api/v1/system/templates/:id`
- List: `GET /api/v1/system/templates`

## Scraper

### Scrape Page
`POST /api/v1/scrape`

Request:
```json
{
  "url": "https://example.com",
  "trace": true,
  "exclude": [],
  "max_pages": 1
}
```
