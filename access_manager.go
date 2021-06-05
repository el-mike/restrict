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
	PolicyDefinition *PolicyDefinition `json:"policyDefinition"`
}

// AM - singleton variable for storing AccessManager instance.
var AM *AccessManager

// InitAccessManager - initializes AccessManager with provided PolicyDefinition
// and sets singleton instance.
func InitAccessManager(policyDefinition *PolicyDefinition) (*AccessManager, error) {
	if AM != nil {
		return nil, errors.New("AccessManager already initialized!")
	}

	AM = &AccessManager{
		PolicyDefinition: policyDefinition,
	}

	return AM, nil
}

// IsGranted - checks if given Actions are granted to a Role in regard to resourceID
// an entity wants to access.
func (am *AccessManager) IsGranted(roleID string, resourceID string, actions ...Action) (bool, error) {
	role := am.PolicyDefinition.Roles[roleID]

	if role == nil {
		return false, errors.New("Role does not exist!")
	}

	grants := role.Grants[resourceID]
	parents := role.Parents

	if len(grants) == 0 && len(parents) == 0 {
		return false, nil
	}

	for _, action := range actions {
		if am.hasAction(grants, action) {
			return true, nil
		}

		if len(parents) > 0 {
			for _, parent := range parents {
				granted, err := am.IsGranted(parent, resourceID, action)

				if err != nil {
					return false, err
				}

				if granted {
					return true, nil
				}
			}
		}
	}

	return false, nil
}

// hasAction - checks if grants list contain given action.
func (am *AccessManager) hasAction(grants []Action, action Action) bool {
	for _, grant := range grants {
		if grant == action {
			return true
		}
	}

	return false
}
