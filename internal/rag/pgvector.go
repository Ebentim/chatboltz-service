package rag

import (
	"fmt"

	"github.com/alpinesboltltd/boltz-ai/internal/entity"
	"gorm.io/gorm"
)

type PgVectorDB struct {
	db *gorm.DB
}

func NewPgVectorDB(db *gorm.DB) *PgVectorDB {
	return &PgVectorDB{db: db}
}

func (p *PgVectorDB) Store(chunk *entity.DocumentChunk) error {
	return p.db.Create(chunk).Error
}

func (p *PgVectorDB) Search(agentID string, embedding []float32, topK int, threshold float32) ([]entity.RetrievedChunk, error) {
	var results []entity.RetrievedChunk

	query := `
		SELECT content, metadata, 1 - (embedding <=> ?) as score
		FROM document_chunks 
		WHERE agent_id = ? AND 1 - (embedding <=> ?) > ?
		ORDER BY embedding <=> ? 
		LIMIT ?
	`

	rows, err := p.db.Raw(query, embedding, agentID, embedding, threshold, embedding, topK).Rows()
	if err != nil {
		return nil, fmt.Errorf("failed to execute similarity search: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var chunk entity.RetrievedChunk
		if err := rows.Scan(&chunk.Content, &chunk.Metadata, &chunk.Score); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		results = append(results, chunk)
	}

	return results, nil
}

func (p *PgVectorDB) Delete(agentID string) error {
	return p.db.Where("agent_id = ?", agentID).Delete(&entity.DocumentChunk{}).Error
}
