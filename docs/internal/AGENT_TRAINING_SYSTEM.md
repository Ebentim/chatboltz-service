# Agent Training System

## Overview

The agent training system enables training of three types of agents (text, voice, multimodal) with various data sources using a hierarchical provider system: **OpenAI → Google → Cohere**.

## Agent Types

### 1. Text Agents (`AgentType.TextOnly`)
- Optimized for text-based conversations
- Supports text documents, PDFs, web content
- Uses standard text chunking (500 characters)

### 2. Voice Agents (`AgentType.VoiceOnly`) 
- Specialized for voice interactions
- Supports audio transcription via OpenAI Whisper
- Processes video files by extracting audio tracks
- Uses larger chunks (600 characters) for speech patterns

### 3. Multimodal Agents (`AgentType.Multimodal`)
- Handles text, voice, and visual content
- Supports images via OCR and vision models
- Processes all media types
- Most comprehensive training capabilities

## Provider Hierarchy

### Primary: OpenAI
- **Whisper**: Audio/video transcription
- **GPT-4 Vision**: Image analysis and OCR
- **Text Processing**: Document analysis

### Secondary: Google AI
- **Speech-to-Text**: Audio transcription fallback
- **Vision API**: Image OCR and analysis
- **Document AI**: PDF processing

### Tertiary: Cohere
- **Embeddings**: Vector generation (1024 dimensions)
- **Text Analysis**: Content processing

## Training Data Sources

### 1. Text Content
```http
POST /api/v1/agent/{agentId}/train/text
Content-Type: application/json

{
  "title": "Product Documentation",
  "content": "Your text content here..."
}
```

### 2. File Upload
```http
POST /api/v1/agent/{agentId}/train/file
Content-Type: multipart/form-data

file: [uploaded file]
title: "Optional custom title"
```

**Supported File Types:**
- **Documents**: PDF, TXT
- **Images**: JPEG, PNG, GIF, BMP, TIFF, WebP
- **Audio**: MP3, WAV, FLAC, M4A, OGG, AAC
- **Video**: MP4, AVI, MOV, MKV, WebM, FLV

### 3. Web Scraping
```http
POST /api/v1/agent/{agentId}/train/url
Content-Type: application/json

{
  "url": "https://example.com/docs",
  "title": "Website Documentation",
  "trace": true,
  "max_pages": 10
}
```

**Scraping Features:**
- Single page or multi-page crawling
- Same-domain link following
- Content extraction (headings, paragraphs)
- Automatic text cleaning

## Training Workflow

### 1. Document Processing
```
Input → MIME Detection → Media Processing → Text Extraction → Chunking → Embedding → Storage
```

### 2. Media Processing Pipeline
```
Audio/Video → Whisper Transcription → Text Chunks
Images → GPT-4 Vision/OCR → Text Description → Chunks  
PDFs → Text Extraction → Structured Chunks
Web Pages → Scraping → Content Extraction → Chunks
```

### 3. Vector Storage
- **Chunking**: Content split by type-specific rules
- **Embeddings**: Generated via Cohere (1024-dim vectors)
- **Storage**: PostgreSQL with pgvector extension
- **Indexing**: Optimized for similarity search

## API Endpoints

### Training Operations
- `POST /agent/{agentId}/train/text` - Train with text content
- `POST /agent/{agentId}/train/file` - Train with file upload
- `POST /agent/{agentId}/train/url` - Train with web content
- `GET /agent/{agentId}/training/documents` - List training documents
- `GET /agent/{agentId}/training/stats` - Training statistics
- `DELETE /agent/{agentId}/training` - Clear all training data

### Knowledge Base Queries
- `POST /agent/{agentId}/training/query` - RAG query for testing

### Legacy Migration
- `POST /agent/{agentId}/training/migrate` - Migrate old training data

## Configuration

### Required Environment Variables
```env
# Core APIs
OPENAI_API_KEY=sk-...          # Primary media processor
GOOGLE_API_KEY=...             # Secondary media processor  
COHERE_API_KEY=...             # Embeddings and fallback

# Database
DATABASE_URL=postgresql://...   # With pgvector extension
VECTOR_DB_TYPE=pgvector        # or pinecone

# Optional
PINECONE_API_KEY=...           # If using Pinecone
PINECONE_INDEX_NAME=...        # Pinecone index name
```

### Database Setup
```sql
-- Enable pgvector extension
CREATE EXTENSION IF NOT EXISTS vector;

-- Run migration
-- See: internal/repository/migrations/20241202_create_training_system.sql
```

## Usage Examples

### Training a Text Agent
```bash
# Upload documentation
curl -X POST "http://localhost:8080/api/v1/agent/agent-123/train/text" \
  -H "Authorization: Bearer $JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "API Documentation", 
    "content": "Our API supports REST endpoints..."
  }'
```

### Training with Audio
```bash
# Upload audio file
curl -X POST "http://localhost:8080/api/v1/agent/agent-123/train/file" \
  -H "Authorization: Bearer $JWT_TOKEN" \
  -F "file=@meeting-recording.mp3" \
  -F "title=Team Meeting Notes"
```

### Training with Website
```bash
# Scrape website content
curl -X POST "http://localhost:8080/api/v1/agent/agent-123/train/url" \
  -H "Authorization: Bearer $JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "url": "https://docs.example.com",
    "trace": true,
    "max_pages": 20
  }'
```

### Query Knowledge Base
```bash
# Test RAG retrieval
curl -X POST "http://localhost:8080/api/v1/agent/agent-123/training/query" \
  -H "Authorization: Bearer $JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "query": "How do I authenticate with the API?",
    "top_k": 5,
    "threshold": 0.7
  }'
```

## Performance Optimizations

### Chunking Strategy
- **Text**: 500 characters with 50 character overlap
- **Audio/Video**: 600 characters (speech patterns)
- **PDF**: 800 characters (document structure)
- **Images**: Single chunk per image
- **FAQ**: Q&A pair chunking

### Provider Fallbacks
1. **OpenAI fails** → Try Google AI
2. **Google AI fails** → Use Cohere (text only)
3. **All fail** → Return error with details

### Caching & Performance
- Document processing results cached
- Embedding generation batched
- Vector similarity search optimized
- Concurrent processing for multiple files

## Monitoring & Analytics

### Training Statistics
```json
{
  "total_documents": 45,
  "document_types": {
    "text": 20,
    "pdf": 10,
    "audio": 8,
    "image": 5,
    "video": 2
  },
  "total_chunks": 1250,
  "last_trained": "2024-12-02T10:30:00Z"
}
```

### Error Handling
- Provider fallback logging
- Failed document tracking
- Processing time metrics
- Storage usage monitoring

## Security & Privacy

### Data Protection
- All training data encrypted at rest
- API key rotation supported
- User data isolation by agent ID
- Automatic cleanup options

### Access Control
- JWT authentication required
- Agent ownership validation
- Rate limiting on training endpoints
- File size and type restrictions

## Migration from Legacy System

The system supports migrating from the old `TrainingData` JSONB format:

```bash
curl -X POST "http://localhost:8080/api/v1/agent/agent-123/training/migrate" \
  -H "Authorization: Bearer $JWT_TOKEN"
```

This converts legacy training texts to the new RAG system with proper chunking and embeddings.