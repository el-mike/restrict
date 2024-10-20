package restrict

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type storageAdapterMock struct {
	mock.Mock
}

func (m *storageAdapterMock) LoadPolicy() (*PolicyDefinition, error) {
	args := m.Called()

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*PolicyDefinition), args.Error(1)
}

func (m *storageAdapterMock) SavePolicy(policy *PolicyDefinition) error {
	args := m.Called(policy)

	return args.Error(0)
}

type policyManagerSuite struct {
	suite.Suite

	testError error
}

func (s *policyManagerSuite) SetupSuite() {
	s.testError = errors.New("testError")
}

func TestPolicyManagerSuite(t *testing.T) {
	suite.Run(t, new(policyManagerSuite))
}

func (s *policyManagerSuite) TestNewPolicyManager() {
	testPolicy := getBasicPolicy()

	// Failing adapter
	testAdapter := new(storageAdapterMock)
	testAdapter.On("LoadPolicy").Return(nil, s.testError).Once()

	_, err := NewPolicyManager(testAdapter, false)

	assert.Error(s.T(), err)

	// Working adapter
	testAdapter = new(storageAdapterMock)
	testAdapter.On("LoadPolicy").Return(testPolicy, nil).Once()

	manager, err := NewPolicyManager(testAdapter, true)

	assert.Nil(s.T(), err)

	assert.IsType(s.T(), new(PolicyManager), manager)
	assert.Equal(s.T(), true, manager.autoUpdate)

	testAdapter.AssertNumberOfCalls(s.T(), "LoadPolicy", 1)
}

func (s *policyManagerSuite) TestLoadPolicy() {
	testPolicy := getBasicPolicy()

	// Working adapter
	testAdapter := new(storageAdapterMock)
	testAdapter.On("LoadPolicy").Return(testPolicy, nil).Once()

	manager, err := NewPolicyManager(testAdapter, false)

	assert.Nil(s.T(), err)

	testAdapter.On("LoadPolicy").Return(testPolicy, nil).Once()

	err = manager.LoadPolicy()

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), testPolicy, manager.policy)
	// We expect 2, since NewPolicyManager calls LoadPolicy as well.
	testAdapter.AssertNumberOfCalls(s.T(), "LoadPolicy", 2)

	// Failing adapter
	testAdapter.On("LoadPolicy").Return(nil, s.testError).Once()

	err = manager.LoadPolicy()

	assert.Error(s.T(), err)
}

func (s *policyManagerSuite) TestLoadPolicy_ApplyPresets() {
	testPolicy := getBasicPolicy()

	testAdapter := new(storageAdapterMock)
	testAdapter.On("LoadPolicy").Return(testPolicy, nil)

	manager, _ := NewPolicyManager(testAdapter, false)

	testPolicy.PermissionPresets = PermissionPresets{
		"testPreset1": &Permission{
			Action: "test-action-1",
		},
		"testPreset2": &Permission{
			Action: "test-action-2",
		},
	}

	testPermissions := Permissions{
		&Permission{Preset: "testPreset1"},
		&Permission{Preset: "testPreset2"},
	}

	testPolicy.Roles[basicRoleOneName].Grants[basicResourceOneName] = testPermissions

	err := manager.LoadPolicy()

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), "test-action-1", testPermissions[0].Action)
	assert.Equal(s.T(), "test-action-2", testPermissions[1].Action)
}

func (s *policyManagerSuite) TestLoadPolicy_ApplyPresetFailure() {
	testPolicy := getBasicPolicy()

	testAdapter := new(storageAdapterMock)
	testAdapter.On("LoadPolicy").Return(testPolicy, nil)

	manager, _ := NewPolicyManager(testAdapter, false)

	// applyPreset error handling - missing preset for Permission
	testPolicy.Roles[basicRoleOneName].Grants[basicResourceOneName][0].Preset = "incorrect-preset"

	err := manager.LoadPolicy()

	assert.Error(s.T(), err)
	assert.IsType(s.T(), new(PermissionPresetNotFoundError), err)
}

func (s *policyManagerSuite) TestSavePolicy() {
	testPolicy := getBasicPolicy()

	testAdapter := new(storageAdapterMock)
	testAdapter.On("LoadPolicy").Return(testPolicy, nil)

	manager, _ := NewPolicyManager(testAdapter, false)

	// Failing adapter
	testAdapter.On("SavePolicy", mock.Anything).Return(s.testError).Once()

	err := manager.SavePolicy()

	assert.Error(s.T(), err)
	testAdapter.AssertNumberOfCalls(s.T(), "SavePolicy", 1)
	testAdapter.AssertCalled(s.T(), "SavePolicy", testPolicy)

	// Working adapter
	testAdapter.On("SavePolicy", mock.Anything).Return(nil).Once()

	err = manager.SavePolicy()

	assert.Nil(s.T(), err)
}

