package restrict

import (
	"encoding/json"

	"gopkg.in/yaml.v3"
)

// Permission - describes an Action that can be performed in regards to
// some resource, with specified conditions.
type Permission struct {
	// Permission's name. Should be unique in the scope of Resource it belongs to.
	// If no Name is specified, it will be set to Action.
	Name string `json:"name,omitempty" yaml:"name,omitempty"`
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

// marshalablePermission - helper alias type for handling marshaling/unmarshaling
// of JSON and YAML structures.
type marshalablePermission Permission

// UnmarshalJSON - implementation of json.Marshaler.UnmarshalJSON, for creating Permission
// object from JSON. It takes care of initializating empty Conditions map.
func (p *Permission) UnmarshalJSON(jsonData []byte) error {
	var perm = marshalablePermission{
		Conditions: Conditions{},
	}

	if err := json.Unmarshal(jsonData, &perm); err != nil {
		return err
	}

	p.Name = perm.Name
	p.Action = perm.Action
	p.Preset = perm.Preset
	p.Conditions = perm.Conditions
	p.ExtendPresetConditions = perm.ExtendPresetConditions

	if err := p.ResolveName(); err != nil {
		return err
	}
	return nil
}

// UnmarshalYAML - implementation of yaml.Marshaler.UnmarshalYAML, for creating Permission
// object from YAML. It takes care of initializating empty Conditions map.
func (p *Permission) UnmarshalYAML(value *yaml.Node) error {
	var perm = marshalablePermission{
		Conditions: Conditions{},
	}

	if err := value.Decode(&perm); err != nil {
		return err
	}

	p.Name = perm.Name
	p.Action = perm.Action
	p.Preset = perm.Preset
	p.Conditions = perm.Conditions
	p.ExtendPresetConditions = perm.ExtendPresetConditions

	if err := p.ResolveName(); err != nil {
		return err
	}

	return nil
}

// ResolveName - sets Permission's Name based on it's state.
func (p *Permission) ResolveName() error {
	// If there was no Name specified, check if Permission has preset - if yes,
	// set preset as name. Otherwise, fall back to Action.
	if p.Name == "" {
		if p.Preset != "" {
			p.Name = p.Preset
		} else {
			p.Name = p.Action
		}
	}

	// if name could not be resolved, return an error.
	if p.Name == "" {
		return NewMissingPermissionNameError(p)
	}

	return nil
}

// MergePreset - merges preset values into Permission.
func (p *Permission) MergePreset(preset *Permission) {
	if preset == nil {
		return
	}

	p.Action = preset.Action

	// If given Permission should extend preset's Conditions, we merge both
	// Condition maps. Otherwise, we just reassign it.
	if p.ExtendPresetConditions {
		for key, condition := range preset.Conditions {
			p.Conditions[key] = condition
		}
	} else {
		p.Conditions = preset.Conditions
	}

	// We set Preset value to zero, to prevent subsequent merges while updating
	// policies.
	p.Preset = ""
}
