# Acceptance Criteria - ChatBoltz Service

## 1. Workspaces
- [ ] Database schema includes `workspaces` and `workspace_members` tables.
- [ ] API endpoint `POST /workspaces` creates a new workspace.
- [ ] API endpoint `GET /workspaces` lists user's workspaces.
- [ ] Agents are associated with a `workspace_id`.

## 2. Agents
- [ ] Database is seeded with "Virtual Assistant", "SDR", "BDR", "Customer Service" templates.
- [ ] API endpoint `POST /agents` allows creating an agent from a template.
- [ ] Agent execution loop functions for "Virtual Assistant" tasks.

## 3. Training
- [ ] API endpoint `POST /agents/{id}/training` accepts PDF, URL, Text.
- [ ] System successfully extracts text from PDF and URL.
- [ ] System stores embeddings in vector store (or mock for MVP).

## 4. Integrations
- [ ] Backend supports "Website Widget" communication.
- [ ] Backend has structure to support "Zoom", "Teams" in the future.

## 5. Testing
- [ ] Unit tests for Workspace logic.
- [ ] Integration tests for Agent creation.
