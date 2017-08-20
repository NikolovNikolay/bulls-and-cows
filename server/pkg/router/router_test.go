package router

import (
	"net/http"
	"testing"
)

func TestNew(t *testing.T) {
	r := New()

	if r.R == nil {
		t.Error("Router did not initialize properly")
	}
}

type testService struct{}

func (s testService) Method() string {
	return "GET"
}

func (s testService) Endpoint() string {
	return "/test"
}

func (s testService) Handle(w http.ResponseWriter, response *http.Request) {
	w.WriteHeader(200)
}

func TestServiceReg(t *testing.T) {
	rtr := New()
	r := rtr.RegisterService(testService{})
	meths, e := r.GetMethods()
	if e != nil || len(meths) > 1 || len(meths) < 1 || meths[0] != "GET" {
		t.Error("Incorrectly added a new service to router")
	}
}
