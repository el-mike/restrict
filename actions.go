package restrict

const (
	// NoopAction - empty action.
	NoopAction string = ""
	// Create - action for creating a resource.
	Create string = "create"
	// Read - action for reading resource of given type.
	Read string = "read"
	// Update - action for updating resource of given type.
	Update string = "update"
	// Delete - action for deleting resource of given type.
	Delete string = "delete"
	// CRUD - action encompassing all CRUD actions.
	CRUD string = "crud"
)

// Actions - alias type for a slice of Actions.
type Actions []string
