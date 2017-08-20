package socket

import (
	"testing"
	"time"

	"github.com/NikolovNikolay/bulls-and-cows/server/pkg/game"
	"github.com/NikolovNikolay/bulls-and-cows/server/pkg/player"
	"github.com/NikolovNikolay/bulls-and-cows/server/pkg/utils"
	"github.com/googollee/go-socket.io"
	"gopkg.in/mgo.v2/bson"
)

func TestSocket(t *testing.T) {

	var socket *Socket

	t.Run("Test socket", func(t *testing.T) {
		t.Run("Instantiate", func(t *testing.T) {
			db := utils.GetDBSession()
			io, e := socketio.NewServer(nil)
			if e != nil {
				t.Error("Could not instantiate socket.io server")
			}

			socket = New(io, db, true)
			if socket.Socket == nil {
				t.Error("Did not instantiate socket.io correctly")
			}
		})

		t.Run("Init", func(t *testing.T) {
			e := socket.Init()
			if e != nil {
				t.Error("Could not init socket.io correctly")
			}
		})

		t.Run("Create game", func(t *testing.T) {
			success := socket.createGameHandler(nil)("room1", 1)
			if success == "" {
				t.Error("Could not create game with socket.io")
			}
		})

		t.Run("Join game", func(t *testing.T) {
			gID := bson.NewObjectId()
			hostName := "Joe"
			// First create a new player
			p := player.New(
				hostName,
				false,
				utils.DBNameTest,
				socket.db)

			p.LogIn(&gID)

			// Then create a new game and update the id
			g := game.New(3)
			g.ID = gID
			g.AddPlayer(p)

			// Remove all players with the host name
			_, e := socket.db.DB(utils.DBNameTest).C(utils.DBCPlayers).RemoveAll(bson.M{"name": hostName})
			if e != nil {
				t.Error("Could not delete all Joe`s records in DB")
			}

			// Add above objects to DB
			e = g.Add(utils.DBNameTest, socket.db)
			if e != nil {
				t.Error("Could not add game to DB")
			}
			e = p.Add(utils.DBNameTest, socket.db)
			if e != nil {
				t.Error("Could not add player to DB")
			}

			// Now join
			success := socket.joinGameHandler(nil)("Joe", "Bill")
			if success == "" {
				t.Error("Could not init socket.io correctly")
			}
		})

		t.Run("Get all games", func(t *testing.T) {
			socket.roomMap["Joe"] = 2
			socket.roomMap["Bill"] = 1
			socket.roomMap["Jo"] = 1
			socket.roomMap["Steven"] = 1
			games := socket.getGames(nil)("Steven")
			if len(games) != 4 {
				t.Error("Games were not returned as expected")
			}
		})

		t.Run("Set player guess", func(t *testing.T) {
			socket.roomMap["Joe"] = 2
			socket.roomMap["Bill"] = 1
			socket.roomMap["Jo"] = 1
			socket.roomMap["Steven"] = 1
			success := socket.setPlayerGuessNumHandler(nil)("Bill", "2345", "Joe")
			if !success {
				t.Error("Games were not returned as expected")
			}
		})

		t.Run("Increase player count for game", func(t *testing.T) {
			socket.increaseRoomParticipants("Bill", 1)
			socket.increaseRoomParticipants("room", 2)

			if socket.roomMap["Bill"] != 2 {
				t.Error("Could not increase the room participants")
			}
		})

		t.Run("Test ready game and rival numbers", func(t *testing.T) {
			var dbName = getTargetDBName(socket)
			g := game.New(1)
			g.SetNumber(8327)
			g.Start(time.Now().Unix())
			g.End(time.Now().Unix() + 1000)

			e := g.Add(dbName, socket.db)
			if e != nil {
				t.Error("Could not insert game in DB")
			}

			playerOne := player.New("Player one", false, utils.DBNameTest, utils.GetDBSession())
			playerOne.LogIn(&g.ID)
			playerOne.Number = 1234
			playerTwo := player.New("Player two", false, utils.DBNameTest, utils.GetDBSession())
			playerOne.LogIn(&g.ID)
			playerTwo.Number = 9876
			g.AddPlayer(playerOne)
			g.AddPlayer(playerTwo)

			e = g.Update(dbName, socket.db)
			if e != nil {
				t.Error("Could not update game in DB")
			}

			ready := socket.checkIfReadyForP2PStart(g)

			if !ready {
				t.Error("Incorrectly returned ready state of a game")
			}
			num, e := socket.getP2PRivalNumber(playerOne.Name, g)
			if e != nil {
				t.Error("Could not get opponent's number from DB")
			}
			if num != 9876 {
				t.Error("Opponent number got incorrectly")
			}

			success := socket.setPlayerGuessNumHandler(nil)("", "7645", "Player one")
			if success {
				g, e = game.FindByID(g.ID.Hex(), dbName, socket.db)
				if e != nil {
					t.Error("Could not get game by id from DB")
				}
				if g.Players[0].Number != 7645 {
					t.Error("Player number was not set correctly")
				}
			}

			gp := socket.makeGuessHandler(nil)("", "1234", "Player two", g.ID.Hex())
			if gp.Win != true && gp.BC.Bulls != 4 {
				t.Error("Make guess handler result did not get correct results")
			}
		})
	})
}
