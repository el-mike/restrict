package adapters

import (
	"fmt"

	"github.com/el-Mike/restrict"
)

const (
	basicRoleName        = "BasicRole"
	basicResourceOneName = "BasicResourceOne"
)

const (
	createAction = "create"
	readAction   = "read"
)

const BasicConditionOne = "BASIC_CONDITION_ONE"
const BasicConditionTwo = "BASIC_CONDITION_TWO"

func getBasicPolicyJSONString() string {
	return fmt.Sprintf(`{
		"roles": {
			"%s": {
				"id": "%s",
				"description": "Basic role",
				"grants": {
					"%s": [
						{ "action": "%s" },
						{ "action: "%s" }
					]
				}
			}
		}
	}`, basicRoleName,
		basicRoleName,
		basicResourceOneName,
		createAction,
		readAction,
	)
}

func getBasicPolicyYAMLString() string {
	return fmt.Sprintf(`
		roles:
			%s:
				id: %s,
				description: "Basic role",
				grants:
					%s:
						- action: %s
						- action: %s
	`, basicRoleName,
		basicRoleName,
		basicResourceOneName,
		createAction,
		readAction,
	)
}

func getBasicRole() *restrict.Role {
	return &restrict.Role{
		ID:          basicRoleName,
		Description: "Basic Role",
		Grants: restrict.GrantsMap{
			basicResourceOneName: {
				&restrict.Permission{Action: createAction},
				&restrict.Permission{Action: readAction},
			},
		},
	}
}

func getBasicPolicy() *restrict.PolicyDefinition {
	return &restrict.PolicyDefinition{
		Roles: restrict.Roles{
			basicRoleName: getBasicRole(),
		},
	}
}

func getEmptyPolicy() *restrict.PolicyDefinition {
	return &restrict.PolicyDefinition{}
}
