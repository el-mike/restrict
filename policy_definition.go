package restrict

// PolicyDefinition - describes a model of roles and grants that
// are defined for the application.
type PolicyDefinition struct {
	Resources []string
	Actions   []Action
	Roles     map[string]*Role
}
