# ChatBoltz Service - Requirements Document

## 1. Overview
Backend service for Boltz AI, providing agent orchestration, data management, and integrations.

## 2. Core Modules

### 2.1. Workspaces
- **Data Model**: Support `workspaces` and `workspace_members`.
- **API**: Endpoints for CRUD operations on workspaces.
- **Auth**: Middleware to ensure users access only their workspaces.

### 2.2. Agents
- **Default Agents**: Seed database with:
  1. Virtual Assistant
  2. SDR
  3. BDR
  4. Customer Service
- **Custom Agents**: API to create agents with custom prompts.
- **Orchestration**: Level 4 autonomy loop (Perceive -> Plan -> Act -> Verify).

### 2.3. Training (RAG)
- **Ingestion**: Support PDF, URL, Text, QnA.
- **Vector Store**: Store embeddings for retrieval.
- **Retrieval**: Context-aware retrieval for agent responses.

### 2.4. Integrations
- **Platform Connectors**: Abstraction layer for external platforms.
- **Meeting Integration**: Support for joining video calls (LiveKit/Zoom/etc.).

## 3. MVP Requirements (Tonight)
- **Focus**: Virtual Assistant Agent.
- **Features**:
  - Workspace API.
  - Seed "Virtual Assistant" agent.
  - Basic RAG pipeline (Text/URL).
  - Integration stubs for future expansion.
