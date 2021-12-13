package tests

import (
	"fmt"

	"github.com/el-Mike/restrict"
)

const hasUserConditionType = "BELONGS_TO"

type hasUserCondition struct{}

func (c *hasUserCondition) Type() string {
	return hasUserConditionType
}

func (c *hasUserCondition) Check(request *restrict.AccessRequest) error {
	user, ok := request.Subject.(*User)
	if !ok {
		return restrict.NewConditionNotSatisfiedError(c, request, fmt.Errorf("Subject has to be a User"))
	}

	conversation, ok := request.Resource.(*Conversation)
	if !ok {
		return restrict.NewConditionNotSatisfiedError(c, request, fmt.Errorf("Resource has to be a Conversation"))
	}

	for _, userId := range conversation.Participants {
		if userId == user.ID {
			return nil
		}
	}

	return restrict.NewConditionNotSatisfiedError(c, request, fmt.Errorf("User does not belong to Conversation with ID: %s", conversation.ID))
}
