package restrict

// ActionType - string describing an operation represented
// by given Action.
type ActionType string

const (
	Noop   ActionType = ""       // Noop - empty action.
	Create ActionType = "CREATE" // Create - action for creating a resource.
	Read   ActionType = "READ"   // Read - action for reading resource of given type.
	Update ActionType = "UPDATE" // Update - action for updating resource of given type.
	Delete ActionType = "DELETE" // Delete - action for deleting resource of given type.
	CRUD   ActionType = "CRUD"   // CRUD - action encompassing all CRUD actions.
)

// Action - describes an action that can be done in regard to
// given resource.
type Action struct {
	ActionType ActionType `json:"actionType"`
	Modifier   Modifier   `json:"modifier"`
}
