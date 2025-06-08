# GDMC - Go Daily Message Creator

🤖 A global CLI tool that generates professional messages from your git commits using Gemini AI.

## Features

- 📊 **Professional Reports**: Generate status reports from git commits
- 🗣️ **Meeting Transcripts**: Create standup meeting updates
- ⏰ **Multiple Intervals**: Daily, weekly, monthly analysis
- 🎯 **Customizable Templates**: Fully configurable message templates
- 🌍 **Global Tool**: Works in any git repository
- 🔧 **Configurable**: Save preferences and API keys

## Installation

### Prerequisites

- Go 1.18+ installed
- Git installed and accessible
- Terminal/shell access

### Step 1: Clone and Install

```bash
# Clone the repository
git clone https://github.com/fatihbahadir/go-daily-message-creator.git
cd go-daily-message-creator

# Install globally
go install .
```

### Step 2: Setup PATH (if needed)

If `gdmc` command is not found after installation, add Go's bin directory to your PATH:

```bash
# For zsh (default on macOS)
echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.zshrc
source ~/.zshrc

# For bash
echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.bashrc
source ~/.bashrc

# Check if it works
gdmc --help
```

### Step 3: Verify Installation

```bash
# Should show help message
gdmc --help

# Check version and available commands
gdmc generate --help
gdmc config --help
```

## Setup

### 1. Get your Gemini API Key

Get your API key from [Google AI Studio](https://makersuite.google.com/app/apikey)

### 2. Configure gdmc

```bash
# Set your default author email (use your git email)
gdmc config set author "your-email@example.com"

# Set your Gemini API key
gdmc config set api_key "your-gemini-api-key"

# Alternative: Set as environment variable
export GEMINI_API_KEY="your-gemini-api-key"

# Verify configuration
gdmc config show
```

## Usage

Navigate to any git repository and run gdmc:

```bash
# Basic usage (uses config defaults)
cd /path/to/your/project
gdmc generate

# Specific examples
gdmc generate --author="john@company.com" -t=report --i=weekly
gdmc generate --author="jane@team.com" --t=transcript --i=daily
gdmc generate --author="dev@startup.com" --t=summary --i=monthly

# Show help
gdmc --help
gdmc generate --help
```

## Message Types

### 📋 Report (Professional Status Report)
```bash
gdmc generate -t=report
```
Generates structured status reports with:
- Summary of accomplishments
- Key changes and improvements
- Technical details
- Project impact
- Next steps

### 🗣️ Transcript (Meeting Format)
```bash
gdmc generate -t=transcript
```
Creates standup meeting updates with:
- What was accomplished
- Current focus areas
- Next priorities
- Blockers and notes

### 📝 Summary (Concise Overview)
```bash
gdmc generate -t=summary
```
Provides brief summaries highlighting the most important changes.

## Time Intervals

- **daily**: Yesterday midnight to now
- **weekly**: Last 7 days
- **monthly**: Last 30 days

## Configuration Management

```bash
# View current configuration
gdmc config show

# Set configuration values
gdmc config set author "your@email.com"
gdmc config set default_type "report"
gdmc config set api_key "your-gemini-key"
```

Configuration is stored in `~/.config/gdmc/config.json`

## Real-World Examples

### Daily Standup Report
```bash
cd ~/projects/my-awesome-app
gdmc generate -t=transcript -i=daily
```

Output example:
```
**What I accomplished**: Implemented user authentication system
**Current focus**: Working on API rate limiting 
**Next priorities**: Database optimization and caching
**Blockers/Notes**: Waiting for security review on auth module
```

### Weekly Team Update
```bash
cd ~/work/backend-service
gdmc generate --author="team-lead@company.com" -t=report -i=weekly
```

### Monthly Summary for Multiple Contributors
```bash
cd ~/open-source/project
gdmc generate --author="contributor@domain.com" -t=summary -i=monthly
```

## Template Customization

Templates are fully customizable through the config file. Each template supports variables:
- `{{.Commits}}`: List of git commits
- `{{.Interval}}`: Time period name
- `{{.Author}}`: Author email

## Git Settings

Configure git behavior in your config:
- **Include/exclude merges**: Control merge commit inclusion
- **Branch selection**: Choose which branches to analyze
- **Path exclusion**: Exclude specific files/directories

## Troubleshooting

### Command not found: gdmc
```bash
# Check if Go bin is in PATH
echo $PATH | grep go

# Add Go bin to PATH
export PATH=$PATH:$(go env GOPATH)/bin

# Make it permanent
echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.zshrc
```

### No commits found
- Ensure you're in a git repository: `git status`
- Check author email matches your git config: `git config user.email`
- Verify the time interval has commits: `git log --oneline --since="yesterday"`

### API errors
- Verify API key: `gdmc config show`
- Check internet connection
- Ensure you have Gemini API access

## Requirements

- Go 1.18+ 
- Git installed and accessible
- Valid Gemini API key
- Internet connection for AI generation
- Terminal with shell access (bash/zsh/fish)

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

- 🐛 **Issues**: Report bugs or request features
- 📖 **Documentation**: Check the help commands (`gdmc --help`)
- 💬 **Discussions**: Share use cases and tips

---

Made with ❤️ for developers who want to communicate their work better. 