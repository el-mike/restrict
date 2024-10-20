package restrict

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type policyProviderMock struct {
	mock.Mock
}

func (m *policyProviderMock) GetRole(name string) (*Role, error) {
	args := m.Called(name)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*Role), args.Error(1)
}

type accessManagerSuite struct {
	suite.Suite

	testError error
}

func (s *accessManagerSuite) SetupSuite() {
	s.testError = errors.New("testError")
}

func TestAccessManagerSuite(t *testing.T) {
	suite.Run(t, new(accessManagerSuite))
}

func (s *accessManagerSuite) TestNewAccessManager() {
	testPolicyProvider := new(policyProviderMock)

	manager := NewAccessManager(testPolicyProvider)

	assert.NotNil(s.T(), manager)
	assert.IsType(s.T(), new(AccessManager), manager)
}

func (s *accessManagerSuite) TestAuthorize_MalformedRequest() {
	testPolicyProvider := new(policyProviderMock)
	testPolicyProvider.On("GetRole", mock.Anything).Return(nil, nil)

	manager := NewAccessManager(testPolicyProvider)

	testResource := new(resourceMock)

	testRequest := &AccessRequest{
		Subject:  nil,
		Resource: testResource,
	}

	err := manager.Authorize(testRequest)

	assert.IsType(s.T(), new(RequestMalformedError), err)

	testSubject := new(subjectMock)

	testRequest.Subject = testSubject
	testRequest.Resource = nil

	err = manager.Authorize(testRequest)

	assert.IsType(s.T(), new(RequestMalformedError), err)
}

func (s *accessManagerSuite) TestAuthorize_MalformedSubjectOrResource() {
	testPolicyProvider := new(policyProviderMock)
	testPolicyProvider.On("GetRole", mock.Anything).Return(getBasicRoleOne(), nil)

	manager := NewAccessManager(testPolicyProvider)

	testSubject := new(subjectMock)
	testResource := new(resourceMock)

	// Failing Subject, working Resource.
	testSubject.On("GetRoles").Return([]string{}).Once()
	testResource.On("GetResourceName").Return(basicResourceOneName).Once()

	testRequest := &AccessRequest{
		Subject:  testSubject,
		Resource: testResource,
	}

	err := manager.Authorize(testRequest)

	assert.IsType(s.T(), new(RequestMalformedError), err)

	testSubject.AssertNumberOfCalls(s.T(), "GetRoles", 1)
	testResource.AssertNumberOfCalls(s.T(), "GetResourceName", 1)

	// Working Subject, failing Resource.
	testSubject.On("GetRoles").Return(getBasicRolesSet())
	testResource.On("GetResourceName").Return("").Once()

	err = manager.Authorize(testRequest)

	assert.IsType(s.T(), new(RequestMalformedError), err)

	testResource.On("GetResourceName").Return(basicResourceOneName).Once()

	err = manager.Authorize(testRequest)

	// Note that err is nil because we actually supplied providerMock with correct Role.
	// Otherwise, err would still be set, but with different type.
	assert.Nil(s.T(), err)
}

func (s *accessManagerSuite) TestAuthorize_NoPermissions() {
	// Failing GetRole check.
	testPolicyProvider := new(policyProviderMock)
	testPolicyProvider.On("GetRole", mock.Anything).Return(nil, s.testError).Once()

	manager := NewAccessManager(testPolicyProvider)

	testSubject := new(subjectMock)
	testResource := new(resourceMock)

	testSubject.On("GetRoles").Return(getBasicRolesSet())
	testResource.On("GetResourceName").Return(basicResourceOneName)

	testRequest := &AccessRequest{
		Subject:  testSubject,
		Resource: testResource,
		Actions:  []string{"testAction"},
	}

	err := manager.Authorize(testRequest)

	testPolicyProvider.AssertNumberOfCalls(s.T(), "GetRole", 1)
	assert.Error(s.T(), err)

	testRole := getBasicRoleOne()

	// Empty grants check.
	testRole.Grants = nil
	testPolicyProvider.On("GetRole", mock.Anything).Return(testRole, nil)

	err = manager.Authorize(testRequest)

	assert.IsType(s.T(), new(AccessDeniedError), err)

	// 0 length grants.
	testRole.Grants = GrantsMap{}

	err = manager.Authorize(testRequest)

	assert.IsType(s.T(), new(AccessDeniedError), err)
}

