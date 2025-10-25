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

```bash
agent-builder
```

Check version:
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

## Releasing

To create a new release:

1. Tag the commit:
   ```bash
   git tag v0.1.0
   ```

2. Push the tag:
   ```bash
   git push origin v0.1.0
   ```

This will automatically:
- Build binaries for multiple platforms
- Create a GitHub Release
- Update the Homebrew tap

### First Release Setup

Before your first release, you need to:

1. Create a personal access token with `repo` permissions
2. Add it as a repository secret named `HOMEBREW_TAP_GITHUB_TOKEN`
3. Create the homebrew tap repository: `github.com/doji-co/homebrew-tap`

## About ADK

This tool helps build [ADK (Agent Development Kit)](https://google.github.io/adk-docs/agents/multi-agents) multi-agents.

## License

MIT
