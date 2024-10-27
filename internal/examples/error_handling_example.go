package examples

import (
	"fmt"
	"log"

	"github.com/el-mike/restrict/v2"
	"github.com/el-mike/restrict/v2/adapters"
)

var errorHandlingExamplePolicy = &restrict.PolicyDefinition{
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

func mainErrorHandling() {
	// Create an instance of PolicyManager, which will be responsible for handling given PolicyDefinition.
	// You can use one of the built-in persistence adapters (in-memory or json/yaml file adapters), or provide your own.
	policyManager, err := restrict.NewPolicyManager(adapters.NewInMemoryAdapter(errorHandlingExamplePolicy), true)
	if err != nil {
		log.Fatal(err)
	}

	manager := restrict.NewAccessManager(policyManager)
	err = manager.Authorize(&restrict.AccessRequest{
		Subject:  &User{},
		Resource: &Conversation{},
		Actions:  []string{"read"},
	})

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
}
