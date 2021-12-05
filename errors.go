package restrict

import "fmt"

// RoleNotFoundError - thrown when there is operation called for a Role
// that does not exist.
type RoleNotFoundError struct {
	roleID string
}

// NewRoleNotFoundError - returns new RoleNotFoundError instance.
func NewRoleNotFoundError(roleID string) *RoleNotFoundError {
	return &RoleNotFoundError{
		roleID: roleID,
	}
}

// Error - error interface implementation.
func (e *RoleNotFoundError) Error() string {
	return fmt.Sprintf("Role with ID: \"%s\" has not been found", e.roleID)
}

// RoleAlreadyExistsError - thrown when new role is being added with
// ID that already exists in policy.
type RoleAlreadyExistsError struct {
	roleID string
}

// NewRoleAlreadyExistsError - returns new RoleAlreadyExistsError instance.
func NewRoleAlreadyExistsError(roleID string) *RoleAlreadyExistsError {
	return &RoleAlreadyExistsError{
		roleID: roleID,
	}
}

// Error - error interface implementation.
func (e *RoleAlreadyExistsError) Error() string {
	return fmt.Sprintf("Role with ID: \"%s\" already exists", e.roleID)
}

// NoAvailablePermissionsError - thrown when no Permissions are available for given role.
type NoAvailablePermissionsError struct {
	roleID string
}

// NewNoAvailablePermissionsError - returns new NoAvailablePermissionsError instance.
func NewNoAvailablePermissionsError(roleID string) *NoAvailablePermissionsError {
	return &NoAvailablePermissionsError{
		roleID: roleID,
	}
}

// Error - error interface implementation.
func (e *NoAvailablePermissionsError) Error() string {
	return fmt.Sprintf("No permissions are available for role: \"%s\"", e.roleID)
}

// PermissionPresetNotFoundError - thrown when Permission specifies preset which is not
// defined in PermissionPresets on PolicyDefinition.
type PermissionPresetNotFoundError struct {
	name string
}

// NewPermissionPresetNotFoundError - returns new PermissionPresetNotFoundError instance.
func NewPermissionPresetNotFoundError(name string) *PermissionPresetNotFoundError {
	return &PermissionPresetNotFoundError{
		name: name,
	}
}

// Error - error interface implementation.
func (e *PermissionPresetNotFoundError) Error() string {
	return fmt.Sprintf("Permission preset: \"%s\" has not been found.", e.name)
}

// PermissionPresetAlreadyExistsError - thrown when new Permission preset is being added
// with a name (key) that already exists.
type PermissionPresetAlreadyExistsError struct {
	name string
}

// NewPermissionPresetAlreadyExistsError - returns new PermissionPresetAlreadyExistsError instance.
func NewPermissionPresetAlreadyExistsError(name string) *PermissionPresetAlreadyExistsError {
	return &PermissionPresetAlreadyExistsError{
		name: name,
	}
}

func (e *PermissionPresetAlreadyExistsError) Error() string {
	return fmt.Sprintf("Permission preset with name: \"%s\" already exists", e.name)
}

// AccessDeniedError - thrown when AccessRequest could not be satisfied due to
// insufficient permissions setup.
type AccessDeniedError struct {
	action  string
	request *AccessRequest
	reason  error
}

// NewAccessDeniedError - returns new AccessDeniedError instance.
func NewAccessDeniedError(request *AccessRequest, action string, reason error) *AccessDeniedError {
	return &AccessDeniedError{
		request: request,
		action:  action,
		reason:  reason,
	}
}

// Error - error interface implementation.
func (e *AccessDeniedError) Error() string {
	return fmt.Sprintf("Access denied for action: \"%s\". Reason: %v", e.action, e.reason.Error())
}

// FailedRequest - returns an AccessRequest for which access has been denied.
func (e *AccessDeniedError) FailedRequest() *AccessRequest {
	return e.request
}

// Reason - returns underlying reason (an error) for denying the access.
func (e *AccessDeniedError) Reason() error {
	return e.reason
}