func (s *accessManagerSuite) TestAuthorize_ActionsWithoutConditions() {
	testRole := getBasicRoleOne()

	testPolicyProvider := new(policyProviderMock)
	testPolicyProvider.On("GetRole", mock.Anything).Return(testRole, nil)

	manager := NewAccessManager(testPolicyProvider)

	testSubject := new(subjectMock)
	testResource := new(resourceMock)

	testSubject.On("GetRoles").Return(getBasicRolesSet())
	testResource.On("GetResourceName").Return(basicResourceOneName)

	// Action does not exist on role.
	testRequest := &AccessRequest{
		Subject:  testSubject,
		Resource: testResource,
		Actions:  []string{deleteAction},
	}

	err := manager.Authorize(testRequest)

	assert.IsType(s.T(), new(AccessDeniedError), err)
	assert.IsType(s.T(), new(PermissionError), err.(*AccessDeniedError).Errors.First())

	// One of the actions does not exist on role.
	testRequest.Actions = []string{createAction, deleteAction}
	err = manager.Authorize(testRequest)

	assert.IsType(s.T(), new(AccessDeniedError), err)
	assert.IsType(s.T(), new(PermissionError), err.(*AccessDeniedError).Errors.First())

	// One of the actions is empty string.
	testRequest.Actions = []string{createAction, ""}
	err = manager.Authorize(testRequest)

	assert.IsType(s.T(), new(RequestMalformedError), err)

	// Action exists on role.
	testRequest.Actions = []string{createAction}
	err = manager.Authorize(testRequest)

	assert.Nil(s.T(), err)
}

func (s *accessManagerSuite) TestAuthorize_ActionsWithConditions() {
	testConditionedAction := "conditioned-action"

	// Working Condition.
	testWorkingCondition := new(conditionMock)
	testWorkingCondition.On("Check", mock.Anything).Return(nil)

	testConditionedPermission := &Permission{
		Action:     testConditionedAction,
		Conditions: Conditions{testWorkingCondition},
	}
	testRole := getBasicRoleOne()
	testRole.Grants[basicResourceOneName] = append(testRole.Grants[basicResourceOneName], testConditionedPermission)

	testPolicyProvider := new(policyProviderMock)
	testPolicyProvider.On("GetRole", mock.Anything).Return(testRole, nil)

	manager := NewAccessManager(testPolicyProvider)

	testSubject := new(subjectMock)
	testResource := new(resourceMock)

	testSubject.On("GetRoles").Return(getBasicRolesSet())
	testResource.On("GetResourceName").Return(basicResourceOneName)

	testRequest := &AccessRequest{
		Subject:  testSubject,
		Resource: testResource,
		Actions:  []string{testConditionedAction},
	}

	err := manager.Authorize(testRequest)

	assert.Nil(s.T(), err)

	// Failing Condition
	testFailingCondition := new(conditionMock)
	testConditionError := NewConditionNotSatisfiedError(testFailingCondition, testRequest, s.testError)

	testFailingCondition.On("Check", mock.Anything).Return(testConditionError)

	testConditionedPermission.Conditions = Conditions{testFailingCondition}

	err = manager.Authorize(testRequest)
	permissionErr := err.(*AccessDeniedError).Errors.First()

	assert.IsType(s.T(), new(AccessDeniedError), err)
	assert.IsType(s.T(), new(PermissionError), permissionErr)
	assert.IsType(s.T(), new(ConditionNotSatisfiedError), permissionErr.ConditionError)

	// AND - should expect all Conditions to be satisfied
	testConditionedPermission.Conditions = Conditions{testWorkingCondition, testWorkingCondition, testFailingCondition}

	err = manager.Authorize(testRequest)
	permissionErr = err.(*AccessDeniedError).Errors.First()

	assert.IsType(s.T(), new(AccessDeniedError), err)
	assert.IsType(s.T(), new(PermissionError), permissionErr)
	assert.IsType(s.T(), new(ConditionNotSatisfiedError), permissionErr.ConditionError)

	// OR - should expect one of Permissions to be granted
	testConditionedPermission.Conditions = Conditions{testWorkingCondition, testFailingCondition}

	secondTestConditionedPermission := &Permission{
		Action:     testConditionedAction,
		Conditions: Conditions{testWorkingCondition},
	}

	testRole.Grants[basicResourceOneName] = append(testRole.Grants[basicResourceOneName], secondTestConditionedPermission)

	err = manager.Authorize(testRequest)

	assert.Nil(s.T(), err)
}

func (s *accessManagerSuite) TestAuthorize_UnknownConditionError() {
	testConditionedAction := "conditioned-action"
	testRole := getBasicRoleOne()

	testParentRole := getBasicParentRole()

	testRole.Parents = []string{testParentRole.ID}

	testPolicyProvider := new(policyProviderMock)
	testPolicyProvider.On("GetRole", basicRoleOneName).Return(testRole, nil)
	testPolicyProvider.On("GetRole", basicParentRoleName).Return(testParentRole, nil)

	// Failing Condition
	testFailingCondition := new(conditionMock)
	testConditionError := errors.New("Custom error")

	testFailingCondition.On("Check", mock.Anything).Return(testConditionError)

	testConditionedPermission := &Permission{
		Action:     testConditionedAction,
		Conditions: Conditions{testFailingCondition},
	}

	testRole.Grants[basicResourceOneName] = append(testRole.Grants[basicResourceOneName], testConditionedPermission)

	manager := NewAccessManager(testPolicyProvider)

	testSubject := new(subjectMock)
	testResource := new(resourceMock)

	testSubject.On("GetRoles").Return(getBasicRolesSet())
	testResource.On("GetResourceName").Return(basicResourceOneName)

	testRequest := &AccessRequest{
		Subject:  testSubject,
		Resource: testResource,
		Actions:  []string{testConditionedAction},
	}

	err := manager.Authorize(testRequest)

	assert.Equal(s.T(), testConditionError, err)

	// Unknown error on parents
	testRole.Grants[basicResourceOneName] = Permissions{}
	testParentRole.Grants[basicResourceOneName] = append(testParentRole.Grants[basicResourceOneName], testConditionedPermission)

	err = manager.Authorize(testRequest)

	assert.Equal(s.T(), testConditionError, err)
}

