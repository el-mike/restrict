package restrict

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type equalConditionSuite struct {
	suite.Suite
}

func TestEqualConditionSuite(t *testing.T) {
	suite.Run(t, new(equalConditionSuite))
}

func (s *equalConditionSuite) TestType_Equal() {
	testCondition := &EqualCondition{}

	assert.Equal(s.T(), EqualConditionType, testCondition.Type())
}

func (s *equalConditionSuite) TestType_NotEqual() {
	testCondition := &NotEqualCondition{}

	assert.Equal(s.T(), NotEqualConditionType, testCondition.Type())
}

func (s *equalConditionSuite) TestCheck_Equal() {
	testSubject := new(subjectMock)

	testSubject.FieldOne = "testValue"
	testSubject.FieldTwo = 1

	testResource := new(resourceMock)

	testResource.FieldOne = "testValue"
	testResource.FieldTwo = 2

	testRequest := &AccessRequest{
		Subject:  testSubject,
		Resource: testResource,
	}

	// Failing descriptors - missing field
	testCondition := &EqualCondition{
		Left:  &ValueDescriptor{Source: SubjectField},
		Right: &ValueDescriptor{Source: ResourceField},
	}

	err := testCondition.Check(testRequest)

	assert.Error(s.T(), err)

	testCondition.Left = &ValueDescriptor{
		Source: SubjectField,
		Field:  "FieldOne",
	}

	err = testCondition.Check(testRequest)

	assert.Error(s.T(), err)

	testCondition.Right = &ValueDescriptor{
		Source: ResourceField,
		Field:  "FieldOne",
	}

	err = testCondition.Check(testRequest)

	assert.Nil(s.T(), err)

	// Error - different types
	testCondition.Right = &ValueDescriptor{
		Source: ResourceField,
		Field:  "FieldTwo",
	}

	err = testCondition.Check(testRequest)

	assert.IsType(s.T(), new(ConditionNotSatisfiedError), err)

	// Error - different values
	testCondition.Left = &ValueDescriptor{
		Source: SubjectField,
		Field:  "FieldTwo",
	}

	err = testCondition.Check(testRequest)

	assert.IsType(s.T(), new(ConditionNotSatisfiedError), err)
}

func (s *equalConditionSuite) TestCheck_NotEqual() {
	testSubject := new(subjectMock)

	testSubject.FieldOne = "testValue"
	testSubject.FieldTwo = 1

	testResource := new(resourceMock)

	testResource.FieldOne = "testValue"
	testResource.FieldTwo = 2

	testRequest := &AccessRequest{
		Subject:  testSubject,
		Resource: testResource,
	}

	// Failing descriptors - missing field
	testCondition := &NotEqualCondition{
		Left:  &ValueDescriptor{Source: SubjectField},
		Right: &ValueDescriptor{Source: ResourceField},
	}

	err := testCondition.Check(testRequest)

	assert.Error(s.T(), err)

	testCondition.Left = &ValueDescriptor{
		Source: SubjectField,
		Field:  "FieldOne",
	}

	err = testCondition.Check(testRequest)

	assert.Error(s.T(), err)

	testCondition.Right = &ValueDescriptor{
		Source: ResourceField,
		Field:  "FieldOne",
	}

	// Values are the same
	err = testCondition.Check(testRequest)

	assert.IsType(s.T(), new(ConditionNotSatisfiedError), err)

	// Satisfied - different types
	testCondition.Right = &ValueDescriptor{
		Source: ResourceField,
		Field:  "FieldTwo",
	}

	err = testCondition.Check(testRequest)

	assert.Nil(s.T(), err)

	// Satisfied - different values
	testCondition.Left = &ValueDescriptor{
		Source: SubjectField,
		Field:  "FieldTwo",
	}

	err = testCondition.Check(testRequest)

	assert.Nil(s.T(), err)
}
