# anc (add-n-commit)

A beautiful, interactive CLI tool for generating AI-powered git commit messages using OpenAI's API.

![Go Version](https://img.shields.io/badge/Go-1.24.4-blue)
![License](https://img.shields.io/badge/license-MIT-green)

## Features

- 🎨 **Beautiful TUI**: Interactive terminal UI built with Bubble Tea
- 📝 **Multiple Generation Modes**:
  - All-in-one summary
  - File-by-file summaries
  - Custom prompt with additional context
- 🔐 **Secure Key Management**: Store and manage your OpenAI API key safely
- 📁 **Granular File Selection**: Choose exactly which files to stage and commit
- ✏️ **Message Editing**: Review and edit generated messages before committing
- ⚙️ **Configurable**: Customize prompts, model, temperature, and more

## Installation

### From Source

```bash
git clone https://github.com/oconnorjohnson/add-n-commit.git
cd add-n-commit
go build -o anc
sudo mv anc /usr/local/bin/  # Or add to your PATH
```

### Prerequisites

- Go 1.24.4 or higher
- Git
- OpenAI API key

## Quick Start

1. **Set your OpenAI API key**:

   ```bash
   anc --set-key sk-your-api-key-here
   ```

2. **Make some changes to your code**

3. **Run anc**:

   ```bash
   anc
   ```

4. **Follow the interactive prompts** to:
   - Select files to stage
   - Choose commit message generation mode
   - Review and edit the generated message
   - Commit your changes

## Usage

### Interactive Mode (Default)

Simply run `anc` without any arguments to enter the interactive TUI:

```bash
anc
```

### Command Line Options

```bash
anc [options]

OPTIONS:
    --set-key <key>    Set the OpenAI API key
    --show-key         Show the current OpenAI API key (masked)
    --delete-key       Delete the stored OpenAI API key
    --config           Open interactive configuration editor
    --version          Show version information
    --help             Show this help message
```

### Configuration

Configuration is stored in `~/.config/anc/config.json`. You can also set the `OPENAI_API_KEY` environment variable.

Use `anc --config` to interactively edit all settings:

- OpenAI API key
- Model (default: o4-mini)
- Default mode (interactive/all/by-file)
- Temperature
- System prompts

## Key Bindings

### File Selection

- `Space`: Toggle file selection
- `a`: Toggle all files
- `Enter`: Continue to mode selection
- `q`: Quit

### Mode Selection

- `Enter`: Select mode
- `q`: Quit

### Message Review

- `Enter`: Commit with current message
- `e`: Edit message
- `r`: Regenerate message
- `q`: Quit

### Message Editing

- `Ctrl+S` or `Ctrl+D`: Save and commit
- `Esc`: Cancel editing

## Examples

### Basic Usage

```bash
# Enter interactive mode
anc
```

### API Key Management

```bash
# Set API key
anc --set-key sk-your-api-key-here

# View current key (masked)
anc --show-key

# Delete stored key
anc --delete-key
```

### Configuration

```bash
# Open configuration editor
anc --config
```

## Comparison with Original Script

This CLI app improves upon the original bash script by adding:

- **Granular file staging**: Select specific files instead of `git add -A`
- **Interactive TUI**: Beautiful interface for all operations
- **Configuration management**: Store settings and customize behavior
- **Better error handling**: Clear error messages and recovery options
- **Message editing**: Review and edit before committing
- **Secure key storage**: Manage API keys without editing scripts

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License - see LICENSE file for details
