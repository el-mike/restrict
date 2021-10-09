package restrict

import (
	"reflect"

	"github.com/el-Mike/restrict/utils"
)

const IS_OWNER_CONDITION_NAME string = "IS_OWNER"

type IsOwnerCondition struct {
	IdentityField  string `json:"identityField,omitempty" yaml:"identityField,omitempty"`
	OwnershipField string `json:"owhershipField,omitempty" yaml:"owhershipField,omitempty"`
}

// Name - returns Condition's name.
func (c *IsOwnerCondition) Name() string {
	return IS_OWNER_CONDITION_NAME
}

func (c *IsOwnerCondition) Check(value interface{}, request *AccessRequest) bool {
	subjectObject := request.Subject
	resourceObject := request.Resource

	if !utils.IsStruct(subjectObject) || !utils.IsStruct(resourceObject) {
		return false
	}

	subject := reflect.ValueOf(subjectObject).Elem().FieldByName(c.IdentityField)
	if !subject.IsValid() {
		return false
	}

	resourceOwner := reflect.ValueOf(resourceObject).Elem().FieldByName(c.OwnershipField)
	if !resourceOwner.IsValid() {
		return false
	}

	return subject.String() == resourceOwner.String()
}
