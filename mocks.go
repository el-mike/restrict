package restrict

import (
	"fmt"

	"github.com/stretchr/testify/mock"
)

const (
	BasicRoleName        = "BasicRole"
	BasicResourceOneName = "BasicResourceOne"
	BasicResourceTwoName = "BasicResourceTwo"
)

const (
	CreateAction = "create"
	ReadAction   = "read"
	UpdateAction = "update"
	DeleteAction = "delete"
)

type SubjectMock struct {
	mock.Mock
	ID string
}

func (m *SubjectMock) GetRole() string {
	args := m.Called()

	// Note that this is not checking if first argument is string, so "" (empty string)
	// can be used when we want to test failing GetRole.
	if args.Get(0) == nil {
		return BasicRoleName
	}

	return args.String(0)
}

type ResourceMock struct {
	mock.Mock

	ID        string
	CreatedBy string
	Type      string
}

func (m *ResourceMock) GetResourceName() string {
	args := m.Called()

	// Note that this is not checking if first argument is string, so "" (empty string)
	// can be used when we want to test failing GetResourceName.
	if args.Get(0) == nil {
		return m.Type
	}

	return args.String(0)
}

func GetBasicRole() *Role {
	return &Role{
		ID:          BasicRoleName,
		Description: "Basic Role",
		Grants: GrantsMap{
			BasicResourceOneName: {
				&Permission{Action: CreateAction},
				&Permission{Action: ReadAction},
			},
		},
	}
}

func GetBasicPolicy() *PolicyDefinition {
	return &PolicyDefinition{
		Roles: Roles{
			BasicRoleName: GetBasicRole(),
		},
	}
}

func GetEmptyPolicy() *PolicyDefinition {
	return &PolicyDefinition{}
}

func GetBasicPolicyJSONString() string {
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
	}`, BasicRoleName,
		BasicRoleName,
		BasicResourceOneName,
		CreateAction,
		ReadAction,
	)
}

func GetBasicPolicyYAMLString() string {
	return fmt.Sprintf(`
		roles:
			%s:
				id: %s,
				description: "Basic role",
				grants:
					%s:
						- action: %s
						- action: %s
	`, BasicRoleName,
		BasicRoleName,
		BasicResourceOneName,
		CreateAction,
		ReadAction,
	)
}
