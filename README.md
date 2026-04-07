# GoSymGym

[![Go Version](https://img.shields.io/badge/Go-1.22%2B-00ADD8?logo=go)](https://golang.org)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
![Go linting and formatting](https://img.shields.io/github/actions/workflow/status/Gitubrr/GoSymGym/build_and_lint.yaml?label=lint&job=lint)
![Build](https://img.shields.io/github/actions/workflow/status/Gitubrr/GoSymGym/build_and_lint.yaml?label=build&job=build)

**GoSymGym** is a CLI tool for getting information about GitHub repositories directly from the terminal. It displays stars, forks, issue count, language, creation and update dates, and a link to the repository.

---

## Installation

### Method 1: Installing with `go install`

```bash
go install github.com/Gitubrr/GoSymGym/cmd/GoSymGym@latest
```
Make sure `$GOPATH/bin` is added to your PATH. This will make the GoSymGym command available globally.

### Method 2: Building from source
Clone the repository:
```bash
git clone https://github.com/Gitubrr/GoSymGym.git
cd GoSymGym
```
Build the binary:
```bash
go build -o GoSymGym ./cmd
```

## Usage
Docks:
```bash
GitHub Repository Info CLI

usage: GoSymGym [-h] -o REPO_OWNER -n REPO_NAME [-t GITHUB_TOKEN] [-T TIMEOUT]

options:
        -h, --help              show this help message and exit
        -o REPO_OWNER           Name of the repository owner
        -n REPO_NAME            Repository name
        -t GITHUB_TOKEN         Personal access token
        -T TIMEOUT              Maximum response time sec
```
### Setting up a token (recommended)
Without a token, the GitHub API only allows 60 requests per hour. With a token, the limit increases to 5,000 requests per hour.

### How to get a token:
1. Go to GitHub → Settings → Developer settings → Personal access tokens → Tokens (classic).
2. Click Generate new token (classic).
3. Select permissions: at least repo to access repositories.
4. Copy the received token.

### How to use the token:
1. Method 1
Pass it via the `-t` flag
2. Method 2 (recommended) 
Set an environment variable (the token does not appear in the command history):
    ```bash
    export GITHUB_TOKEN="your token"
    GoSymGym -o golang -n go
    ```

## Result of work
You should have a similar result:

![workExample](images/workExample.png)
