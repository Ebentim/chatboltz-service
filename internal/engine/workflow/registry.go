package workflow

import (
	"sync"

	"github.com/alpinesboltltd/boltz-ai/internal/engine"
)

type Registry struct {
	mu    sync.RWMutex
	store map[string]engine.Workflow
}

func NewRegistry() *Registry {
	return &Registry{store: make(map[string]engine.Workflow)}
}

func (r *Registry) Register(w engine.Workflow) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.store[w.ID()] = w
}

func (r *Registry) Get(id string) (engine.Workflow, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	w, ok := r.store[id]
	return w, ok
}
