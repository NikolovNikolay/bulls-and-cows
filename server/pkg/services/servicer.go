package services

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/NikolovNikolay/bulls-and-cows/server/pkg/response"
)

// Servicer represents a service executing
// some logic with a handler
type Servicer interface {
	Method() string
	Endpoint() string
	Handle(w http.ResponseWriter, response *http.Request)
}

// DefSendResponseBeh is the common send response behaviour
func DefSendResponseBeh(w http.ResponseWriter, response *response.Response) {
	w.WriteHeader(response.Status)
	if e := json.NewEncoder(w).Encode(response); e != nil {
		log.Println(e)
	}

}
