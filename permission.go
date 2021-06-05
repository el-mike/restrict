package restrict

// Permission - result of checking if Action is granted
// for a given Role in regard of wanted Resource.
type Permission struct {
	ResourceID string   `json:"id"`
	Actions    []Action `json:"actions"`
	Attributes []string `json:"attributes"`
	Granted    bool     `json:"granted"`
}