func (s *policyManagerSuite) TestGetPolicy() {
	testPolicy := getBasicPolicy()

	testAdapter := new(storageAdapterMock)
	testAdapter.On("LoadPolicy").Return(testPolicy, nil)

	manager, _ := NewPolicyManager(testAdapter, false)

	policy := manager.GetPolicy()

	assert.Equal(s.T(), testPolicy, policy)
}

func (s *policyManagerSuite) TestGetRole() {
	testPolicy := getBasicPolicy()

	testAdapter := new(storageAdapterMock)
	testAdapter.On("LoadPolicy").Return(testPolicy, nil)

	manager, _ := NewPolicyManager(testAdapter, false)

	// Incorrect role
	role, err := manager.GetRole("INCORRECT_ROLE")

	assert.Nil(s.T(), role)
	assert.IsType(s.T(), new(RoleNotFoundError), err)

	// Correct role
	role, err = manager.GetRole(basicRoleOneName)

	assert.Equal(s.T(), testPolicy.Roles[basicRoleOneName], role)
	assert.Nil(s.T(), err)
}

func (s *policyManagerSuite) TestAddRole() {
	testPolicy := getBasicPolicy()

	testAdapter := new(storageAdapterMock)
	testAdapter.On("LoadPolicy").Return(testPolicy, nil)
	testAdapter.On("SavePolicy", mock.Anything).Return(nil)

	manager, _ := NewPolicyManager(testAdapter, false)

	// Existing role
	testExistingRole := &Role{
		ID: basicRoleOneName,
	}

	err := manager.AddRole(testExistingRole)

	assert.IsType(s.T(), new(RoleAlreadyExistsError), err)

	// New role
	testNewRole := &Role{ID: "NEW_ROLE"}

	err = manager.AddRole(testNewRole)

	assert.Nil(s.T(), err)
	testAdapter.AssertNumberOfCalls(s.T(), "SavePolicy", 0)

	// With auto update
	testNewRole = &Role{ID: "NEW_ROLE_2"}

	manager.EnableAutoUpdate()

	_ = manager.AddRole(testNewRole)

	// It should still be one
	testAdapter.AssertNumberOfCalls(s.T(), "SavePolicy", 1)
}

func (s *policyManagerSuite) TestUpdateRole() {
	testPolicy := getBasicPolicy()

	testAdapter := new(storageAdapterMock)
	testAdapter.On("LoadPolicy").Return(testPolicy, nil)
	testAdapter.On("SavePolicy", mock.Anything).Return(nil)

	manager, _ := NewPolicyManager(testAdapter, false)

	// New role
	testNewRole := &Role{
		ID: "NEW_ROLE",
	}

	err := manager.UpdateRole(testNewRole)

	assert.IsType(s.T(), new(RoleNotFoundError), err)

	// Existing role
	testExistingRole := &Role{
		ID: basicRoleOneName,
	}

	err = manager.UpdateRole(testExistingRole)

	assert.Nil(s.T(), err)
	testAdapter.AssertNumberOfCalls(s.T(), "SavePolicy", 0)

	// With auto update
	manager.EnableAutoUpdate()

	_ = manager.UpdateRole(testExistingRole)

	// It should still be one
	testAdapter.AssertNumberOfCalls(s.T(), "SavePolicy", 1)
}

func (s *policyManagerSuite) TestUpsertRole() {
	testPolicy := getBasicPolicy()

	testAdapter := new(storageAdapterMock)
	testAdapter.On("LoadPolicy").Return(testPolicy, nil)

	manager, _ := NewPolicyManager(testAdapter, false)

	testNewRole := &Role{
		ID: "NEW_ROLE",
	}

	testExistingRole := &Role{
		ID: basicRoleOneName,
	}

	err := manager.UpsertRole(testNewRole)
	newRole, _ := manager.GetRole(testNewRole.ID)

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), testNewRole, newRole)

	err = manager.UpsertRole(testExistingRole)
	existingRole, _ := manager.GetRole(testExistingRole.ID)

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), testExistingRole, existingRole)
}

