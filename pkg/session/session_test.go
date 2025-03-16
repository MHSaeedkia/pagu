package session

import (
	"testing"
	"time"
)

func TestSession(t *testing.T) {
	// Create a new session
	sess := NewSession()

	// Test setting a state
	key := "user1"
	state := "active"
	sess.SetState(key, state)

	// Test getting the state
	st, exists := sess.GetState(key)
	if !exists {
		t.Errorf("expected state for key %s to exist", key)
	}
	if st.State != state {
		t.Errorf("expected state to be %s, got %s", state, st.State)
	}

	// Test getting all states
	allStates := sess.GetAllStates()
	if len(allStates) != 1 {
		t.Errorf("expected 1 state, got %d", len(allStates))
	}
	if allStates[key].State != state {
		t.Errorf("expected state %s, got %s", state, allStates[key].State)
	}

	// Test updating a state
	newState := "inactive"
	sess.SetState(key, newState)
	st, exists = sess.GetState(key)
	if !exists {
		t.Errorf("expected state for key %s to exist", key)
	}
	if st.State != newState {
		t.Errorf("expected state to be %s, got %s", newState, st.State)
	}

	// Test deleting a state
	sess.DeleteState(key)
	_, exists = sess.GetState(key)
	if exists {
		t.Errorf("expected state for key %s to be deleted", key)
	}

	// Test getting all states after deletion
	allStates = sess.GetAllStates()
	if len(allStates) != 0 {
		t.Errorf("expected no states, got %d", len(allStates))
	}
}

func TestSetState_Timestamp(t *testing.T) {
	// Create a new session
	sess := NewSession()

	// Set a state for a key
	key := "user1"
	state := "active"
	sess.SetState(key, state)

	// Retrieve the state and check if the LastUpdate field is set properly
	st, exists := sess.GetState(key)
	if !exists {
		t.Errorf("expected state for key %s to exist", key)
	}

	if st.LastUpdate.IsZero() {
		t.Errorf("expected LastUpdate to be set, but got zero value")
	}

	// Ensure LastUpdate is a recent timestamp
	if time.Since(st.LastUpdate) > time.Second {
		t.Errorf("LastUpdate is too old, expected a recent timestamp")
	}
}
