package game

import "testing"
import "gopkg.in/mgo.v2/bson"
import "github.com/NikolovNikolay/bulls-and-cows/server/pkg/utils"

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

func TestDbOps(t *testing.T) {

	var game *Game
	add := func(t *testing.T) {
		pID := bson.NewObjectId()
		g, e := AddToDb(
			utils.DBNameTest,
			2,
			&pID,
			1234,
			utils.GetDBSession())
		game = g

		if e != nil {
			t.Error("Could not insert game in DB", e)
		}
	}

	findByID := func(t *testing.T) {
		g, e := FindByID(
			game.GameID.Hex(),
			utils.DBNameTest,
			utils.GetDBSession())

		if e != nil {
			t.Error("Error appeared while trying to get game from DB", e)
		}

		if g.GameID.Hex() != game.GameID.Hex() {
			t.Error("Did not get the right game from DB")
		}
	}

	updateByID := func(t *testing.T) {
		game.GuessNumSec = 1234

		e := UpdateByID(
			game,
			utils.DBNameTest,
			utils.GetDBSession())

		if e != nil {
			t.Error("Error appeared while trying to update game in DB", e)
		}

		ng, e := FindByID(
			game.GameID.Hex(),
			utils.DBNameTest,
			utils.GetDBSession())

		if e != nil {
			t.Error("Error appeared while trying to get game from DB", e)
		}

		if ng.GuessNumSec != game.GuessNumSec {
			t.Error("Did not update game DB correctly", game.GameID)
		}
	}

	t.Run("DB ops", func(t *testing.T) {
		t.Run("Add", add)
		t.Run("FindByID", findByID)
		t.Run("UpdateById", updateByID)
	})

}

func TestFindById(t *testing.T) {
	pID := bson.NewObjectId()
	_, e := AddToDb(
		utils.DBNameTest,
		2,
		&pID,
		1234,
		utils.GetDBSession())

	if e != nil {
		t.Error("Could not insert game in DB", e)
	}
}
