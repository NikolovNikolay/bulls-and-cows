package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/NikolovNikolay/bulls-and-cows/server/models"
	"github.com/NikolovNikolay/bulls-and-cows/server/response"
	"github.com/NikolovNikolay/bulls-and-cows/server/utils"
	"github.com/julienschmidt/httprouter"
)

var numGen utils.NumGen
var bcChecker utils.BCCheck

func init() {
	numGen = utils.GetNumGen()
	bcChecker = utils.BCCheck{}
}

// GameController holds methods for managing a game session including http handlers
type GameController struct {
	session *mgo.Session
	pc      PlayerController
}

// NewGameController returns a new instance of GameController
func NewGameController(s *mgo.Session, pc PlayerController) *GameController {
	return &GameController{s, pc}
}

type initPayload struct {
	GameSessionID string `json:"gameID"`
	Guess         string `json:"guess"`
	PlayerName    string `json:"name"`
}

type gamePayload struct {
	BC      *utils.BCCheckResult `json:"bc"`
	Guesses []int                `json:"m"`
	Win     bool                 `json:"win"`
	Time    int64                `json:"t"`
}

// InitHandler initializes a game process
func (gc GameController) InitHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	r.ParseForm()
	response := response.New(200, "", nil)
	payload := &initPayload{}

	// Gettings the player's username
	userName := r.PostFormValue("userName")
	if userName == "" {
		// if no username is set then we autogenerate one
		userName = generateUserName()
	}

	gt, e := validateGameTypeParam(r)
	if e != nil {
		response.Error = e.Error()
		response.Status = http.StatusBadRequest
		sendResponse(w, response)
		return
	}

	dbName := getTargetDbName(r)

	// Searching if a player exists
	p, err := gc.pc.findPlayerByName(userName, utils.DBName)
	if err != nil {
		// if he dows not exist we get an error, then we create it
		p = generateNewPlayer(userName)
		if insErr := gc.pc.createPlayer(p, utils.DBName); insErr != nil {
			log.Fatal(insErr)
			response.Status = http.StatusInternalServerError
			response.Error = insErr.Error()
			sendResponse(w, response)
			return
		}
	}

	if p.Logged == true {
		// if the player exists and he is logged in we return forbidden - the new guy can not continue with this name
		payload.GameSessionID = p.LoggedIn.Hex()
		payload.PlayerName = p.Name
		response.Payload = payload
		sendResponse(w, response)
		return
	}

	g, e := createNewGame(gc, dbName, gt, &p.ID, numGen.Gen())
	if e != nil {
		response.Error = e.Error()
		response.Status = http.StatusBadRequest
		sendResponse(w, response)
		return
	}

	// if the player exists and he is not logged in then its ok
	p.Logged = true
	p.LoggedIn = &g.GameID

	// we update the player in the DB as now logged in and continue
	if e := gc.pc.updatePlayer(p, utils.DBName); e != nil {
		response.Error = e.Error()
		response.Status = http.StatusInternalServerError
		log.Fatal(e)
		sendResponse(w, response)
		return
	}

	payload.GameSessionID = g.GameID.Hex()
	payload.PlayerName = p.Name

	if gt == utils.CVC {
		payload.Guess = strconv.Itoa(g.GuessNum)
	}

	response.Payload = payload
	sendResponse(w, response)
}

// GuessHandler takes the player's guess number and returns the following bulls and cows
func (gc GameController) GuessHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	response := response.New(200, "", nil)
	r.ParseForm()
	guess, e := validateGuessNumberParam(ps)
	if e != nil {
		response.Error = e.Error()
		response.Status = http.StatusBadRequest
		sendResponse(w, response)
		return
	}

	dbName := getTargetDbName(r)
	gID := r.PostFormValue("X-GameID")
	if gID == "" {
		response.Error = "Not a valid guess - not referring to a game"
		response.Status = http.StatusBadRequest
		sendResponse(w, response)
		return
	}

	g, e := gc.getGameByID(gID, dbName)
	if e != nil || g.StartTime == 0 {
		response.Error = "Not a valid guess - not referring to a game"
		response.Status = http.StatusBadRequest
		sendResponse(w, response)
		return
	}

	if g.EndTime != 0 {
		response.Payload = gamePayload{
			BC:      &utils.BCCheckResult{Bulls: 4, Cows: 0},
			Guesses: g.PlayerOneGuesses,
			Win:     true,
			Time:    g.EndTime - g.StartTime}
		sendResponse(w, response)
		return
	}

	br := bcChecker.Check(g.GuessNum, guess)
	guesses := append(g.PlayerOneGuesses, []int{guess}...)

	g.PlayerOneGuesses = guesses

	// check if winner
	var win = false
	now := time.Now().Unix()
	var t = (now - g.StartTime)
	if br.Bulls == 4 {
		win = true
		g.EndTime = now
		p, _ := gc.pc.findPlayerByID(g.PlayerOneID.Hex(), dbName)
		p.Logged = false
		p.LoggedIn = nil
		e = gc.pc.updatePlayer(p, dbName)
	}

	e = gc.updateGameByID(g, dbName)
	if e != nil {
		response.Status = http.StatusInternalServerError
		response.Error = e.Error()
		sendResponse(w, response)
	}

	response.Payload = gamePayload{BC: br, Guesses: g.PlayerOneGuesses, Win: win, Time: t}
	sendResponse(w, response)
}

