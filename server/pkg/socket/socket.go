package socket

import (
	"log"
	"strconv"

	"gopkg.in/mgo.v2"

	"github.com/NikolovNikolay/bulls-and-cows/server/pkg/game"
	"github.com/NikolovNikolay/bulls-and-cows/server/pkg/player"
	"github.com/NikolovNikolay/bulls-and-cows/server/pkg/utils"
	"github.com/googollee/go-socket.io"
)

// Socket gives access to some real-time
// communication via socket.io
type Socket struct {
	Socket  *socketio.Server
	roomMap map[string]int
	db      *mgo.Session
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

// New returns a new instance of SocketController
func New(socket *socketio.Server, db *mgo.Session) *Socket {
	return &Socket{
		Socket:  socket,
		roomMap: make(map[string]int),
		db:      db}
}

// Init configures the socket.io server
func (s Socket) Init() {
	e := s.Socket.On(evtConnection, func(so socketio.Socket) {
		log.Println(so.Id(), "connected via socket.io")
		var e error
		e = so.Join(roomDefault)
		e = so.On(evtCreateGame, s.createGameHandler(so))
		e = so.On(evtJoinGame, s.joinGameHandler(so))
		e = so.On(evtGetActiveGames, s.getGames(so))
		e = so.On(evtMakeGuess, s.setPlayerGuessNumHandler(so))
		if e != nil {
			log.Println("There was an error trying to register custom socket events")
		}
	})
	if e != nil {
		log.Println("Could not connect to socket")
	}
	e = s.Socket.On(evtError, func(so socketio.Socket, err error) {
		log.Println("error:", err)
	})
	if e != nil {
		log.Println("Could not register error handler for socket")
	}
	e = s.Socket.On(evtDisconnect, func(so socketio.Socket, err error) {
		log.Println("disconnect:", err)
	})
	if e != nil {
		log.Println("Could not register disconnect handler for socket")
	}
}

func (s *Socket) createGameHandler(
	so socketio.Socket) func(data string) bool {

	return func(room string) bool {
		if s.roomMap[room] > 0 {
			return false
		}

		s.increaseRoomParticipants(room)
		dbName := utils.DBName
		db := utils.GetDBSession()

		e := so.BroadcastTo(
			roomDefault,
			evtUpdActiveGames,
			[1]string{room})
		if e != nil {
			log.Printf("Could not broadcast to '%s' room", room)
		}

		e = so.Leave(roomDefault)
		if e != nil {
			log.Printf("Could not leave the default room")
		}
		e = so.Join(room)
		if e != nil {
			log.Printf("Could not join '%s' room", room)
		}

		p := player.New(
			room,
			false,
			utils.DBName,
			utils.GetDBSession())

		e = p.Add(dbName, db)
		if e != nil {
			log.Println("Could n insert new player to DB")
			return false
		}

		g := game.New(
			utils.P2P,
			dbName,
			db)

		p.LogIn(&g.ID)
		e = p.Update(dbName, db)
		if e != nil {
			log.Println("Could not update player in DB")
			return false
		}

		return true
	}
}

func (s *Socket) joinGameHandler(
	so socketio.Socket) func(a, b string) bool {

	return func(room, rival string) bool {
		roomFull := false
		for k, v := range s.roomMap {
			if k == room && v == 2 {
				roomFull = true
			}
		}

		if roomFull {
			return false
		}

		dbName := utils.DBName
		db := utils.GetDBSession()
		s.increaseRoomParticipants(room)
		e := so.BroadcastTo(roomDefault, evtJoinedAGame, rival)
		if e != nil {
			log.Printf("Could not broadcast event to '%s' room", roomDefault)
		}
		e = so.Join(room)
		if e != nil {
			log.Printf("Could not join '%s", room)
			return false
		}

		host, e := player.FindByName(room, utils.DBName, s.db)
		if e != nil {
			log.Printf("Could not find host of '%s'", room)
			return false
		}

		gID := host.LoggedIn
		g, e := game.FindByID(gID.Hex(), utils.DBName, s.db)
		if e != nil {
			log.Printf("Could not get game with host '%s'", host.Name)
			return false
		}
		p := player.New(
			rival,
			false,
			dbName,
			db)
		p.LogIn(gID)

		e = p.Add(dbName, db)
		if e != nil {
			log.Printf("Could not add player '%s' to DB when attempting to join", p.Name)
		}

		g.AddPlayer(p)
		e = g.Update(dbName, db)
		if e != nil {
			log.Printf("Could not update game '%s' in DB", g.ID)
		}
		e = so.BroadcastTo(room, evtConfirmJoinGame)
		if e != nil {
			log.Printf("Could not broadcast join to '%s' room", room)
		}

		return true
	}
}

func (s *Socket) getGames(
	so socketio.Socket) func(a string) []string {

	return func(playerName string) []string {
		avrooms := make([]string, 1)

		for k, v := range s.roomMap {
			if v < 2 && k != playerName {
				avrooms = append(avrooms, k)
			}
		}

		return avrooms
	}
}

func (s *Socket) setPlayerGuessNumHandler(
	so socketio.Socket) func(a, b string) bool {
	return func(guess, playerName string) bool {
		dbName := utils.DBName
		db := utils.GetDBSession()
		p, e := player.FindByName(playerName, dbName, db)
		if e != nil {
			log.Printf("Could not find player '%s'", playerName)
			return false
		}
		g, e := game.FindByID(p.LoggedIn.Hex(), dbName, db)
		if e != nil {
			log.Printf("Could not find game '%s'", g.ID)
			return false
		}

		if g.Players[0].ID == p.ID {
			g.Players[0].Number, _ = strconv.Atoi(guess)
		} else {
			g.Players[1].Number, _ = strconv.Atoi(guess)
		}

		e = g.Update(dbName, db)
		if e != nil {
			log.Printf("Could not update game '%s'", g.ID)
			return false
		}

		return true
	}
}

func (s *Socket) increaseRoomParticipants(room string) {
	for k, v := range s.roomMap {
		if k == room {
			s.roomMap[k] = v + 1
		}
	}
}
