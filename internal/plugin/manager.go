package plugin

import (
	"sync"
)

// HookFunc is a callback that plugins can register
type HookFunc func(args ...interface{}) (interface{}, error)

// Event names for common hook points
const (
	EventPostCreated   = "post.created"
	EventThreadCreated = "thread.created"
	EventUserRegistered = "user.registered"
	EventBeforePostSave = "post.before_save"
	EventAfterPostSave  = "post.after_save"
)

// Manager stores and dispatches hooks
type Manager struct {
	hooks map[string][]HookFunc
	mu    sync.RWMutex
}

var global = &Manager{hooks: make(map[string][]HookFunc)}

// Global returns the global hook manager
func Global() *Manager { return global }

// Register adds a hook for an event
func (m *Manager) Register(event string, fn HookFunc) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.hooks[event] = append(m.hooks[event], fn)
}

// Dispatch runs all hooks for an event
func (m *Manager) Dispatch(event string, args ...interface{}) []interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var results []interface{}
	for _, fn := range m.hooks[event] {
		r, _ := fn(args...)
		results = append(results, r)
	}
	return results
}

// Clear removes all hooks (useful for testing)
func (m *Manager) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.hooks = make(map[string][]HookFunc)
}
