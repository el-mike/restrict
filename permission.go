package restrict

import (
	"encoding/json"

	"gopkg.in/yaml.v3"
)

// Permission - describes an Action that can be performed in regards to
// some resource, with specified conditions.
type Permission struct {
	Action     string     `json:"action" yaml:"action"`
	Conditions Conditions `json:"conditions" yaml:"conditions"`
}

// marshalablePermission - helper type for handling marshaling/unmarshaling
// of JSON and YAML structures.
type marshalablePermission struct {
	Action     string     `json:"action" yaml:"action"`
	Conditions Conditions `json:"conditions" yaml:"conditions"`
}

// UnmarshalJSON - implementation of json.Marshaler.UnmarshalJSON, for creating Permission
// object from JSON. It takes care of initializating empty Conditions map.
func (p *Permission) UnmarshalJSON(jsonData []byte) error {
	var perm = marshalablePermission{
		Conditions: Conditions{},
	}

	if err := json.Unmarshal(jsonData, &perm); err != nil {
		return err
	}

	p.Action = perm.Action
	p.Conditions = perm.Conditions

	return nil
}

// UnmarshalJSON - implementation of yaml.Marshaler.UnmarshalJSON, for creating Permission
// object from YAML. It takes care of initializating empty Conditions map.
func (p *Permission) UnmarshalYAML(value *yaml.Node) error {
	var perm = marshalablePermission{
		Conditions: Conditions{},
	}

	if err := value.Decode(&perm); err != nil {
		return err
	}

	p.Action = perm.Action
	p.Conditions = perm.Conditions

	return nil
}
