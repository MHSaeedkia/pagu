package session

import (
	"maps"
	"sync"
	"time"
)

type State struct {
	State      string
	LastUpdate time.Time
}
type Session struct {
	mu     sync.Mutex
	States map[string]State
}

type SessionInterface interface {
	SetState(key, state string)
	GetState(key string) (State, bool)
	DeleteState(key string)
	GetAllStates() map[string]State
}

// NewSession creates and initializes a new session
func NewSession() SessionInterface {
	return &Session{
		States: make(map[string]State),
	}
}

// SetState updates or adds a new state for a given key
func (s *Session) SetState(key, state string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.States[key] = State{
		State:      state,
		LastUpdate: time.Now(),
	}
}

// GetState retrieves the state for a given key
func (s *Session) GetState(key string) (State, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	state, exists := s.States[key]
	return state, exists
}

// DeleteState removes a state from the session
func (s *Session) DeleteState(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.States, key)
}

// GetAllStates returns a copy of all states
func (s *Session) GetAllStates() map[string]State {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Return a copy to prevent external modification
	copyStates := make(map[string]State, len(s.States))
	maps.Copy(copyStates, s.States)
	return copyStates
}
