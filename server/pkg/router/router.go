package router

import (
	"github.com/NikolovNikolay/bulls-and-cows/server/pkg/services"
	"github.com/gorilla/mux"
)

// BCRouter represents a custom
// wrapper over gorilla
type BCRouter struct {
	R *mux.Router
}

// New returns a new router instance
func New() *BCRouter {
	return &BCRouter{R: mux.NewRouter()}
}

// RegisterService takes a Service interface and
// registers it to the internal mux
func (b BCRouter) RegisterService(s services.Servicer) *mux.Route {
	return b.R.HandleFunc(
		s.Endpoint(),
		s.Handle).Methods(s.Method())
}
