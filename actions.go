package restrict

const (
	Noop   string = ""       // Noop - empty action.
	Create string = "CREATE" // Create - action for creating a resource.
	Read   string = "READ"   // Read - action for reading resource of given type.
	Update string = "UPDATE" // Update - action for updating resource of given type.
	Delete string = "DELETE" // Delete - action for deleting resource of given type.
	CRUD   string = "CRUD"   // CRUD - action encompassing all CRUD actions.
)
