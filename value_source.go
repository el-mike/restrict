package restrict

import (
	"bytes"
	"encoding/json"

	"gopkg.in/yaml.v3"
)

// ValueSource - enum type for source of value for given ValueDescriptor.
type ValueSource int

const (
	// NoopValueSource - zero value for ValueSource, useful when marshaling/unmarshaling.
	noopValueSource ValueSource = iota
	// SubjectField - value that comes from Subject's field.
	SubjectField
	// ResourceField - value taht comes from Resource's field.
	ResourceField
	// ContextField - value that comes from Context's field.
	ContextField
	// Explicit - value set explicitly in PolicyDefinition.
	Explicit
)

var byValue = map[ValueSource]string{
	SubjectField:  "SubjectField",
	ResourceField: "ResourceField",
	ContextField:  "ContextField",
	Explicit:      "Explicit",
}

var byName = map[string]ValueSource{
	"SubjectField":  SubjectField,
	"ResourceField": ResourceField,
	"ContextField":  ContextField,
	"Explicit":      Explicit,
}

// String - Stringer implementation.
func (vs ValueSource) String() string {
	return byValue[vs]
}

// MarshalJSON - marshals a ValueSource enum into its name as string.
func (vs ValueSource) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)

	buffer.WriteString(byValue[vs])
	buffer.WriteString(`"`)

	return buffer.Bytes(), nil
}

// MarshalYAML - marshals a ValueSource enum into its name as string.
func (vs ValueSource) MarshalYAML() (interface{}, error) {
	return byValue[vs], nil
}

// UnmarshalJSON - unmarshals a string into ValueSource.
func (vs *ValueSource) UnmarshalJSON(jsonData []byte) error {
	var sourceName string

	if err := json.Unmarshal(jsonData, &sourceName); err != nil {
		return err
	}

	*vs = byName[sourceName]

	return nil
}

// UnmarshalYAML - unmarshals a string into ValueSource.
func (vs *ValueSource) UnmarshalYAML(value *yaml.Node) error {
	var sourceName string

	if err := value.Decode(&sourceName); err != nil {
		return err
	}

	*vs = byName[sourceName]

	return nil
}
