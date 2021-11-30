package adapters

import (
	"testing"

	"github.com/el-Mike/restrict/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type inMemoryAdapterSuite struct {
	suite.Suite
}

func TestInMemoryAdapterSuite(t *testing.T) {
	suite.Run(t, new(inMemoryAdapterSuite))
}

func (s *inMemoryAdapterSuite) TestNewInMemoryAdapter() {
	testPolicy := mocks.GetBasicPolicy()

	adapter := NewInMemoryAdapter(testPolicy)

	assert.NotNil(s.T(), adapter)
	assert.IsType(s.T(), new(InMemoryAdapter), adapter)
}

func (s *inMemoryAdapterSuite) TestLoadPolicy() {
	testPolicy := mocks.GetBasicPolicy()

	adapter := NewInMemoryAdapter(testPolicy)

	policy, err := adapter.LoadPolicy()

	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), policy)
	assert.Equal(s.T(), policy, testPolicy)
}

func (s *inMemoryAdapterSuite) TestSavePolicy() {
	emptyPolicy := mocks.GetEmptyPolicy()
	testPolicy := mocks.GetBasicPolicy()

	adapter := NewInMemoryAdapter(emptyPolicy)

	assert.NotNil(s.T(), adapter)

	policy, err := adapter.LoadPolicy()

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), policy, emptyPolicy)

	err = adapter.SavePolicy(testPolicy)

	assert.Nil(s.T(), err)

	policy, err = adapter.LoadPolicy()

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), policy, testPolicy)
}
