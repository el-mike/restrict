package restrict

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type valueDescriptorSuite struct {
	suite.Suite
}

func TestValueDescriptorSuite(t *testing.T) {
	suite.Run(t, new(valueDescriptorSuite))
}

func (s *valueDescriptorSuite) TestGetValue_Explicit() {
	testRequest := &AccessRequest{}

	testDescriptor := &ValueDescriptor{
		Source: Explicit,
	}

	value, err := testDescriptor.GetValue(testRequest)

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), nil, value)

	testDescriptor.Value = 1

	value, err = testDescriptor.GetValue(testRequest)

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 1, value)
}

func (s *valueDescriptorSuite) TestGetValue_Subject() {
	testSubject := new(subjectMock)

	testSubject.FieldOne = "testValue"
	testSubject.FieldThree = []int{}

	testRequest := &AccessRequest{
		Subject: testSubject,
	}

	testDescriptor := &ValueDescriptor{
		Source: SubjectField,
	}

	value, err := testDescriptor.GetValue(testRequest)

	assert.IsType(s.T(), new(ValueDescriptorMalformedError), err)
	assert.Nil(s.T(), value)

	testDescriptor.Field = "IncorrectField"

	value, err = testDescriptor.GetValue(testRequest)

	assert.IsType(s.T(), new(ValueDescriptorMalformedError), err)
	assert.Nil(s.T(), value)

	testDescriptor.Field = "FieldOne"

	value, err = testDescriptor.GetValue(testRequest)

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), testSubject.FieldOne, value)

	testDescriptor.Field = "FieldThree"

	value, err = testDescriptor.GetValue(testRequest)

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), testSubject.FieldThree, value)
}

func (s *valueDescriptorSuite) TestGetValue_Resource() {
	testResource := new(resourceMock)

	testResource.FieldOne = "testValue"
	testResource.FieldThree = []int{}

	testRequest := &AccessRequest{
		Resource: testResource,
	}

	testDescriptor := &ValueDescriptor{
		Source: ResourceField,
		Field:  "FieldOne",
	}

	value, err := testDescriptor.GetValue(testRequest)

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), testResource.FieldOne, value)
}

func (s *valueDescriptorSuite) TestGetValue_Context() {
	testContext := Context{
		"FieldOne": "testValue",
		"FieldTwo": 2,
	}

	testRequest := &AccessRequest{
		Context: testContext,
	}

	testDescriptor := &ValueDescriptor{
		Source: ContextField,
		Field:  "FieldOne",
	}

	value, err := testDescriptor.GetValue(testRequest)

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), testContext["FieldOne"], value)

	testDescriptor.Field = "FieldTwo"

	value, err = testDescriptor.GetValue(testRequest)

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), testContext["FieldTwo"], value)
}

func (s *valueDescriptorSuite) TestGetValue_MissingSource() {
	testDescriptor := &ValueDescriptor{
		Source: noopValueSource,
		Field:  "TestField",
	}

	value, err := testDescriptor.GetValue(&AccessRequest{})

	assert.Nil(s.T(), value)
	assert.IsType(s.T(), new(ValueDescriptorMalformedError), err)
}
