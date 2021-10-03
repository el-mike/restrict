package restrict

type GrantsMap map[string][]*Permission

// Role - a set of granted permissions, that can be
// assign to a user.
type Role struct {
	ID          string    `json:"id"`
	Description string    `json:"description"`
	Grants      GrantsMap `json:"grants"`
	Parents     []string  `json:"parents"`
}

type Roles map[string]*Role
