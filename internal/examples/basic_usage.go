package examples

import (
	"fmt"
	"log"

	"github.com/el-Mike/restrict"
	"github.com/el-Mike/restrict/adapters"
)

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
		Actions:        []string{"read", "delete"},
		SkipConditions: false,
	}); err != nil {
		fmt.Print(err) // Access denied for action: "delete". Reason: Permission for action: "delete" is not granted for Resource: "Conversation"
	}

	if accessError, ok := err.(*restrict.AccessDeniedError); ok {
		// Error() implementation. Returns a message in a form:
		// Access denied for action: "...". Reason: Permission for action: "..." is not granted for Resource: "..."
		fmt.Print(accessError)
		// Returns an AccessRequest that failed.
		fmt.Print(accessError.FailedRequest())
		// Returns underlying error which was the reason of a failure.
		fmt.Print(accessError.Reason())

		// If the reason of an AccessDeniedError was failed Condition,
		// this helper method returns it directly. Otherwise, nil will be returned.
		failedCondition := accessError.FailedCondition()

		// You can later cast the Condition to the type you want.
		if emptyCondition, ok := failedCondition.(*restrict.EmptyCondition); failedCondition != nil && ok {
			fmt.Print(emptyCondition.ID)
		}
	}
}
