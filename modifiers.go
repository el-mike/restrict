package restrict

// Modifier - describes a modifier that can be applied to Action.
type Modifier string

const (
	Own       Modifier = "OWN"
	BelongsTo Modifier = "BELONGS_TO"
)
