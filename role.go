package restrict

// Permissions - alias type for slice of Permissions.
type Permissions []*Permission

// GrantsMap - alias type for map of Permission slices.
type GrantsMap map[string]Permissions

// Role - a set of granted permissions, that can be
// assign to a user.
type Role struct {
	ID          string    `json:"id" yaml:"id"`
	Description string    `json:"description,omitempty" yaml:"description,omitempty"`
	Grants      GrantsMap `json:"grants" yaml:"grants"`
	Parents     []string  `json:"parents,omitempty" yaml:"parents,omitempty"`
}

// Roles - alias type for map of Roles.
type Roles map[string]*Role
