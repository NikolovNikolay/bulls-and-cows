package game

import (
	"testing"
	"time"

	"github.com/NikolovNikolay/bulls-and-cows/server/pkg/player"
	"github.com/NikolovNikolay/bulls-and-cows/server/pkg/utils"
)

func TestNewGameObject(t *testing.T) {

	allVal := func(t *testing.T) {
		g := New(1)
		if g.GameType != 1 {
			t.Error("Could not initialize a new game object")
		}
	}

	t.Run("Create new game", func(t *testing.T) {
		t.Run("All valid params", allVal)
	})
}

func TestDbOps(t *testing.T) {

	var game *Game
	add := func(t *testing.T) {
		g := New(1)
		g.SetNumber(8327)
		g.Start(time.Now().Unix())
		g.End(time.Now().Unix() + 1000)
		e := g.Add(utils.DBNameTest, utils.GetDBSession())

		game = g

		if e != nil {
			t.Error("Could not insert game in DB", e)
		}
	}

	findByID := func(t *testing.T) {
		g, e := FindByID(
			game.ID.Hex(),
			utils.DBNameTest,
			utils.GetDBSession())

		if e != nil {
			t.Error("Error appeared while trying to get game from DB", e)
		}

		if g.ID.Hex() != game.ID.Hex() {
			t.Error("Did not get the right game from DB")
		}
	}

	updateByID := func(t *testing.T) {
		game.Number = 1234

		e := game.Update(
			utils.DBNameTest,
			utils.GetDBSession())

		if e != nil {
			t.Error("Error appeared while trying to update game in DB", e)
		}

		ng, e := FindByID(
			game.ID.Hex(),
			utils.DBNameTest,
			utils.GetDBSession())

		if e != nil {
			t.Error("Error appeared while trying to get game from DB", e)
		}

		if ng.Number != game.Number {
			t.Error("Did not update game DB correctly", game.ID)
		}
	}

	t.Run("DB ops", func(t *testing.T) {
		t.Run("Add", add)
		t.Run("FindByID", findByID)
		t.Run("UpdateById", updateByID)
	})

}

func TestFindById(t *testing.T) {
	g := New(2)
	g.AddPlayer(player.New(
		"Test User",
		false,
		utils.DBNameTest,
		utils.GetDBSession()))
	e := g.Add(utils.DBNameTest, utils.GetDBSession())

	if e != nil {
		t.Error("Could not insert game in DB", e)
	}

	fr, e := FindByID(g.ID.Hex(), utils.DBNameTest, utils.GetDBSession())
	if e != nil || fr.ID.Hex() != g.ID.Hex() {
		t.Error("Could not find the proper game doc in DB")
	}
}
