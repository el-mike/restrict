package tests

import (
	"fmt"
	"log"
	"testing"

	"github.com/el-mike/restrict"
	"github.com/el-mike/restrict/adapters"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type integrationSuite struct {
	suite.Suite

	testUserId string
}

func TestPoliciesSuite(t *testing.T) {
	suite.Run(t, new(integrationSuite))
}

func (s *integrationSuite) SetupSuite() {
	s.testUserId = "testUser1"

	//nolint
	restrict.RegisterConditionFactory(hasUserConditionType, func() restrict.Condition {
		return new(hasUserCondition)
	})

	//nolint
	restrict.RegisterConditionFactory(greatherThanType, func() restrict.Condition {
		return new(greaterThanCondition)
	})
}

func (s *integrationSuite) TestRestrict_JSON() {
	jsonAdapter := adapters.NewFileAdapter("test_policy.json", adapters.JSONFile)

	policyManager, err := restrict.NewPolicyManager(jsonAdapter, true)
	if err != nil {
		log.Fatal(err)
	}

	s.testPolicy(policyManager)
}

func (s *integrationSuite) TestRestrict_YAML() {
	yamlAdapter := adapters.NewFileAdapter("test_policy.yaml", adapters.YAMLFile)

	policyManager, err := restrict.NewPolicyManager(yamlAdapter, true)
	if err != nil {
		log.Fatal(err)
	}

	s.testPolicy(policyManager)
}

func (s *integrationSuite) TestRestrict() {
	policyManager, err := restrict.NewPolicyManager(adapters.NewInMemoryAdapter(PolicyOne), true)
	if err != nil {
		log.Fatal(err)
	}

	s.testPolicy(policyManager)
}

func (s *integrationSuite) testPolicy(policyManager *restrict.PolicyManager) {
	user := &User{
		ID:   s.testUserId,
		Role: UserRole,
	}

	conversation := &Conversation{
		ID:           "testConversation1",
		CreatedBy:    "otherUser1",
		Participants: []string{},
		Active:       true,
	}

	manager := restrict.NewAccessManager(policyManager)

	// "read" is not granted - User does not belong to the Conversation.
	err := manager.Authorize(&restrict.AccessRequest{
		Subject:  user,
		Resource: conversation,
		Actions:  []string{"read"},
	})

	assert.IsType(s.T(), new(restrict.AccessDeniedError), err)
	assert.IsType(s.T(), new(restrict.ConditionNotSatisfiedError), err.(*restrict.AccessDeniedError).Reason())

	// "read" granted - User belongs to the Conversation.
	conversation.Participants = []string{s.testUserId}

	err = manager.Authorize(&restrict.AccessRequest{
		Subject:  user,
		Resource: conversation,
		Actions:  []string{"read"},
	})

	assert.Nil(s.T(), err)

	// "update" granted - User owns the conversation.
	conversation.CreatedBy = s.testUserId

	err = manager.Authorize(&restrict.AccessRequest{
		Subject:  user,
		Resource: conversation,
		Actions:  []string{"update"},
	})

	fmt.Print(err)

	assert.Nil(s.T(), err)

	// "modify" is not granted.
	err = manager.Authorize(&restrict.AccessRequest{
		Subject:  user,
		Resource: conversation,
		Actions:  []string{"read", "modify"},
	})

	assert.IsType(s.T(), new(restrict.AccessDeniedError), err)
	assert.IsType(s.T(), new(restrict.PermissionNotGrantedError), err.(*restrict.AccessDeniedError).Reason())

	// "delete" condition not satisfied - Conversation must be unactive.
	err = manager.Authorize(&restrict.AccessRequest{
		Subject:  user,
		Resource: conversation,
		Actions:  []string{"delete"},
	})

	assert.IsType(s.T(), new(restrict.AccessDeniedError), err)

	conditionErr := err.(*restrict.AccessDeniedError).Reason()
	assert.IsType(s.T(), new(restrict.ConditionNotSatisfiedError), conditionErr)

	condition := err.(*restrict.AccessDeniedError).FailedCondition().(*restrict.EmptyCondition)
	assert.Equal(s.T(), "deleteActive", condition.ID)

	// "delete" condition not satisfied - Conversation has to have less than 100 messages.
	conversation.Active = false
	conversation.MessagesCount = 110

	err = manager.Authorize(&restrict.AccessRequest{
		Subject:  user,
		Resource: conversation,
		Actions:  []string{"delete"},
		Context: restrict.Context{
			"Max": 100,
		},
	})

	assert.IsType(s.T(), new(restrict.AccessDeniedError), err)

	// User CAN read itself
	err = manager.Authorize(&restrict.AccessRequest{
		Subject:  user,
		Resource: user,
		Actions:  []string{"read"},
	})

	assert.Nil(s.T(), err)

	// User can NOT create other users
	err = manager.Authorize(&restrict.AccessRequest{
		Subject:  user,
		Resource: &User{},
		Actions:  []string{"create"},
	})

	assert.IsType(s.T(), new(restrict.AccessDeniedError), err)

	// Admin CAN create other users
	admin := &User{
		ID:   "admin1",
		Role: AdminRole,
	}

	err = manager.Authorize(&restrict.AccessRequest{
		Subject:  admin,
		Resource: &User{},
		Actions:  []string{"create"},
	})

	assert.Nil(s.T(), err)

	// Admin CAN create Converation because inherits from User.
	err = manager.Authorize(&restrict.AccessRequest{
		Subject:  admin,
		Resource: &Conversation{},
		Actions:  []string{"create"},
	})

	assert.Nil(s.T(), err)

	// Admin CAN read any Converation, because it has unconditional read permission
	// (along with conditional one inherited from User).
	err = manager.Authorize(&restrict.AccessRequest{
		Subject:  admin,
		Resource: conversation,
		Actions:  []string{"read"},
	})

	assert.Nil(s.T(), err)
}
