# restrict

[![Go Report Card](https://goreportcard.com/badge/github.com/el-Mike/gochat)](https://goreportcard.com/report/github.com/el-Mike/restrict)
![License](https://img.shields.io/github/license/el-Mike/restrict)
[![master](https://github.com/el-Mike/restrict/actions/workflows/master.yml/badge.svg)](https://github.com/el-Mike/restrict/actions/workflows/master.yml)

Restrict is a authorization library that provides a hybrid of RBAC and ABAC models, allowing to define simple role-based policies while using more fine-grained control when needed.

## Table of contents
* [Installation](#installation)
* [Concepts](#concepts)
* [Basic usage](#basic-usage)
* [Policy](#policy)

## Installation
To install the library, run:
```
go get github.com/el-Mike/restrict
```

## Concepts
Restrict helps with building simple yet powerful access policies in declarative way. In order to do that, we introduce following concepts:
* **Subject** - an entity that wants to perform some actions. Needs to implement `Subject` interface and provide unique role name. Subject is usually any kind of user in your application. 
* **Resource** - an entity that is a target of the actions. Needs to implement `Resource` interface and provide unique resource name. Resource can be implemented by any entity or object in your domain.
* **Action** - an arbitrary operation that can be performed on given Resource.
* **Context** - a map of values containing any additional data needed to validate the access rights.

Restrict uses those informations to determine whether an access can be granted.

## Basic usage
```go
type User struct {
	ID string
}

// Subject interface implementation.
func (u *User) GetRole() string {
	return "User"
}

// Example entity with some fields.
type Conversation struct {
	ID           string
	CreatedBy    string
	Participants []string
	Active       bool
}

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
	// Create an instance of PolicyManager, which will be responsible for handling given PolicyDefinition.
	// You can use one of the built-in persistence adapters (in-memory or json/yaml file adapters), or provide your own.
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

## Policy
Policy is the description of access rules that should be enforced in given system. It consists of a Roles map, each with a set of Permissions granted per Resource, and Permission presets, that can be reused under various Roles and Resources. Here is an example of a policy:

```go
var policy = &restrict.PolicyDefinition{
	// A map of reusable Permissions. Key corresponds to a preset's name, which can
	// be later used to apply it.
	PermissionPresets: restrict.PermissionPresets{
		"updateOwn": &restrict.Permission{
			// An action that given Permission allows to perform.
			Action: "update",
			// Optional Conditions that when defined, need to be satisfied in order
			// to allow the access.
			Conditions: restrict.Conditions{
				// EqualCondition requires two values (described by ValueDescriptors)
				// to be equal in order to grant the access.
				// In this example we want to check if Conversation.CreatedBy and User.ID
				// are the same, meaning that Conversation was created by given User.
				&restrict.EqualCondition{
					// Optional ID helpful when we need to identify the exact Condition that failed
					// when checking the access.
					ID: "isOwner",
					// First value to compare.
					Left: &restrict.ValueDescriptor{
						Source: restrict.ResourceField,
						Field:  "CreatedBy",
					},
					// Second value to compare.
					Right: &restrict.ValueDescriptor{
						Source: restrict.SubjectField,
						Field:  "ID",
					},
				},
			},
		},
	},
	// A map of Roles. Key corresponds to a Role that Subjects in your system can belong to.
	Roles: restrict.Roles{
		"User": {
			// Optional, human readable description.
			Description: "This is a simple User role, with permissions for basic chat operations.",
			// Map of Permissions per Resource.
			// Grants map can be nil, meaning given Role has no Permissions (but can still inherit some).
			Grants: restrict.GrantsMap{
				"Conversation": {
					// Subject "User" can "read" any "Conversation".
					&restrict.Permission{Action: "read"},
					// Subject "User" can "create" a "Conversation".
					&restrict.Permission{Action: "create"},
					// Subject "User" can "update" ONLY a "Coversation" that was
					// created by it. Check "updateOwn" preset definition above.
					&restrict.Permission{Preset: "updateOwn"},
					// Subject "User" can "delete" ONLY inactive "Conversation".
					&restrict.Permission{
						Action: "delete",
						Conditions: restrict.Conditions{
							// EmptyCondition requires a value (described by ValueDescriptor)
							// to be empty (falsy) in order to grant the access.
							// In this example, we want Conversation.Active to be false.
							&restrict.EmptyCondition{
								ID: "deleteActive",
								Value: &restrict.ValueDescriptor{
									Source: restrict.ResourceField,
									Field:  "Active",
								},
							},
						},
					},
				},
			},
		},
		"Admin": {
			Description: "This is an Admin role, with permissions to manage Users.",
			// "Admin" can do everything "User" can.
			Parents: []string{"User"},
			// AND can also perform other operations that User itself
			// is not allowed to do.
			Grants: restrict.GrantsMap{
				// Please note that in order to make this work,
				// User needs to implement Resource interface.
				"User": {
					// Subject "Admin" can create a "User".
					&restrict.Permission{Action: "create"},
				},
			},
		},
	},
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

