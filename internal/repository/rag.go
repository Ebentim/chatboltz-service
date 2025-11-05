package repository

import (
	"fmt"

	"github.com/alpinesboltltd/boltz-ai/internal/entity"
	"gorm.io/gorm"
)

type RAGRepository struct {
	db *gorm.DB
}

func NewRAGRepository(db *gorm.DB) *RAGRepository {
	return &RAGRepository{db: db}
}

func (r *RAGRepository) CreateTrainingDocument(doc *entity.TrainingDocument) error {
	return r.db.Create(doc).Error
}

func (r *RAGRepository) GetTrainingDocumentsByAgentID(agentID string) ([]entity.TrainingDocument, error) {
	var docs []entity.TrainingDocument
	err := r.db.Where("agent_id = ?", agentID).Preload("Chunks").Find(&docs).Error
	return docs, err
}

func (r *RAGRepository) UpdateTrainingDocument(doc *entity.TrainingDocument) error {
	return r.db.Save(doc).Error
}

func (r *RAGRepository) DeleteTrainingDocument(docID string) error {
	return r.db.Delete(&entity.TrainingDocument{}, "id = ?", docID).Error
}

func (r *RAGRepository) StoreChunks(chunks []entity.DocumentChunk) error {
	return r.db.CreateInBatches(chunks, 100).Error
}

func (r *RAGRepository) SearchSimilar(agentID string, embedding []float32, topK int, threshold float32) ([]entity.RetrievedChunk, error) {
	var results []entity.RetrievedChunk

	query := `
		SELECT dc.content, dc.metadata, dc.document_id, td.document_type, 1 - (dc.embedding <=> ?) as score
		FROM document_chunks dc
		JOIN training_documents td ON dc.document_id = td.id
		WHERE dc.agent_id = ? AND td.is_active = true AND 1 - (dc.embedding <=> ?) > ?
		ORDER BY dc.embedding <=> ? 
		LIMIT ?
	`

	rows, err := r.db.Raw(query, embedding, agentID, embedding, threshold, embedding, topK).Rows()
	if err != nil {
		return nil, fmt.Errorf("failed to execute similarity search: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var chunk entity.RetrievedChunk
		if err := rows.Scan(&chunk.Content, &chunk.Metadata, &chunk.DocumentID, &chunk.DocumentType, &chunk.Score); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		results = append(results, chunk)
	}

	return results, nil
}

func (r *RAGRepository) DeleteChunksByDocumentID(docID string) error {
	return r.db.Where("document_id = ?", docID).Delete(&entity.DocumentChunk{}).Error
}

func (r *RAGRepository) DeleteChunksByAgentID(agentID string) error {
	return r.db.Where("agent_id = ?", agentID).Delete(&entity.DocumentChunk{}).Error
}
