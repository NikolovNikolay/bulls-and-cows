package services

import (
	"errors"
	"net/http"
	"time"

	"github.com/NikolovNikolay/bulls-and-cows/server/pkg/game"
	"github.com/NikolovNikolay/bulls-and-cows/server/pkg/response"
	"github.com/gorilla/mux"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// GetGameService returns data for a specific game from the DB
type GetGameService struct {
	db *mgo.Session
}

// NewGetGameService returns a new GetGameService instance
func NewGetGameService(db *mgo.Session) GetGameService {
	return GetGameService{db: db}
}

type gamePayload GuessPayload

// Endpoint returns the endpoint of the service
func (gg GetGameService) Endpoint() string {
	return "/api/game/{gameID}"
}

// Method returns the method of the service
func (gg GetGameService) Method() string {
	return "GET"
}

// Handle is the handle function used to register in the mux
func (gg GetGameService) Handle(w http.ResponseWriter, r *http.Request) {
	response := response.New(200, "", nil)

	vars := mux.Vars(r)
	gameID := vars["gameID"]
	if gameID == "" || !bson.IsObjectIdHex(gameID) {
		response.Status = http.StatusBadRequest
		response.Error = errors.New("Invalid gameID parameter").Error()
		DefSendResponseBeh(w, response)
		return
	}
	dbName := getTargetDbName(r)

	g, e := game.FindByID(gameID, dbName, gg.db)
	if e != nil || g.StartTime == 0 {
		response.Error = "Not a valid guess - not referring to a game"
		response.Status = http.StatusBadRequest
		DefSendResponseBeh(w, response)
		return
	}

	response.Payload = gamePayload{
		BC:      nil,
		Guesses: g.PlayerOneGuesses,
		Win:     g.EndTime != 0,
		Time:    time.Now().Unix() - g.StartTime}
	DefSendResponseBeh(w, response)
}
