// restrict package provides lightweight implementation of RBAC
// authorization model.
package restrict

import (
	"errors"
)

// AccessManager - manages all of the defined Permissions and Roles,
// provides an interface to perform authorization checks and add/remove
// Permissions and Roles.
type AccessManager struct {
	policyManager *PolicyManager
}

// InitAccessManager - initializes AccessManager with provided PolicyDefinition
// and sets singleton instance.
func NewAccessManager(policyManager *PolicyManager) *AccessManager {
	return &AccessManager{
		policyManager: policyManager,
	}
}

// IsGranted - checks if given Actions are granted to a Role in regard to resourceID
// an entity wants to access.
func (am *AccessManager) IsGranted(request *AccessRequest) (bool, error) {
	policy := am.policyManager.GetPolicy()

	role := policy.Roles[request.Role]

	if role == nil {
		return false, errors.New("Role does not exist!")
	}

	grants := role.Grants[request.Resource]
	parents := role.Parents

	isGranted := true

	if len(grants) == 0 && len(parents) == 0 {
		return false, nil
	}

	for _, action := range request.Actions {
		isGranted = isGranted && am.validateAction(grants, action, request)

		if len(parents) > 0 {
			for _, parent := range parents {
				request.Role = parent
				granted, err := am.IsGranted(request)

				if err != nil {
					return false, err
				}

				if granted {
					return true, nil
				}
			}
		}
	}

	return isGranted, nil
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
	for key, condition := range permission.Conditions {
		if satisfied := condition.Check(request.Context[key], request); !satisfied {
			return false
		}
	}

	return true
}
