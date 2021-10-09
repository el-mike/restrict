package restrict

type Context map[string]interface{}

type AccessRequest struct {
	Role       string
	ResourceID string
	Actions    []string
	Subject    interface{}
	Resource   interface{}
	Context    Context
}
