package router

import (
	"github.com/NikolovNikolay/bulls-and-cows/server/services"
	"github.com/gorilla/mux"
)

type BCRouter struct {
	R *mux.Router
}

// New returns a new router instance
func New() BCRouter {
	return BCRouter{R: mux.NewRouter()}
}

// RegisterService takes a Service interface and
// registers it to the internal mux
func (b BCRouter) RegisterService(s services.Service) {
	b.R.HandleFunc(
		s.Endpoint(),
		s.Handle).Methods(s.Method())
}
