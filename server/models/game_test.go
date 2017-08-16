package models

import "testing"
import "gopkg.in/mgo.v2/bson"

func TestNewGameObject(t *testing.T) {

	allVal := func(t *testing.T) {
		fID := bson.NewObjectId()
		sID := bson.NewObjectId()
		_, e := NewGame(
			bson.NewObjectId(),
			1,
			&fID,
			&sID,
			1234,
			3467)

		if e != nil {
			t.Error("Could not pass with valid game params", e)
		}
	}

	nilFPlayer := func(t *testing.T) {
		_, e := NewGame(
			bson.NewObjectId(),
			1,
			nil,
			nil,
			1234,
			3467)

		if e == nil {
			t.Error("Should not pass with nil player one ID")
		}
	}

	t.Run("Create new game", func(t *testing.T) {
		t.Run("All valid params", allVal)
		t.Run("No main player ID", nilFPlayer)
	})
}
