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

func init() {
	numGen = utils.GetNumGen()
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

type initResponse struct {
	GameSessionID string `json:"gameID"`
}

// InitHandler initializes a game process
func (gc GameController) InitHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	r.ParseForm()
	response := response.New(200, "", nil)

	// Gettings the player's username
	userName := r.Form.Get("userName")
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
		response.Status = http.StatusForbidden
		response.Error = errors.New("Player with that name has already logged").Error()
		sendResponse(w, response)
		return
	}

	// if the player exists and he is not logged in then its ok
	p.Logged = true

	// we update the player in the DB as now logged in and continue
	if e := gc.pc.updatePlayer(p, utils.DBName); e != nil {
		response.Error = e.Error()
		response.Status = http.StatusInternalServerError
		log.Fatal(e)
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

	response.Payload = initResponse{GameSessionID: g.GameID.Hex()}
	sendResponse(w, response)
}

// Helper function

func sendResponse(w http.ResponseWriter, response *response.Response) {
	w.WriteHeader(response.Status)
	json.NewEncoder(w).Encode(response)
}

func createNewGame(gc GameController, dbName string, gt *int, pID *bson.ObjectId, guessNum int) (*models.Game, error) {
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
		ID:     bson.NewObjectId(),
		Name:   name,
		Wins:   0,
		Logged: false}
}

func getTargetDbName(r *http.Request) string {
	th := r.Header.Get("x-test")
	isTesting := th != ""
	if isTesting {
		return utils.DBNameTest
	}

	return utils.DBName
}

func validateGameTypeParam(r *http.Request) (*int, error) {
	gameType := r.Form.Get("gameType")
	if gameType == "" {
		return nil, errors.New("Missing parameter game type")
	}

	gt, e := strconv.Atoi(gameType)
	if e != nil {
		return nil, errors.New("Could not parse game type parameter")
	}
	return &gt, nil
}
