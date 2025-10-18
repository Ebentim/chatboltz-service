package usecase

import (
	"sync"
	"time"

	"github.com/alpinesboltltd/boltz-ai/internal/entity"
)

type AgentConfig struct {
	Agent             entity.Agent
	Behavior          entity.AgentBehavior
	SystemInstruction entity.SystemInstruction
	PromptTemplate    entity.PromptTemplate
	LoadedAt          time.Time
}

type AgentCache struct {
	cache map[string]*AgentConfig
	mutex sync.RWMutex
	ttl   time.Duration
	repo  AgentRepository
}

type AgentRepository interface {
	GetAgentByID(id string) (*entity.Agent, error)
	GetAgentBehavior(agentID string) (*entity.AgentBehavior, error)
	GetSystemInstruction(id string) (*entity.SystemInstruction, error)
	GetPromptTemplate(id string) (*entity.PromptTemplate, error)
}

func NewAgentCache(repo AgentRepository, ttl time.Duration) *AgentCache {
	cache := &AgentCache{
		cache: make(map[string]*AgentConfig),
		ttl:   ttl,
		repo:  repo,
	}
	go cache.cleanup()
	return cache
}

func (c *AgentCache) GetAgentConfig(agentID string) (*AgentConfig, error) {
	c.mutex.RLock()
	config, exists := c.cache[agentID]
	c.mutex.RUnlock()

	if exists && time.Since(config.LoadedAt) < c.ttl {
		return config, nil
	}

	return c.loadAgentConfig(agentID)
}

func (c *AgentCache) loadAgentConfig(agentID string) (*AgentConfig, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	agent, err := c.repo.GetAgentByID(agentID)
	if err != nil {
		return nil, err
	}

	behavior, err := c.repo.GetAgentBehavior(agentID)
	if err != nil {
		return nil, err
	}

	var sysInst entity.SystemInstruction
	if behavior.SystemInstructionId != nil && *behavior.SystemInstructionId != "" {
		inst, err := c.repo.GetSystemInstruction(*behavior.SystemInstructionId)
		if err == nil {
			sysInst = *inst
		}
	}

	var promptTmpl entity.PromptTemplate
	if behavior.PromptTemplateId != nil && *behavior.PromptTemplateId != "" {
		tmpl, err := c.repo.GetPromptTemplate(*behavior.PromptTemplateId)
		if err == nil {
			promptTmpl = *tmpl
		}
	}

	config := &AgentConfig{
		Agent:             *agent,
		Behavior:          *behavior,
		SystemInstruction: sysInst,
		PromptTemplate:    promptTmpl,
		LoadedAt:          time.Now(),
	}

	c.cache[agentID] = config
	return config, nil
}

func (c *AgentCache) cleanup() {
	ticker := time.NewTicker(c.ttl)
	defer ticker.Stop()

	for range ticker.C {
		c.mutex.Lock()
		for id, config := range c.cache {
			if time.Since(config.LoadedAt) > c.ttl {
				delete(c.cache, id)
			}
		}
		c.mutex.Unlock()
	}
}
