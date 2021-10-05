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
	action     string
}

// NewPermissionNotFoundError - returns new PermissionNotFoundError instance.
func NewPermissionNotFoundError(resourceID, action string) *PermissionNotFoundError {
	return &PermissionNotFoundError{
		resourceID: resourceID,
		action:     action,
	}
}

// Error - error interface implementation.
func (e *PermissionNotFoundError) Error() string {
	return fmt.Sprintf("Permission for action: %s dot not exist for resource: %s", e.action, e.resourceID)
}

// PermissionAlreadyExistsError - thrown when new permision is being added
// with an action that already exists for given resource.
type PermissionAlreadyExistsError struct {
	resourceID string
	action     string
}

// NewPermissionAlreadyExistsError - returns new PermissionAlreadyExistsError instance.
func NewPermissionAlreadyExistsError(resourceID, action string) *PermissionAlreadyExistsError {
	return &PermissionAlreadyExistsError{
		resourceID: resourceID,
		action:     action,
	}
}

// Error - error interface implementation.
func (e *PermissionAlreadyExistsError) Error() string {
	return fmt.Sprintf("Permission for action: %s already exists for resource: %s", e.action, e.resourceID)
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

// FailedRequest - returns AccessRequest for which access has been denied.
func (e *AccessDeniedError) FailedRequest() *AccessRequest {
	return e.request
}
