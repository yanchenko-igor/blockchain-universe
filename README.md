# Blockchain Universe Agent MVP

A production-ready MVP implementation of an autonomous agent operating within the Blockchain Universe (BU) concept, where reality is defined solely by blockchain events.

## Overview

This agent:
- Exists exclusively within a blockchain-based reality
- Makes autonomous decisions using LLM reasoning
- Creates cryptographically signed events
- Maintains causal relationships between events
- Operates without reference to external physical reality

## Features

- ✅ **Production-ready architecture** with proper error handling
- ✅ **Structured logging** with configurable log levels
- ✅ **Configuration management** via YAML
- ✅ **Cryptographic signatures** using Ed25519
- ✅ **LLM integration** with OpenAI-compatible APIs (Ollama, etc.)
- ✅ **Event chain validation** and verification
- ✅ **Graceful shutdown** handling
- ✅ **Docker support** for containerized deployment
- ✅ **Comprehensive testing** setup

## Project Structure

```
blockchain-universe/
├── cmd/
│   └── agent/
│       └── main.go              # Application entry point
├── internal/
│   ├── agent/
│   │   └── agent.go             # Agent logic and decision-making
│   ├── blockchain/
│   │   └── blockchain.go        # Event management and verification
│   ├── config/
│   │   └── config.go            # Configuration handling
│   └── llm/
│       └── client.go            # LLM API client
├── pkg/
│   └── logger/
│       └── logger.go            # Structured logging
├── config.yaml                  # Configuration file
├── go.mod                       # Go module definition
├── Dockerfile                   # Container definition
├── Makefile                     # Build automation
└── README.md                    # This file
```

## Prerequisites

- Go 1.21 or higher
- Ollama (or OpenAI-compatible LLM API)
- Make (optional, for build automation)
- Docker (optional, for containerized deployment)

## Installation

1. **Clone the repository:**
```bash
git clone https://github.com/yanchenko.igor/blockchain-universe.git
cd blockchain-universe
```

2. **Install dependencies:**
```bash
make install-deps
```

3. **Set up Ollama (if using local LLM):**
```bash
# Install Ollama
curl -fsSL https://ollama.com/install.sh | sh

# Pull a model
ollama pull llama3.2
```

4. **Configure the application:**
```bash
cp config.yaml config.local.yaml
# Edit config.local.yaml with your settings
```

## Configuration

Edit `config.yaml` to customize the agent:

```yaml
agent:
  decision_interval: 30s      # How often to make decisions
  max_event_chain: 100        # Max depth for event chains

llm:
  api_endpoint: "http://localhost:11434/v1/completions"
  api_key: ""                 # Leave empty for local Ollama
  model: "llama3.2"
  max_tokens: 150
  temperature: 0.7
  timeout_seconds: 30
```

## Usage

### Running Locally

**Build and run:**
```bash
make run
```

**Run with debug logging:**
```bash
make run-debug
```

**Run directly:**
```bash
go run cmd/agent/main.go -config config.yaml -log-level info
```

### Running with Docker

**Build Docker image:**
```bash
make docker-build
```

**Run container:**
```bash
docker run --rm \
  -v $(pwd)/config.yaml:/app/config.yaml \
  --network host \
  bu-agent:latest
```

### Command-line Options

- `-config`: Path to configuration file (default: `config.yaml`)
- `-log-level`: Log level - debug, info, warn, error (default: `info`)

## Development

### Building

```bash
make build
```

### Testing

```bash
# Run tests
make test

# Run tests with coverage
make test-coverage
```

### Code Quality

```bash
# Format code
make fmt

# Lint code
make lint
```

## Architecture

### Event Structure

Each event in the Blockchain Universe contains:
- **Data**: Type, description, payload, and timestamp
- **Parents**: References to previous events (causal chain)
- **Signature**: Ed25519 cryptographic signature
- **Author**: Public key of the creating agent

### Decision Flow

1. Agent reads recent blockchain events
2. Constructs context prompt for LLM
3. LLM suggests next action based on BU principles
4. Agent creates and signs new event
5. Event is validated and added to blockchain
6. Process repeats at configured interval

### LLM System Prompt

The agent uses a specialized system prompt that enforces the Blockchain Universe worldview:
- No external physical reality
- Only blockchain events exist
- Time, space, and matter are redefined in blockchain terms
- All reasoning must be based on event chains

## API Compatibility

The LLM client is compatible with:
- **Ollama** (local, recommended for development)
- **OpenAI API**
- **Anthropic Claude** (via compatible proxy)
- Any OpenAI-compatible completion endpoint

## Extending the MVP

### Adding HTTP API

To expose agent status via HTTP:

```go
// In cmd/agent/main.go
import "net/http"

http.HandleFunc("/stats", func(w http.ResponseWriter, r *http.Request) {
    stats := agentInstance.GetStats()
    json.NewEncoder(w).Encode(stats)
})
go http.ListenAndServe(":8080", nil)
```

### Adding Persistence

To persist events to disk:

```go
// In internal/blockchain/blockchain.go
func (bc *Blockchain) SaveToFile(path string) error {
    // Serialize events to JSON
    // Write to file
}

func (bc *Blockchain) LoadFromFile(path string) error {
    // Read from file
    // Deserialize and validate events
}
```

### Multi-Agent Communication

To enable agent-to-agent communication:

1. Add network layer (gRPC or HTTP)
2. Implement event synchronization protocol
3. Add conflict resolution for concurrent events
4. Implement consensus mechanism

## Troubleshooting

### LLM Connection Issues

```bash
# Check Ollama is running
curl http://localhost:11434/api/version

# Check model is available
ollama list
```

### Event Verification Failures

- Ensure system time is synchronized
- Check event parent hashes exist
- Verify signature generation is consistent

### High Memory Usage

- Reduce `max_event_chain` in config
- Implement event pruning
- Add persistence and load events on-demand

## Security Considerations

- **Private keys** are generated in-memory and not persisted (implement secure storage for production)
- **Event signatures** ensure authenticity and integrity
- **No external input validation** in MVP (add for production)
- **Rate limiting** not implemented (add for public deployment)

## License

MIT License - see LICENSE file for details

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes with tests
4. Submit a pull request

## Roadmap

- [ ] Event persistence layer
- [ ] Multi-agent networking
- [ ] Web dashboard for visualization
- [ ] Event pruning and archival
- [ ] Consensus mechanisms
- [ ] Smart contract-like event rules
- [ ] Performance optimizations
- [ ] Comprehensive benchmarks

## Support

For issues and questions:
- Open an issue on GitHub
- Check existing documentation
- Review Ollama documentation for LLM setup

---

**Note**: This is an MVP implementation focused on demonstrating the Blockchain Universe concept. For production deployment, additional security hardening, monitoring, and operational features should be added.