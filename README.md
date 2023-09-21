<div align="center">
    <img src=".github/grove_logo.png" alt="Grove logo" width="600"/>
    <h1>Portal Database Client</h1>
    <big>Go database client for interacting with the Portal HTTP Database</big>
    <div>
    <br/>
        <a href="https://github.com/pokt-foundation/db-client/pulse"><img src="https://img.shields.io/github/last-commit/pokt-foundation/db-client.svg"/></a>
        <a href="https://github.com/pokt-foundation/db-client/pulls"><img src="https://img.shields.io/github/issues-pr/pokt-foundation/db-client.svg"/></a>
        <a href="https://github.com/pokt-foundation/db-client/issues"><img src="https://img.shields.io/github/issues-closed/pokt-foundation/db-client.svg"/></a>
    </div>
</div>
<br/>

# Usage

This client should be installed in any Grove Go backend repo that needs to interact with PHD. The interface to be used depends on the interactions required by the repo:

- `IDBReader`: read-only
- `IDBWrite`: write-only
- `IDBClient`: read & write

# Publishing

This client will automatically publish when a Pull Request is merged to the `main` branch. The tag versioning system follows the Semantic Release standard and will be updated as such based on the commit messages in the merged branch.

# Development

This client should be updated to reflect any changes to PHD endpoints (including updating the tests file), published and then the necessary repos updated.

## Pre-Commit Installation

Before starting development work on this repo, `pre-commit` must be installed.

In order to do so, run the command **`make init-pre-commit`** from the repository root.

Once this is done, the following checks will be performed on every commit to the repo and must pass before the commit is allowed:

### 1. Basic checks

- **check-yaml** - Checks YAML files for errors
- **check-merge-conflict** - Ensures there are no merge conflict markers
- **end-of-file-fixer** - Adds a newline to end of files
- **trailing-whitespace** - Trims trailing whitespace
- **no-commit-to-branch** - Ensures commits are not made directly to `main`

### 2. Go-specific checks

- **go-fmt** - Runs `gofmt`
- **go-imports** - Runs `goimports`
- **golangci-lint** - run `golangci-lint run ./...`
- **go-critic** - run `gocritic check ./...`
- **go-build** - run `go build`
- **go-mod-tidy** - run `go mod tidy -v`

### 3. Detect Secrets

Will detect any potential secrets or sensitive information before allowing a commit.

- Test variables that may resemble secrets (random hex strings, etc.) should be prefixed with `test_`
- The inline comment `pragma: allowlist secret` may be added to a line to force acceptance of a false positive

## Packages in Use

- [Mockery](https://github.com/vektra/mockery) - Generates a mock from the Driver interface for testing purposes.

## Generating code

**Before committing any code to the repo, run the default Make target (`make`)**

This will generate up to date mocks of the `IDBReader`, `IDBWriter` & `IDBClient` interfaces for testing purposes.
