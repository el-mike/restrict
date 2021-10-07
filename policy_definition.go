package restrict

type PermissionPresets map[string]*Permission

// PolicyDefinition - describes a model of roles and grants that
// are defined for the application.
type PolicyDefinition struct {
	PermissionPresets PermissionPresets `json:"permissionPresets" yaml:"permissionPresets"`
	Roles             Roles             `json:"roles" yaml:"roles"`
	IdentityField     string            `json:"identityField" yaml:"identityField"`
	OwnershipField    string            `json:"ownershipField" yaml:"ownershipField"`
}