// FailedCondition - helper function for retrieving underlying failed Condition.
func (e *AccessDeniedError) FailedCondition() Condition {
	if e.reason != nil {
		if conditionErr, ok := e.reason.(*ConditionNotSatisfiedError); ok {
			return conditionErr.condition
		}
	}

	return nil
}

// RequestMalformedError - thrown when AccessRequest is no correct or
// does not contain all necessary information.
type RequestMalformedError struct {
	request *AccessRequest
	reason  error
}

// NewRequestMalformedError - returns new SubjectNotDefinedError instance.
func NewRequestMalformedError(request *AccessRequest, reason error) *RequestMalformedError {
	return &RequestMalformedError{
		request: request,
		reason:  reason,
	}
}

// Error - error interface implementation.
func (e *RequestMalformedError) Error() string {
	return fmt.Sprintf("Subject is not defined!")
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

// NewConditionFactoryAlreadyExistsError - returns new ConditionFactoryAlreadyExistsError instance.
func NewConditionFactoryAlreadyExistsError(conditionName string) *ConditionFactoryAlreadyExistsError {
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

// NewConditionFactoryNotFoundError - returns new ConditionFactoryNotFoundError instance.
func NewConditionFactoryNotFoundError(conditionName string) *ConditionFactoryNotFoundError {
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
	descriptor ValueDescriptor
	reason     error
}

// NewValueDescriptorMalformedError - returns new ValueDescriptorMalformedError instance.
func NewValueDescriptorMalformedError(descriptor *ValueDescriptor, reason error) *ValueDescriptorMalformedError {
	return &ValueDescriptorMalformedError{
		descriptor: *descriptor,
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
	return &e.descriptor
}

// ConditionNotSatisfiedError - thrown when given Condition was not satisfied due to
// insufficient privileges for given AccessRequest.
type ConditionNotSatisfiedError struct {
	condition Condition
	request   *AccessRequest
	reason    error
}

// NewConditionNotSatisfiedError - returns new ConditionNotSatisfiedError instance.
func NewConditionNotSatisfiedError(condition Condition, request *AccessRequest, reason error) *ConditionNotSatisfiedError {
	return &ConditionNotSatisfiedError{
		condition: condition,
		request:   request,
		reason:    reason,
	}
}

// Error - error interface implementation.
func (e *ConditionNotSatisfiedError) Error() string {
	return fmt.Sprintf("Condition: \"%v\" was not satisfied! %s", e.condition.Type(), e.reason.Error())
}

// Reason - returns underlying reason (an error) of failing Condition.
func (e *ConditionNotSatisfiedError) Reason() error {
	return e.reason
}

// FailedCondition - returns failed Condition.
func (e *ConditionNotSatisfiedError) FailedCondition() Condition {
	return e.condition
}

// FailedRequest - returns failed AccessRequest.
func (e *ConditionNotSatisfiedError) FailedRequest() *AccessRequest {
	return e.request
}

// PermissionNotGrantedError - thrown when Permission grant for action was not found
// for given Resource.
type PermissionNotGrantedError struct {
	action       string
	resourceName string
}

// NewPermissionNotGrantedError - returns new ActionNotFoundError instance.
func NewPermissionNotGrantedError(action string, resourceName string) *PermissionNotGrantedError {
	return &PermissionNotGrantedError{
		action:       action,
		resourceName: resourceName,
	}
}

// Error - error interface implementation.
func (e *PermissionNotGrantedError) Error() string {
	return fmt.Sprintf("Permission for action: \"%v\" is not granted for Resource: \"%v\"", e.action, e.resourceName)
}

// RoleInheritanceCycleError - thrown when circular inheritance is detected.
type RoleInheritanceCycleError struct {
	roles []string
}

// NewRoleInheritanceCycleError - returns new RoleInheritanceCycleError instance.
func NewRoleInheritanceCycleError(roles []string) *RoleInheritanceCycleError {
	return &RoleInheritanceCycleError{
		roles: roles,
	}
}

// Error - error interface implementation.
func (e *RoleInheritanceCycleError) Error() string {
	message := fmt.Sprintf("Role inheritance cycle has been detected: ")

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
