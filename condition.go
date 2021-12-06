package restrict

import (
	"encoding/json"

	"gopkg.in/yaml.v3"
)

// Condition - additional requirement that needs to be satisfied
// to grant given permission.
type Condition interface {
	// Type - returns Condition's type.
	Type() string

	// Check - returns true if Condition is satisfied by
	// given request, false otherwise.
	Check(request *AccessRequest) error
}

// Conditions - alias type for Conditions array.
type Conditions []Condition

// AppendCondition - adds new Condition to Conditions slice.
func (cs *Conditions) appendCondition(condition Condition) {
	if cs == nil {
		cs = &Conditions{}
	}

	*cs = append(*cs, condition)
}

// jsonMarshalableCondition - helper type for handling marshaling/unmarshaling
// of JSON structures.
type jsonMarshalableCondition struct {
	Type    string          `json:"type"`
	Options json.RawMessage `json:"options,omitempty"`
}

// yamlMarshalableCondition - helper type for handling marshaling/unmarshaling
// of YAML structures.
type yamlMarshalableCondition struct {
	Type    string    `yaml:"type"`
	Options yaml.Node `yaml:"options,omitempty"`
}

// MarshalJSON - marshals a map of Conditions to JSON data.
func (cs Conditions) MarshalJSON() ([]byte, error) {
	result := []*jsonMarshalableCondition{}

	for _, condition := range cs {
		options, err := json.Marshal(condition)
		if err != nil {
			return nil, err
		}

		result = append(result, &jsonMarshalableCondition{
			Type:    condition.Type(),
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
			Type:    condition.Type(),
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
	var jsonValue []jsonMarshalableCondition

	if err := json.Unmarshal(jsonData, &jsonValue); err != nil {
		return err
	}

	for _, jsonCondition := range jsonValue {
		factory := ConditionFactories[jsonCondition.Type]

		if factory == nil {
			return newConditionFactoryNotFoundError(jsonCondition.Type)
		}

		condition := factory()

		if len(jsonCondition.Options) > 0 {
			if err := json.Unmarshal(jsonCondition.Options, condition); err != nil {
				return err
			}
		}

		cs.appendCondition(condition)
	}

	return nil
}

// UnmarshalYAML - unmarshals a YAML-coded map of Conditions.
func (cs *Conditions) UnmarshalYAML(value *yaml.Node) error {
	var yamlValue []yamlMarshalableCondition

	if err := value.Decode(&yamlValue); err != nil {
		return err
	}

	for _, yamlCondition := range yamlValue {
		// Guard for conditions maps being empty - YAML will still
		// create a Nodes from them.
		if yamlCondition.Type == "" {
			continue
		}

		factory := ConditionFactories[yamlCondition.Type]

		if factory == nil {
			return newConditionFactoryNotFoundError(yamlCondition.Type)
		}

		condition := factory()

		if len(yamlCondition.Options.Content) > 0 {
			if err := yamlCondition.Options.Decode(condition); err != nil {
				return err
			}
		}

		cs.appendCondition(condition)
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
	IsEqualConditionType: func() Condition {
		return new(IsEqualCondition)
	},
}

// RegisterConditionFactory - adds a new ConditionFactory under given name. If given name
// is already taken, an error is returned.
func RegisterConditionFactory(name string, factory ConditionFactory) error {
	if ConditionFactories[name] != nil {
		return newConditionFactoryAlreadyExistsError(name)
	}

	ConditionFactories[name] = factory
	return nil
}
