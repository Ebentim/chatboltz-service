package usecase

import (
	"sync"
	"time"

	"github.com/alpinesboltltd/boltz-ai/internal/entity"
	"github.com/alpinesboltltd/boltz-ai/internal/repository"
)

type AgentConfig struct {
	Agent             entity.Agent
	Behavior          entity.AgentBehavior
	SystemInstruction entity.SystemInstruction
	PromptTemplate    entity.PromptTemplate
	LoadedAt          time.Time
}

type AgentCache struct {
	cache      map[string]*AgentConfig
	mutex      sync.RWMutex
	ttl        time.Duration
	agentRepo  repository.AgentRepositoryInterface
	systemRepo repository.SystemRepositoryInterface
}

func NewAgentCache(agentRepo repository.AgentRepositoryInterface, systemRepo repository.SystemRepositoryInterface, ttl time.Duration) *AgentCache {
	cache := &AgentCache{
		cache:      make(map[string]*AgentConfig),
		ttl:        ttl,
		agentRepo:  agentRepo,
		systemRepo: systemRepo,
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

	agent, err := c.agentRepo.GetAgent(agentID)
	if err != nil {
		return nil, err
	}

	behavior, err := c.agentRepo.GetAgentBehavior(agentID)
	if err != nil {
		return nil, err
	}

	var sysInst entity.SystemInstruction
	if behavior.SystemInstructionId != nil && *behavior.SystemInstructionId != "" {
		inst, err := c.systemRepo.GetSystemInstruction(*behavior.SystemInstructionId)
		if err == nil {
			sysInst = *inst
		}
	}

	var promptTmpl entity.PromptTemplate
	if behavior.PromptTemplateId != nil && *behavior.PromptTemplateId != "" {
		tmpl, err := c.systemRepo.GetPromptTemplate(*behavior.PromptTemplateId)
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
