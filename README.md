# restrict

[![Go Report Card](https://goreportcard.com/badge/github.com/el-Mike/gochat)](https://goreportcard.com/report/github.com/el-Mike/restrict)
![License](https://img.shields.io/github/license/el-mike/restrict)
[![Tests](https://github.com/el-mike/restrict/actions/workflows/v2.yml/badge.svg)](https://github.com/el-mike/restrict/actions/workflows/v2.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/el-mike/restrict.svg)](https://pkg.go.dev/github.com/el-mike/restrict/v2)
![Release](https://img.shields.io/github/v/release/el-mike/restrict)

Restrict is an authorization library that provides a hybrid of RBAC and ABAC models, allowing to define simple role-based policies while also offering fine-grained control when needed. It helps you enforce your access policies away from the business logic, and express them in a convenient way.
### [API Reference](https://pkg.go.dev/github.com/el-mike/restrict/v2#section-documentation)
### [Changelog](https://github.com/el-mike/restrict/blob/v2/CHANGELOG.md)

## Table of contents
* [Installation](#installation)
* [Concepts](#concepts)
* [Basic usage](#basic-usage)
* [Policy](#policy)
* [Access Request](#access-request)
* [Access Manager](#access-manager)
* [Validation and errors](#validation-and-errors)
  * [Validation strategy](#validation-strategy)
* [Conditions](#conditions)
	* [Built-in Conditions](#built-in-conditions)
		* [Empty Condition](#empty-condition)
		* [Equal Condition](#equal-condition)
	* [Value Descriptor](#value-descriptor)
	* [Composition](#composition)
	* [Custom Conditions](#custom-conditions)
* [Presets](#presets)
* [PolicyManager and persistence](#policymanager-and-persistence)
	* [Storage adapter](#storage-adapter)
	* [Built-in Adapters](#built-in-adapters)
	* [Policy management](#policy-management)
* [Examples](#examples)
	* [Middleware function](#middleware-function)
* [Roadmap](#roadmap)

## Installation
To install the library, run:
```
go get github.com/el-mike/restrict/v2
```
**Go version 1.15+ is required!**  
Restrict follows [semantic versioning](https://semver.org/), so any changes will be applied according to its principles.

## Concepts
Restrict helps with building simple yet powerful access policies in a declarative way. In order to do that, we introduce following concepts:
* **Subject** - an entity that wants to perform some actions. Needs to implement `Subject` interface and provide a list of assigned roles. Subject is usually any kind of user or client in your domain. 
* **Resource** - an entity that is a target of an action. Needs to implement `Resource` interface and provide unique resource name. Resource can be implemented by any entity or object in your domain.
* **Action** - an arbitrary operation that can be performed on given Resource.
* **Context** - a map of values containing any additional data needed to validate the access rights.
* **Condition** - requirement that needs to be satisfied in order to grant the access. There are a couple of built-in Conditions, but any custom Condition can be added, as long as it implements `Condition` interface. Conditions are the way to express more granular control.

## Basic usage
```go
type User struct {
	ID string
}

// Subject interface implementation.
func (u *User) GetRoles() []string {
	return []string{"User"}
}

type Conversation struct {
	ID            string
	CreatedBy     string
	Participants  []string
	MessagesCount int
	Active        bool
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
		Subject:        &User{},
		Resource:       &Conversation{},
		Actions:        []string{"delete"},
	}); err != nil {
		fmt.Println(err) // access denied for Action: "delete" on Resource: "Conversation"
	}
}
```

## Policy
Policy is the description of access rules that should be enforced in given system. It consists of a Roles map, each with a set of Permissions granted per Resource, and Permission presets, that can be reused under various Roles and Resources. Here is an example of a policy:

```go
var policy = &restrict.PolicyDefinition{
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
					// created by it. Check "updateOwn" preset definition below.
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
	// A map of reusable Permissions. Key corresponds is a preset's name.
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
}
```

## Access Request
`AccessRequest` is an object describing a question about the access - can a **Subject** perform given **Actions** on a particular **Resource**.

If you only need RBAC-like functionality, or you want to perform "preliminary" access check (for example to just check if Subject can read given Resource at all, regardless of Conditions), empty Subject/Resource instances will be enough. Otherwise, Subject and Resource should be retrieved prior to authorization and passed along in AccessRequest.

For example, in typical backend application, Subject (User) will likely come from request context. Resource (Conversation) can be fetched from the database. Entire authorization process, along with Subject/Resource retrieval, could take place in middleware function.

Here is an example of AccessRequest:
```go
// ... manager setup

// Create empty instances or provide the correct entities.
user := &User{}
conversation := &Conversation{}

accessRequest := &restrict.AccessRequest{
	// Subject - subject (typically a user) that wants to perform given Actions.
	// Needs to implement Subject interface.
	// Required.
	Subject: user,
	// Resource - resource that given Subject wants to interact with.
	// Needs to implement Resource interface.
	// Required.
	Resource: conversation,
	// Actions - list of operations Subject wants to perform on given Resource.
	// Required.
	Actions:  []string{"read", "create"},
	// Context - map of any additional values needed while checking the Conditions.
	// Optional.
	Context: restrict.Context{
		"SomeField": "someValue",
	},
	// SkipConditions - allows to skip Conditions while checking the access.
	// Optional, default: false.
	SkipConditions: false,
	// CompleteValidation - when true, validation will not return early, and all possible errors
	// will be returned, including all Conditions checks.
	// Optional, default: false.
	CompleteValidation bool
}

// If the access is granted, err will be nil - otherwise,
// an error will be returned containing information about the failure.
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

## Access Manager
`AccessManager` is responsible for the actual validation. Once set up with proper `PolicyManager` instance (see [PolicyManager and persistence](#policymanager-and-persistence) for details), you can use its `Authorize` method in order to check given `AccessRequest`. `Authorize` returns an error if access is not granted, and `nil` otherwise (meaning there is no error and the access is granted).

```go
var policy = &restrict.PolicyDefinition{
	// ... policy details
}

adapter := adapters.NewInMemoryAdapter(policy)
policyMananger, err := restrict.NewPolicyManager(adapter, true)
if err != nil {
	log.Fatal(err)
}

manager := restrict.NewAccessManager(policyMananger)

accessRequest := &restrict.AccessRequest{
	// ... request details
}

// 
err := manager.Authorize(accessRequest)
```

## Validation and errors
Since `Authorize` method depends on various operations, including external ones provided in a form of Conditions, its return type is a general `error` type. However, when error is caused by actual policy validation (i.e. Permission is not granted or Conditions were not satisfied), `Authorize` returns an instance of `AccessDeniedError`, which provides a lot of information and context about the reason behind denied access. It facilitates easy error handling and debugging.

If the error is caused by policy validation, it will be returned according to [validation strategy](#validation-strategy) used in the request. Any other error will be returned immediately.

`AccessDeniedError` can contain one or more underlying reasons - errors of type `PermissionError`, that can be accessed via `Reasons` property. Each `PermissionError` can contain zero or more `ConditionNotSatisfiedError` errors, accessible via `ConditionErrors` property.  

The error returned can be tested using type assertion.
```go
err := manager.Authorize(accessRequest)
if accessError, ok := err.(*restrict.AccessDeniedError); ok {
	// Error() implementation. Returns a message in a form: "access denied for Action/s: ... on Resource: ..."
	fmt.Println(accessError)
	// Returns an AccessRequest that failed.
	fmt.Println(accessError.Request)
	// Returns first reason for the denied access.
	// Especially helpful in fail-early mode, where there will only be one Reason.
	fmt.Println(accessError.FirstReason())
	
	// Reasons property will hold all errors that caused the access to be denied.
	for _, permissionErr := range accessError.Reasons {
		fmt.Println(permissionErr)
		fmt.Println(permissionErr.Action)
		fmt.Println(permissionErr.RoleName)
		fmt.Println(permissionErr.ResourceName)
		
		// Returns first ConditionNotSatisfied error for given PermissionError, if any was returned for given PermissionError.
		// Especially helpful in fail-early mode, where there will only be one failed Condition.
		fmt.Println(permissionErr.FirstConditionError())
		
		// ConditionErrors property will hold all ConditionNotSatisfied errors.
		for _, conditionErr := range permissionErr.ConditionErrors {
			fmt.Println(conditionErr)
			fmt.Println(conditionErr.Reason)
			
			// Every ConditionNotSatisfied contains an instance of Condition that returned it,
			// so it can be tested using type assertion to get more details about failed Condition.
			if emptyCondition, ok := conditionErr.Condition.(*restrict.EmptyCondition); ok {
				fmt.Println(emptyCondition.ID)
			}
		}
	}
}
```

### Validation strategy
Access can be validated in one of two modes - fail-early or complete. Note that **by default, restrict will run in fail-early mode**. To use complete validation, set `CompleteValidation` to true on `AccessRequest` being authorized.

In **fail-early mode**, validation will be aborted as soon as any policy validation error is encountered. `AccessDeniedError` will only contain one reason for failure, even when Subject with multiple Roles is being tested, or when `AccessRequest` contains multiple Actions. `FirstReason()` and `FirstConditionError()` can be used to make error handling easier in fail-early mode.

```go
err := manager.Authorize(accessRequest)
if accessError, ok := err.(*restrict.AccessDeniedError); ok {
	reason := accessError.FirstReason()
	fmt.Println(reason)
	
	conditionErr := reason.FirstConditionError()
	if conditionErr != nil {
		fmt.Println(conditionErr)
		
		// If needed, ConditionNotSatisfied can be introspected further. 
		if emptyCondition, ok := conditionErr.Condition.(*restrict.EmptyCondition); ok {
			fmt.Println(emptyCondition.ID)
		}
	}
}
```

In **complete mode**, validation will not be aborted on first encountered error - instead, it will continue until all Roles and Actions are checked, gathering the errors along the way. If given `Permission` have more than one `Condition` attached, they will all be checked as well.
This mode is helpful when full information about access denial is needed - for example when debugging, or when checking user's access as an administrator.

`Reasons` property on `AccessDeniedError` provides some additional utility methods, that make it easier to handle the errors:
- `GetByRoleName` will return only the reasons related to checking given Role
- `GetByAction` will return only the reasons related to checking given Action
- `GetFailedActions` will return all Actions for which access has been denied

```go
accessRequest := &restrict.AccessRequest{
	// ...
	CompleteValidation: true,
}

err := manager.Authorize(accessRequest)
if accessError, ok := err.(*restrict.AccessDeniedError); ok {
	// Get only the Reasons related to Role user, that were returned for "read" action.
	userReasons := accessError.Reasons.GetByRoleName("user").GetByAction("read")
	
	for _, permissionErr := range userReasons {
		fmt.Println(permissionErr)
		
		// ConditionErrors property will hold all ConditionNotSatisfied errors.
		for _, conditionErr := range permissionErr.ConditionErrors {
		fmt.Println(conditionErr)
		
		// If needed, ConditionNotSatisfied can be introspected further. 
		if emptyCondition, ok := conditionErr.Condition.(*restrict.EmptyCondition); ok {
			fmt.Println(emptyCondition.ID)
			}
		}
	}
}
```
Please note than in complete mode, **errors from parent Roles will not be gathered!** If a check for given Role and all its parents returned an error, only the error from the original Role will be returned as a reason.

You can find more comprehensive error handling in [this example](https://github.com/el-mike/restrict/blob/v2/internal/examples/error_handling_example.go).

## Conditions
Conditions allows to define more specific access control. It's similar to ABAC model, where more than just associated role needs to be validated in order to grant the access. For example, a Subject can only update the Resource if it was created by it. Such a requirement can be expressed with Restrict as a Condition. If the Condition is not satisfied, access will not be granted, even if Subject does have the required Role.

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
Please note that if you want to get `AccessDeniedError` from `Authorize` method when your custom Condition is not satisfied, you should return `ConditionNotSatisfiedError` - any other error will be treated like internal error and will be returned directly. You can use `NewConditionNotSatisfiedError` to create a new instance if it.

Sticking to previous examples with `User` and `Conversation`, we can consider a case where we want to allow the `User` to read a `Conversation` only if it participates in it. Such a Condition could look like this:
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

## Presets
Preset is simply a Permission with unique name, that you can reuse across your PolicyDefinition. The main reason behind introducing presets is saving the necessity of defining the same Conditions for different Permissions, as in many cases the same action will have identical Conditions for various Resources.

Let's consider following example:
```go
var policy = &restrict.PolicyDefinition{
	Roles: restrict.Roles{
		"User": {
			Grants: restrict.GrantsMap{
				"Conversation": {
					&restrict.Permission{Preset: "updateOwn"},
				},
				"Message": {
					&restrict.Permission{Preset: "updateOwn"},
				}
			},
		},
	},
	PermissionPresets: restrict.PermissionPresets{
		"updateOwn": &restrict.Permission{
			Action: "update",
			Conditions: restrict.Conditions{
				&restrict.EqualCondition{
					// ... condition details
				},
			},
		},
	},
}
```
In this case, we can express that `User` can update only its own `Conversation` or `Message`, without the need for repeating Conditions definition.

But what in case we need the same Conditions, but for different actions? We can just define an action name of Permission itself, along the preset:
```go
var policy = &restrict.PolicyDefinition{
	Roles: restrict.Roles{
		"User": {
			Grants: restrict.GrantsMap{
				"Conversation": {
					&restrict.Permission{
						Action: "update",
						Preset: "accessOwn",
					},
					&restrict.Permission{
						Action: "delete",
						Preset: "accessOwn",
					},
				},
				"Message": {
					&restrict.Permission{
						Action: "update",
						Preset: "accessOwn",
					},
					&restrict.Permission{
						Action: "delete",
						Preset: "accessOwn",
					},
				},
			},
		},
	},
	PermissionPresets: restrict.PermissionPresets{
		// Note that this preset does not have an Action anymore,
		// but it can - Permission's Action just overrides preset's Action.
		"accessOwn": &restrict.Permission{
			Conditions: restrict.Conditions{
				&restrict.EqualCondition{
					// ... condition details
				},
			},
		},
	},
}
```
Now we can reuse the same Conditions for different actions. 

## PolicyManager and persistence
`PolicyManager` provides thread-safe, runtime policy management, that allows to easily retrieve and manipulate your policy. It is used by `AccessManager` to retrieve Permissions for given role when checking `AccessRequest`. You can create `PolicyManager` like so:
```go
myStorageAdapter := // ... create adapter 

// Second argument let's you set auto-update feature. If set to true,
// any change made via PolicyManager will be automatically saved with StorageAdapter.
// You can later disable/enable auto-update with DisableAutoUpdate() and EnableAutoUpdate() methods.
policyManager, err := restrict.NewPolicyManager(myStorageAdapter, true)
```

### Storage adapter
`PolicyManager` relies on `StorageAdapter` instance, which is an entity providing perstistence logic for PolicyDefinition. Restrict ships with two built-in, ready to go StorageAdapters, but you can easily provide your own, by implementing following interface:
```go
type StorageAdapter interface {
	// LoadPolicy - loads and returns PolicyDefinition from underlying
	// storage provider.
	LoadPolicy() (*PolicyDefinition, error)

	// SavePolicy - saves PolicyDefinition in underlying storage provider.
	SavePolicy(policy *PolicyDefinition) error
}
```

All of the Restrict's models are JSON and YAML compliant, so you can marshal/unmarshal PolicyDefinition in those formats.

### Built-in Adapters

#### InMemoryAdapter
Simple, in-memory storage for PolicyDefinition. You can create and use it like so:
```go
inMemoryAdapter := adapters.NewInMemoryAdapter(policy)

policyManager, err := restrict.NewPolicyManager(inMemoryAdapter, true)
```
`InMemoryAdapter` will keep PolicyDefinition object directly in memory. Using `InMemoryAdapter`, you will propably keep your PolicyDefinition in .go files. Please note that when using `InMemoryAdapter`, calling `inMemoryAdapter.SavePolicy(policy)` does NOT save it permanently, therefore any changes you've made with `PolicyManager` will be lost once program exits.

#### FileAdapter
`FileAdapter` uses file system to persit the PolicyDefinition. You can use JSON or YAML files. Here is how to use it:
```go
fileAdapter := adapters.NewFileAdapter("filename.json", adapters.JSONFile)
// alternatively, to use YAML file:
fileAdapter := adapters.NewFileAdapter("filename.yml", adapters.YAMLFile)

policyManager, err := restrict.NewPolicyManager(fileAdapter, true)
```
`FileAdapter` will load the PolicyDefinition from given file, and keep it in sync with any changes when calling `fileAdapter.SavePolicy(policy)`. You can also easily trasform your in-memory PolicyDefinition into JSON/YAML one, like so:
```go
policy := &restrict.PolicyDefinition{
	// ... policy details
}

// assuming "filename.json" file does not exist or is empty
fileAdapter := adapters.NewFileAdapter("filename.json", adapters.JSONFile)

err := fileAdapter.SavePolicy(policy)
if err != nil {
	// ... error handling
}
```
Please refer to:
* [JSON policy](https://github.com/el-mike/restrict/blob/v2/internal/examples/policy_example.json)
* [YAML policy](https://github.com/el-mike/restrict/blob/v2/internal/examples/policy_example.yaml)

To see examples of JSON/YAML policies.

### Policy management
`PolicyManager` provides a set of methods that will help you manage your policy in a dynamic way. You can manipulate it in runtime, or create custom tools in order to add and remove Roles, grant and revoke Permissions or manage presets. Full list of `PolicyManager`'s methods can be found here:

[PolicyManager docs](https://pkg.go.dev/github.com/el-mike/restrict#PolicyManager)

## Examples

### Middleware function
You can easily use Restrict as a middleware function in your backend application. Here is an example of using Restrict inside Gin framework's handler function:
```go
func WithAuth(
	actions []string,
	resource restrict.Resource,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := authenticateUser(c)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Set(UserContextKey, user)

		// If any previous handler fetched concrete Resource, we want to
		// override it here, so we can test it against Conditions if necessary.
		if contextResource, ok := c.Get(ResourceContextKey); ok {
			if res, ok := contextResource.(restrict.Resource); ok {
				resource = res
			}
		}

		err = accessManager.Authorize(&restrict.AccessRequest{
			Subject:  user,
			Resource: resource,
			Actions:  actions,
		})

		// If error is related to insufficient privileges, we want to
		// return appropriate response.
		if accessError, ok := err.(*restrict.AccessDeniedError); ok {
			log.Print(accessError)

			c.AbortWithError(http.StatusForbidden, accessError)
			return
		}

		// If not, some other error occurred.
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	}
}
```
You can see entire working example **[HERE](https://github.com/el-mike/restrict-middleware-example)**

## Roadmap
Check this [Project](https://github.com/el-mike/restrict/projects/3) for upcoming features.

## Contributing
### Prerequisites

1. Install [golangci-lint](https://golangci-lint.run/usage/install/)
2. Set your IDE to use golangci-lint ([instructions](https://golangci-lint.run/usage/integrations/))
3. Run `git config core.hooksPath .githooks` to wire up project's git hooks

Please note that test suite uses external library [testify](github.com/stretchr/testify), and minimal version of Go required to run the test suite is **1.17**. 

### Conventions

This repository follows [ConventionalCommits](https://www.conventionalcommits.org/en/v1.0.0/) specification for creating commit messages. There is `prepare-commit-msg` hook set up to ensure following those rules. Branch names should also reflect the type of work it contains - one of following should be used:
* `feat/<task-description>`
* `fix/<task-description>`
* `chore/<task-description>`
