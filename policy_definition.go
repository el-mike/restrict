package restrict

// PolicyDefinition - describes a model of Roles and Permissions that
// are defined for the domain.
type PolicyDefinition struct {
	// PermissionPresets - a map of Permission presets.
	PermissionPresets PermissionPresets `json:"permissionPresets" yaml:"permissionPresets"`
	// Roles - collection of Roles used in the domain.
	Roles Roles `json:"roles" yaml:"roles"`
}
