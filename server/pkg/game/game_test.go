package game

import "testing"
import "gopkg.in/mgo.v2/bson"

func TestNewGameObject(t *testing.T) {

	allVal := func(t *testing.T) {
		fID := bson.NewObjectId()
		sID := bson.NewObjectId()
		_, e := New(
			bson.NewObjectId(),
			1,
			&fID,
			&sID,
			1234,
			3467)
		if e != nil {
			t.Error("Could not create a new game object", e)
		}
	}

	invalidPOneID := func(t *testing.T) {
		_, e := New(
			bson.NewObjectId(),
			1,
			nil,
			nil,
			1234,
			3467)

		if e == nil {
			t.Error("Should not be able to create a new game without host ID")
		}
	}

	t.Run("Create new game", func(t *testing.T) {
		t.Run("All valid params", allVal)
		t.Run("Should fail if invalid host ID", invalidPOneID)
	})
}
