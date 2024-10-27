package examples

import (
	"fmt"
	"log"

	"github.com/el-mike/restrict/v2"
	"github.com/el-mike/restrict/v2/adapters"
)

var basicUsagePolicy = &restrict.PolicyDefinition{
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
	policyManager, err := restrict.NewPolicyManager(adapters.NewInMemoryAdapter(basicUsagePolicy), true)
	if err != nil {
		log.Fatal(err)
	}

	manager := restrict.NewAccessManager(policyManager)

	if err = manager.Authorize(&restrict.AccessRequest{
		Subject:  &User{},
		Resource: &Conversation{},
		Actions:  []string{"delete"},
	}); err != nil {
		fmt.Println(err) // access denied for Action: "delete" on Resource: "Conversation"
	}
}
