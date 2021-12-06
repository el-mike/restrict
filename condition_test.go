package restrict

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gopkg.in/yaml.v3"
)

const marshalableConditionMockName = "TEST_CONDITION"
const invalidMarshalableConditionMockName = "INVALID_TEST_CONDITION"

type marshalableConditionMock struct {
	TestPropertyOne int    `json:"testPropertyOne" yaml:"testPropertyOne"`
	TestPropertyTwo string `json:"testPropertyTwo" yaml:"testPropertyTwo"`
}

func (m *marshalableConditionMock) Type() string {
	return marshalableConditionMockName
}

func (m *marshalableConditionMock) Check(request *AccessRequest) error {
	return nil
}

type invalidMarshalableConditionMock struct {
	TestProperty int `json:"testProperty" yaml:"testProperty"`
}

func (m *invalidMarshalableConditionMock) Type() string {
	return invalidMarshalableConditionMockName
}

func (m *invalidMarshalableConditionMock) Check(request *AccessRequest) error {
	return nil
}

func (m *invalidMarshalableConditionMock) MarshalJSON() ([]byte, error) {
	return nil, errors.New("testError")
}

func (m *invalidMarshalableConditionMock) MarshalYAML() (interface{}, error) {
	return nil, errors.New("testError")
}

func (m *invalidMarshalableConditionMock) UnarshalJSON(jsonDate []byte) error {
	return errors.New("testError")
}

func (m *invalidMarshalableConditionMock) UnmarshalYAML(value *yaml.Node) error {
	return errors.New("testError")
}

type conditionsSuite struct {
	suite.Suite
}

func resetFactories() {
	ConditionFactories[marshalableConditionMockName] = nil
	ConditionFactories[invalidMarshalableConditionMockName] = nil
}

func TestConditionsSuite(t *testing.T) {
	suite.Run(t, new(conditionsSuite))
}

func (s *conditionsSuite) BeforeTest(_, _ string) {
	//nolint
	RegisterConditionFactory(marshalableConditionMockName, func() Condition {
		return &marshalableConditionMock{}
	})

	//nolint
	RegisterConditionFactory(invalidMarshalableConditionMockName, func() Condition {
		return &invalidMarshalableConditionMock{}
	})
}

func (s *conditionsSuite) AfterTest(_, _ string) {
	resetFactories()
}

func (s *conditionsSuite) TestRegisterConditionFactory() {
	resetFactories()

	factory := func() Condition {
		return &marshalableConditionMock{}
	}

	err := RegisterConditionFactory(marshalableConditionMockName, factory)

	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), ConditionFactories[marshalableConditionMockName])

	err = RegisterConditionFactory(marshalableConditionMockName, factory)

	assert.IsType(s.T(), new(ConditionFactoryAlreadyExistsError), err)
}

func (s *conditionsSuite) TestUnmarshalJSON() {
	// Since ConditionsFactories is global, we make sure that the mocked Conditions's
	// factory is nil.
	ConditionFactories[marshalableConditionMockName] = nil

	conditionsData := []byte(`[
		{
			"type": "TEST_CONDITION",
			"options": {
				"testPropertyOne": 1,
				"testPropertyTwo": "testString1"
			}
		},
		{
			"type": "TEST_CONDITION",
			"options": {
				"testPropertyOne": 2,
				"testPropertyTwo": "testString2"
			}
		}
	]`)

	assert.True(s.T(), json.Valid(conditionsData))

	testConditions := Conditions{}

	err := testConditions.UnmarshalJSON(conditionsData)

	assert.IsType(s.T(), new(ConditionFactoryNotFoundError), err)

	//nolint
	RegisterConditionFactory(marshalableConditionMockName, func() Condition {
		return &marshalableConditionMock{}
	})

	testConditions = Conditions{}

	err = testConditions.UnmarshalJSON(conditionsData)

	assert.Nil(s.T(), err)

	assert.IsType(s.T(), new(marshalableConditionMock), testConditions[0])
	assert.IsType(s.T(), new(marshalableConditionMock), testConditions[1])

	testConditionOne := testConditions[0].(*marshalableConditionMock)
	testConditionTwo := testConditions[1].(*marshalableConditionMock)

	assert.Equal(s.T(), 1, testConditionOne.TestPropertyOne)
	assert.Equal(s.T(), "testString1", testConditionOne.TestPropertyTwo)

	assert.Equal(s.T(), 2, testConditionTwo.TestPropertyOne)
	assert.Equal(s.T(), "testString2", testConditionTwo.TestPropertyTwo)
}

func (s *conditionsSuite) TestUnmarshalJSON_InvalidData() {
	//nolint
	RegisterConditionFactory(marshalableConditionMockName, func() Condition {
		return &marshalableConditionMock{}
	})

	invalidConditionsData := []byte(`[
		{
			type": "TEST_CONDITION"
		}
	]`)

	testConditions := Conditions{}

	err := testConditions.UnmarshalJSON(invalidConditionsData)

	assert.Error(s.T(), err)

	// "testPropertyOne" has string instead of int.
	invalidConditionsData = []byte(`[
		{
			"type": "TEST_CONDITION",
			"options": {
				"testPropertyOne": "2"
			}
		}
	]`)

	err = json.Unmarshal(invalidConditionsData, &testConditions)

	assert.Error(s.T(), err)
}

