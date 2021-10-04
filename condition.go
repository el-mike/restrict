package restrict

import (
	"encoding/json"
	"errors"
	"fmt"

	"gopkg.in/yaml.v3"
)

// Condition - additional requirement that needs to be satisfied
// to grant given permission.
type Condition interface {
	// Name - returns Condition's name, which is it's unique identifier.
	Name() string

	// Check - returns true if Condition is satisfied by
	// given request, false otherwise.
	Check(interface{}, *AccessRequest) bool
}

// Conditions - alias type for Conditions map.
type Conditions map[string]Condition

// jsonMarshalableCondition - helper type for handling marshaling/unmarshaling
// of JSON structures.
type jsonMarshalableCondition struct {
	Name    string          `json:"name"`
	Options json.RawMessage `json:"options"`
}

// yamlMarshalableCondition - helper type for handling marshaling/unmarshaling
// of YAML structures.
type yamlMarshalableCondition struct {
	Name    string    `yaml:"name"`
	Options yaml.Node `yaml:"options"`
}

// MarshalJSON - marshals a map of Conditions to JSON data.
func (cs Conditions) MarshalJSON() ([]byte, error) {
	result := make(map[string]*jsonMarshalableCondition, len(cs))

	for key, condition := range cs {
		options, err := json.Marshal(condition)
		if err != nil {
			return []byte{}, err
		}

		result[key] = &jsonMarshalableCondition{
			Name:    condition.Name(),
			Options: json.RawMessage(options),
		}
	}

	return json.Marshal(result)
}

// MarshalYAML - marshals a map of Conditions to YAML data.
func (cs Conditions) MarshalYAML() (interface{}, error) {
	result := make(map[string]*yamlMarshalableCondition, len(cs))

	for key, condition := range cs {
		options := yaml.Node{}

		if err := options.Encode(condition); err != nil {
			return nil, err
		}

		result[key] = &yamlMarshalableCondition{
			Name:    condition.Name(),
			Options: options,
		}
	}

	output := yaml.Node{}

	if err := output.Encode(result); err != nil {
		return nil, err
	}

	return output, nil
}

// UnmarshalJSON - unmarshals a JSON-coded map of Conditions.
func (cs Conditions) UnmarshalJSON(jsonData []byte) error {
	if cs == nil {
		return errors.New("Cannot unmarshal nil value")
	}

	var jsonValue map[string]jsonMarshalableCondition

	if err := json.Unmarshal(jsonData, &jsonValue); err != nil {
		return err
	}

	for key, jsonCondition := range jsonValue {
		factory := ConditionFactories[jsonCondition.Name]

		if factory == nil {
			return fmt.Errorf("No factory found for Condition: %v", jsonCondition.Name)
		}

		condition := factory()

		if len(jsonCondition.Options) > 0 {
			if err := json.Unmarshal(jsonCondition.Options, condition); err != nil {
				return err
			}
		}

		cs[key] = condition
	}

	return nil
}

// UnmarshalYAML - unmarshals a YAML-coded map of Conditions.
func (cs Conditions) UnmarshalYAML(value *yaml.Node) error {
	if cs == nil {
		return errors.New("Cannot unmarshal nil value")
	}

	var yamlValue map[string]yamlMarshalableCondition

	if err := value.Decode(&yamlValue); err != nil {
		return err
	}

	for key, yamlCondition := range yamlValue {
		// Guard for conditions maps being empty - YAML will still
		// create a Nodes from them.
		if yamlCondition.Name == "" {
			continue
		}

		factory := ConditionFactories[yamlCondition.Name]

		if factory == nil {
			return fmt.Errorf("No factory found for Condition: %v", yamlCondition.Name)
		}

		condition := factory()

		if len(yamlCondition.Options.Content) > 0 {
			if err := yamlCondition.Options.Decode(condition); err != nil {
				return err
			}
		}

		cs[key] = condition
	}

	return nil
}

// ConditionFactory - factory function for Conditions.
type ConditionFactory func() Condition

// ConditionFatoriesMap - map of Condition factories.
type ConditionFatoriesMap = map[string]ConditionFactory

// ConditionFactories - stores a map of functions responsible for
// creating new Conditions, based on it's names.
var ConditionFactories = ConditionFatoriesMap{
	IS_EQUAL_CONDITION_NAME: func() Condition {
		return new(IsEqualCondition)
	},
}
