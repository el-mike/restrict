package restrict

// GrantsMap - alias type for map of Permission slices.
type GrantsMap map[string]Permissions

// Role - describes privileges of a Role's members.
type Role struct {
	// ID - unique identifier of the Role.
	ID string `json:"id" yaml:"id"`
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
