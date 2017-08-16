package player

import (
	"testing"

	"github.com/NikolovNikolay/bulls-and-cows/server/pkg/utils"

	"gopkg.in/mgo.v2/bson"
)

func TestNew(t *testing.T) {
	pID := bson.NewObjectId()
	New(pID, "Test", true)
}

func TestDbOps(t *testing.T) {
	pID := bson.NewObjectId()
	newPlayer := New(pID, "Test User", false)

	add := func(t *testing.T) {
		e := AddToDB(
			newPlayer,
			utils.DBNameTest,
			utils.GetDBSession())

		if e != nil {
			t.Error("Could not insert new player in DB", e)
		}
	}

	findByName := func(t *testing.T) {
		p, e := FindByName(
			newPlayer.Name,
			utils.DBNameTest,
			utils.GetDBSession())

		if e != nil {
			t.Error("Could not get player from DB", e)
			return
		}

		if p.Name != newPlayer.Name {
			t.Error("Did not return correct player record from db")
		}
	}

	findByInvalidID := func(t *testing.T) {
		_, e := FindByID(
			"newPlayer.Name",
			utils.DBNameTest,
			utils.GetDBSession())

		if e == nil {
			t.Error("Should not return player from DB with invalid id")
			return
		}

		if e != nil && e.Error() != "Invalid player id" {
			t.Error("Did not return correct player record from db")
		}
	}

	findByID := func(t *testing.T) {
		_, e := FindByID(
			newPlayer.ID.Hex(),
			utils.DBNameTest,
			utils.GetDBSession())

		if e != nil {
			t.Error("Could not get player from DB")
			return
		}
	}

	update := func(t *testing.T) {
		newPlayer.Logged = true
		e := Update(
			newPlayer,
			utils.DBNameTest,
			utils.GetDBSession())

		if e != nil {
			t.Error("Could not update player in DB", e)
			return
		}

		g, e := FindByID(
			newPlayer.ID.Hex(),
			utils.DBNameTest,
			utils.GetDBSession())

		if e != nil {
			t.Error("Could not find player in DB", e)
			return
		}

		if g.Logged != true {
			t.Error("Player was not updated correctly in DB")
		}
	}

	t.Run("Db ops", func(t *testing.T) {
		t.Run("Add", add)
		t.Run("Find by name", findByName)
		t.Run("Find by invalid ID", findByInvalidID)
		t.Run("Find by ID", findByID)
		t.Run("Update", update)
	})
}
