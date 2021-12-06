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
	testPolicyProvider.On("GetRole", mock.Anything).Return(getBasicRole(), nil)

	manager := NewAccessManager(testPolicyProvider)

	testSubject := new(subjectMock)
	testResource := new(resourceMock)

	// Failing Subject, working Resource.
	testSubject.On("GetRole").Return("").Once()
	testResource.On("GetResourceName").Return(basicResourceOneName).Once()

	testRequest := &AccessRequest{
		Subject:  testSubject,
		Resource: testResource,
	}

	err := manager.Authorize(testRequest)

	assert.IsType(s.T(), new(RequestMalformedError), err)

	testSubject.AssertNumberOfCalls(s.T(), "GetRole", 1)
	testResource.AssertNumberOfCalls(s.T(), "GetResourceName", 1)

	// Working Subject, failing Resource.
	testSubject.On("GetRole").Return(basicRoleName)
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

	testSubject.On("GetRole").Return(basicRoleName)
	testResource.On("GetResourceName").Return(basicResourceOneName)

	testRequest := &AccessRequest{
		Subject:  testSubject,
		Resource: testResource,
	}

	err := manager.Authorize(testRequest)

	testPolicyProvider.AssertNumberOfCalls(s.T(), "GetRole", 1)
	assert.Error(s.T(), err)

	testRole := getBasicRole()

	// Empty grants check.
	testRole.Grants = nil
	testPolicyProvider.On("GetRole", mock.Anything).Return(testRole, nil)

	err = manager.Authorize(testRequest)

	assert.IsType(s.T(), new(NoAvailablePermissionsError), err)

	// 0 length grants.
	testRole.Grants = GrantsMap{}

	err = manager.Authorize(testRequest)

	assert.IsType(s.T(), new(NoAvailablePermissionsError), err)

	// 0 length grants with parents.
	testRole.Parents = []string{"TestParentOne"}

	err = manager.Authorize(testRequest)

	// Err is nil, because no actions are passed, meaning that Request is granted.
	assert.Nil(s.T(), err)

}

func (s *accessManagerSuite) TestAuthorize_ActionsWithoutConditions() {
	testRole := getBasicRole()

	testPolicyProvider := new(policyProviderMock)
	testPolicyProvider.On("GetRole", mock.Anything).Return(testRole, nil)

	manager := NewAccessManager(testPolicyProvider)

	testSubject := new(subjectMock)
	testResource := new(resourceMock)

	testSubject.On("GetRole").Return(basicRoleName)
	testResource.On("GetResourceName").Return(basicResourceOneName)

	// Action does not exist on role.
	testRequest := &AccessRequest{
		Subject:  testSubject,
		Resource: testResource,
		Actions:  []string{deleteAction},
	}

	err := manager.Authorize(testRequest)

	assert.IsType(s.T(), new(AccessDeniedError), err)
	assert.IsType(s.T(), new(PermissionNotGrantedError), err.(*AccessDeniedError).Reason())

	// One of the actions does not exists on role.
	testRequest.Actions = []string{createAction, deleteAction}

	err = manager.Authorize(testRequest)

	assert.IsType(s.T(), new(AccessDeniedError), err)
	assert.IsType(s.T(), new(PermissionNotGrantedError), err.(*AccessDeniedError).Reason())

	// One of the actions is empty string.
	testRequest.Actions = []string{createAction, ""}

	err = manager.Authorize(testRequest)

	assert.IsType(s.T(), new(RequestMalformedError), err)
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
	testRole := getBasicRole()
	testRole.Grants[basicResourceOneName] = append(testRole.Grants[basicResourceOneName], testConditionedPermission)

	testPolicyProvider := new(policyProviderMock)
	testPolicyProvider.On("GetRole", mock.Anything).Return(testRole, nil)

	manager := NewAccessManager(testPolicyProvider)

	testSubject := new(subjectMock)
	testResource := new(resourceMock)

	testSubject.On("GetRole").Return(basicRoleName)
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

	assert.IsType(s.T(), new(AccessDeniedError), err)
	assert.IsType(s.T(), new(ConditionNotSatisfiedError), err.(*AccessDeniedError).Reason())

	// AND - should expect all Conditions to be satisfied
	testConditionedPermission.Conditions = Conditions{testWorkingCondition, testWorkingCondition, testFailingCondition}

	err = manager.Authorize(testRequest)

	assert.IsType(s.T(), new(AccessDeniedError), err)
	assert.IsType(s.T(), new(ConditionNotSatisfiedError), err.(*AccessDeniedError).Reason())

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

func (s *accessManagerSuite) TestAuthorize_ActionsOnParents() {
	testRole := getBasicRole()
	testParentRole := getBasicParentRole()

	testRole.Parents = []string{testParentRole.ID}

	testPolicyProvider := new(policyProviderMock)
	testPolicyProvider.On("GetRole", basicRoleName).Return(testRole, nil)
	testPolicyProvider.On("GetRole", basicParentRoleName).Return(testParentRole, nil)

	manager := NewAccessManager(testPolicyProvider)

	testSubject := new(subjectMock)
	testResource := new(resourceMock)

	testSubject.On("GetRole").Return(basicRoleName)
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
	assert.IsDecreasing(s.T(), new(PermissionNotGrantedError), err.(*AccessDeniedError).Reason())

	testGrantParentRoleName := "BasicGrandParent"
	testGrandParentRole := getBasicParentRole()

	testGrandParentRole.ID = testGrantParentRoleName
	testGrandParentRole.Grants[basicResourceOneName] =
		append(testGrandParentRole.Grants[basicResourceOneName], &Permission{Action: deleteAction})

	testPolicyProvider.On("GetRole", testGrantParentRoleName).Return(testGrandParentRole, nil)

	testParentRole.Parents = []string{testGrantParentRoleName}

	// Action exist on grand parent.
	testRequest.Actions = []string{deleteAction}

	err = manager.Authorize(testRequest)

	assert.Nil(s.T(), err)

	// Ignore inhertiance cycle when permission is granted beforehand
	testGrandParentRole.Parents = []string{testRole.ID}

	err = manager.Authorize(testRequest)

	assert.Nil(s.T(), err)

	// Detect inhertiance cycle when permission is not granted beforehand
	testRequest.Actions = []string{"NewAction"}

	err = manager.Authorize(testRequest)

	assert.IsType(s.T(), new(RoleInheritanceCycleError), err)
}
