package router

import "github.com/julienschmidt/httprouter"

// New returns a new router instance
func New() *httprouter.Router {
	return httprouter.New()
}
