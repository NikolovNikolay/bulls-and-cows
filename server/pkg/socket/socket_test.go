package socket

import "testing"
import "github.com/googollee/go-socket.io"
import "github.com/NikolovNikolay/bulls-and-cows/server/pkg/utils"
import "github.com/NikolovNikolay/bulls-and-cows/server/pkg/game"
import "gopkg.in/mgo.v2/bson"
import "github.com/NikolovNikolay/bulls-and-cows/server/pkg/player"
import "gopkg.in/mgo.v2"

func TestSocket(t *testing.T) {

	var socket *Socket
	var db *mgo.Session

	t.Run("Test socket", func(t *testing.T) {
		t.Run("Instantiate", func(t *testing.T) {
			db = utils.GetDBSession()
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
			success := socket.createGameHandler(nil)("room1")
			if success == false {
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
				db)

			p.LogIn(&gID)

			// Then create a new game and update the id
			g := game.New(3)
			g.ID = gID
			g.AddPlayer(p)

			// Remove all players with the host name
			_, e := db.DB(utils.DBNameTest).C(utils.DBCPlayers).RemoveAll(bson.M{"name": hostName})
			if e != nil {
				t.Error("Could not delete all Joe`s records in DB")
			}

			// Add above objects to DB
			e = g.Add(utils.DBNameTest, db)
			if e != nil {
				t.Error("Could not add game to DB")
			}
			e = p.Add(utils.DBNameTest, db)
			if e != nil {
				t.Error("Could not add player to DB")
			}

			// Now join
			success := socket.joinGameHandler(nil)("Joe", "Bill")
			if success == false {
				t.Error("Could not init socket.io correctly")
			}
		})
		t.Run("Get all games", func(t *testing.T) {
			socket.roomMap["Joe"] = 2
			socket.roomMap["Bill"] = 1
			socket.roomMap["Jo"] = 1
			socket.roomMap["Steven"] = 1
			games := socket.getGames(nil)("Steven")
			if len(games) != 2 {
				t.Error("Games were not returned as expected")
			}
		})
		t.Run("Set player guess", func(t *testing.T) {
			socket.roomMap["Joe"] = 2
			socket.roomMap["Bill"] = 1
			socket.roomMap["Jo"] = 1
			socket.roomMap["Steven"] = 1
			success := socket.setPlayerGuessNumHandler(nil)("2345", "Joe")
			if !success {
				t.Error("Games were not returned as expected")
			}
		})
		t.Run("Increase player count for game", func(t *testing.T) {
			socket.increaseRoomParticipants("Bill")
			socket.increaseRoomParticipants("room")

			if socket.roomMap["Bill"] != 2 {
				t.Error("Could not increase the room participants")
			}
		})

	})
}

func TestSetPlayerGuess(t *testing.T) {

}
