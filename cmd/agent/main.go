package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/yanchenko.igor/blockchain-universe/internal/agent"
	"github.com/yanchenko.igor/blockchain-universe/internal/blockchain"
	"github.com/yanchenko.igor/blockchain-universe/internal/config"
	"github.com/yanchenko.igor/blockchain-universe/internal/llm"
	"github.com/yanchenko.igor/blockchain-universe/pkg/logger"
)

var (
	configPath = flag.String("config", "config.yaml", "Path to configuration file")
	logLevel   = flag.String("log-level", "info", "Log level (debug, info, warn, error)")
)

func main() {
	flag.Parse()

	// Initialize logger
	log := logger.New(*logLevel)
	log.Info("Starting Blockchain Universe Agent...")

	// Load configuration
	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatal("Failed to load configuration", "error", err)
	}

	// Create context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize blockchain
	bc := blockchain.New(log)

	// Initialize LLM client
	llmClient, err := llm.NewClient(cfg.LLM, log)
	if err != nil {
		log.Fatal("Failed to initialize LLM client", "error", err)
	}

	// Initialize agent
	agentInstance, err := agent.New(cfg.Agent, bc, llmClient, log)
	if err != nil {
		log.Fatal("Failed to initialize agent", "error", err)
	}

	log.Info("Agent initialized", "public_key", agentInstance.PublicKeyHex())

	// Start agent in background
	go func() {
		if err := agentInstance.Start(ctx); err != nil {
			log.Error("Agent error", "error", err)
			cancel()
		}
	}()

	// Create initial event
	if err := agentInstance.CreateInitialEvent(ctx); err != nil {
		log.Error("Failed to create initial event", "error", err)
	}

	// Run decision loop
	ticker := time.NewTicker(cfg.Agent.DecisionInterval)
	defer ticker.Stop()

	// Setup graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	log.Info("Agent running. Press Ctrl+C to stop.")

	for {
		select {
		case <-ctx.Done():
			log.Info("Context cancelled, shutting down...")
			return
		case <-sigChan:
			log.Info("Received shutdown signal")
			cancel()
			// Allow some time for graceful shutdown
			time.Sleep(2 * time.Second)
			return
		case <-ticker.C:
			if err := agentInstance.MakeDecision(ctx); err != nil {
				log.Error("Decision error", "error", err)
			}
		}
	}
}