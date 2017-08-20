package utils

import "testing"

func TestDbInit(t *testing.T) {
	s, e := initMongo()
	if e != nil || s == nil {
		t.Error("Could not init mongo connection")
	}
}
