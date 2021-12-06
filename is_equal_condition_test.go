package restrict

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type isEqualConditionSuite struct {
	suite.Suite
}

func TestIsEqualConditionSuite(t *testing.T) {
	suite.Run(t, new(isEqualConditionSuite))
}

func (s *isEqualConditionSuite) TestType() {
	testCondition := &IsEqualCondition{}

	assert.Equal(s.T(), IsEqualConditionType, testCondition.Type())
}

func (s *isEqualConditionSuite) TestCheck() {
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
	testCondition := &IsEqualCondition{
		Value:  &ValueDescriptor{Source: SubjectField},
		Equals: &ValueDescriptor{Source: ResourceField},
	}

	err := testCondition.Check(testRequest)

	assert.Error(s.T(), err)

	testCondition.Value = &ValueDescriptor{
		Source: SubjectField,
		Field:  "FieldOne",
	}

	err = testCondition.Check(testRequest)

	assert.Error(s.T(), err)

	testCondition.Equals = &ValueDescriptor{
		Source: ResourceField,
		Field:  "FieldOne",
	}

	err = testCondition.Check(testRequest)

	assert.Nil(s.T(), err)

	// Error - different types
	testCondition.Equals = &ValueDescriptor{
		Source: ResourceField,
		Field:  "FieldTwo",
	}

	err = testCondition.Check(testRequest)

	assert.IsType(s.T(), new(ConditionNotSatisfiedError), err)

	// Error - different values
	testCondition.Value = &ValueDescriptor{
		Source: SubjectField,
		Field:  "FieldTwo",
	}

	err = testCondition.Check(testRequest)

	assert.IsType(s.T(), new(ConditionNotSatisfiedError), err)
}
