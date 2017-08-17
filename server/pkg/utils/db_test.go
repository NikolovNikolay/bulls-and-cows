package utils

import "testing"

func TestDbInit(t *testing.T) {
	s := initMongo()
	if s == nil {
		t.Error("Could not init mongo connection")
	}
}
