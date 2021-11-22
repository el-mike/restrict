package restrict

const (
	Noop   string = ""       // Noop - empty action.
	Create string = "create" // Create - action for creating a resource.
	Read   string = "read"   // Read - action for reading resource of given type.
	Update string = "update" // Update - action for updating resource of given type.
	Delete string = "delete" // Delete - action for deleting resource of given type.
	CRUD   string = "crud"   // CRUD - action encompassing all CRUD actions.
)

// Actions - alias type for a slice of Actions.
type Actions []string
