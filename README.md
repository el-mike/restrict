# restrict

[![Go Report Card](https://goreportcard.com/badge/github.com/el-Mike/gochat)](https://goreportcard.com/report/github.com/el-Mike/restrict)
![License](https://img.shields.io/github/license/el-Mike/restrict)
[![release](https://github.com/el-Mike/restrict/actions/workflows/release.yml/badge.svg)](https://github.com/el-Mike/restrict/actions/workflows/release.yml)

Restrict is a authorization library that provides a hybrid of RBAC and ABAC models, allowing to define simple role-based policies while using more fine-grained control when needed.

## Table of contents
* [Installation](#installation)
* [Concepts](#concepts)
* [Basic usage](#basic-usage)

## Installation
To install the library, run:
```
go get github.com/el-Mike/restrict
```

## Concepts
Restrict helps with building simple yet powerful access policies in declarative way. In order to do that, we introduce following concepts:
* **Subject** - an entity that wants to perform some actions. Needs to implement Subject interface and provide unique role name. Subject is usually any kind of user in your application. 
* **Resource** - an entity that is a target of the actions. Needs to implement Resource interface and provide unique resource name. Resource can be implemented by any entity or object in your domain.
* **Action** - an arbitrary operation that can be performed on given Resource.
* **Context** - a map of values containing any additional data needed to validate the access rights.

Restrict uses those informations to determine whether an access can be granted.

## Basic usage
```go
type User struct{}

// Subject interface implementation.
func (u *User) GetRole() string {
	return "User"
}

type Conversation struct{}

// Resource interface implementation.
func (c *Conversation) GetResourceName() string {
	return "Conversation"
}

var policy = &restrict.PolicyDefinition{
	Roles: restrict.Roles{
		"User": {
			Grants: restrict.GrantsMap{
				"Conversation": {
					&restrict.Permission{Action: "read"},
					&restrict.Permission{Action: "create"},
				},
			},
		},
	},
}

func main() {
	policyMananger, err := restrict.NewPolicyManager(adapters.NewInMemoryAdapter(policy), true)
	if err != nil {
		log.Fatal(err)
	}

	manager := restrict.NewAccessManager(policyMananger)

	if err = manager.Authorize(&restrict.AccessRequest{
		Subject:  &User{},
		Resource: &Conversation{},
		Actions:  []string{"read", "delete"},
	}); err != nil {
		fmt.Print(err) // Access denied for action: "delete". Reason: Permission for action: "delete" is not granted for Resource: "Conversation"

	}
}
```

## Development
### Prerequisites

1. Install [golangci-lint](https://golangci-lint.run/usage/install/)
2. Set your IDE to use golangci-lint ([instructions](https://golangci-lint.run/usage/integrations/))
3. Install [python3](https://www.python.org/download/releases/3.0/)
4. Run `git config core.hooksPath .githooks` to wire up project's git hooks

### Conventions

This repository follows [ConventionalCommits](https://www.conventionalcommits.org/en/v1.0.0/) specification for creating commit messages. There is `prepare-commit-msg` hook set up to ensure following those rules. Branch names should also reflect the type of work it contains - one of following should be used:
* `feature/<task-description>`
* `bugfix/<task-description>`
* `chore/<task-description>`

