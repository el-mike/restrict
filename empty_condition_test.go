package restrict

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type emptyConditionSuite struct {
	suite.Suite
}

func TestEmptyConditionSuite(t *testing.T) {
	suite.Run(t, new(emptyConditionSuite))
}

func (s *emptyConditionSuite) TestType_Empty() {
	testCondition := &EmptyCondition{}

	assert.Equal(s.T(), EmptyConditionType, testCondition.Type())
}

func (s *emptyConditionSuite) TestCheck_Empty() {
	testSubject := new(subjectMock)

	testSubject.FieldOne = "testValue"
	testSubject.FieldTwo = 1
	testSubject.FieldThree = nil

	testRequest := &AccessRequest{
		Subject: testSubject,
	}

	// Failing descriptors - missing field
	testCondition := &EmptyCondition{
		Value: &ValueDescriptor{Source: SubjectField},
	}

	err := testCondition.Check(testRequest)

	assert.Error(s.T(), err)

	// Not empty
	testCondition.Value = &ValueDescriptor{
		Source: SubjectField,
		Field:  "FieldOne",
	}

	err = testCondition.Check(testRequest)

	assert.IsType(s.T(), new(ConditionNotSatisfiedError), err)

	// Missing value field
	testCondition.Value = &ValueDescriptor{
		Source: SubjectField,
		Field:  "FieldThree",
	}

	err = testCondition.Check(testRequest)

	assert.Nil(s.T(), err)

	// 0 string value
	testSubject.FieldOne = ""
	testCondition.Value = &ValueDescriptor{
		Source: SubjectField,
		Field:  "FieldOne",
	}

	err = testCondition.Check(testRequest)

	assert.Nil(s.T(), err)

	// 0 int value
	testSubject.FieldTwo = 0
	testCondition.Value = &ValueDescriptor{
		Source: SubjectField,
		Field:  "FieldTwo",
	}

	err = testCondition.Check(testRequest)

	assert.Nil(s.T(), err)
}

func (s *emptyConditionSuite) TestType_NotEmpty() {
	testCondition := &NotEmptyCondition{}

	assert.Equal(s.T(), NotEmptyConditionType, testCondition.Type())
}

func (s *emptyConditionSuite) TestCheck_NotEmpty() {
	testSubject := new(subjectMock)

	testSubject.FieldOne = "testValue"
	testSubject.FieldTwo = 1
	testSubject.FieldThree = nil

	testRequest := &AccessRequest{
		Subject: testSubject,
	}

	// Failing descriptors - missing field
	testCondition := &NotEmptyCondition{
		Value: &ValueDescriptor{Source: SubjectField},
	}

	err := testCondition.Check(testRequest)

	assert.Error(s.T(), err)

	// Empty
	testSubject.FieldOne = ""
	testCondition.Value = &ValueDescriptor{
		Source: SubjectField,
		Field:  "FieldOne",
	}

	err = testCondition.Check(testRequest)

	assert.IsType(s.T(), new(ConditionNotSatisfiedError), err)

	// Not empty
	testSubject.FieldOne = "testValue"
	testCondition.Value = &ValueDescriptor{
		Source: SubjectField,
		Field:  "FieldOne",
	}

	err = testCondition.Check(testRequest)

	assert.Nil(s.T(), err)
}