func (s *conditionsSuite) TestUnmarshalYAML() {
	// Since ConditionsFactories is global, we make sure that the mocked Conditions's
	// factory is nil.
	ConditionFactories[marshalableConditionMockName] = nil

	// Note that one of the Conditions has empty type - UnmarshalYAML should omit
	// this Condition.
	conditionsData := []byte(`
- type: TEST_CONDITION
  options:
    testPropertyOne: 1
    testPropertyTwo: "testString1"

- type: TEST_CONDITION
  options:
    testPropertyOne: 2
    testPropertyTwo: "testString2"

- type:
  options:
    testPropertyOne: 2
    testPropertyTwo: "testString2"
`)

	testConditions := Conditions{}

	err := yaml.Unmarshal(conditionsData, &testConditions)

	assert.IsType(s.T(), new(ConditionFactoryNotFoundError), err)

	//nolint
	RegisterConditionFactory(marshalableConditionMockName, func() Condition {
		return &marshalableConditionMock{}
	})

	err = yaml.Unmarshal(conditionsData, &testConditions)

	assert.Nil(s.T(), err)

	assert.IsType(s.T(), new(marshalableConditionMock), testConditions[0])
	assert.IsType(s.T(), new(marshalableConditionMock), testConditions[1])

	testConditionOne := testConditions[0].(*marshalableConditionMock)
	testConditionTwo := testConditions[1].(*marshalableConditionMock)

	assert.Equal(s.T(), 1, testConditionOne.TestPropertyOne)
	assert.Equal(s.T(), "testString1", testConditionOne.TestPropertyTwo)

	assert.Equal(s.T(), 2, testConditionTwo.TestPropertyOne)
	assert.Equal(s.T(), "testString2", testConditionTwo.TestPropertyTwo)
}

func (s *conditionsSuite) TestUnmarshalYAML_InvalidData() {
	// Missing hyphen at the beginning
	conditionsData := []byte(`
- type: TEST_CONDITION
  options:
    testPropertyOne: "1"
    testPropertyTwo: "testString1"
`)

	testConditions := Conditions{}
	err := yaml.Unmarshal(conditionsData, &testConditions)

	assert.Error(s.T(), err)
}

func (s *conditionsSuite) TestMarshalJSON() {
	//nolint
	RegisterConditionFactory(marshalableConditionMockName, func() Condition {
		return &marshalableConditionMock{}
	})

	testConditionOne := &marshalableConditionMock{
		TestPropertyOne: 1,
		TestPropertyTwo: "testString1",
	}

	testConditionTwo := &marshalableConditionMock{
		TestPropertyOne: 2,
		TestPropertyTwo: "testString2",
	}

	testConditions := Conditions{
		testConditionOne,
		testConditionTwo,
	}

	conditionsJSON, err := testConditions.MarshalJSON()

	assert.Nil(s.T(), err)
	assert.True(s.T(), json.Valid(conditionsJSON))

	testConditions = Conditions{}

	err = testConditions.UnmarshalJSON(conditionsJSON)

	assert.Nil(s.T(), err)

	testConditionOne = testConditions[0].(*marshalableConditionMock)
	testConditionTwo = testConditions[1].(*marshalableConditionMock)

	assert.Equal(s.T(), 1, testConditionOne.TestPropertyOne)
	assert.Equal(s.T(), "testString1", testConditionOne.TestPropertyTwo)

	assert.Equal(s.T(), 2, testConditionTwo.TestPropertyOne)
	assert.Equal(s.T(), "testString2", testConditionTwo.TestPropertyTwo)
}

func (s *conditionsSuite) TestMarshalJSON_InvalidCondition() {
	testCondition := &invalidMarshalableConditionMock{
		TestProperty: 1,
	}

	testConditions := Conditions{
		testCondition,
	}

	conditionsJSON, err := testConditions.MarshalJSON()

	assert.Nil(s.T(), conditionsJSON)
	assert.Error(s.T(), err)
}

func (s *conditionsSuite) TestMarshalYAML() {
	//nolint
	RegisterConditionFactory(marshalableConditionMockName, func() Condition {
		return &marshalableConditionMock{}
	})

	testConditionOne := &marshalableConditionMock{
		TestPropertyOne: 1,
		TestPropertyTwo: "testString1",
	}

	testConditionTwo := &marshalableConditionMock{
		TestPropertyOne: 2,
		TestPropertyTwo: "testString2",
	}

	testConditions := Conditions{
		testConditionOne,
		testConditionTwo,
	}

	conditionsYAML, err := yaml.Marshal(testConditions)

	assert.Nil(s.T(), err)

	testConditions = Conditions{}

	err = yaml.Unmarshal(conditionsYAML, &testConditions)

	assert.Nil(s.T(), err)

	testConditionOne = testConditions[0].(*marshalableConditionMock)
	testConditionTwo = testConditions[1].(*marshalableConditionMock)

	assert.Equal(s.T(), 1, testConditionOne.TestPropertyOne)
	assert.Equal(s.T(), "testString1", testConditionOne.TestPropertyTwo)

	assert.Equal(s.T(), 2, testConditionTwo.TestPropertyOne)
	assert.Equal(s.T(), "testString2", testConditionTwo.TestPropertyTwo)
}

func (s *conditionsSuite) TestMarshalYAML_InvalidCondition() {
	testCondition := &invalidMarshalableConditionMock{
		TestProperty: 1,
	}

	testConditions := Conditions{
		testCondition,
	}

	conditionsYAML, err := yaml.Marshal(testConditions)

	assert.Nil(s.T(), conditionsYAML)
	assert.Error(s.T(), err)
}

func (s *conditionsSuite) TestFactories() {
	isEqualConditionFactory := ConditionFactories[IsEqualConditionType]

	assert.IsType(s.T(), new(IsEqualCondition), isEqualConditionFactory())
}
