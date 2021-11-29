# restrict

[![Go Report Card](https://goreportcard.com/badge/github.com/el-Mike/gochat)](https://goreportcard.com/report/github.com/el-Mike/restrict)
![License](https://img.shields.io/github/license/el-Mike/restrict)
[![dev](https://github.com/el-Mike/restrict/actions/workflows/go.yml/badge.svg?branch=develop)](https://github.com/el-Mike/restrict/actions/workflows/go.yml)
[![release](https://github.com/el-Mike/restrict/actions/workflows/release.yml/badge.svg)](https://github.com/el-Mike/restrict/actions/workflows/release.yml)

Restrict is a authorization library that provides a simple RBAC model, while allowing to use a more fine-grained access control when needed, similar to ABAC. 
## Prerequisites

1. Install [golangci-lint](https://golangci-lint.run/usage/install/)
2. Set your IDE to use golangci-lint ([instructions](https://golangci-lint.run/usage/integrations/))
3. Install [python3](https://www.python.org/download/releases/3.0/)
4. Run `git config core.hooksPath .githooks` to wire up project's git hooks

## Conventions

This repository follows [ConventionalCommits](https://www.conventionalcommits.org/en/v1.0.0/) specification for creating commit messages. There is `prepare-commit-msg` hook set up to ensure following those rules. Branch names should also reflect the type of work it contains - one of following should be used:
* `feature/<task-description>`
* `bugfix/<task-description>`
* `chore/<task-description>`

