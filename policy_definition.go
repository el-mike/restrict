package restrict

// PolicyDefinition - describes a model of roles and grants that
// are defined for the application.
type PolicyDefinition struct {
	Roles          Roles
	IdentityField  string
	OwnershipField string
}
