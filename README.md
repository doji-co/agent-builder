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

### Create Command

The `create` command offers two options:

```bash
agent-builder create
```

You'll be prompted to choose:

#### Option 1: Starter Project

Creates a complete multi-agent system from scratch. The CLI guides you through:

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

**Generated structure:**
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

**Running your project:**
```bash
cd your-project
pip install -r requirements.txt

# Run with command line
python main.py "Your prompt here"

# Or use ADK web interface
adk web
```

#### Option 2: Single Agent

Creates a single agent folder in the current directory. Perfect for adding new sub-agents to an existing project.

The CLI prompts for:
- Agent name
- Agent type (LLM or Tool)
- Instruction
- Output key
- Model selection

**Generated structure:**
```
./agent_name/
└── agent.py
```

**To use the agent in your project:**
```python
# In your orchestrator's agent.py
from agent_name.agent import agent as agent_name

agent = SequentialAgent(
    name="YourOrchestrator",
    sub_agents=[..., agent_name],  # Add the new agent
)
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
