package restrict

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type permissionSuite struct {
	suite.Suite

	testPresetName string
	testAction     string
}

func TestPermissionSuite(t *testing.T) {
	suite.Run(t, new(permissionSuite))
}

func (s *permissionSuite) SetupSuite() {
	s.testAction = "test-action"
	s.testPresetName = "testPreset"
}

func (s *permissionSuite) TestMergePreset() {
	// Do not extend conditions
	testCondition := new(conditionMock)

	testPreset := &Permission{
		Action: s.testAction,
		Conditions: Conditions{
			testCondition,
			testCondition,
		},
	}

	testPermission := &Permission{}

	assert.NotPanics(s.T(), func() {
		testPermission.mergePreset(nil)
	})

	testPermission.mergePreset(testPreset)

	assert.Equal(s.T(), testPreset.Action, testPermission.Action)
	assert.Equal(s.T(), "", testPermission.Preset)
	assert.Equal(s.T(), 2, len(testPermission.Conditions))

	assert.Equal(s.T(), testCondition, testPermission.Conditions[0])

	// Extend conditions
	testPermission = &Permission{
		ExtendPresetConditions: true,
		Conditions: Conditions{
			testCondition,
		},
	}

	testPermission.mergePreset(testPreset)

	assert.Equal(s.T(), 3, len(testPermission.Conditions))
}
