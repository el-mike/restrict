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
		ID:    s.testUserId,
		Roles: []string{UserRole},
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

	permissionErr := err.(*restrict.AccessDeniedError).FirstReason()
	conditionErr := permissionErr.FirstConditionError()

	assert.IsType(s.T(), new(restrict.PermissionError), permissionErr)
	assert.IsType(s.T(), new(restrict.ConditionNotSatisfiedError), conditionErr)

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
	assert.IsType(s.T(), new(restrict.PermissionError), err.(*restrict.AccessDeniedError).FirstReason())

	// "delete" condition not satisfied - Conversation must be inactive.
	err = manager.Authorize(&restrict.AccessRequest{
		Subject:  user,
		Resource: conversation,
		Actions:  []string{"delete"},
	})

	assert.IsType(s.T(), new(restrict.AccessDeniedError), err)

	permissionErr = err.(*restrict.AccessDeniedError).FirstReason()
	conditionErr = permissionErr.FirstConditionError()

	assert.IsType(s.T(), new(restrict.ConditionNotSatisfiedError), conditionErr)

	condition := conditionErr.Condition.(*restrict.EmptyCondition)
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
		ID:    "admin1",
		Roles: []string{AdminRole},
	}

	err = manager.Authorize(&restrict.AccessRequest{
		Subject:  admin,
		Resource: &User{},
		Actions:  []string{"create"},
	})

	assert.Nil(s.T(), err)

	// Admin CAN create Conversation because inherits from User.
	err = manager.Authorize(&restrict.AccessRequest{
		Subject:  admin,
		Resource: &Conversation{},
		Actions:  []string{"create"},
	})

	assert.Nil(s.T(), err)

	// Admin CAN read any Conversation, because it has unconditional read permission
	// (along with conditional one inherited from User).
	err = manager.Authorize(&restrict.AccessRequest{
		Subject:  admin,
		Resource: conversation,
		Actions:  []string{"read"},
	})

	assert.Nil(s.T(), err)
}

func (s *integrationSuite) TestRestrict_CompleteValidation() {
	policyManager, err := restrict.NewPolicyManager(adapters.NewInMemoryAdapter(PolicyOne), true)
	if err != nil {
		log.Fatal(err)
	}

	manager := restrict.NewAccessManager(policyManager)

	user := &User{
		ID:    s.testUserId,
		Roles: []string{UserRole},
	}

	conversation := &Conversation{
		ID:            "testConversation1",
		CreatedBy:     "otherUser1",
		Participants:  []string{},
		Active:        true,
		MessagesCount: 20,
	}

	err = manager.Authorize(&restrict.AccessRequest{
		Subject:            user,
		Resource:           conversation,
		Actions:            []string{"read", "update", "delete"},
		CompleteValidation: true,
		Context: restrict.Context{
			"Max": 10,
		},
	})

	assert.IsType(s.T(), new(restrict.AccessDeniedError), err)

	accessErr := err.(*restrict.AccessDeniedError)

	// Expect 3 permission errors, one per each Action.
	assert.Equal(s.T(), 3, len(accessErr.Reasons))

	permissionErrors := accessErr.Reasons

	// "read" action failing with unsatisfied hasUserCondition (readWhereBelongs preset).
	assert.Equal(s.T(), "read", permissionErrors[0].Action)
	assert.Equal(s.T(), 1, len(permissionErrors[0].ConditionErrors))
	assert.IsType(s.T(), new(hasUserCondition), permissionErrors[0].ConditionErrors[0].Condition)

	// "update" action failing with unsatisfied EqualCondition (updateOwn preset).
	assert.Equal(s.T(), "update", permissionErrors[1].Action)
	assert.Equal(s.T(), 1, len(permissionErrors[1].ConditionErrors))
	assert.IsType(s.T(), new(restrict.EqualCondition), permissionErrors[1].ConditionErrors[0].Condition)

	// "delete" action failing with unsatisfied EmptyCondition and greaterThanCondition.
	assert.Equal(s.T(), "delete", permissionErrors[2].Action)
	assert.Equal(s.T(), 2, len(permissionErrors[2].ConditionErrors))
	assert.IsType(s.T(), new(restrict.EmptyCondition), permissionErrors[2].ConditionErrors[0].Condition)
	assert.IsType(s.T(), new(greaterThanCondition), permissionErrors[2].ConditionErrors[1].Condition)
}
