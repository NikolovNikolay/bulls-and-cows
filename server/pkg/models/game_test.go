package models

import "testing"
import "gopkg.in/mgo.v2/bson"

func TestNewGameObject(t *testing.T) {

	allVal := func(t *testing.T) {
		fID := bson.NewObjectId()
		sID := bson.NewObjectId()
		NewGame(
			bson.NewObjectId(),
			1,
			&fID,
			&sID,
			1234,
			3467)
	}

	t.Run("Create new game", func(t *testing.T) {
		t.Run("All valid params", allVal)
		// t.Run("SecondGen", testFunc)
		// t.Run("ThirdGen", testFunc)
		// t.Run("FourthGen", testFunc)
	})
}
