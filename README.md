# Agent Builder

Agent Builder is a Go CLI for scaffolding Google ADK Python agent projects and reusable sub-agent packages.

## Installation

### Homebrew

```bash
brew install doji-co/tap/agent-builder
```

### Manual

Download the latest release from the GitHub releases page.

## Usage

### Create a Project

```bash
agent-builder create
```

The starter-project flow generates a current ADK package layout:

```text
your-project/
├── requirements.txt
├── README.md
└── your_project/
    ├── __init__.py
    ├── .env.example
    ├── agent.py
    └── sub_agents/
        ├── __init__.py
        ├── researcher/
        │   ├── __init__.py
        │   └── agent.py
        └── formatter/
            ├── __init__.py
            └── agent.py
```

Generated projects target stable ADK `1.32.x`, pin `google-adk==1.32.0`, default to `gemini-2.5-flash`, and offer current Gemini 3 / 2.5 model options during scaffolding.

### Run a Generated Project

```bash
cd your-project
python3 -m venv .venv
source .venv/bin/activate
pip install -r requirements.txt
cp your_project/.env.example your_project/.env
adk run your_project
adk web --port 8000
```

### Create a Single Agent Package

The single-agent flow generates a reusable Python package with `__init__.py` and `agent.py` so it can be dropped into an existing ADK project's `sub_agents/` tree. It supports both `LlmAgent` scaffolds and `BaseAgent` custom-agent scaffolds.

### Version

```bash
agent-builder version
```

### Patterns

```bash
agent-builder patterns
```

## Development

### Prerequisites

- Go 1.26+

### Local Checks

```bash
go test ./...
go test -cover ./...
go vet ./...
go build ./...
```

## Releasing

```bash
git tag v0.1.0
git push origin v0.1.0
```

This publishes release binaries and updates the Homebrew tap.

## About ADK

The generated Python output follows the current ADK package conventions centered on `agent.py`, `__init__.py`, `root_agent`, `adk run`, and `adk web`.

## License

MIT
