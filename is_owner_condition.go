package restrict

// IsOwnerConditionName - IsOwnerCondition's identifier.
const IsOwnerConditionName string = "IS_OWNER"

// IsOwnerCondition - Condition for testing whether given Subject is an owner of given Resource.
// Subject and Resource need to implement IdentifiableSubject and OwnableResource interfaces
// respectively.
type IsOwnerCondition struct{}

// Name - returns Condition's name.
func (c *IsOwnerCondition) Name() string {
	return IsOwnerConditionName
}

// Check - returns true if Condition is satisfied, false otherwise.
func (c *IsOwnerCondition) Check(request *AccessRequest) error {
	subjectObject := request.Subject
	resourceObject := request.Resource

	if subjectObject == nil || resourceObject == nil {
		return NewConditionNotSatisfiedError(c, request)
	}

	subject, ok := subjectObject.(IdentifiableSubject)
	if !ok {
		return NewConditionNotSatisfiedError(c, request)
	}

	resource, ok := resourceObject.(OwnableResource)
	if !ok {
		return NewConditionNotSatisfiedError(c, request)
	}

	isOwner := subject.GetId() == resource.GetOwner()

	if !isOwner {
		return NewConditionNotSatisfiedError(c, request)
	}

	return nil
}
