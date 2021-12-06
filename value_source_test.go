package restrict

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gopkg.in/yaml.v3"
)

type valueSourcesWrapper struct {
	Subject  ValueSource `json:"subject" yaml:"subject"`
	Resource ValueSource `json:"resource" yaml:"resource"`
	Context  ValueSource `json:"context" yaml:"context"`
	Explicit ValueSource `json:"explicit" yaml:"explicit"`
}

type valueSourceSuiteMock struct {
	suite.Suite
}

func TestValueSourceSuite(t *testing.T) {
	suite.Run(t, new(valueSourceSuiteMock))
}

func (s *valueSourceSuiteMock) TestValueSourceString() {
	assert.Equal(s.T(), "SubjectField", SubjectField.String())
	assert.Equal(s.T(), "ResourceField", ResourceField.String())
	assert.Equal(s.T(), "ContextField", ContextField.String())
	assert.Equal(s.T(), "Explicit", Explicit.String())
	assert.Equal(s.T(), "", noopValueSource.String())
}

func (s *valueSourceSuiteMock) TestUnmarshalJSON() {
	valueSourceData := []byte(`{
		"subject": "SubjectField",
		"resource": "ResourceField",
		"context": "ContextField",
		"explicit": "Explicit"
	}`)

	assert.True(s.T(), json.Valid(valueSourceData))

	testValueSources := &valueSourcesWrapper{}

	err := json.Unmarshal(valueSourceData, testValueSources)

	assert.Nil(s.T(), err)

	s.assertValueSourcesWrapper(testValueSources)
}

func (s *valueSourceSuiteMock) TestUnmarshalYAML() {
	valueSourceData := []byte(`
subject: "SubjectField"
resource: "ResourceField"
context: "ContextField"
explicit: "Explicit"
`)

	testValueSources := &valueSourcesWrapper{}

	err := yaml.Unmarshal(valueSourceData, testValueSources)

	assert.Nil(s.T(), err)

	s.assertValueSourcesWrapper(testValueSources)
}

func (s *valueSourceSuiteMock) TestMarshalJSON() {
	testValueSources := &valueSourcesWrapper{
		Subject:  SubjectField,
		Resource: ResourceField,
		Context:  ContextField,
		Explicit: Explicit,
	}

	valueSourcesJSON, err := json.Marshal(testValueSources)

	assert.Nil(s.T(), err)
	assert.True(s.T(), json.Valid(valueSourcesJSON))

	testValueSources = &valueSourcesWrapper{}

	err = json.Unmarshal(valueSourcesJSON, testValueSources)

	assert.Nil(s.T(), err)

	s.assertValueSourcesWrapper(testValueSources)
}

func (s *valueSourceSuiteMock) TestMarshalYAML() {
	testValueSources := &valueSourcesWrapper{
		Subject:  SubjectField,
		Resource: ResourceField,
		Context:  ContextField,
		Explicit: Explicit,
	}

	valueSourcesYAML, err := yaml.Marshal(testValueSources)

	assert.Nil(s.T(), err)

	testValueSources = &valueSourcesWrapper{}

	err = yaml.Unmarshal(valueSourcesYAML, testValueSources)

	assert.Nil(s.T(), err)

	s.assertValueSourcesWrapper(testValueSources)
}

func (s *valueSourceSuiteMock) assertValueSourcesWrapper(testValueSources *valueSourcesWrapper) {
	assert.Equal(s.T(), testValueSources.Subject, SubjectField)
	assert.Equal(s.T(), testValueSources.Resource, ResourceField)
	assert.Equal(s.T(), testValueSources.Context, ContextField)
	assert.Equal(s.T(), testValueSources.Explicit, Explicit)
}
