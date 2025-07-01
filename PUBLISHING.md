# Publishing Guide for add-n-commit

This guide covers how to distribute `anc` (add-n-commit) to users.

## Distribution Methods

### 1. Direct Go Install (Easiest)

Users with Go installed can install directly:

```bash
go install github.com/oconnorjohnson/add-n-commit@latest
```

### 2. GitHub Releases (Automated)

We use GoReleaser to automatically build and publish binaries for multiple platforms.

To create a new release:

```bash
# Tag a new version
git tag -a v0.1.0 -m "Release v0.1.0"
git push origin v0.1.0
```

GitHub Actions will automatically:

- Build binaries for Linux, macOS (Intel & ARM), and Windows
- Create a GitHub release with downloadable binaries
- Generate checksums

### 3. Homebrew (macOS/Linux)

After setting up your tap repository:

1. Create a tap repository: `homebrew-tap`
2. The GoReleaser config will automatically update the formula

Users can then install with:

```bash
brew tap oconnorjohnson/tap
brew install add-n-commit
```

### 4. Manual Installation

Users can download pre-built binaries from GitHub Releases:

```bash
# macOS (Apple Silicon)
curl -L https://github.com/oconnorjohnson/add-n-commit/releases/latest/download/add-n-commit_Darwin_arm64.tar.gz | tar xz
sudo mv anc /usr/local/bin/

# macOS (Intel)
curl -L https://github.com/oconnorjohnson/add-n-commit/releases/latest/download/add-n-commit_Darwin_x86_64.tar.gz | tar xz
sudo mv anc /usr/local/bin/

# Linux
curl -L https://github.com/oconnorjohnson/add-n-commit/releases/latest/download/add-n-commit_Linux_x86_64.tar.gz | tar xz
sudo mv anc /usr/local/bin/
```

## Installation Instructions for README

Add this to your README.md:

````markdown
## Installation

### Using Go

```bash
go install github.com/oconnorjohnson/add-n-commit@latest
```
````

### Using Homebrew (macOS/Linux)

```bash
brew tap oconnorjohnson/tap
brew install add-n-commit
```

### Download Binary

Download the latest binary for your platform from [releases](https://github.com/oconnorjohnson/add-n-commit/releases).

````

## Setting up Homebrew Tap

1. Create a new repository called `homebrew-tap`
2. GoReleaser will automatically create/update the formula when you create a release
3. The formula will be at `homebrew-tap/Formula/add-n-commit.rb`

## Version Management

Add version info to your main.go:
```go
var (
    version = "dev"
    commit  = "none"
    date    = "unknown"
)
````

Then use it in a version command or --version flag.

## Pre-release Checklist

- [ ] Update version in code if needed
- [ ] Update CHANGELOG.md
- [ ] Run tests: `go test ./...`
- [ ] Build locally: `go build -o anc`
- [ ] Test the binary
- [ ] Commit all changes
- [ ] Tag the release: `git tag -a v0.1.0 -m "Release v0.1.0"`
- [ ] Push the tag: `git push origin v0.1.0`
