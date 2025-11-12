package blockchain

import (
	"crypto/ed25519"
	"crypto/sha3"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/yanchenko.igor/blockchain-universe/pkg/logger"
)

// Event represents a blockchain event
type Event struct {
	Data struct {
		Type        string            `json:"type"`
		Description string            `json:"description"`
		Payload     map[string]string `json:"payload"`
		Timestamp   string            `json:"timestamp"`
	} `json:"data"`
	Parents      []string `json:"parents"`
	Signature    string   `json:"signature"`
	AuthorPubKey string   `json:"author_pubkey"`
}

// AgentInfo stores information about known agents
type AgentInfo struct {
	PubKey        string
	LastEventHash string
	LastSeen      time.Time
}

// Blockchain manages events and agents
type Blockchain struct {
	events map[string]*Event
	agents map[string]*AgentInfo
	mu     sync.RWMutex
	log    logger.Logger
}

// New creates a new Blockchain instance
func New(log logger.Logger) *Blockchain {
	return &Blockchain{
		events: make(map[string]*Event),
		agents: make(map[string]*AgentInfo),
		log:    log,
	}
}

// CreateEvent creates and signs a new event
func (bc *Blockchain) CreateEvent(
	eventType, description string,
	payload map[string]string,
	parents []string,
	pub ed25519.PublicKey,
	priv ed25519.PrivateKey,
) (*Event, error) {
	event := &Event{}
	event.Data.Type = eventType
	event.Data.Description = description
	event.Data.Payload = payload
	event.Data.Timestamp = time.Now().UTC().Format(time.RFC3339)
	event.Parents = parents
	event.AuthorPubKey = hex.EncodeToString(pub)

	// Sign the event
	signature, err := bc.signEvent(event, priv)
	if err != nil {
		return nil, fmt.Errorf("failed to sign event: %w", err)
	}
	event.Signature = signature

	return event, nil
}

// AddEvent adds an event to the blockchain
func (bc *Blockchain) AddEvent(event *Event) error {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	// Verify event signature
	if err := bc.verifyEvent(event); err != nil {
		return fmt.Errorf("event verification failed: %w", err)
	}

	hash := bc.HashEvent(event)
	bc.events[hash] = event

	// Update agent info
	bc.agents[event.AuthorPubKey] = &AgentInfo{
		PubKey:        event.AuthorPubKey,
		LastEventHash: hash,
		LastSeen:      time.Now(),
	}

	bc.log.Debug("Event added", "hash", hash, "type", event.Data.Type)
	return nil
}

// GetEvent retrieves an event by hash
func (bc *Blockchain) GetEvent(hash string) (*Event, bool) {
	bc.mu.RLock()
	defer bc.mu.RUnlock()

	event, exists := bc.events[hash]
	return event, exists
}

// GetRecentEvents returns the N most recent events
func (bc *Blockchain) GetRecentEvents(limit int) []*Event {
	bc.mu.RLock()
	defer bc.mu.RUnlock()

	events := make([]*Event, 0, len(bc.events))
	for _, event := range bc.events {
		events = append(events, event)
	}

	// Sort by timestamp (simplified - in production use proper sorting)
	if len(events) > limit {
		events = events[len(events)-limit:]
	}

	return events
}

// GetAgents returns all known agents
func (bc *Blockchain) GetAgents() map[string]*AgentInfo {
	bc.mu.RLock()
	defer bc.mu.RUnlock()

	agents := make(map[string]*AgentInfo, len(bc.agents))
	for k, v := range bc.agents {
		agents[k] = v
	}
	return agents
}

// HashEvent computes the hash of an event
func (bc *Blockchain) HashEvent(event *Event) string {
	eventBytes, _ := json.Marshal(event.Data)
	hash := sha3.Sum512(eventBytes)
	return hex.EncodeToString(hash[:])
}

// signEvent signs an event with a private key
func (bc *Blockchain) signEvent(event *Event, priv ed25519.PrivateKey) (string, error) {
	eventHash := bc.HashEvent(event)
	signature := ed25519.Sign(priv, []byte(eventHash))
	return hex.EncodeToString(signature), nil
}

// verifyEvent verifies an event's signature
func (bc *Blockchain) verifyEvent(event *Event) error {
	pubKeyBytes, err := hex.DecodeString(event.AuthorPubKey)
	if err != nil {
		return fmt.Errorf("invalid public key: %w", err)
	}

	signatureBytes, err := hex.DecodeString(event.Signature)
	if err != nil {
		return fmt.Errorf("invalid signature: %w", err)
	}

	eventHash := bc.HashEvent(event)
	if !ed25519.Verify(pubKeyBytes, []byte(eventHash), signatureBytes) {
		return fmt.Errorf("signature verification failed")
	}

	return nil
}

// GetEventChain returns the chain of events leading to a specific event
func (bc *Blockchain) GetEventChain(hash string, maxDepth int) []*Event {
	bc.mu.RLock()
	defer bc.mu.RUnlock()

	chain := make([]*Event, 0)
	visited := make(map[string]bool)

	var traverse func(string, int)
	traverse = func(h string, depth int) {
		if depth >= maxDepth || visited[h] {
			return
		}

		event, exists := bc.events[h]
		if !exists {
			return
		}

		visited[h] = true
		chain = append(chain, event)

		for _, parent := range event.Parents {
			traverse(parent, depth+1)
		}
	}

	traverse(hash, 0)
	return chain
}