// GetGameDataHandler returns data for a game session
func (gc GameController) GetGameDataHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	response := response.New(200, "", nil)
	r.ParseForm()
	gameID := ps.ByName("gameID")
	if gameID == "" || !bson.IsObjectIdHex(gameID) {
		response.Status = http.StatusBadRequest
		response.Error = errors.New("Invalid gameID parameter").Error()
		sendResponse(w, response)
		return
	}
	dbName := getTargetDbName(r)
	g, e := gc.getGameByID(gameID, dbName)
	if e != nil || g.StartTime == 0 {
		response.Error = "Not a valid guess - not referring to a game"
		response.Status = http.StatusBadRequest
		sendResponse(w, response)
		return
	}

	response.Payload = gamePayload{BC: nil, Guesses: g.PlayerOneGuesses, Win: g.EndTime != 0, Time: time.Now().Unix() - g.StartTime}
	sendResponse(w, response)
}

// Helper function

func (gc GameController) getGameByID(gameID string, dbName string) (models.Game, error) {
	game := models.Game{}
	var err error
	if bson.IsObjectIdHex(gameID) {
		err = gc.session.DB(dbName).C(utils.DBCGames).FindId(bson.ObjectIdHex(gameID)).One(&game)
	} else {
		err = errors.New("Invalid gameID")
	}

	return game, err
}

func (gc GameController) updateGameByID(g models.Game, dbName string) error {
	if bson.IsObjectIdHex(g.GameID.Hex()) {
		e := gc.session.DB(dbName).C(utils.DBCGames).UpdateId(g.GameID, g)
		if e != nil {
			return e
		}

		return nil
	}
	return errors.New("Invalid gameID")

}

func sendResponse(w http.ResponseWriter, response *response.Response) {
	w.WriteHeader(response.Status)
	json.NewEncoder(w).Encode(response)
}

func createNewGame(gc GameController, dbName string, gt int, pID *bson.ObjectId, guessNum int) (*models.Game, error) {
	gameID := bson.NewObjectId()
	game, e := models.NewGame(
		gameID,
		gt,
		pID,
		nil,
		guessNum)

	if e != nil {
		return nil, e
	}

	if er := gc.session.DB(dbName).C(utils.DBCGames).Insert(game); er != nil {
		return nil, er
	}

	return game, nil
}

func generateUserName() string {
	rand.Seed(time.Now().UnixNano())
	randPostfix := rand.Intn(10000)
	return fmt.Sprintf("user%d", randPostfix)
}

func generateNewPlayer(name string) models.Player {
	return models.Player{
		ID:       bson.NewObjectId(),
		Name:     name,
		Wins:     0,
		Logged:   false,
		LoggedIn: nil}
}

func getTargetDbName(r *http.Request) string {
	th := r.Header.Get("x-test")
	isTesting := th != ""
	if isTesting {
		return utils.DBNameTest
	}

	return utils.DBName
}

func validateGameTypeParam(r *http.Request) (int, error) {
	gameType := r.PostFormValue("gameType")
	if gameType == "" {
		return 0, errors.New("Missing parameter for game type")
	}

	gt, e := strconv.Atoi(gameType)
	if e != nil {
		return 0, errors.New("Could not parse game type parameter")
	}
	return gt, nil
}

func validateGuessNumberParam(ps httprouter.Params) (int, error) {
	g := ps.ByName("guess")
	if g == "" {
		return -1, errors.New("Missing parameter guess")
	}

	if !bcChecker.ValidateMadeGuess(g) {
		return -1, errors.New("Invalid guess number")
	}

	return strconv.Atoi(g)
}
