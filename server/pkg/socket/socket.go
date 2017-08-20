package socket

import (
	"errors"
	"log"
	"strconv"
	"time"

	"github.com/NikolovNikolay/bulls-and-cows/server/pkg/guess"

	"gopkg.in/mgo.v2"

	"github.com/NikolovNikolay/bulls-and-cows/server/pkg/game"
	"github.com/NikolovNikolay/bulls-and-cows/server/pkg/player"
	"github.com/NikolovNikolay/bulls-and-cows/server/pkg/services"
	"github.com/NikolovNikolay/bulls-and-cows/server/pkg/utils"
	"github.com/googollee/go-socket.io"
)

// Socket gives access to some real-time
// communication via socket.io
type Socket struct {
	Socket      *socketio.Server
	roomMap     map[string]int
	roomTypeMap map[string]int
	db          *mgo.Session
	inTest      bool
	bcChecker   utils.BCCheck
}

const (
	evtConnection      = "connection"
	evtDisconnect      = "disconnect"
	evtError           = "error"
	evtCreateGame      = "creater"
	evtJoinGame        = "joinr"
	evtJoinedAGame     = "joinmy"
	evtGetActiveGames  = "getavr"
	evtInputGuess      = "inputguess"
	evtUpdActiveGames  = "updater"
	evtConfirmJoinGame = "confjoin"
	evtStartP2P        = "startp2p"
	evtMakeGuess       = "makeguess"
	evtGameEnd         = "gameend"

	roomDefault = "defroom"

	// SockIoEndpoint is socket.io server route
	SockIoEndpoint = "/socket.io/"
)

// New returns a new instance of SocketController
func New(socket *socketio.Server, db *mgo.Session, inTest bool) *Socket {
	s := Socket{
		Socket: socket,
		db:     db,
		inTest: inTest}
	s.bcChecker = utils.BCCheck{}
	s.roomMap = make(map[string]int)
	s.roomTypeMap = make(map[string]int)
	return &s
}

// Init configures the socket.io server - sets the
// custom event listeners
func (s Socket) Init() error {
	e := s.Socket.On(evtConnection, func(so socketio.Socket) {
		log.Println(so.Id(), "connected via socket.io")
		var e error
		e = so.Join(roomDefault)
		e = so.On(evtCreateGame, s.createGameHandler(so))
		e = so.On(evtJoinGame, s.joinGameHandler(so))
		e = so.On(evtGetActiveGames, s.getGames(so))
		e = so.On(evtInputGuess, s.setPlayerGuessNumHandler(so))
		e = so.On(evtMakeGuess, s.makeGuessHandler(so))
		if e != nil {
			log.Println("There was an error trying to register custom socket events")
		}
	})
	if e != nil {
		log.Println("Could not connect to socket")
		return e
	}
	e = s.Socket.On(evtError, func(so socketio.Socket, err error) {
		log.Println("error:", err)
	})
	if e != nil {
		log.Println("Could not register error handler for socket")
		return e
	}
	e = s.Socket.On(evtDisconnect, func(so socketio.Socket, err error) {
		log.Println("disconnect:", err)
	})
	if e != nil {
		log.Println("Could not register disconnect handler for socket")
		return e
	}

	return nil
}

