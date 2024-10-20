package restrict

import (
	"github.com/stretchr/testify/mock"
)

const (
	basicRoleOneName     = "BasicRoleOne"
	basicRoleTwoName     = "BasicRoleTwo"
	basicParentRoleName  = "BasicParentRole"
	basicResourceOneName = "BasicResourceOne"
	basicResourceTwoName = "BasicResourceTwo"
)

const (
	createAction = "create"
	readAction   = "read"
	updateAction = "update"
	deleteAction = "delete"
)

const basicConditionOne = "BASIC_CONDITION_ONE"

func getBasicRolesSet() []string {
	return []string{basicRoleOneName}
}

type subjectMock struct {
	mock.Mock

	ID         string
	FieldOne   string
	FieldTwo   int
	FieldThree []int
}

func (m *subjectMock) GetRoles() []string {
	args := m.Called()

	// If not specified exactly when using the mock, it will return an array with single role.
	// This is helpful for tests where roles don't matter than much.
	if args.Get(0) == nil {
		return getBasicRolesSet()
	}

	return args.Get(0).([]string)
}

type resourceMock struct {
	mock.Mock

	ID         string
	CreatedBy  string
	Type       string
	FieldOne   string
	FieldTwo   int
	FieldThree []int
}

func (m *resourceMock) GetResourceName() string {
	args := m.Called()

	// Note that this is not checking if first argument is string, so "" (empty string)
	// can be used when we want to test failing GetResourceName.
	if args.Get(0) == nil {
		return m.Type
	}

	return args.String(0)
}

type conditionMock struct {
	mock.Mock
}

func (m *conditionMock) Type() string {
	args := m.Called()

	if args.Get(0) == nil {
		return basicConditionOne
	}
	return args.String(0)
}

func (m *conditionMock) Check(request *AccessRequest) error {
	args := m.Called()

	return args.Error(0)
}

func getBasicRoleOne() *Role {
	return &Role{
		ID:          basicRoleOneName,
		Description: "Basic Role",
		Grants: GrantsMap{
			basicResourceOneName: {
				&Permission{Action: createAction},
				&Permission{Action: readAction},
			},
		},
	}
}

func getBasicRoleTwo() *Role {
	return &Role{
		ID:          basicRoleTwoName,
		Description: "Basic Role",
		Grants: GrantsMap{
			basicResourceOneName: {
				&Permission{Action: createAction},
				&Permission{Action: readAction},
				&Permission{Action: updateAction},
				&Permission{Action: deleteAction},
			},
		},
	}
}

func getBasicParentRole() *Role {
	role := getBasicRoleOne()

	role.ID = basicParentRoleName
	role.Description = "Basic Parent Role"
	role.Grants[basicResourceOneName] = append(role.Grants[basicResourceOneName], &Permission{Action: updateAction})

	return role
}

func getBasicPolicy() *PolicyDefinition {
	return &PolicyDefinition{
		Roles: Roles{
			basicRoleOneName: getBasicRoleOne(),
		},
	}
}
