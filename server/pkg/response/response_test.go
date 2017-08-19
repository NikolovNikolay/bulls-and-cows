package response

import (
	"errors"
	"testing"
)

func TestNew(t *testing.T) {
	r1 := New(200, "test one payload", nil)
	r2 := New(500, "", errors.New("internal error"))

	if r1.Status != 200 && r1.Payload != "test one payload" {
		t.Error("Did not initialize response object properly", r1)
	}
	if r2.Status != 500 && r2.Error != "internal error" {
		t.Error("Did not initialize response object properly", r2)
	}
}
