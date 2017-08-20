package services

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"gopkg.in/mgo.v2"

	"github.com/NikolovNikolay/bulls-and-cows/server/pkg/game"
	"github.com/NikolovNikolay/bulls-and-cows/server/pkg/guess"
	"github.com/NikolovNikolay/bulls-and-cows/server/pkg/player"
	"github.com/NikolovNikolay/bulls-and-cows/server/pkg/response"
	"github.com/NikolovNikolay/bulls-and-cows/server/pkg/utils"
)

const (
	paramGuessKey  = "guessNum"
	paramGameIDKey = "gameID"
)

// GuessService returns results to client,
// that contains how many bulls and cows
// his number has
type GuessService struct {
	bcChecker utils.BCCheck
	db        *mgo.Session
}

// GuessPayload is the payload of the GuessService
type GuessPayload struct {
	BC      *utils.BCCheckResult `json:"bc"`
	Guesses []*guess.Guess       `json:"m"`
	Win     bool                 `json:"win"`
	Time    int64                `json:"t"`
}

// NewGuessService returns a new instance of GuessService
func NewGuessService(db *mgo.Session) *GuessService {
	return &GuessService{
		bcChecker: utils.BCCheck{},
		db:        db}
}

// Endpoint returns the endpoint of the service
func (gs GuessService) Endpoint() string {
	return "/api/guess/{guessNum:[0-9]+}"
}

// Method returns the method of the service
func (gs GuessService) Method() string {
	return "PUT"
}

// Handle is the handle function used to register in the mux
func (gs GuessService) Handle(w http.ResponseWriter, r *http.Request) {
	response := response.New(200, "", nil)
	parseForm(r)
	newGuess, e := gs.validateNumber(getVarFromRequest(r, paramGuessKey))
	if e != nil {
		response.Error = e.Error()
		response.Status = http.StatusBadRequest
		DefSendResponseBeh(w, response)
		return
	}

	dbName := getTargetDbName(r)
	db := utils.GetDBSession()
	gID := r.PostFormValue(paramGameIDKey)
	if gID == "" {
		response.Error = "Request not referring to a game"
		response.Status = http.StatusBadRequest
		DefSendResponseBeh(w, response)
		return
	}

	g, e := game.FindByID(gID, dbName, gs.db)
	if e != nil || g.StartTime == 0 {
		response.Error = "Request not referring to a game"
		response.Status = http.StatusBadRequest
		DefSendResponseBeh(w, response)
		return
	}
	if g.EndTime != 0 {
		response.Payload = GuessPayload{
			BC:      &utils.BCCheckResult{Bulls: 4, Cows: 0},
			Guesses: g.Players[0].Guesses,
			Win:     true,
			Time:    g.EndTime - g.StartTime}
		DefSendResponseBeh(w, response)
		return
	}

	br := gs.bcChecker.Check(g.Number, newGuess)
	dbGuess := guess.New(newGuess, br.Bulls, br.Cows)
	guesses := append(g.Players[0].Guesses, []*guess.Guess{&dbGuess}...)
	g.Players[0].Guesses = guesses

	// check if winner
	var win = false
	now := time.Now().Unix()
	var t = (now - g.StartTime)
	if br.Bulls == 4 {
		win = true
		g.End(now)
		p, _ := player.FindByID(g.Players[0].ID.Hex(), dbName, gs.db)
		p.LogOut()
		e = p.Update(dbName, db)
	}

	e = g.Update(dbName, db)
	if e != nil {
		response.Status = http.StatusInternalServerError
		response.Error = e.Error()
		DefSendResponseBeh(w, response)
	}

	response.Payload = GuessPayload{
		BC:      br,
		Guesses: g.Players[0].Guesses,
		Win:     win,
		Time:    t}
	DefSendResponseBeh(w, response)
}

func (gs GuessService) validateNumber(guess string) (int, error) {
	if guess == "" {
		return -1, errors.New("Missing parameter guess")
	}
	if !gs.bcChecker.ValidateMadeGuess(guess) {
		return -1, errors.New("Invalid guess number")
	}

	return strconv.Atoi(guess)
}