func (s *policyManagerSuite) TestDeleteRole() {
	testPolicy := getBasicPolicy()

	testAdapter := new(storageAdapterMock)
	testAdapter.On("LoadPolicy").Return(testPolicy, nil)
	testAdapter.On("SavePolicy", mock.Anything).Return(nil)

	manager, _ := NewPolicyManager(testAdapter, false)

	// Incorrect role, without auto update
	err := manager.DeleteRole("INCORRECT_ROLE")

	assert.IsType(s.T(), new(RoleNotFoundError), err)
	testAdapter.AssertNumberOfCalls(s.T(), "SavePolicy", 0)

	// Correct role, auto update
	manager.EnableAutoUpdate()

	err = manager.DeleteRole(basicRoleOneName)

	assert.Nil(s.T(), err)
	testAdapter.AssertNumberOfCalls(s.T(), "SavePolicy", 1)

	_, err = manager.GetRole(basicRoleOneName)

	assert.IsType(s.T(), new(RoleNotFoundError), err)
}

func (s *policyManagerSuite) TestAddPermission() {
	testPolicy := getBasicPolicy()

	testAdapter := new(storageAdapterMock)
	testAdapter.On("LoadPolicy").Return(testPolicy, nil)
	testAdapter.On("SavePolicy", mock.Anything).Return(nil)

	manager, _ := NewPolicyManager(testAdapter, false)

	testPermission := &Permission{
		Action: createAction,
	}

	assert.Nil(s.T(), testPolicy.Roles[basicRoleOneName].Grants[basicResourceTwoName])

	// Incorrect role, without auto update
	err := manager.AddPermission("INCORRECT_ROLE", basicResourceTwoName, testPermission)

	assert.IsType(s.T(), new(RoleNotFoundError), err)

	err = manager.AddPermission(basicRoleOneName, basicResourceTwoName, testPermission)

	assert.Nil(s.T(), err)
	testAdapter.AssertNumberOfCalls(s.T(), "SavePolicy", 0)

	// With auto update
	manager.EnableAutoUpdate()

	_ = manager.AddPermission(basicRoleOneName, basicResourceTwoName, testPermission)

	testAdapter.AssertNumberOfCalls(s.T(), "SavePolicy", 1)

	// Try preset application
	testPermission.Preset = "incorrect-preset"

	err = manager.AddPermission(basicRoleOneName, basicResourceTwoName, testPermission)

	assert.IsType(s.T(), new(PermissionPresetNotFoundError), err)
}

func (s *policyManagerSuite) TestDeletePermission() {
	testPolicy := getBasicPolicy()

	testAdapter := new(storageAdapterMock)
	testAdapter.On("LoadPolicy").Return(testPolicy, nil)
	testAdapter.On("SavePolicy", mock.Anything).Return(nil)

	manager, _ := NewPolicyManager(testAdapter, false)

	assert.Equal(s.T(), len(testPolicy.Roles[basicRoleOneName].Grants[basicResourceOneName]), 2)

	// Incorrect role, without auto update
	err := manager.DeletePermission("INCORRECT_ROLE", basicResourceOneName, createAction)

	assert.IsType(s.T(), new(RoleNotFoundError), err)
	testAdapter.AssertNumberOfCalls(s.T(), "SavePolicy", 0)

	// Correct role, auto update
	manager.EnableAutoUpdate()

	err = manager.DeletePermission(basicRoleOneName, basicResourceOneName, createAction)

	assert.Nil(s.T(), err)
	testAdapter.AssertNumberOfCalls(s.T(), "SavePolicy", 1)

	role, _ := manager.GetRole(basicRoleOneName)

	assert.Equal(s.T(), len(role.Grants[basicResourceOneName]), 1)

	// Incorrect action (do nothing)
	err = manager.DeletePermission(basicRoleOneName, basicResourceOneName, "incorrect-action")

	assert.Nil(s.T(), err)
}

func (s *policyManagerSuite) TestAddPermissionPreset() {
	testPolicy := getBasicPolicy()

	testAdapter := new(storageAdapterMock)
	testAdapter.On("LoadPolicy").Return(testPolicy, nil)
	testAdapter.On("SavePolicy", mock.Anything).Return(nil)

	manager, _ := NewPolicyManager(testAdapter, false)

	assert.Nil(s.T(), testPolicy.PermissionPresets)

	testPresetName := "TestPreset"

	testPreset := &Permission{
		Action: "test-action-2",
	}

	// Preset exists
	testPolicy.PermissionPresets = PermissionPresets{
		testPresetName: testPreset,
	}

	err := manager.AddPermissionPreset(testPresetName, testPreset)

	assert.IsType(s.T(), new(PermissionPresetAlreadyExistsError), err)

	// Presets set to nil, Preset does not exist, without auto update
	testPolicy.PermissionPresets = nil

	err = manager.AddPermissionPreset(testPresetName, testPreset)

	policy := manager.GetPolicy()

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), testPreset, policy.PermissionPresets[testPresetName])

	// With auto update
	manager.EnableAutoUpdate()

	//nolint
	manager.AddPermissionPreset("SecondTestPreset", testPreset)

	testAdapter.AssertNumberOfCalls(s.T(), "SavePolicy", 1)
}

