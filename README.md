# GitHub Followers Manager

A powerful Go-based CLI tool to manage your GitHub following relationships. Follow followers of specific users who don't already follow you, or unfollow users who don't follow you back.

[![Go Version](https://img.shields.io/badge/go-1.24+-blue?logo=go)](https://golang.org/)
[![Go Reference](https://pkg.go.dev/badge/github.com/jvcByte/gh_followers.svg)](https://pkg.go.dev/github.com/jvcByte/gh_followers)
[![Release](https://img.shields.io/github/v/release/jvcByte/gh_followers)](https://github.com/jvcByte/gh_followers/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/jvcByte/gh_followers)](https://goreportcard.com/report/github.com/jvcByte/gh_followers)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## Features

- **Smart Following**: Follow followers of specific users who don't follow you back
- **Cleanup**: Unfollow users who don't follow you back
- **Efficient**: Handles large follower lists with early-exit optimization
- **Safe**: Interactive confirmation before bulk operations
- **Fast**: Concurrent processing with configurable worker pool
- **Secure**: OAuth2 token-based GitHub API authentication
- **Flexible**: Environment variables and command-line flags
- **Rate-Limited**: Configurable delays to respect GitHub API limits
- **Automated**: GitHub Actions workflows for daily operations

## Table of Contents

- [Installation](#installation)
- [Configuration](#configuration)
- [Usage](#usage)
- [Automated Workflows](#automated-workflows)
- [Project Structure](#project-structure)
- [Security](#security)
- [Troubleshooting](#troubleshooting)
- [Contributing](#contributing)
- [License](#license)

## Installation

### Prerequisites

- Go 1.24 or higher
- GitHub account
- GitHub Personal Access Token with `user:follow` scope

### Option 1: Build from Source

```bash
git clone https://github.com/jvcByte/gh_followers.git
cd gh_followers
go build -o gh_followers ./cmd
```

### Option 2: Install via Go

```bash
go install github.com/jvcByte/gh_followers/cmd@latest
```

## Configuration

### GitHub Personal Access Token

1. Go to GitHub Settings → Developer settings → Personal access tokens → Tokens (classic)
2. Generate a new token with the `user:follow` scope
3. Copy the token

### Environment Variables

Create a `.env` file in the project root:

```env
# Required
GH_USERNAME=your_github_username
GH_TOKEN=ghp_your_token_here

# Optional (defaults shown)
WORKER_COUNT=1      # Concurrent workers (1-5 recommended)
QUEUE_SIZE=3        # Task queue buffer size
TIME_DELAY_MS=2000  # Delay between API calls (1000-3000 recommended)
```

> Copy `.env.example` to `.env` and fill in your credentials.

## Usage

### Follow Command

Follow followers of a specific user who don't already follow you:

```bash
# Follow up to 5 followers of torvalds
./github-followers follow torvalds --limit 5

# Follow without confirmation
./github-followers follow username --force

# Follow with custom limit
./github-followers follow username --limit 10
```

**Flags:**
- `--limit, -l`: Maximum number of users to follow (default: 0 = no limit)
- `--force, -f`: Skip confirmation prompt

### Unfollow Command

Unfollow users who don't follow you back:

```bash
# Show users and confirm before unfollowing
./github-followers unfollow

# Unfollow without confirmation
./github-followers unfollow --force
```

**Flags:**
- `--force, -f`: Skip confirmation prompt

### Help

```bash
./github-followers --help
./github-followers follow --help
./github-followers unfollow --help
```

## Automated Workflows

The project includes GitHub Actions workflows for automated daily operations:

### Daily Follow Workflow

Automatically follows 1 user from each of 30 top developer accounts daily (30 users/day total):

- **Schedule**: Daily at 8 AM UTC
- **Trigger**: Manual via workflow_dispatch
- **Log**: Maintains `follow-log.json` with daily activity

### Daily Unfollow Workflow

Automatically unfollows users who don't follow you back (with 3-day grace period):

- **Schedule**: Daily at 9 AM UTC (1 hour after follow workflow)
- **Logic**: Only unfollows users followed 3+ days ago who haven't followed back
- **Trigger**: Manual via workflow_dispatch (with optional dry-run mode)

### Setup for Automated Workflows

1. Fork this repository
2. Go to Settings → Secrets and variables → Actions
3. Add secrets:
   - `GH_USERNAME`: Your GitHub username
   - `GH_TOKEN`: Your personal access token with `user:follow` scope
4. Enable GitHub Actions in your fork

## Project Structure

```
gh_followers/
├── .github/
│   └── workflows/       # GitHub Actions workflows
├── cmd/
│   └── main.go         # Application entry point
├── internal/
│   ├── cli/            # CLI commands (follow, unfollow)
│   ├── config/         # Configuration management
│   ├── git_hub_manager/# GitHub API client
│   ├── helper/         # Utility functions
│   └── worker/         # Concurrent worker pool
├── .env.example        # Example configuration
├── go.mod              # Go module definition
└── README.md           # This file
```

## Security

- Tokens require only `user:follow` scope (no repo access)
- Tokens loaded from environment variables only
- Never commit `.env` to version control (included in `.gitignore`)
- Use GitHub secrets for automated workflows

## Rate Limiting

GitHub API limits:
- 5,000 requests/hour for authenticated requests
- Each follow/unfollow = 1 API call
- Built-in configurable delays via `TIME_DELAY_MS`

**Recommended settings for large operations:**
```env
WORKER_COUNT=1
TIME_DELAY_MS=3000
```

## Troubleshooting

### Authentication Failed
- Verify `GH_TOKEN` is correct and has `user:follow` scope
- Check token hasn't expired

### Rate Limit Exceeded
- Increase `TIME_DELAY_MS` value
- Reduce `WORKER_COUNT` to 1
- Wait for rate limit reset (1 hour)

### Command Hangs on Large Accounts
- Use `--limit` flag to fetch only needed followers
- Example: `./github-followers follow torvalds --limit 10`

## Contributing

Contributions welcome! Please:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Disclaimer

This tool is provided as-is without warranties. Use at your own risk. Always review user lists before confirming bulk operations.

---

If you find this tool useful, please ⭐ the repository!
