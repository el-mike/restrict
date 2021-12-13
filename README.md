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
* [Access Request](#access-request)
* [Conditions](#conditions)
	* [Built-in Conditions](#built-in-conditions)
		* [Empty Condition](#empty-condition)
		* [Equal Condition](#equal-condition)
	* [Value Descriptor](#value-descriptor)
	* [Composition](#composition)
	* [Custom Conditions](#custom-conditions)

## Installation
To install the library, run:
```
go get github.com/el-Mike/restrict
```

## Concepts
Restrict helps with building simple yet powerful access policies in declarative way. In order to do that, we introduce following concepts:
* **Subject** - an entity that wants to perform some actions. Needs to implement `Subject` interface and provide unique role name. Subject is usually any kind of user or client in your domain. 
* **Resource** - an entity that is a target of the actions. Needs to implement `Resource` interface and provide unique resource name. Resource can be implemented by any entity or object in your domain.
* **Action** - an arbitrary operation that can be performed on given Resource.
* **Context** - a map of values containing any additional data needed to validate the access rights.
* **Condition** - requirement that needs to be satisfied in order to grant the access. There are couple of built-in Conditions, but any custom Condition can be added, as long as it implements `Condition` interface. Conditions are the way to express more granular control.

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

## Access Request
`AccessRequest` is an object describing a question about the access - can **Subject** perform given **Actions** on particular **Resource**.

If you only need RBAC-like functionality, or you want to perform "preliminary" access check (for example to just check if Subject can read given Resource at all, regardless of Conditions), empty Subject/Resource instances will be enough. Otherwise, Subject and Resource should be retrieved prior to authorization and passed along in AccessRequest.

For example, in typical backend application, Subject (User) will likely come from request context. Resource (Conversation) can be fetched from the database. Entire authorization process, along with Subject/Resource retrieval, could take place in middleware function.

Here is an example of AccessRequest:
```go
// ... manager setup

// Create empty instances or provide the correct entitites.
user := &User{}
conversation := &Conversation{}

accessRequest := &restrict.AccessRequest{
	// Required - who wants to perform the actions. It has to be
	// an instance of Subject interface. 
	Subject: user,
	// Required - on which Resource actions will be performed. It has to be
	// an instance of Resource interface.
	Resource: conversation,
	// Required - operations that given Subject wants to perform.
	Actions:  []string{"read", "create"},
	// Optional - a map of additional, external values that can be
	// accessed by Conditions. Values can be of any type.
	Context: restrict.Context{
		"SomeField": "someValue",
	},
	// Optional, lets you to skip Conditions checking.
	// Default: false.
	SkipConditions: false,
}

// If the access is granted, err will be nil - otherwise,
// an error will be returned containing an information about the failure.
err = manager.Authorize(accessRequest)
```

Alternatively to empty Subject/Resource instances, there are two helper functions - `UseSubject()` and `UseResource()`, that can be useful when you don't want to create empty instances or given Subject/Resource in not represented by any type in your domain. In this case, you can use:
```go
accessRequest := &restrict.AccessRequest{
	Subject:  restrict.UseSubject("User"),
	Resource: restrict.UseResource("Conversation"),
	Actions:  []string{"read", "create"},
}
```

## Conditions
Conditions allows to define more specific access control. It's similar to ABAC model, where more than just associated role need to be validated in order to grant the access. For example, a Subject can only update the Resource if it was created by it. Such a requirement can be expressed in Restrict as a Condition. If the Condition is not satsfied, access will not be granted, even if Subject does have the required Role.

Restrict ships with couple of built-in Conditions, but any number of custom Conditions can be added.

### Built-in Conditions

#### Equal Condition
`EqualCondition` and `NotEqualCondition` allow to check if two values, described by ValueDescriptors, are equal or not.
```go
&restrict.Permission{
	Action: "update",
	Conditions: restrict.Conditions{
		&restrict.EqualCondition{ // or &restrict.NotEqualCondition
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
```

#### Empty Condition
`EmptyCondition` and `NotEmptyCondition` allow to check if value described by ValueDescriptor is empty (not defined) or not empty (defined).
```go
&restrict.Permission{
	// An action that given Permission allows to perform.
	Action: "delete",
	// Optional Conditions that when defined, need to be satisfied in order
	// to allow the access.
	Conditions: restrict.Conditions{
		&restrict.EmptyCondition{ // or &restrict.NotEmptyCondition
			// Optional - helps with identifying failing Condition when checking an error
			// returned from .Authorized method.
			ID: "deleteActive",
			// Value to be checked.
			Value: &restrict.ValueDescriptor{
				Source: restrict.ResourceField,
				Field:  "Active",
			},
		},
	},
}
```

### Value Descriptor
`ValueDescriptor` is an object describing the value that needs to be retrieved from `AccessRequest` and tested by given Condition. `ValueDescriptor` allows to check various attributes without coupling your domain's entities to the library itself or forcing you to implement arbitrary interfaces. It uses reflection to get needed values.

`ValueDescriptor` needs to define value's source, which can be one of the predefined `ValueSource` enum type: `SubjectField`, `ResourceField`, `ContextField` or `Explicit`, and `Field` or `Value`, based on chosen source.
```go
type exampleCondition struct {
	ValueFromSubject *restrict.ValueDescriptor
	ValueFromResource *restrict.ValueDescriptor
	ValueFromContext *restrict.ValueDescriptor
	ExplicitValue *restrict.ValueDescriptor
}

condition := &exampleCondition{
	// This value will be taken from Subject's "SomeField" passed in AccessRequest.
	ValueFromSubject: &restrict.ValueDescriptor{
		// Required, ValueSource enum.
		Source: restrict.SubjectField,
		// Optional, string.
		Field: "SomeField",
	},
	// This value will be taken from Resource's "SomeField" passed in AccessRequest.
	ValueFromResource: &restrict.ValueDescriptor{
		Source: restrict.ResourceField,
		Field: "SomeField",
	},
	// This value will be taken from Context's "SomeField" passed in AccessRequest.
	ValueFromContext: &restrict.ValueDescriptor{
		Source: restrict.ContextField,
		Field: "SomeField",
	},
	// This value will be set explicitly to 10 - please note that we are using "Value"
	// instead of "Field" here. "Value" can be of any type.
	ExplicitValue: &restrict.ValueDescriptor{
		Source: restrict.Explicit,
		// Optional, interface{}.
		Value: 10,
	},
}
```

### Composition
Conditions can be composed in various ways, adding some flexibility to your policy. Let's consider following example:
```go
&restrict.Permission{
	Action: "delete",
	Conditions: restrict.Conditions{
		&restrict.EmptyCondition{
			ID: "ConditionOne",
			// ... Conditions details
		},
		&restrict.NotEmptyCondition{
			ID: "ConditionTwo",
			// ... Conditions details
		},
		&restrict.EqualCondition{
			ID: "ConditionThree",
			// ... Conditions details
		},
	},
}
```
In this case, 3 different **Conditions** need to be satisfied in order to grant permission for "delete" action. This way of defining Conditions works as **AND** operator - if just one of them fails, permission is not granted.

But we could also define Conditions like so:
```go
&restrict.Permission{
	Action: "delete",
	Conditions: restrict.Conditions{
		&restrict.EmptyCondition{
			ID: "ConditionOne",
			// ... Conditions details
		},
	},
},
&restrict.Permission{
	Action: "delete",
	Conditions: restrict.Conditions{
		&restrict.NotEmptyCondition{
			ID: "ConditionTwo",
			// ... Conditions details
		},
	},
},
&restrict.Permission{
	Action: "delete",
	Conditions: restrict.Conditions{
		&restrict.EqualCondition{
			ID: "ConditionThree",
			// ... Conditions details
		},
	},
}
```
We have 3 different **Permissions** with the same name but different sets of Conditions, effectively making it an **OR** operation - just one set of the Conditions needs to be satisfied in order to grant permission for "delete" action.

### Custom Conditions
You can add any number of Conditions to match requirements of your access policy. Condition needs to implement `Condition` interface:
```go
type Condition interface {
	// Type - returns Condition's type. Type needs to be unique
	// amongst other Conditions.
	Type() string

	// Check - returns true if Condition is satisfied by
	// given AccessRequest, false otherwise.
	Check(request *AccessRequest) error
}

```
For example, sticking to previous examples with `User` and `Conversation`, we can consider a case where we want to allow the `User` to read a `Conversation` only if it participates in it. Such a Condition could look like this:
```go
// Type is spelled with upper case - it's not necessary, but built-in Conditions
// follow this convention, to make a distinction between Condition type and other
// tokens, like preset or role name.
const hasUserConditionType = "BELONGS_TO"

type hasUserCondition struct{}

func (c *hasUserCondition) Type() string {
	return hasUserConditionType
}

func (c *hasUserCondition) Check(request *restrict.AccessRequest) error {
	user, ok := request.Subject.(*User)
	if !ok {
		return restrict.NewConditionNotSatisfiedError(c, request, fmt.Errorf("Subject has to be a User"))
	}

	conversation, ok := request.Resource.(*Conversation)
	if !ok {
		return restrict.NewConditionNotSatisfiedError(c, request, fmt.Errorf("Resource has to be a Conversation"))
	}

	for _, userId := range conversation.Participants {
		if userId == user.ID {
			return nil
		}
	}

	return restrict.NewConditionNotSatisfiedError(c, request, fmt.Errorf("User does not belong to Conversation with ID: %s", conversation.ID))
}

// ... policy definition - "User" -> "Grants"
ConversationResource: {
	// ... other permissions
	&restrict.Permission{
		Action: "read",
		Conditions: restrict.Conditions{
			&hasUserCondition{},
		},
	},
},

// ... check
user := &User{ID: "user-one"}
conversation := &Conversation{Participants: []string{"user-one"}}

err := manager.Authorize(&restrict.AccessRequest{
	Subject:  user,
	Resource: conversation,
	Actions:  []string{"read"},
})
// err is nil - "user-one" belongs to conversation's Participants slice.
```
Or you could want to allow to delete a `Conversation` only when it has less than 100 messages. In this case, you could create more generalized `Condition`, using `ValueDescriptor`, and pass `Max` value via Context:
```go
const greatherThanType = "GREATER_THAN"

type greaterThanCondition struct {
	// Please note that this field needs to have json/yaml tags if
	// you are using JSON/YAML based persistence.
	Value *restrict.ValueDescriptor `json:"value" yaml:"value"`
}

func (c *greaterThanCondition) Type() string {
	return greatherThanType
}

func (c *greaterThanCondition) Check(request *restrict.AccessRequest) error {
	value, err := c.Value.GetValue(request)
	if err != nil {
		return err
	}

	intValue, ok := value.(int)
	if !ok {
		return restrict.NewConditionNotSatisfiedError(c, request, fmt.Errorf("Value has to be an integer"))
	}

	intMax, ok := request.Context["Max"].(int)
	if !ok {
		return restrict.NewConditionNotSatisfiedError(c, request, fmt.Errorf("Max has to be an integer"))
	}

	if intValue > intMax {
		return restrict.NewConditionNotSatisfiedError(c, request, fmt.Errorf("Value is greater than max"))
	}

	return nil
}

// ... policy definition
ConversationResource: {
	// ... other permissions
	&restrict.Permission{
		Action: "delete",
		Conditions: restrict.Conditions{
			&greaterThanCondition{
				Value: &restrict.ValueDescriptor{
					Source: restrict.ResourceField,
					Field:  "MessagesCount",
				},
			},
		},
	},
},

// ... check
user := &User{}
conversation := &Conversation{MessagesCount: 90}

err = manager.Authorize(&restrict.AccessRequest{
	Subject:  user,
	Resource: conversation,
	Actions:  []string{"update"},
	Context: restrict.Context{
		"Max": 100,
	},
})
// err is nil - conversation has less than 100 messages.
```
You could also provide `Max` value as explicit value (see [Value Descriptor](#value-descriptor) section) and set it in your PolicyDefinition.

All of the checking logic is up to you - restrict only provides some building blocks and ensures that your Conditions will be used as specified in your policy.

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

