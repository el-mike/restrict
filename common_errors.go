package restrict

import "fmt"

// RoleNotFoundError - thrown when there is an operation called for a Role
// that does not exist.
type RoleNotFoundError struct {
	roleID string
}

// newRoleNotFoundError - returns new RoleNotFoundError instance.
func newRoleNotFoundError(roleID string) *RoleNotFoundError {
	return &RoleNotFoundError{
		roleID: roleID,
	}
}

// Error - error interface implementation.
func (e *RoleNotFoundError) Error() string {
	return fmt.Sprintf("Role with ID: \"%s\" has not been found", e.roleID)
}

// RoleAlreadyExistsError - thrown when new Role is being added with
// ID that already exists in the PolicyDefinition.
type RoleAlreadyExistsError struct {
	roleID string
}

// newRoleAlreadyExistsError - returns new RoleAlreadyExistsError instance.
func newRoleAlreadyExistsError(roleID string) *RoleAlreadyExistsError {
	return &RoleAlreadyExistsError{
		roleID: roleID,
	}
}

// Error - error interface implementation.
func (e *RoleAlreadyExistsError) Error() string {
	return fmt.Sprintf("Role with ID: \"%s\" already exists", e.roleID)
}

// PermissionPresetNotFoundError - thrown when Permission specifies a preset which is not
// defined in PermissionPresets on PolicyDefinition.
type PermissionPresetNotFoundError struct {
	name string
}

// newPermissionPresetNotFoundError - returns new PermissionPresetNotFoundError instance.
func newPermissionPresetNotFoundError(name string) *PermissionPresetNotFoundError {
	return &PermissionPresetNotFoundError{
		name: name,
	}
}

// Error - error interface implementation.
func (e *PermissionPresetNotFoundError) Error() string {
	return fmt.Sprintf("Permission preset: \"%s\" has not been found.", e.name)
}

// PermissionPresetAlreadyExistsError - thrown when a new Permission preset is being added
// with a name (key) that already exists.
type PermissionPresetAlreadyExistsError struct {
	name string
}

// newPermissionPresetAlreadyExistsError - returns new PermissionPresetAlreadyExistsError instance.
func newPermissionPresetAlreadyExistsError(name string) *PermissionPresetAlreadyExistsError {
	return &PermissionPresetAlreadyExistsError{
		name: name,
	}
}

func (e *PermissionPresetAlreadyExistsError) Error() string {
	return fmt.Sprintf("Permission preset with name: \"%s\" already exists", e.name)
}

// RequestMalformedError - thrown when AccessRequest is not correct or
// does not contain all necessary information.
type RequestMalformedError struct {
	request *AccessRequest
	reason  error
}

// newRequestMalformedError - returns new SubjectNotDefinedError instance.
func newRequestMalformedError(request *AccessRequest, reason error) *RequestMalformedError {
	return &RequestMalformedError{
		request: request,
		reason:  reason,
	}
}

// Error - error interface implementation.
func (e *RequestMalformedError) Error() string {
	return "Subject is not defined!"
}

// Reason - returns underlying reason (an error) of malformed Request.
func (e *RequestMalformedError) Reason() error {
	return e.reason
}

// FailedRequest - returns an AccessRequest for which access has been denied.
func (e *RequestMalformedError) FailedRequest() *AccessRequest {
	return e.request
}

// ConditionFactoryAlreadyExistsError - thrown when ConditionFactory is being added under a name
// that's already set in ConditionFactories map.
type ConditionFactoryAlreadyExistsError struct {
	conditionName string
}

// newConditionFactoryAlreadyExistsError - returns new ConditionFactoryAlreadyExistsError instance.
func newConditionFactoryAlreadyExistsError(conditionName string) *ConditionFactoryAlreadyExistsError {
	return &ConditionFactoryAlreadyExistsError{
		conditionName: conditionName,
	}
}

// Error - error interface implementation.
func (e *ConditionFactoryAlreadyExistsError) Error() string {
	return fmt.Sprintf("ConditionFactory for Condition: \"%v\" already exists!", e.conditionName)
}

// ConditionFactoryNotFoundError - thrown when ConditionFactory is not found while
// unmarshaling a Permission.
type ConditionFactoryNotFoundError struct {
	conditionName string
}

// newConditionFactoryNotFoundError - returns new ConditionFactoryNotFoundError instance.
func newConditionFactoryNotFoundError(conditionName string) *ConditionFactoryNotFoundError {
	return &ConditionFactoryNotFoundError{
		conditionName: conditionName,
	}
}

// Error - error interface implementation.
func (e *ConditionFactoryNotFoundError) Error() string {
	return fmt.Sprintf("ConditionFactory not found for Condition: \"%v\"", e.conditionName)
}

// ValueDescriptorMalformedError - thrown when malformed ValueDescriptor is being resolved.
type ValueDescriptorMalformedError struct {
	descriptor *ValueDescriptor
	reason     error
}

// newValueDescriptorMalformedError - returns new ValueDescriptorMalformedError instance.
func newValueDescriptorMalformedError(descriptor *ValueDescriptor, reason error) *ValueDescriptorMalformedError {
	return &ValueDescriptorMalformedError{
		descriptor: descriptor,
		reason:     reason,
	}
}

// Error - error interface implementation.
func (e *ValueDescriptorMalformedError) Error() string {
	return fmt.Sprintf("ValueDescriptor could not be resolved! Reason: %s", e.reason.Error())
}

// Reason - returns underlying reason (an error) of malformed ValueDescriptor.
func (e *ValueDescriptorMalformedError) Reason() error {
	return e.reason
}

// FailedDescriptor - returns failed ValueDescriptor.
func (e *ValueDescriptorMalformedError) FailedDescriptor() *ValueDescriptor {
	return e.descriptor
}

// RoleInheritanceCycleError - thrown when circular Role inheritance is detected.
type RoleInheritanceCycleError struct {
	roles []string
}

// newRoleInheritanceCycleError - returns new RoleInheritanceCycleError instance.
func newRoleInheritanceCycleError(roles []string) *RoleInheritanceCycleError {
	return &RoleInheritanceCycleError{
		roles: roles,
	}
}

// Error - error interface implementation.
func (e *RoleInheritanceCycleError) Error() string {
	message := "Role inheritance cycle has been detected: "

	for i, role := range e.roles {
		if i > 0 {
			message += " -> "
		}

		message += fmt.Sprintf("\"%s\"", role)
	}

	// We want to add the first role at the end, to indicate the cycle.
	message += fmt.Sprintf(" -> \"%s\"", e.roles[0])

	return message
}
