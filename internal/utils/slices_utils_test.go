package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type slicesUtilsSuite struct {
	suite.Suite
}

func TestSlicesUtilsSuite(t *testing.T) {
	suite.Run(t, new(slicesUtilsSuite))
}

func (s *slicesUtilsSuite) TestStringSliceContains() {
	testSlice := []string{"one", "two", "three", "three"}

	assert.True(s.T(), StringSliceContains(testSlice, "one"))
	assert.False(s.T(), StringSliceContains(testSlice, "four"))
}
