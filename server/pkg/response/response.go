package response

// Response represents a custom response object
// holding some metadata, needed for returning the main
// response to the client
type Response struct {
	Status  int         `json:"s"`
	Payload interface{} `json:"p"`
	Error   string      `json:"e"`
}

// New returns new custom response object
func New(status int, payload interface{}, e error) *Response {
	r := &Response{}

	if status != 0 {
		r.Status = status
	}
	if payload != nil {
		r.Payload = payload
	}
	if e != nil {
		r.Error = e.Error()
	}

	return r
}
