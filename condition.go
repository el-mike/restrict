package restrict

import (
	"encoding/json"
	"errors"

	"gopkg.in/yaml.v3"
)

// Condition - additional requirement that needs to be satisfied
// to grant given permission.
type Condition interface {
	// Name - returns Condition's name, which is it's unique identifier.
	Name() string

	// Check - returns true if Condition is satisfied by
	// given request, false otherwise.
	Check(request *AccessRequest) error
}

// Conditions - alias type for Conditions array.
type Conditions []Condition

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
	result := []*jsonMarshalableCondition{}

	for _, condition := range cs {
		options, err := json.Marshal(condition)
		if err != nil {
			return []byte{}, err
		}

		result = append(result, &jsonMarshalableCondition{
			Name:    condition.Name(),
			Options: json.RawMessage(options),
		})
	}

	return json.Marshal(result)
}

// MarshalYAML - marshals a map of Conditions to YAML data.
func (cs Conditions) MarshalYAML() (interface{}, error) {
	result := []*yamlMarshalableCondition{}

	for _, condition := range cs {
		options := yaml.Node{}

		if err := options.Encode(condition); err != nil {
			return nil, err
		}

		result = append(result, &yamlMarshalableCondition{
			Name:    condition.Name(),
			Options: options,
		})
	}

	output := yaml.Node{}

	if err := output.Encode(result); err != nil {
		return nil, err
	}

	return output, nil
}

// UnmarshalJSON - unmarshals a JSON-coded map of Conditions.
func (cs *Conditions) UnmarshalJSON(jsonData []byte) error {
	if cs == nil {
		return errors.New("Cannot unmarshal nil value")
	}

	var jsonValue []jsonMarshalableCondition

	if err := json.Unmarshal(jsonData, &jsonValue); err != nil {
		return err
	}

	for _, jsonCondition := range jsonValue {
		factory := ConditionFactories[jsonCondition.Name]

		if factory == nil {
			return NewConditionFactoryNotFoundError(jsonCondition.Name)
		}

		condition := factory()

		if len(jsonCondition.Options) > 0 {
			if err := json.Unmarshal(jsonCondition.Options, condition); err != nil {
				return err
			}
		}

		*cs = append(*cs, condition)
	}

	return nil
}

// UnmarshalYAML - unmarshals a YAML-coded map of Conditions.
func (cs *Conditions) UnmarshalYAML(value *yaml.Node) error {
	if cs == nil {
		return errors.New("Cannot unmarshal nil value")
	}

	var yamlValue []yamlMarshalableCondition

	if err := value.Decode(&yamlValue); err != nil {
		return err
	}

	for _, yamlCondition := range yamlValue {
		// Guard for conditions maps being empty - YAML will still
		// create a Nodes from them.
		if yamlCondition.Name == "" {
			continue
		}

		factory := ConditionFactories[yamlCondition.Name]

		if factory == nil {
			return NewConditionFactoryNotFoundError(yamlCondition.Name)
		}

		condition := factory()

		if len(yamlCondition.Options.Content) > 0 {
			if err := yamlCondition.Options.Decode(condition); err != nil {
				return err
			}
		}

		*cs = append(*cs, condition)
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
	IsEqualConditionName: func() Condition {
		return new(IsEqualCondition)
	},
	IsOwnerConditionName: func() Condition {
		return new(IsOwnerCondition)
	},
}

// RegisterConditionFactory - adds a new ConditionFactory under given name. If given name
// is already taken, an error is returned.
func RegisterConditionFactory(name string, factory ConditionFactory) error {
	if ConditionFactories[name] != nil {
		return NewConditionFactoryAlreadyExistsError(name)
	}

	ConditionFactories[name] = factory
	return nil
}
