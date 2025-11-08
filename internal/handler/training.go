package handler

import (
	"io"
	"net/http"

	"github.com/alpinesboltltd/boltz-ai/internal/entity"
	appErrors "github.com/alpinesboltltd/boltz-ai/internal/errors"
	"github.com/alpinesboltltd/boltz-ai/internal/usecase"
	"github.com/gin-gonic/gin"
)

// TrainingHandler handles HTTP requests for agent training operations
type TrainingHandler struct {
	trainingUsecase *usecase.TrainingUseCase
}

// NewTrainingHandler creates a new training handler
func NewTrainingHandler(trainingUsecase *usecase.TrainingUseCase) *TrainingHandler {
	return &TrainingHandler{
		trainingUsecase: trainingUsecase,
	}
}

// TrainWithText trains an agent with text content
func (h *TrainingHandler) TrainWithText(c *gin.Context) {
	agentID := c.Param("agentId")
	if agentID == "" {
		appErrors.HandleError(c, appErrors.NewValidationError("Agent ID is required"), "TrainWithText")
		return
	}

	var req struct {
		Title   string `json:"title" binding:"required"`
		Content string `json:"content" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		appErrors.HandleError(c, appErrors.NewValidationError("Invalid request format"), "TrainWithText")
		return
	}

	err := h.trainingUsecase.ProcessDocument(agentID, req.Title, req.Content, entity.DocumentTypeText, nil)
	if err != nil {
		appErrors.HandleError(c, err, "TrainWithText")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Text training completed successfully"})
}

// TrainWithFile trains an agent with uploaded file
func (h *TrainingHandler) TrainWithFile(c *gin.Context) {
	agentID := c.Param("agentId")
	if agentID == "" {
		appErrors.HandleError(c, appErrors.NewValidationError("Agent ID is required"), "TrainWithFile")
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		appErrors.HandleError(c, appErrors.NewValidationError("File is required"), "TrainWithFile")
		return
	}
	defer file.Close()

	title := c.PostForm("title")
	if title == "" {
		title = header.Filename
	}

	fileData, err := io.ReadAll(file)
	if err != nil {
		appErrors.HandleError(c, appErrors.NewInternalError("Failed to read file", err.Error()), "TrainWithFile")
		return
	}

	mimeType := header.Header.Get("Content-Type")
	err = h.trainingUsecase.ProcessFileWithMimeDetection(agentID, title, fileData, mimeType, nil)
	if err != nil {
		appErrors.HandleError(c, err, "TrainWithFile")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File training completed successfully"})
}

// GetTrainingDocuments retrieves all training documents for an agent
func (h *TrainingHandler) GetTrainingDocuments(c *gin.Context) {
	agentID := c.Param("agentId")
	if agentID == "" {
		appErrors.HandleError(c, appErrors.NewValidationError("Agent ID is required"), "GetTrainingDocuments")
		return
	}

	documents, err := h.trainingUsecase.GetAgentDocuments(agentID)
	if err != nil {
		appErrors.HandleError(c, err, "GetTrainingDocuments")
		return
	}

	c.JSON(http.StatusOK, gin.H{"documents": documents})
}

// DeleteTrainingData removes all training data for an agent
func (h *TrainingHandler) DeleteTrainingData(c *gin.Context) {
	agentID := c.Param("agentId")
	if agentID == "" {
		appErrors.HandleError(c, appErrors.NewValidationError("Agent ID is required"), "DeleteTrainingData")
		return
	}

	err := h.trainingUsecase.DeleteAgentTraining(agentID)
	if err != nil {
		appErrors.HandleError(c, err, "DeleteTrainingData")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Training data deleted successfully"})
}

// QueryKnowledgeBase performs RAG query on agent's knowledge base
func (h *TrainingHandler) QueryKnowledgeBase(c *gin.Context) {
	agentID := c.Param("agentId")
	if agentID == "" {
		appErrors.HandleError(c, appErrors.NewValidationError("Agent ID is required"), "QueryKnowledgeBase")
		return
	}

	var req struct {
		Query     string  `json:"query" binding:"required"`
		TopK      int     `json:"top_k,omitempty"`
		Threshold float32 `json:"threshold,omitempty"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		appErrors.HandleError(c, appErrors.NewValidationError("Invalid request format"), "QueryKnowledgeBase")
		return
	}

	ragQuery := entity.RAGQuery{
		Query:     req.Query,
		AgentID:   agentID,
		TopK:      req.TopK,
		Threshold: req.Threshold,
	}

	// Get RAG service from training usecase
	response, err := h.trainingUsecase.QueryKnowledgeBase(ragQuery)
	if err != nil {
		appErrors.HandleError(c, err, "QueryKnowledgeBase")
		return
	}

	c.JSON(http.StatusOK, response)
}

// MigrateLegacyTraining migrates legacy training data to new RAG system
func (h *TrainingHandler) MigrateLegacyTraining(c *gin.Context) {
	agentID := c.Param("agentId")
	if agentID == "" {
		appErrors.HandleError(c, appErrors.NewValidationError("Agent ID is required"), "MigrateLegacyTraining")
		return
	}

	err := h.trainingUsecase.TrainAgentFromLegacyData(agentID)
	if err != nil {
		appErrors.HandleError(c, err, "MigrateLegacyTraining")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Legacy training data migrated successfully"})
}

// TrainWithURL trains an agent with content from URL
func (h *TrainingHandler) TrainWithURL(c *gin.Context) {
	agentID := c.Param("agentId")
	if agentID == "" {
		appErrors.HandleError(c, appErrors.NewValidationError("Agent ID is required"), "TrainWithURL")
		return
	}

	var req struct {
		URL      string `json:"url" binding:"required"`
		Title    string `json:"title,omitempty"`
		Trace    bool   `json:"trace,omitempty"`
		MaxPages int    `json:"max_pages,omitempty"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		appErrors.HandleError(c, appErrors.NewValidationError("Invalid request format"), "TrainWithURL")
		return
	}

	title := req.Title
	if title == "" {
		title = "Content from " + req.URL
	}

	err := h.trainingUsecase.ProcessURL(agentID, req.URL, title, req.Trace, req.MaxPages)
	if err != nil {
		appErrors.HandleError(c, err, "TrainWithURL")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "URL training completed successfully"})
}

// GetTrainingStats returns training statistics for an agent
func (h *TrainingHandler) GetTrainingStats(c *gin.Context) {
	agentID := c.Param("agentId")
	if agentID == "" {
		appErrors.HandleError(c, appErrors.NewValidationError("Agent ID is required"), "GetTrainingStats")
		return
	}

	documents, err := h.trainingUsecase.GetAgentDocuments(agentID)
	if err != nil {
		appErrors.HandleError(c, err, "GetTrainingStats")
		return
	}

	stats := map[string]interface{}{
		"total_documents": len(documents),
		"document_types":  make(map[string]int),
		"total_chunks":    0,
	}

	docTypes := stats["document_types"].(map[string]int)
	totalChunks := 0

	for _, doc := range documents {
		docTypes[string(doc.DocumentType)]++
		totalChunks += len(doc.Chunks)
	}

	stats["total_chunks"] = totalChunks

	c.JSON(http.StatusOK, stats)
}
