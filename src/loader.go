package gorulesengine

import (
	"context"
	"sync"
	"time"
)

// RuleProvider defines the interface for fetching rules from external sources.
type RuleProvider interface {
	// FetchRules retrieves a list of rules from the source.
	FetchRules(ctx context.Context) ([]*Rule, error)
}

// HotReloader manages the periodic reloading of rules from a provider.
type HotReloader struct {
	engine   *Engine
	provider RuleProvider
	interval time.Duration
	stopChan chan struct{}
	wg       sync.WaitGroup
	mu       sync.Mutex
	running  bool
	onUpdate func([]*Rule)
	onError  func(error)
}

// NewHotReloader creates a new hot reloader for an engine.
func NewHotReloader(engine *Engine, provider RuleProvider, interval time.Duration) *HotReloader {
	return &HotReloader{
		engine:   engine,
		provider: provider,
		interval: interval,
		stopChan: make(chan struct{}),
	}
}

// OnUpdate sets a callback to be called whenever rules are updated.
func (h *HotReloader) OnUpdate(callback func([]*Rule)) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.onUpdate = callback
}

// OnError sets a callback to be called whenever an error occurs during reload.
func (h *HotReloader) OnError(callback func(error)) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.onError = callback
}

// Start begins the periodic reloading process.
func (h *HotReloader) Start(ctx context.Context) {
	h.mu.Lock()
	if h.running {
		h.mu.Unlock()
		return
	}
	h.running = true
	h.mu.Unlock()

	h.wg.Add(1)
	go func() {
		defer h.wg.Done()
		ticker := time.NewTicker(h.interval)
		defer ticker.Stop()

		// Initial fetch
		h.reload(ctx)

		for {
			select {
			case <-ticker.C:
				h.reload(ctx)
			case <-h.stopChan:
				return
			case <-ctx.Done():
				return
			}
		}
	}()
}

// Stop stops the hot reloader.
func (h *HotReloader) Stop() {
	h.mu.Lock()
	if !h.running {
		h.mu.Unlock()
		return
	}
	close(h.stopChan)
	h.running = false
	h.mu.Unlock()
	h.wg.Wait()
}

func (h *HotReloader) reload(ctx context.Context) {
	rules, err := h.provider.FetchRules(ctx)
	if err != nil {
		h.mu.Lock()
		if h.onError != nil {
			h.onError(err)
		}
		h.mu.Unlock()
		return
	}

	// If rules is nil, it means StatusNotModified or no rules returned
	if rules == nil {
		return
	}

	// Hot swap rules in engine
	// Note: We need a thread-safe way to swap rules in the engine.
	h.engine.SetRules(rules)

	h.mu.Lock()
	if h.onUpdate != nil {
		h.onUpdate(rules)
	}
	h.mu.Unlock()
}
