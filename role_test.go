package restrict

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gopkg.in/yaml.v3"
)

type roleSuite struct {
	suite.Suite
}

func TestRoleSuite(t *testing.T) {
	suite.Run(t, new(roleSuite))
}

func (s *roleSuite) TestUnmarshalJSON() {
	testRoleOne := "TestRole1"
	testRoleTwo := "TestRole2"

	rolesData := []byte(fmt.Sprintf(`{
		"%s": {
			"grants": {
				"%s": [
					{"action": "create"}
				]
			}
		},
		"%s": {
			"grants": {
				"%s": [
					{ "action": "update" }
				]
			}
		}
	
	}`, testRoleOne, basicResourceOneName, testRoleTwo, basicResourceOneName))

	assert.True(s.T(), json.Valid(rolesData))

	testRoles := Roles{}

	err := json.Unmarshal(rolesData, &testRoles)

	assert.Nil(s.T(), err)

	assert.IsType(s.T(), new(Role), testRoles[testRoleOne])
	assert.IsType(s.T(), new(Role), testRoles[testRoleTwo])

	assert.Equal(s.T(), testRoleOne, testRoles[testRoleOne].ID)
	assert.Equal(s.T(), testRoleTwo, testRoles[testRoleTwo].ID)
}

func (s *roleSuite) TestUnmarshalJSON_InvalidData() {
	// Array instead of map for grants
	rolesData := []byte(`{
		"TestRole1": {
			"grants": [
				{"action": "create"},
				{"action": "update"}
			]
		}
	}`)

	testRoles := Roles{}

	err := json.Unmarshal(rolesData, &testRoles)

	assert.Error(s.T(), err)
	assert.NotPanics(s.T(), func() { json.Unmarshal(rolesData, &testRoles) }) // nolint
}

func (s *roleSuite) TestUnmarshalYAML() {
	testRoleOne := "TestRole1"
	testRoleTwo := "TestRole2"

	rolesData := []byte(fmt.Sprintf(`
%s:
  grants:
    %s:
      - action: create
%s:
  grants:
    %s:
      - action: update
`, testRoleOne, basicResourceOneName, testRoleTwo, basicResourceOneName))

	testRoles := Roles{}

	err := yaml.Unmarshal(rolesData, &testRoles)

	assert.Nil(s.T(), err)

	assert.IsType(s.T(), new(Role), testRoles[testRoleOne])
	assert.IsType(s.T(), new(Role), testRoles[testRoleTwo])

	assert.Equal(s.T(), testRoleOne, testRoles[testRoleOne].ID)
	assert.Equal(s.T(), testRoleTwo, testRoles[testRoleTwo].ID)
}

func (s *roleSuite) TestUnmarshalYAML_InvalidData() {
	// Array instead of map for grants
	rolesData := []byte(`
  TestRole1:
    grants:
      - action: create
      - action: update
`)

	testRoles := Roles{}

	err := yaml.Unmarshal(rolesData, &testRoles)

	assert.Error(s.T(), err)
	assert.NotPanics(s.T(), func() { yaml.Unmarshal(rolesData, &testRoles) }) // nolint
}
