// Package restrict provides lightweight implementation of RBAC
// authorization model.
package restrict

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

// IsGranted - checks if given AccessRequest can be satisfied given currently loaded policy.
// Returns error if access is not granted or any other problem occurred, nil otherwise.
func (am *AccessManager) IsGranted(request *AccessRequest) error {
	if request.Subject == nil || request.Resource == nil {
		return NewRequestMalformedError(request)
	}

	roleName := request.Subject.GetRole()
	resourceName := request.Resource.GetResourceName()

	if roleName == "" || resourceName == "" {
		return NewRequestMalformedError(request)
	}

	return am.isGranted(request, roleName, resourceName)
}

// isGranted - helper function for decoupling role and resource names retrieval from
// recursive search.
func (am *AccessManager) isGranted(request *AccessRequest, roleName, resourceName string) error {
	role, err := am.policyManager.GetRole(roleName)
	if err != nil {
		return err
	}

	grants := role.Grants[resourceName]
	parents := role.Parents

	// If given role has no permissions granted, and no parents to
	// fall back on, return an error.
	if len(grants) == 0 && len(parents) == 0 {
		return NewNoAvailablePermissionsError(role.ID)
	}

	for _, action := range request.Actions {
		granted := am.validateAction(grants, action, request)

		// If access if not granted for given action on current Role, check if
		// any parent Role can satisfy the request.
		if !granted && len(parents) > 0 {
			for _, parent := range parents {
				parentRequest := &AccessRequest{
					Resource: request.Resource,
					Actions:  []string{action},
					Context:  request.Context,
				}

				if err := am.isGranted(parentRequest, parent, resourceName); err != nil {
					switch err.(type) {
					// If the returned error is one of the below, it just means that
					// access has been denied for some reason.
					case *NoAvailablePermissionsError, *AccessDeniedError:
						granted = false

					// Otherwise, some other problem occurred, and we want to propagate
					// the exception to the caller.
					default:
						return err
					}
				} else {
					// If .IsGranted call with parent Role has returned nil,
					// that means the request is satisfied.
					granted = true
				}
			}
		}

		// If request has not been granted, abort the loop and return an error.
		if !granted {
			return NewAccessDeniedError(request, action)
		}
	}

	return nil
}

// hasAction - checks if grants list contains given action.
func (am *AccessManager) validateAction(permissions []*Permission, action string, request *AccessRequest) bool {
	for _, grant := range permissions {
		if grant.Action == action && am.checkConditions(grant, request) {
			return true
		}
	}

	return false
}

// checkConditions - returns true if all conditions specified for given actions
// are satisfied, false otherwise.
func (am *AccessManager) checkConditions(permission *Permission, request *AccessRequest) bool {
	if permission.Conditions == nil {
		return true
	}

	for _, condition := range permission.Conditions {
		if satisfied := condition.Check(request); !satisfied {
			return false
		}
	}

	return true
}
