package agent

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/yanchenko.igor/blockchain-universe/internal/blockchain"
	"github.com/yanchenko.igor/blockchain-universe/internal/config"
	"github.com/yanchenko.igor/blockchain-universe/internal/llm"
	"github.com/yanchenko.igor/blockchain-universe/pkg/logger"
)

// Agent represents a blockchain universe agent
type Agent struct {
	pubKey     ed25519.PublicKey
	privKey    ed25519.PrivateKey
	blockchain *blockchain.Blockchain
	llmClient  *llm.Client
	config     config.AgentConfig
	log        logger.Logger
	lastEvent  string
}

// New creates a new agent instance
func New(
	cfg config.AgentConfig,
	bc *blockchain.Blockchain,
	llmClient *llm.Client,
	log logger.Logger,
) (*Agent, error) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("failed to generate key pair: %w", err)
	}

	return &Agent{
		pubKey:     pub,
		privKey:    priv,
		blockchain: bc,
		llmClient:  llmClient,
		config:     cfg,
		log:        log,
	}, nil
}

// PublicKeyHex returns the agent's public key as hex string
func (a *Agent) PublicKeyHex() string {
	return hex.EncodeToString(a.pubKey)
}

// Start begins the agent's operation
func (a *Agent) Start(ctx context.Context) error {
	a.log.Info("Agent started", "public_key", a.PublicKeyHex())
	<-ctx.Done()
	a.log.Info("Agent stopped")
	return nil
}

// CreateInitialEvent creates the first event for this agent
func (a *Agent) CreateInitialEvent(ctx context.Context) error {
	event, err := a.blockchain.CreateEvent(
		"initialization",
		"Agent initialization in Blockchain Universe",
		map[string]string{
			"agent_id": a.PublicKeyHex()[:16],
			"state":    "active",
			"version":  "1.0.0",
		},
		[]string{},
		a.pubKey,
		a.privKey,
	)
	if err != nil {
		return fmt.Errorf("failed to create initial event: %w", err)
	}

	if err := a.blockchain.AddEvent(event); err != nil {
		return fmt.Errorf("failed to add initial event: %w", err)
	}

	a.lastEvent = a.blockchain.HashEvent(event)
	a.log.Info("Initial event created", "hash", a.lastEvent)

	return nil
}

// MakeDecision uses LLM to decide on next action
func (a *Agent) MakeDecision(ctx context.Context) error {
	// Build context from blockchain state
	prompt := a.buildPrompt()

	a.log.Debug("Requesting LLM decision", "prompt_length", len(prompt))

	// Get decision from LLM
	decision, err := a.llmClient.GetCompletion(ctx, prompt)
	if err != nil {
		return fmt.Errorf("failed to get LLM decision: %w", err)
	}

	a.log.Info("LLM decision received", "decision", decision)

	// Create event based on decision
	if err := a.createDecisionEvent(ctx, decision); err != nil {
		return fmt.Errorf("failed to create decision event: %w", err)
	}

	return nil
}

// buildPrompt constructs a prompt for the LLM based on current blockchain state
func (a *Agent) buildPrompt() string {
	recentEvents := a.blockchain.GetRecentEvents(5)
	agents := a.blockchain.GetAgents()

	prompt := "Current Blockchain Universe state:\n\n"

	// Add recent events
	prompt += fmt.Sprintf("Recent events (%d):\n", len(recentEvents))
	for i, event := range recentEvents {
		prompt += fmt.Sprintf("%d. [%s] %s - %s\n",
			i+1,
			event.Data.Type,
			event.Data.Description,
			event.Data.Timestamp,
		)
	}

	// Add known agents
	prompt += fmt.Sprintf("\nKnown agents (%d):\n", len(agents))
	for pubKey, info := range agents {
		prompt += fmt.Sprintf("- Agent %s (last seen: %s)\n",
			pubKey[:16],
			info.LastSeen.Format(time.RFC3339),
		)
	}

	// Add my last event
	if a.lastEvent != "" {
		prompt += fmt.Sprintf("\nMy last event hash: %s\n", a.lastEvent)
	}

	prompt += "\nWhat should be the next event in the Blockchain Universe? " +
		"Provide a brief description (max 100 characters) for the event."

	return prompt
}

// createDecisionEvent creates an event based on LLM decision
func (a *Agent) createDecisionEvent(ctx context.Context, decision string) error {
	parents := []string{}
	if a.lastEvent != "" {
		parents = append(parents, a.lastEvent)
	}

	event, err := a.blockchain.CreateEvent(
		"state_change",
		decision,
		map[string]string{
			"agent_id": a.PublicKeyHex()[:16],
			"action":   "llm_decision",
		},
		parents,
		a.pubKey,
		a.privKey,
	)
	if err != nil {
		return fmt.Errorf("failed to create event: %w", err)
	}

	if err := a.blockchain.AddEvent(event); err != nil {
		return fmt.Errorf("failed to add event: %w", err)
	}

	a.lastEvent = a.blockchain.HashEvent(event)
	a.log.Info("Decision event created", "hash", a.lastEvent, "description", decision)

	return nil
}

// GetStats returns current agent statistics
func (a *Agent) GetStats() map[string]interface{} {
	return map[string]interface{}{
		"public_key":       a.PublicKeyHex(),
		"last_event_hash":  a.lastEvent,
		"total_events":     len(a.blockchain.GetRecentEvents(1000)),
		"known_agents":     len(a.blockchain.GetAgents()),
	}
}