package restrict

// Permission - describes an Action that can be performed in regards to
// some Resource, with specified Conditions.
type Permission struct {
	// Action that will be allowed to perform if the Permission is granted, and Conditions
	// are satisfied.
	Action string `json:"action,omitempty" yaml:"action,omitempty"`
	// Conditions that need to be satisfied in order to allow the subject perform given Action.
	Conditions Conditions `json:"conditions,omitempty" yaml:"conditions,omitempty"`
	// Preset allows to extend Permission defined in PolicyDefinition.
	Preset string `json:"preset,omitempty" yaml:"preset,omitempty"`
	// ExtendPresetConditions specifies if preset's Conditions should be extended with
	// Permission's own Conditions, or should they be overridden.
	ExtendPresetConditions bool `json:"extendPresetConditions,omitempty" yaml:"extendPresetConditions,omitempty"`
}

// Permissions - alias type for slice of Permissions.
type Permissions []*Permission

// mergePreset - merges preset values into Permission.
func (p *Permission) mergePreset(preset *Permission) {
	if preset == nil {
		return
	}

	p.Action = preset.Action

	// If given Permission should extend preset's Conditions, we merge both
	// Condition maps. Otherwise, we just reassign it.
	if p.ExtendPresetConditions {
		p.Conditions = append(p.Conditions, preset.Conditions...)
	} else {
		p.Conditions = preset.Conditions
	}

	// We set Preset value to zero, to prevent subsequent merges while updating
	// policies.
	p.Preset = ""
}

// PermissionPresets - a map of reusable Permissions. Map key serves as a preset's name,
// that can be later referenced by Permission.
// Presets are applied when policy is loaded.
type PermissionPresets map[string]*Permission
