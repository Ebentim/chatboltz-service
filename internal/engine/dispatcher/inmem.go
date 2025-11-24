package dispatcher

import (
	"context"
	"sync"

	"github.com/alpinesboltltd/boltz-ai/internal/engine"
)

type InMemDispatcher struct {
	mu   sync.RWMutex
	subs map[string][]chan engine.OutboxEvent
}

func NewInMemDispatcher() *InMemDispatcher {
	return &InMemDispatcher{subs: make(map[string][]chan engine.OutboxEvent)}
}

func (d *InMemDispatcher) Dispatch(ctx context.Context, ev engine.OutboxEvent) error {
	d.mu.RLock()
	defer d.mu.RUnlock()
	if chs, ok := d.subs[ev.EventType]; ok {
		for _, ch := range chs {
			select {
			case ch <- ev:
			default:
			}
		}
	}
	return nil
}

func (d *InMemDispatcher) Subscribe(eventType string) (<-chan engine.OutboxEvent, error) {
	ch := make(chan engine.OutboxEvent, 100)
	d.mu.Lock()
	d.subs[eventType] = append(d.subs[eventType], ch)
	d.mu.Unlock()
	return ch, nil
}
