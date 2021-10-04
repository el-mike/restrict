package restrict

// PolicyDefinition - describes a model of roles and grants that
// are defined for the application.
type PolicyDefinition struct {
	Roles          Roles  `json:"roles" yaml:"roles"`
	IdentityField  string `json:"identityField" yaml:"identityField"`
	OwnershipField string `json:"ownershipField" yaml:"ownershipField"`
}
