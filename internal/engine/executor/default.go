package executor

import (
	"context"
	"encoding/json"
	"log"

	"github.com/alpinesboltltd/boltz-ai/internal/engine"
	"github.com/alpinesboltltd/boltz-ai/internal/entity"
	"github.com/alpinesboltltd/boltz-ai/internal/provider/smtp"
	"github.com/alpinesboltltd/boltz-ai/internal/rag"
	"github.com/google/uuid"
)

type DefaultExecutor struct {
	llmFunc   func(ctx context.Context, input []byte) (string, error)
	smtpClient *smtp.Client
	store     engine.StateStore
	ragSvc    *rag.RAGService
}

// NewDefaultExecutor accepts optional dependencies. llmFunc can be nil for placeholder behavior.
func NewDefaultExecutor(llmFunc func(ctx context.Context, input []byte) (string, error), smtpClient *smtp.Client, store engine.StateStore, ragSvc *rag.RAGService) *DefaultExecutor {
	return &DefaultExecutor{llmFunc: llmFunc, smtpClient: smtpClient, store: store, ragSvc: ragSvc}
}

func (e *DefaultExecutor) RunStep(ctx context.Context, step *engine.WorkflowStepRecord) (engine.StepResult, error) {
	switch step.StepName {
	case "fetch_ticket":
		// For now, simply echo the input as the fetched ticket payload.
		// In a full implementation this would call a ticketing/CRM service.
		var payload interface{}
		if err := json.Unmarshal(step.Input, &payload); err != nil {
			log.Printf("executor: fetch_ticket invalid input: %v", err)
			return engine.StepResult{Success: false}, err
		}
		out, _ := json.Marshal(map[string]interface{}{"ticket": payload})
		return engine.StepResult{Success: true, Output: out}, nil

	case "retrieve_context":
		// Placeholder: in production this would run vector search / DB lookups.
		// If a RAG service is available, call it with the provided query.
		if e.ragSvc != nil {
			var in map[string]string
			_ = json.Unmarshal(step.Input, &in)
			query := in["query"]
			agentID := in["agent_id"]
			// fallback: if query empty, use raw input as string
			if query == "" {
				query = string(step.Input)
			}
			ragQuery := entity.RAGQuery{Query: query, AgentID: agentID, TopK: 5}
			resp, err := e.ragSvc.Query(ragQuery)
			if err != nil {
				log.Printf("executor: rag query error: %v", err)
				return engine.StepResult{Success: false}, err
			}
			out, _ := json.Marshal(resp)
			return engine.StepResult{Success: true, Output: out}, nil
		}
		// Return a simple context object when RAG is not configured.
		ctxObj := map[string]string{"context": "no_additional_context_available"}
		out, _ := json.Marshal(ctxObj)
		return engine.StepResult{Success: true, Output: out}, nil

	case "draft_response":
		// Use llmFunc if provided to generate draft text from step.Input
		if e.llmFunc != nil {
			resp, err := e.llmFunc(ctx, step.Input)
			if err != nil {
				log.Printf("executor: llm draft error: %v", err)
				return engine.StepResult{Success: false}, err
			}
			// store LLM result as JSON string
			out, _ := json.Marshal(map[string]string{"draft": resp})
			return engine.StepResult{Success: true, Output: out}, nil
		}
		// fallback placeholder
		out, _ := json.Marshal(map[string]string{"draft": "(llm disabled)"})
		return engine.StepResult{Success: true, Output: out}, nil

	case "send_response":
		// Expect input to contain to/subject/body/html; enqueue an outbox event
		var payload map[string]string
		if err := json.Unmarshal(step.Input, &payload); err != nil {
			log.Printf("executor: invalid send_response input: %v", err)
			return engine.StepResult{Success: false}, err
		}
		// build outbox payload
		p, _ := json.Marshal(payload)
		ev := &engine.OutboxEvent{
			ID:        uuid.NewString(),
			EventType: "email_send",
			Payload:   p,
			State:     "pending",
			Published: false,
		}
		if err := e.store.EnqueueEvent(ctx, ev); err != nil {
			log.Printf("executor: enqueue outbox error: %v", err)
			return engine.StepResult{Success: false}, err
		}
		return engine.StepResult{Success: true, Output: []byte(`{"enqueued":true}`)}, nil

	case "human_review":
		// Notify a human agent for review by enqueueing an email to the agent.
		var payload map[string]string
		if err := json.Unmarshal(step.Input, &payload); err != nil {
			log.Printf("executor: human_review invalid input: %v", err)
			return engine.StepResult{Success: false}, err
		}
		agentEmail := payload["agent_email"]
		draft := payload["draft"]
		ticketID := payload["ticket_id"]
		subject := "CSR Review Required"
		body := "A draft response requires your review.\n\nTicket: " + ticketID + "\n\nDraft:\n" + draft
		mail := map[string]string{"to": agentEmail, "subject": subject, "body": body}
		b, _ := json.Marshal(mail)
		ev := &engine.OutboxEvent{ID: uuid.NewString(), EventType: "email_send", Payload: b, State: "pending", Published: false}
		if err := e.store.EnqueueEvent(ctx, ev); err != nil {
			log.Printf("executor: enqueue human review email error: %v", err)
			return engine.StepResult{Success: false}, err
		}
		return engine.StepResult{Success: true, Output: []byte(`{"enqueued_review":true}`)}, nil

	default:
		// generic placeholder success for unknown steps
		return engine.StepResult{Success: true, Output: []byte(`{"status":"ok"}`)}, nil
	}
}