func (s *accessManagerSuite) TestAuthorize_ActionsOnParents() {
	testRole := getBasicRoleOne()
	testParentRole := getBasicParentRole()

	testRole.Parents = []string{testParentRole.ID}

	testPolicyProvider := new(policyProviderMock)
	testPolicyProvider.On("GetRole", basicRoleOneName).Return(testRole, nil)
	testPolicyProvider.On("GetRole", basicParentRoleName).Return(testParentRole, nil)

	manager := NewAccessManager(testPolicyProvider)

	testSubject := new(subjectMock)
	testResource := new(resourceMock)

	testSubject.On("GetRoles").Return(getBasicRolesSet())
	testResource.On("GetResourceName").Return(basicResourceOneName)

	// Action exist on parent.
	testRequest := &AccessRequest{
		Subject:  testSubject,
		Resource: testResource,
		Actions:  []string{updateAction},
	}

	err := manager.Authorize(testRequest)

	assert.Nil(s.T(), err)

	// Action does not exist on parent.
	testRequest.Actions = []string{deleteAction}

	err = manager.Authorize(testRequest)

	assert.IsType(s.T(), new(AccessDeniedError), err)
	assert.IsDecreasing(s.T(), new(PermissionError), err.(*AccessDeniedError).Errors.First())

	testGrantParentRoleName := "BasicGrandParent"
	testGrandParentRole := getBasicParentRole()

	testGrandParentRole.ID = testGrantParentRoleName
	testGrandParentRole.Grants[basicResourceOneName] =
		append(testGrandParentRole.Grants[basicResourceOneName], &Permission{Action: deleteAction})

	testPolicyProvider.On("GetRole", testGrantParentRoleName).Return(testGrandParentRole, nil)

	testParentRole.Parents = []string{testGrantParentRoleName}

	// Action exist on grandparent.
	testRequest.Actions = []string{deleteAction}

	err = manager.Authorize(testRequest)

	assert.Nil(s.T(), err)

	// Ignore inheritance cycle when permission is granted beforehand.
	testGrandParentRole.Parents = []string{testRole.ID}

	err = manager.Authorize(testRequest)

	assert.Nil(s.T(), err)

	// Detect inheritance cycle when permission is not granted beforehand.
	testRequest.Actions = []string{"NewAction"}

	err = manager.Authorize(testRequest)

	assert.IsType(s.T(), new(RoleInheritanceCycleError), err)
}

func (s *accessManagerSuite) TestAuthorize_MultipleRoles() {
	testRoleOne := getBasicRoleOne()
	testRoleTwo := getBasicRoleTwo()

	testMissingAction := "missing-action"

	testPolicyProvider := new(policyProviderMock)
	testPolicyProvider.On("GetRole", basicRoleOneName).Return(testRoleOne, nil)
	testPolicyProvider.On("GetRole", basicRoleTwoName).Return(testRoleTwo, nil)

	manager := NewAccessManager(testPolicyProvider)

	testSubject := new(subjectMock)
	testResource := new(resourceMock)

	testSubject.On("GetRoles").Return([]string{basicRoleOneName, basicRoleTwoName})
	testResource.On("GetResourceName").Return(basicResourceOneName)

	// Action does not exist on neither role.
	testRequest := &AccessRequest{
		Subject:  testSubject,
		Resource: testResource,
		Actions:  []string{testMissingAction, "delete"},
	}

	err := manager.Authorize(testRequest)
	accessError := err.(*AccessDeniedError)

	assert.IsType(s.T(), new(AccessDeniedError), err)

	// Should have one error per each Role that is missing the Permission for testMissingAction.
	// Missing "delete" action on testRoleOne won't be included, as it fails early.
	// @TODO: Adjust when complete validation option is introduced.
	assert.True(s.T(), len(accessError.Errors) == 2)

	roleOneErrors := accessError.Errors.GetByRoleName(basicRoleOneName)
	assert.True(s.T(), len(roleOneErrors) == 1)
	assert.True(s.T(), roleOneErrors[0].Action == testMissingAction)

	roleTwoErrors := accessError.Errors.GetByRoleName(basicRoleTwoName)
	assert.True(s.T(), len(roleTwoErrors) == 1)
	assert.True(s.T(), roleTwoErrors[0].Action == testMissingAction)

	assert.True(s.T(), len(accessError.Errors.GetByAction(testMissingAction)) == 2)

	// Action exists on one of the roles.
	testRequest.Actions = []string{deleteAction}
	err = manager.Authorize(testRequest)

	assert.Nil(s.T(), err)
}