// Triggered when the client calls to host and
// create a new game
func (s *Socket) createGameHandler(
	so socketio.Socket) func(data string, gt int) string {

	return func(room string, gameType int) string {
		if s.roomMap[room] > 0 {
			return ""
		}
		s.increaseRoomParticipants(room, gameType)

		dbName := getTargetDBName(s)
		var e error
		if !s.inTest {
			// we broadcast the client in the lets say it lobby
			// for the available rooms to join
			e = so.BroadcastTo(
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
		}

		p := player.New(
			room,
			false,
			utils.DBName,
			utils.GetDBSession())
		e = p.Add(dbName, s.db)

		if e != nil {
			log.Println("Could n insert new player to DB")
			return ""
		}
		g := game.New(utils.P2P)
		p.LogIn(&g.ID)
		g.AddPlayer(p)
		e = g.Add(dbName, s.db)

		if e != nil {
			log.Println("Could not save game in DB")
			return ""
		}

		e = p.Update(dbName, s.db)
		if e != nil {
			log.Println("Could not update player in DB")
			return ""
		}

		return g.ID.Hex()
	}
}

// Triggered when the client calls to join a game
func (s *Socket) joinGameHandler(
	so socketio.Socket) func(a, b string) string {

	return func(room, rival string) string {
		if (s.roomMap[room] == 2 && s.roomTypeMap[room] == utils.P2P) || room == rival {
			return ""
		}

		dbName := getTargetDBName(s)
		s.increaseRoomParticipants(room, 0)
		if !s.inTest {
			e := so.Join(room)
			if e != nil {
				log.Printf("Could not join '%s", room)
				return ""
			}
			e = so.BroadcastTo(room, evtJoinedAGame, rival)
			if e != nil {
				log.Printf("Could not broadcast event to '%s' room", roomDefault)
			}
		}

		host, e := player.FindByName(room, dbName, s.db)
		if e != nil {
			log.Printf("Could not find host of '%s'", room)
			return ""
		}
		gID := host.LoggedIn
		g, e := game.FindByID(gID.Hex(), dbName, s.db)
		if e != nil {
			log.Printf("Could not get game with host '%s'", host.Name)
			return ""
		}
		p := player.New(
			rival,
			false,
			dbName,
			s.db)
		p.LogIn(gID)
		e = p.Add(dbName, s.db)
		if e != nil {
			log.Printf("Could not add player '%s' to DB when attempting to join", p.Name)
		}
		g.AddPlayer(p)
		e = g.Update(dbName, s.db)
		if e != nil {
			log.Printf("Could not update game '%s' in DB", g.ID)
		}
		if !s.inTest {
			e = so.BroadcastTo(room, evtConfirmJoinGame)
			if e != nil {
				log.Printf("Could not broadcast join to '%s' room", room)
			}

		}
		return g.ID.Hex()
	}
}

// Triggered from the client and returns the
// available for a player rooms to join
func (s *Socket) getGames(
	so socketio.Socket) func(a string) []string {

	return func(playerName string) []string {
		avrooms := make([]string, 1)
		for k, v := range s.roomMap {
			if v < 2 && k != playerName {
				avrooms = append(avrooms, []string{k}...)
			}
		}

		return avrooms
	}
}

// Triggered from the client when attempting to set
// a its number for the other player to guess
func (s *Socket) setPlayerGuessNumHandler(
	so socketio.Socket) func(a, b, c string) bool {

	return func(room, rawGuess, playerName string) bool {
		valid := s.bcChecker.ValidateMadeGuess(rawGuess)
		if !valid {
			log.Printf("Invalid guess")
			return false
		}
		parsedNum, e := strconv.Atoi(rawGuess)
		if e != nil {
			log.Printf("Could not parse the guess num")
			return false
		}

		dbName := getTargetDBName(s)
		p, e := player.FindByName(playerName, dbName, s.db)
		if e != nil {
			log.Printf("Could not find player '%s'", playerName)
			return false
		}
		g, e := game.FindByID(p.LoggedIn.Hex(), dbName, s.db)
		if e != nil {
			log.Printf("Could not find game '%s'", g.ID)
			return false
		}
		for _, player := range g.Players {
			if player.ID == p.ID {
				player.Number = parsedNum
			}
		}

		var ready bool
		if ready = s.checkIfReadyForP2PStart(g); ready && !s.inTest {
			// All players have set their numbers, so we return the host
			// name, which is with index 0
			e = so.BroadcastTo(room, evtStartP2P, g.Players[0].Name)
			if e != nil {
				log.Printf("Could not broadcast player name to room %s", room)
			}
		}
		if ready {
			g.Start(time.Now().Unix())
		}
		e = g.Update(dbName, s.db)
		if e != nil {
			log.Printf("Could not update game '%s'", g.ID)
			return false
		}
		if ready {
			return true
		}
		return true
	}
}

// Triggered when a player sends a number, attempting to guess
// hit opponent's number
func (s *Socket) makeGuessHandler(
	so socketio.Socket) func(a, b, c, d string) *services.GuessPayload {
	return func(room, rawGuess, playerName, gameID string) *services.GuessPayload {
		valid := s.bcChecker.ValidateMadeGuess(rawGuess)
		if !valid {
			log.Printf("Could not parse guess")
			return nil
		}
		guessInt, e := strconv.Atoi(rawGuess)
		if e != nil {
			log.Printf("Could not parse guess")
			return nil
		}
		dbName := getTargetDBName(s)
		g, e := game.FindByID(gameID, dbName, s.db)
		if e != nil {
			log.Printf("Could not get game from DB")
			return nil
		}
		rn, e := s.getP2PRivalNumber(playerName, g)
		if e != nil {
			log.Printf(e.Error())
			return nil
		}
		res := s.bcChecker.Check(rn, guessInt)
		var pl *player.Player
		for _, p := range g.Players {
			if p.Name == playerName {
				pl = p
			}
		}

		pl.Guesses = append(pl.Guesses, []*guess.Guess{&guess.Guess{Bc: res, Number: guessInt}}...)

		win := false
		if res.Bulls == 4 {
			win = true
			g.End(time.Now().Unix())
		}
		e = g.Update(dbName, s.db)
		if e != nil {
			log.Printf("Could not get update game in DB")
			return nil
		}
		gp := services.GuessPayload{
			BC:      res,
			Guesses: pl.Guesses,
			Win:     win,
			Time:    g.EndTime - g.StartTime}
		if win && !s.inTest {
			e = so.BroadcastTo(room, evtGameEnd)
			if e != nil {
				log.Printf("Could not broadcast player name to room %s", room)
			}
		}

		return &gp
	}
}

// checkIfReadyForP2PStart checks the number of
// all players in a game. If all of them have
// input their own guesses, the game can start
func (s *Socket) checkIfReadyForP2PStart(g *game.Game) bool {
	ready := 0
	for _, p := range g.Players {
		if p.Number != 0 {
			ready = ready + 1
		}
	}
	if ready == len(g.Players) {
		return true
	}

	return false
}

func (s *Socket) increaseRoomParticipants(room string, gameType int) {
	s.roomMap[room] = s.roomMap[room] + 1
	if gameType > 0 {
		s.roomTypeMap[room] = gameType
	}
}

func (s *Socket) getP2PRivalNumber(currentPlayer string, g *game.Game) (int, error) {
	for _, p := range g.Players {
		if p.Name != currentPlayer {
			return p.Number, nil
		}
	}

	return 0, errors.New("Could not find rival's number")
}

func getTargetDBName(s *Socket) string {

	dbName := utils.DBName
	if s.inTest {
		dbName = utils.DBNameTest
	}

	return dbName
}
