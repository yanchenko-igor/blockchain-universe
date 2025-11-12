package blockchain

import (
	"crypto/ed25519"
	"crypto/rand"
	"testing"

	"github.com/yourusername/blockchain-universe/pkg/logger"
)

func TestCreateEvent(t *testing.T) {
	log := logger.New("error")
	bc := New(log)
	pub, priv, _ := ed25519.GenerateKey(rand.Reader)

	event, err := bc.CreateEvent(
		"test_event",
		"Test event description",
		map[string]string{"key": "value"},
		[]string{},
		pub,
		priv,
	)

	if err != nil {
		t.Fatalf("Failed to create event: %v", err)
	}

	if event.Data.Type != "test_event" {
		t.Errorf("Expected type 'test_event', got '%s'", event.Data.Type)
	}

	if event.Data.Description != "Test event description" {
		t.Errorf("Expected description 'Test event description', got '%s'", event.Data.Description)
	}

	if event.Data.Payload["key"] != "value" {
		t.Errorf("Expected payload key 'value', got '%s'", event.Data.Payload["key"])
	}

	if event.Signature == "" {
		t.Error("Event signature should not be empty")
	}
}

func TestAddEvent(t *testing.T) {
	log := logger.New("error")
	bc := New(log)
	pub, priv, _ := ed25519.GenerateKey(rand.Reader)

	event, _ := bc.CreateEvent(
		"test_event",
		"Test event description",
		map[string]string{},
		[]string{},
		pub,
		priv,
	)

	err := bc.AddEvent(event)
	if err != nil {
		t.Fatalf("Failed to add event: %v", err)
	}

	hash := bc.HashEvent(event)
	retrieved, exists := bc.GetEvent(hash)

	if !exists {
		t.Error("Event should exist after being added")
	}

	if retrieved.Data.Type != event.Data.Type {
		t.Errorf("Retrieved event type mismatch")
	}
}

func TestEventVerification(t *testing.T) {
	log := logger.New("error")
	bc := New(log)
	pub, priv, _ := ed25519.GenerateKey(rand.Reader)

	event, _ := bc.CreateEvent(
		"test_event",
		"Test event description",
		map[string]string{},
		[]string{},
		pub,
		priv,
	)

	// Valid event should verify
	err := bc.verifyEvent(event)
	if err != nil {
		t.Errorf("Valid event failed verification: %v", err)
	}

	// Tampered event should fail verification
	event.Data.Description = "Tampered description"
	err = bc.verifyEvent(event)
	if err == nil {
		t.Error("Tampered event should fail verification")
	}
}

func TestGetRecentEvents(t *testing.T) {
	log := logger.New("error")
	bc := New(log)
	pub, priv, _ := ed25519.GenerateKey(rand.Reader)

	// Add multiple events
	for i := 0; i < 5; i++ {
		event, _ := bc.CreateEvent(
			"test_event",
			"Test event description",
			map[string]string{},
			[]string{},
			pub,
			priv,
		)
		bc.AddEvent(event)
	}

	recent := bc.GetRecentEvents(3)
	if len(recent) != 3 {
		t.Errorf("Expected 3 recent events, got %d", len(recent))
	}
}

func TestEventChain(t *testing.T) {
	log := logger.New("error")
	bc := New(log)
	pub, priv, _ := ed25519.GenerateKey(rand.Reader)

	// Create first event
	event1, _ := bc.CreateEvent(
		"event1",
		"First event",
		map[string]string{},
		[]string{},
		pub,
		priv,
	)
	bc.AddEvent(event1)
	hash1 := bc.HashEvent(event1)

	// Create second event with first as parent
	event2, _ := bc.CreateEvent(
		"event2",
		"Second event",
		map[string]string{},
		[]string{hash1},
		pub,
		priv,
	)
	bc.AddEvent(event2)
	hash2 := bc.HashEvent(event2)

	// Get event chain
	chain := bc.GetEventChain(hash2, 10)

	if len(chain) != 2 {
		t.Errorf("Expected chain length 2, got %d", len(chain))
	}
}

func TestAgentTracking(t *testing.T) {
	log := logger.New("error")
	bc := New(log)
	pub, priv, _ := ed25519.GenerateKey(rand.Reader)

	event, _ := bc.CreateEvent(
		"test_event",
		"Test event description",
		map[string]string{},
		[]string{},
		pub,
		priv,
	)
	bc.AddEvent(event)

	agents := bc.GetAgents()
	if len(agents) != 1 {
		t.Errorf("Expected 1 agent, got %d", len(agents))
	}

	pubKeyHex := event.AuthorPubKey
	agent, exists := agents[pubKeyHex]
	if !exists {
		t.Error("Agent should be tracked")
	}

	if agent.PubKey != pubKeyHex {
		t.Error("Agent public key mismatch")
	}
}

func BenchmarkCreateEvent(b *testing.B) {
	log := logger.New("error")
	bc := New(log)
	pub, priv, _ := ed25519.GenerateKey(rand.Reader)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bc.CreateEvent(
			"test_event",
			"Test event description",
			map[string]string{"key": "value"},
			[]string{},
			pub,
			priv,
		)
	}
}

func BenchmarkHashEvent(b *testing.B) {
	log := logger.New("error")
	bc := New(log)
	pub, priv, _ := ed25519.GenerateKey(rand.Reader)
	event, _ := bc.CreateEvent(
		"test_event",
		"Test event description",
		map[string]string{},
		[]string{},
		pub,
		priv,
	)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bc.HashEvent(event)
	}
}