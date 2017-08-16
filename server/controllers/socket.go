package controllers

import (
	"log"
	"strconv"

	"gopkg.in/mgo.v2/bson"

	"github.com/NikolovNikolay/bulls-and-cows/server/models"
	"github.com/NikolovNikolay/bulls-and-cows/server/utils"
	"github.com/googollee/go-socket.io"
)

// SocketController gives access to some real-time
// communication via socket.io
type SocketController struct {
	Socket  *socketio.Server
	Gc      GameController
	Pc      PlayerController
	roomMap map[string]int
}

const (
	evtConnection      = "connection"
	evtDisconnect      = "disconnect"
	evtError           = "error"
	evtCreateGame      = "creater"
	evtJoinGame        = "joinr"
	evtJoinedAGame     = "joinmy"
	evtGetActiveGames  = "getavr"
	evtMakeGuess       = "inputguess"
	evtUpdActiveGames  = "updater"
	evtConfirmJoinGame = "confjoin"

	roomDefault = "defroom"

	// SockIoEndpoint is socket.io server route
	SockIoEndpoint = "/socket.io/"
)

// NewSocketController returns a new
// instance of SocketController
func NewSocketController(
	socket *socketio.Server,
	gc GameController,
	pc PlayerController) SocketController {
	return SocketController{
		Gc:      gc,
		Pc:      pc,
		Socket:  socket,
		roomMap: make(map[string]int)}
}

// Init configures the socket.io server
func (sc SocketController) Init() {
	sc.Socket.On(evtConnection, func(so socketio.Socket) {
		log.Println(so.Id(), "connected via socket.io")
		so.Join(roomDefault)

		so.On(evtCreateGame, sc.createGameHandler(so))
		so.On(evtJoinGame, sc.joinGameHandler(so))
		so.On(evtGetActiveGames, sc.getGames(so))
		so.On(evtMakeGuess, sc.setPlayerGuessNumHandler(so))
	})
	sc.Socket.On(evtError, func(so socketio.Socket, err error) {
		log.Println("error:", err)
	})
	sc.Socket.On(evtDisconnect, func(so socketio.Socket, err error) {
		log.Println("disconnect:", err)
	})
}

func (sc *SocketController) createGameHandler(
	so socketio.Socket) func(data string) bool {

	return func(room string) bool {
		if sc.roomMap[room] > 0 {
			return false
		}

		sc.increaseRoomParticipants(room)

		so.BroadcastTo(roomDefault, evtUpdActiveGames, [1]string{room})
		so.Leave(roomDefault)
		so.Join(room)

		p := models.Player{
			ID:     bson.NewObjectId(),
			Logged: true,
			Name:   room}
		sc.Pc.CreatePlayer(p, utils.DBName)

		g, _ := sc.Gc.CreateNewGame(
			utils.DBName,
			utils.P2P,
			&p.ID,
			0)

		p.LoggedIn = &g.GameID
		sc.Pc.UpdatePlayer(p, utils.DBName)

		return true
	}
}

func (sc *SocketController) joinGameHandler(
	so socketio.Socket) func(a, b string) bool {

	return func(room, rival string) bool {
		roomFull := false
		for k, v := range sc.roomMap {
			if k == room && v == 2 {
				roomFull = true
			}
		}

		if roomFull {
			return false
		}

		sc.increaseRoomParticipants(room)
		so.BroadcastTo(roomDefault, evtJoinedAGame, rival)
		so.Join(room)

		host, _ := sc.Pc.FindPlayerByName(room, utils.DBName)
		gID := host.LoggedIn
		g, _ := sc.Gc.GetGameByID(gID.Hex(), utils.DBName)
		p := models.Player{
			ID:       bson.NewObjectId(),
			Logged:   true,
			LoggedIn: gID,
			Name:     rival}
		sc.Pc.CreatePlayer(p, utils.DBName)
		g.PlayerTwoID = &p.ID
		sc.Gc.UpdateGameByID(g, utils.DBName)
		so.BroadcastTo(room, evtConfirmJoinGame)

		return true
	}
}

func (sc *SocketController) getGames(
	so socketio.Socket) func(a string) []string {

	return func(playerName string) []string {
		avrooms := make([]string, 1)

		for k, v := range sc.roomMap {
			if v < 2 && k != playerName {
				avrooms = append(avrooms, k)
			}
		}

		return avrooms
	}
}

func (sc *SocketController) setPlayerGuessNumHandler(
	so socketio.Socket) func(a, b string) bool {
	return func(guess, playerName string) bool {
		p, _ := sc.Pc.FindPlayerByName(playerName, utils.DBName)
		g, _ := sc.Gc.GetGameByID(p.LoggedIn.Hex(), utils.DBName)

		if g.PlayerOneID == &p.ID {
			g.GuessNum, _ = strconv.Atoi(guess)
		} else {
			g.GuessNumSec, _ = strconv.Atoi(guess)
		}

		sc.Gc.UpdateGameByID(g, utils.DBName)

		return true
	}
}

// Utils

func (sc *SocketController) increaseRoomParticipants(room string) {
	for k, v := range sc.roomMap {
		if k == room {
			sc.roomMap[k] = v + 1
		}
	}
}
