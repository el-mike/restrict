package restrict

// PermissionPresets - alias type for map of Permission pointers.
type PermissionPresets map[string]*PermissionPreset

// PolicyDefinition - describes a model of roles and grants that
// are defined for the application.
type PolicyDefinition struct {
	PermissionPresets PermissionPresets `json:"permissionPresets" yaml:"permissionPresets"`
	Roles             Roles             `json:"roles" yaml:"roles"`
}
