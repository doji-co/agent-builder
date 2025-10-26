# Agent Builder

ADK Multi-Agent Builder CLI - A tool to help build ADK (Agent Development Kit) agents.

## Installation

### Homebrew (macOS/Linux)

```bash
brew install doji-co/tap/agent-builder
```

### Manual Installation

Download the latest release from the [releases page](https://github.com/doji-co/agent-builder/releases).

## Usage

### Create a New Multi-Agent Project

Run the interactive CLI to create a new multi-agent system:

```bash
agent-builder create
```

The CLI will guide you through:
1. **Project name** - Name for your multi-agent project
2. **Orchestrator details** - The root agent that coordinates sub-agents
   - Name
   - Orchestration pattern (Sequential, Parallel, LLM-Coordinated, or Loop)
   - Description
   - Model selection (gemini-2.5-flash, gemini-2.5-pro, or gemini-2.5-flash-lite)
3. **Sub-agents** - Individual agents that perform specific tasks
   - Name
   - Type (LLM or Tool)
   - Instruction
   - Output key
   - Model

### Generated Project Structure

```
your-project/
├── orchestrator_name/
│   └── agent.py       # Orchestrator agent
├── sub_agent_1/
│   └── agent.py       # Sub-agent 1
├── sub_agent_2/
│   └── agent.py       # Sub-agent 2
├── main.py            # Example usage
├── requirements.txt   # Python dependencies
└── README.md          # Project documentation
```

### Running Your Agent

After creating your project:

```bash
cd your-project
pip install -r requirements.txt

# Run with command line
python main.py "Your prompt here"

# Or use ADK web interface
adk web
```

### Check Version

```bash
agent-builder version
```

## Development

### Prerequisites

- Go 1.23 or higher

### Building

```bash
go build -o agent-builder .
```

### Running

```bash
./agent-builder
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

### Releasing (for maintainers)

To create a new release, tag the commit and push:

```bash
git tag v0.1.0
git push origin v0.1.0
```

This will automatically:
- Build binaries for multiple platforms
- Create a GitHub Release
- Update the Homebrew tap

## About ADK

This tool helps build [ADK (Agent Development Kit)](https://google.github.io/adk-docs/agents/multi-agents) multi-agents.

## License

MIT
