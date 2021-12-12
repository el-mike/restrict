package restrict

import (
	"github.com/stretchr/testify/mock"
)

const (
	basicRoleName        = "BasicRole"
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

type subjectMock struct {
	mock.Mock

	ID         string
	FieldOne   string
	FieldTwo   int
	FieldThree []int
}

func (m *subjectMock) GetRole() string {
	args := m.Called()

	// Note that this is not checking if first argument is string, so "" (empty string)
	// can be used when we want to test failing GetRole.
	if args.Get(0) == nil {
		return basicRoleName
	}

	return args.String(0)
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

func getBasicRole() *Role {
	return &Role{
		ID:          basicRoleName,
		Description: "Basic Role",
		Grants: GrantsMap{
			basicResourceOneName: {
				&Permission{Action: createAction},
				&Permission{Action: readAction},
			},
		},
	}
}

func getBasicParentRole() *Role {
	role := getBasicRole()

	role.ID = basicParentRoleName
	role.Description = "Basic Parent Role"
	role.Grants[basicResourceOneName] = append(role.Grants[basicResourceOneName], &Permission{Action: updateAction})

	return role
}

func getBasicPolicy() *PolicyDefinition {
	return &PolicyDefinition{
		Roles: Roles{
			basicRoleName: getBasicRole(),
		},
	}
}
