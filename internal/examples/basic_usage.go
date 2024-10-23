package examples

import (
	"fmt"
	"log"

	"github.com/el-mike/restrict"
	"github.com/el-mike/restrict/adapters"
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
	policyManager, err := restrict.NewPolicyManager(adapters.NewInMemoryAdapter(policy), true)
	if err != nil {
		log.Fatal(err)
	}

	manager := restrict.NewAccessManager(policyManager)

	if err = manager.Authorize(&restrict.AccessRequest{
		Subject:        &User{},
		Resource:       &Conversation{},
		Actions:        []string{"read"},
		SkipConditions: false,
	}); err != nil {
		fmt.Println(err) // access denied for Action: "read" on Resource: "Conversation". Reason: ...
	}

	if accessError, ok := err.(*restrict.AccessDeniedError); ok {
		// Error() implementation. Returns a message in a form:
		// Access denied for Action/s: "...". Reason: Permission for action: "..." is not granted for Resource: "..."
		fmt.Println(accessError)
		// Returns an AccessRequest that failed.
		fmt.Println(accessError.Request)
		// We can use FirstReason() to get the first encountered PermissionError.
		// Especially helpful in fail-early mode, where there will only be one Reason.
		fmt.Println(accessError.FirstReason())

		// We can use Reasons property to loop over all PermissionErrors.
		for _, permissionErr := range accessError.Reasons {
			fmt.Println(permissionErr)
			fmt.Println(permissionErr.Action)
			fmt.Println(permissionErr.RoleName)
			fmt.Println(permissionErr.ResourceName)

			// If the reason of a Permission was failed Condition, it will be stored in ConditionErrors slice.
			// Especially helpful in fail-early mode, where there will only be one failed Condition.
			conditionErr := permissionErr.FirstConditionError()

			if conditionErr != nil {
				// It can be later cast to the type you want.
				if emptyCondition, ok := conditionErr.Condition.(*restrict.EmptyCondition); ok {
					fmt.Print(emptyCondition.ID)
				}
			}

		}
	}
}
