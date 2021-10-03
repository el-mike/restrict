package restrict

// Permission - describes an Action that can be performed in regards to
// some resource, with specified conditions.
type Permission struct {
	Action     string     `json:"action"`
	Conditions Conditions `json:"conditions"`
}
