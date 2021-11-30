package mocks

import "github.com/el-Mike/restrict"

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
	ID string
}

func (sm *SubjectMock) GetRole() string {
	return BasicRoleName
}

type ResourceMock struct {
	ID        string
	CreatedBy string
	Type      string
}

func (rm *ResourceMock) GetResourceName() string {
	return rm.Type
}

func GetBasicRole() *restrict.Role {
	return &restrict.Role{
		ID:          BasicRoleName,
		Description: "Basic Role",
		Grants: restrict.GrantsMap{
			BasicResourceOneName: {
				&restrict.Permission{Action: CreateAction},
				&restrict.Permission{Action: ReadAction},
			},
		},
	}
}

func GetBasicPolicy() *restrict.PolicyDefinition {
	return &restrict.PolicyDefinition{
		Roles: restrict.Roles{
			BasicRoleName: GetBasicRole(),
		},
	}
}

func GetEmptyPolicy() *restrict.PolicyDefinition {
	return &restrict.PolicyDefinition{}
}
