package services

import (
	"encoding/json"
	"net/http"

	"github.com/NikolovNikolay/bulls-and-cows/server/response"
)

// Service represents a service executing
// some logic with a handler
type Service interface {
	Method() string
	Endpoint() string
	Handle(w http.ResponseWriter, response *http.Request)
}

// DefSendResponseBeh is the common send response behaviour
func DefSendResponseBeh(w http.ResponseWriter, response *response.Response) {
	w.WriteHeader(response.Status)
	json.NewEncoder(w).Encode(response)
}
