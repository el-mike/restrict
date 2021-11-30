// Package restrict provides lightweight implementation of RBAC
// authorization model.
package restrict

import "fmt"

// AccessManager - manages all of the defined Permissions and Roles,
// provides an interface to perform authorization checks and add/remove
// Permissions and Roles.
type AccessManager struct {
	// Instance of PolicyManager, responsible for managing currently loaded policy.
	policyManager *PolicyManager
}

// NewAccessManager - initializes AccessManager with provided PolicyDefinition
// and sets singleton instance.
func NewAccessManager(policyManager *PolicyManager) *AccessManager {
	return &AccessManager{
		policyManager: policyManager,
	}
}

// Authorize - checks if given AccessRequest can be satisfied given currently loaded policy.
// Returns error if access is not granted or any other problem occurred, nil otherwise.
func (am *AccessManager) Authorize(request *AccessRequest) error {
	if request.Subject == nil || request.Resource == nil {
		return NewRequestMalformedError(request, fmt.Errorf("Subject or Resource not defined"))
	}

	roleName := request.Subject.GetRole()
	resourceName := request.Resource.GetResourceName()

	if roleName == "" || resourceName == "" {
		return NewRequestMalformedError(request, fmt.Errorf("Missing roleName or resourceName"))
	}

	return am.authorize(request, roleName, resourceName)
}

// authorize - helper function for decoupling role and resource names retrieval from
// recursive search.
func (am *AccessManager) authorize(request *AccessRequest, roleName, resourceName string) error {
	role, err := am.policyManager.GetRole(roleName)
	if err != nil {
		return NewAccessDeniedError(request, "", err)
	}

	grants := role.Grants[resourceName]
	parents := role.Parents

	// If given role has no permissions granted, and no parents to
	// fall back on, return an error.
	if len(grants) == 0 && len(parents) == 0 {
		return NewNoAvailablePermissionsError(role.ID)
	}

	for _, action := range request.Actions {
		authorizeError := am.validateAction(grants, action, request)

		// If access if not granted for given action on current Role, check if
		// any parent Role can satisfy the request.
		if authorizeError != nil && len(parents) > 0 {
			for _, parent := range parents {
				parentRequest := &AccessRequest{
					Resource: request.Resource,
					Actions:  []string{action},
					Context:  request.Context,
				}

				if err := am.authorize(parentRequest, parent, resourceName); err != nil {
					switch err.(type) {
					// If the returned error is one of the below, it just means that
					// access has been denied for some reason.
					case *NoAvailablePermissionsError,
						*ConditionNotSatisfiedError,
						*ActionNotFoundError,
						*AccessDeniedError:
						authorizeError = err

					// Otherwise, some other problem occurred, and we want to propagate
					// the exception to the caller.
					default:
						return err
					}
				} else {
					// If .Authorize call with parent Role has returned nil,
					// that means the request is satisfied.
					authorizeError = nil
				}
			}
		}

		// If request has not been granted, abort the loop and return an error.
		if authorizeError != nil {
			return NewAccessDeniedError(request, action, authorizeError)
		}
	}

	return nil
}

// hasAction - checks if grants list contains given action.
func (am *AccessManager) validateAction(permissions []*Permission, action string, request *AccessRequest) error {
	for _, grant := range permissions {
		if grant.Action == action {
			if request.SkipConditions {
				return nil
			}

			return am.checkConditions(grant, request)
		}
	}

	return NewActionNotFoundError(action, request.Resource.GetResourceName())
}

// checkConditions - returns nil if all conditions specified for given actions
// are satisfied, error otherwise.
func (am *AccessManager) checkConditions(permission *Permission, request *AccessRequest) error {
	if permission.Conditions == nil {
		return nil
	}

	for _, condition := range permission.Conditions {
		if err := condition.Check(request); err != nil {
			return err
		}
	}

	return nil
}
