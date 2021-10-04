package restrict

type GrantsMap map[string][]*Permission

// Role - a set of granted permissions, that can be
// assign to a user.
type Role struct {
	ID          string    `json:"id" yaml:"id"`
	Description string    `json:"description,omitempty" yaml:"description,omitempty"`
	Grants      GrantsMap `json:"grants" yaml:"grants"`
	Parents     []string  `json:"parents,omitempty" yaml:"parents,omitempty"`
}

type Roles map[string]*Role
