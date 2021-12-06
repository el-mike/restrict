package restrict

import (
	"encoding/json"

	"gopkg.in/yaml.v3"
)

// GrantsMap - alias type for map of Permission slices.
type GrantsMap map[string]Permissions

// Role - describes privileges of a Role's members.
type Role struct {
	// ID - unique identifier of the Role.
	ID string
	// Description - optional description for a Role.
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	// Grants - contains sets of Permissions assigned to Resources.
	Grants GrantsMap `json:"grants" yaml:"grants"`
	// Parents - other Roles that given Role inherits from. If a Permission is granted
	// for a parent, it is also granted for a child.
	Parents []string `json:"parents,omitempty" yaml:"parents,omitempty"`
}

// Roles - alias type for map of Roles.
type Roles map[string]*Role

// UnmarshalJSON - unmarshals a JSON-coded map of Roles.
func (rs Roles) UnmarshalJSON(jsonData []byte) error {
	if rs == nil {
		rs = Roles{}
	}

	var jsonRoles map[string]*Role

	if err := json.Unmarshal(jsonData, &jsonRoles); err != nil {
		return err
	}

	for key, role := range jsonRoles {
		role.ID = key
		rs[key] = role
	}

	return nil
}

// UnmarshalYAML - unmarshals a YAML-coded map of Roles.
func (rs Roles) UnmarshalYAML(value *yaml.Node) error {
	if rs == nil {
		rs = Roles{}
	}

	var yamlRoles map[string]*Role

	if err := value.Decode(&yamlRoles); err != nil {
		return err
	}

	for key, role := range yamlRoles {
		role.ID = key
		rs[key] = role
	}

	return nil
}
