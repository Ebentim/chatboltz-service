package dispatcher

import (
	"context"
	"log"
	"sync"
	"sync/atomic"
	"time"

	"github.com/alpinesboltltd/boltz-ai/internal/engine"
)

// DeliveryTimeout controls how long the dispatcher will wait for a subscriber
// channel to accept an event before considering the delivery dropped. It can
// be adjusted for testing or deployment tuning. Zero means immediate drop.
var DeliveryTimeout = 100 * time.Millisecond

// droppedDeliveries counts how many deliveries have been dropped due to timeouts.
var droppedDeliveries uint64

// DroppedDeliveries returns the number of dropped deliveries observed so far.
func DroppedDeliveries() uint64 { return atomic.LoadUint64(&droppedDeliveries) }

type InMemDispatcher struct {
	mu   sync.RWMutex
	subs map[string][]chan engine.OutboxEvent
	done chan struct{}
}

func NewInMemDispatcher() *InMemDispatcher {
	return &InMemDispatcher{
		subs: make(map[string][]chan engine.OutboxEvent),
		done: make(chan struct{}),
	}
}

// Unsubscribe removes a subscriber channel and closes it. The provided channel
// may be a receive-only channel; we compare by converting stored channels to
// the receive-only type for equality.
func (d *InMemDispatcher) Unsubscribe(eventType string, ch <-chan engine.OutboxEvent) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if chs, ok := d.subs[eventType]; ok {
		for i, c := range chs {
			if (<-chan engine.OutboxEvent)(c) == ch {
				// close the underlying channel and remove it from the slice
				close(c)
				d.subs[eventType] = append(chs[:i], chs[i+1:]...)
				break
			}
		}
	}
}

// Close stops the dispatcher and closes all subscriber channels. It is safe to
// call multiple times; subsequent calls are no-ops.
func (d *InMemDispatcher) Close() error {
	// Try to close the done channel only once
	select {
	case <-d.done:
		// already closed
	default:
		close(d.done)
	}

	d.mu.Lock()
	defer d.mu.Unlock()
	for _, chs := range d.subs {
		for _, ch := range chs {
			// close subscriber channels
			close(ch)
		}
	}
	d.subs = make(map[string][]chan engine.OutboxEvent)
	return nil
}

func (d *InMemDispatcher) Dispatch(ctx context.Context, ev engine.OutboxEvent) error {
	d.mu.RLock()
	defer d.mu.RUnlock()
	if chs, ok := d.subs[ev.EventType]; ok {
		for idx, ch := range chs {
			// Prefer context cancellation if provided
			if ctx != nil {
				select {
				case ch <- ev:
				case <-ctx.Done():
					// context cancelled before delivery; count as dropped
					atomic.AddUint64(&droppedDeliveries, 1)
					log.Printf("dispatcher: delivery cancelled for event=%s subscriber=%d: %v", ev.EventType, idx, ctx.Err())
				case <-time.After(DeliveryTimeout):
					atomic.AddUint64(&droppedDeliveries, 1)
					log.Printf("dispatcher: dropped event=%s for subscriber=%d after timeout=%s", ev.EventType, idx, DeliveryTimeout)
				}
			} else {
				// No context provided, use simple timeout
				select {
				case ch <- ev:
				case <-time.After(DeliveryTimeout):
					atomic.AddUint64(&droppedDeliveries, 1)
					log.Printf("dispatcher: dropped event=%s for subscriber=%d after timeout=%s", ev.EventType, idx, DeliveryTimeout)
				}
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