func (s *policyManagerSuite) TestUpdatePermissionPreset() {
	testPolicy := getBasicPolicy()

	testAdapter := new(storageAdapterMock)
	testAdapter.On("LoadPolicy").Return(testPolicy, nil)
	testAdapter.On("SavePolicy", mock.Anything).Return(nil)

	manager, _ := NewPolicyManager(testAdapter, false)

	testPresetName := "TestPreset"
	testPreset := &Permission{
		Action: "test-action-2",
	}

	// Preset does not exist
	testPolicy.PermissionPresets = PermissionPresets{
		testPresetName: testPreset,
	}

	testIncorrectPreset := &Permission{
		Action: "test-action-2",
	}

	err := manager.UpdatePermissionPreset("IncorrectName", testIncorrectPreset)

	assert.IsType(s.T(), new(PermissionPresetNotFoundError), err)

	// Preset exists
	err = manager.UpdatePermissionPreset(testPresetName, testPreset)

	preset := manager.getPermissionPreset(testPresetName)

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), testPreset.Action, preset.Action)
	testAdapter.AssertNumberOfCalls(s.T(), "SavePolicy", 0)

	// With auto update
	manager.EnableAutoUpdate()
	//nolint
	manager.UpdatePermissionPreset(testPresetName, testPreset)

	testAdapter.AssertNumberOfCalls(s.T(), "SavePolicy", 1)
}

func (s *policyManagerSuite) TestUpsertPermissionPreset() {
	testPolicy := getBasicPolicy()

	testAdapter := new(storageAdapterMock)
	testAdapter.On("LoadPolicy").Return(testPolicy, nil)

	manager, _ := NewPolicyManager(testAdapter, false)

	testPresetName := "TestPreset"
	testPreset := &Permission{
		Action: "test-action-2",
	}

	// Preset does not exist
	err := manager.UpsertPermissionPreset(testPresetName, testPreset)
	preset := manager.getPermissionPreset(testPresetName)

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), testPreset, preset)

	err = manager.UpsertPermissionPreset(testPresetName, testPreset)
	preset = manager.getPermissionPreset(testPresetName)

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), testPreset, preset)
}

func (s *policyManagerSuite) TestDeletePermissionPreset() {
	testPolicy := getBasicPolicy()

	testAdapter := new(storageAdapterMock)
	testAdapter.On("LoadPolicy").Return(testPolicy, nil)
	testAdapter.On("SavePolicy", mock.Anything).Return(nil)

	testPresetName := "TestPreset"
	testPreset := &Permission{
		Action: "test-action-2",
	}

	testPolicy.PermissionPresets = PermissionPresets{
		testPresetName: testPreset,
		"TestPreset2": &Permission{
			Action: "test-action-2",
		},
		"TestPreset3": &Permission{
			Action: "test-action-2",
		},
	}

	manager, _ := NewPolicyManager(testAdapter, false)

	assert.Equal(s.T(), len(testPolicy.PermissionPresets), 3)

	// Incorrect preset, without auto update
	err := manager.DeletePermissionPreset("INCORRECT_PRESET")

	assert.IsType(s.T(), new(PermissionPresetNotFoundError), err)
	testAdapter.AssertNumberOfCalls(s.T(), "SavePolicy", 0)

	// Correct role, auto update
	manager.EnableAutoUpdate()

	err = manager.DeletePermissionPreset(testPresetName)

	assert.Nil(s.T(), err)
	testAdapter.AssertNumberOfCalls(s.T(), "SavePolicy", 1)

	policy := manager.GetPolicy()

	assert.Equal(s.T(), len(policy.PermissionPresets), 2)
}

func (s *policyManagerSuite) TestDisableAutoUpdate() {
	testPolicy := getBasicPolicy()

	testAdapter := new(storageAdapterMock)
	testAdapter.On("LoadPolicy").Return(testPolicy, nil)

	manager, _ := NewPolicyManager(testAdapter, true)

	assert.True(s.T(), manager.autoUpdate)

	manager.DisableAutoUpdate()

	assert.False(s.T(), manager.autoUpdate)

}

func (s *policyManagerSuite) TestEnableAutoUpdate() {
	testPolicy := getBasicPolicy()

	testAdapter := new(storageAdapterMock)
	testAdapter.On("LoadPolicy").Return(testPolicy, nil)

	manager, _ := NewPolicyManager(testAdapter, false)

	assert.False(s.T(), manager.autoUpdate)

	manager.EnableAutoUpdate()

	assert.True(s.T(), manager.autoUpdate)

}
