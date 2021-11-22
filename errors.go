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
	return fmt.Sprintf("Role with ID: %s has not been found", e.roleID)
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
	return fmt.Sprintf("Role with ID: %s already exists", e.roleID)
}

// PermissionNotFoundError - thrown when there is operation called for a Permission
// that does not exist.
type PermissionNotFoundError struct {
	resourceID string
	name       string
}

// NewPermissionNotFoundError - returns new PermissionNotFoundError instance.
func NewPermissionNotFoundError(resourceID, name string) *PermissionNotFoundError {
	return &PermissionNotFoundError{
		resourceID: resourceID,
		name:       name,
	}
}

// Error - error interface implementation.
func (e *PermissionNotFoundError) Error() string {
	return fmt.Sprintf("Permission with name: %s dot not exist for resource: %s", e.name, e.resourceID)
}

// PermissionAlreadyExistsError - thrown when new permision is being added
// with a name that already exists for given resource.
type PermissionAlreadyExistsError struct {
	resourceID string
	name       string
}

// NewPermissionAlreadyExistsError - returns new PermissionAlreadyExistsError instance.
func NewPermissionAlreadyExistsError(resourceID, name string) *PermissionAlreadyExistsError {
	return &PermissionAlreadyExistsError{
		resourceID: resourceID,
		name:       name,
	}
}

// Error - error interface implementation.
func (e *PermissionAlreadyExistsError) Error() string {
	return fmt.Sprintf("Permission with name: %s already exists for resource: %s", e.name, e.resourceID)
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
	return fmt.Sprintf("No permissions are available for role: %s", e.roleID)
}

// MissingPermissionNameError - thrown when Permission without a Name is being added
// or loaded.
type MissingPermissionNameError struct {
	permission *Permission
}

// NewMissingPermissionNameError - returns new MissingPermissionNameError instance.
func NewMissingPermissionNameError(permission *Permission) *MissingPermissionNameError {
	return &MissingPermissionNameError{
		permission: permission,
	}
}

// Error - error interface implementation.
func (e *MissingPermissionNameError) Error() string {
	return fmt.Sprintf("Permission without a name cannot be created")
}

// FailedPermission - returns Permission that could not be created due to missing name.
func (e *MissingPermissionNameError) FailedPermission() *Permission {
	return e.permission
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
	return fmt.Sprintf("Permission preset: %s has not been found.", e.name)
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
	return fmt.Sprintf("Permission preset with name: %s already exists", e.name)
}

// AccessDeniedError - thrown when AccessRequest could not be satisfied due to
// insufficient permissions setup.
type AccessDeniedError struct {
	action  string
	request *AccessRequest
}

// NewAccessDeniedError - returns new AccessDeniedError instance.
func NewAccessDeniedError(request *AccessRequest, action string) *AccessDeniedError {
	return &AccessDeniedError{
		request: request,
		action:  action,
	}
}

// Error - error interface implementation.
func (e *AccessDeniedError) Error() string {
	return fmt.Sprintf("Access denied for action: %s", e.action)
}

// FailedRequest - returns an AccessRequest for which access has been denied.
func (e *AccessDeniedError) FailedRequest() *AccessRequest {
	return e.request
}

// RequestMalformedError - thrown when AccessRequest is no correct or
// does not contain all neccessary information.
type RequestMalformedError struct {
	request *AccessRequest
}

// NewRequestMalformedError - returns new SubjectNotDefinedError instance.
func NewRequestMalformedError(request *AccessRequest) *RequestMalformedError {
	return &RequestMalformedError{
		request: request,
	}
}

// Error - error interface implementation.
func (e *RequestMalformedError) Error() string {
	return fmt.Sprintf("Subject is not defined!")
}

// FailedRequest - returns an AccessRequest for which access has been denied.
func (e *RequestMalformedError) FailedRequest() *AccessRequest {
	return e.request
}
