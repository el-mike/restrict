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

	testResource := new(ResourceMock)

	testRequest := &AccessRequest{
		Subject:  nil,
		Resource: testResource,
	}

	err := manager.Authorize(testRequest)

	assert.IsType(s.T(), new(RequestMalformedError), err)

	testSubject := new(SubjectMock)

	testRequest.Subject = testSubject
	testRequest.Resource = nil

	err = manager.Authorize(testRequest)

	assert.IsType(s.T(), new(RequestMalformedError), err)
}

func (s *accessManagerSuite) TestAuthorize_MalformedSubjectOrResource() {
	testPolicyProvider := new(policyProviderMock)
	testPolicyProvider.On("GetRole", mock.Anything).Return(GetBasicRole(), nil)

	manager := NewAccessManager(testPolicyProvider)

	testSubject := new(SubjectMock)
	testResource := new(ResourceMock)

	// Failing Subject, working Resource.
	testSubject.On("GetRole").Return("").Once()
	testResource.On("GetResourceName").Return(BasicResourceOneName).Once()

	testRequest := &AccessRequest{
		Subject:  testSubject,
		Resource: testResource,
	}

	err := manager.Authorize(testRequest)

	assert.IsType(s.T(), new(RequestMalformedError), err)

	testSubject.AssertNumberOfCalls(s.T(), "GetRole", 1)
	testResource.AssertNumberOfCalls(s.T(), "GetResourceName", 1)

	// Working Subject, failing Resource.
	testSubject.On("GetRole").Return(BasicRoleName)
	testResource.On("GetResourceName").Return("").Once()

	err = manager.Authorize(testRequest)

	assert.IsType(s.T(), new(RequestMalformedError), err)

	testResource.On("GetResourceName").Return(BasicResourceOneName).Once()

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

	testSubject := new(SubjectMock)
	testResource := new(ResourceMock)

	testSubject.On("GetRole").Return(BasicRoleName)
	testResource.On("GetResourceName").Return(BasicResourceOneName)

	testRequest := &AccessRequest{
		Subject:  testSubject,
		Resource: testResource,
	}

	err := manager.Authorize(testRequest)

	testPolicyProvider.AssertNumberOfCalls(s.T(), "GetRole", 1)
	assert.Error(s.T(), err)

	testRole := GetBasicRole()

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
