package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type testStruct struct {
	IntField        int
	StringField     string
	privateIntField int //nolint
}

type typeUtilsSuite struct {
	suite.Suite
}

func TestTypeUtilsSuite(t *testing.T) {
	suite.Run(t, new(typeUtilsSuite))
}

func (s *typeUtilsSuite) TestIsSameType() {
	assert.True(s.T(), IsSameType(1, 2))
	assert.False(s.T(), IsSameType(1, ""))
	assert.False(s.T(), IsSameType(nil, ""))

	aStruct := &testStruct{}
	bStruct := &testStruct{
		IntField: 1,
	}

	assert.True(s.T(), IsSameType(aStruct, bStruct))
	assert.False(s.T(), IsSameType(aStruct, &bStruct))
}

func (s *typeUtilsSuite) TestIsStruct() {
	assert.False(s.T(), IsStruct(1))
	assert.False(s.T(), IsStruct(nil))

	testStruct := testStruct{}

	assert.True(s.T(), IsStruct(testStruct))
	assert.False(s.T(), IsStruct(&testStruct))

	testMap := map[string]int{}

	assert.False(s.T(), IsStruct(testMap))
}

func (s *typeUtilsSuite) TestIsMap() {
	assert.False(s.T(), IsMap(1))
	assert.False(s.T(), IsMap(nil))

	testStruct := testStruct{}

	assert.False(s.T(), IsMap(testStruct))
	assert.False(s.T(), IsMap(&testStruct))

	testMap := map[string]int{}

	assert.True(s.T(), IsMap(testMap))
	assert.False(s.T(), IsMap(&testMap))
}

func (s *typeUtilsSuite) TestGetMapValue() {
	testKey := "testKey"
	testValue := 1

	assert.Nil(s.T(), GetMapValue(1, testKey))

	testMap := map[string]int{
		"testKey": testValue,
	}

	assert.Equal(s.T(), testValue, GetMapValue(testMap, testKey))
	assert.Equal(s.T(), testValue, GetMapValue(&testMap, testKey))
	assert.Zero(s.T(), 0, GetMapValue(testMap, "invalidKey"))
}

func (s *typeUtilsSuite) TestHasField() {
	testStruct := testStruct{
		IntField:    1,
		StringField: "test",
	}

	assert.False(s.T(), HasField(1, "test"))
	assert.True(s.T(), HasField(testStruct, "IntField"))
	assert.True(s.T(), HasField(testStruct, "StringField"))
	assert.True(s.T(), HasField(&testStruct, "IntField"))
	assert.False(s.T(), HasField(testStruct, "InvalidField"))
}

func (s *typeUtilsSuite) TestGetStructFieldValue() {
	testInt := 1
	testString := "test"

	testStruct := testStruct{
		IntField:    testInt,
		StringField: testString,
	}

	assert.Equal(s.T(), testInt, GetStructFieldValue(testStruct, "IntField"))
	assert.Equal(s.T(), testString, GetStructFieldValue(testStruct, "StringField"))

	assert.Equal(s.T(), testInt, GetStructFieldValue(&testStruct, "IntField"))

	assert.Nil(s.T(), GetStructFieldValue(1, "IntField"))
	assert.Nil(s.T(), GetStructFieldValue(testStruct, "InvalidKey"))
	assert.Nil(s.T(), GetStructFieldValue(testStruct, "privateIntField"))
}
