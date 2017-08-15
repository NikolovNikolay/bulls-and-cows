package router

import (
	"github.com/gorilla/mux"
)

// New returns a new router instance
func New() *mux.Router {
	return mux.NewRouter()
}
