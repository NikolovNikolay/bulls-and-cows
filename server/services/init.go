package services

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/NikolovNikolay/bulls-and-cows/server/controllers"

	"gopkg.in/mgo.v2/bson"

	"github.com/NikolovNikolay/bulls-and-cows/server/models"
	"github.com/NikolovNikolay/bulls-and-cows/server/response"
	"github.com/NikolovNikolay/bulls-and-cows/server/utils"
)

// InitService initiates a new game and player in DB
type InitService struct {
	gameControler *controllers.GameController
	numGen        utils.NumGen
}

type initPayload struct {
	GameSessionID string `json:"gameID"`
	Guess         string `json:"guess"`
	PlayerName    string `json:"name"`
}

// NewInitService returns a new instance of InitService
func NewInitService(gc *controllers.GameController) InitService {
	return InitService{gameControler: gc, numGen: utils.GetNumGen()}
}

// Endpoint returns the endpoint of the service
func (is InitService) Endpoint() string {
	return "/api/init"
}

// Method returns the method of the service
func (is InitService) Method() string {
	return "POST"
}

// Handle is the handle function used to register in the mux
func (is InitService) Handle(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	response := response.New(200, "", nil)
	dbName := getTargetDbName(r)

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
		DefSendResponseBeh(w, response)
		return
	}

	// Searching if a player exists
	p, err := is.gameControler.Pc.FindPlayerByName(userName, utils.DBName)
	if err != nil {
		// if he dows not exist we get an error, then we create it
		p = generateNewPlayer(userName)
		if insErr := is.gameControler.Pc.CreatePlayer(p, utils.DBName); insErr != nil {
			log.Fatal(insErr)
			response.Status = http.StatusInternalServerError
			response.Error = insErr.Error()
			DefSendResponseBeh(w, response)
			return
		}
	}

	payload := &initPayload{}
	if p.Logged == true {
		// If the player is logged to a game, then we
		// would like to resume. We send the game session id,
		// so at a later point the front-end can make a
		// request to download this game's data.
		payload.GameSessionID = p.LoggedIn.Hex()
		payload.PlayerName = p.Name
		response.Payload = payload
		DefSendResponseBeh(w, response)
		return
	}

	g, e := createNewGame(is.gameControler, dbName, gt, &p.ID, is.numGen.Gen())
	if e != nil {
		response.Error = e.Error()
		response.Status = http.StatusBadRequest
		DefSendResponseBeh(w, response)
		return
	}

	// if the player exists and he is not logged
	// in a game then we mark him as now logged
	p.Logged = true
	p.LoggedIn = &g.GameID

	// we update the player in DB
	if e := is.gameControler.Pc.UpdatePlayer(p, utils.DBName); e != nil {
		response.Error = e.Error()
		response.Status = http.StatusInternalServerError
		log.Fatal(e)
		DefSendResponseBeh(w, response)
		return
	}

	payload.GameSessionID = g.GameID.Hex()
	payload.PlayerName = p.Name

	// If Comp. vs. Comp mode is selected
	// we return the generated guess also
	if gt == utils.CVC {
		payload.Guess = strconv.Itoa(g.GuessNum)
	}

	response.Payload = payload
	DefSendResponseBeh(w, response)
}

func createNewGame(
	gc *controllers.GameController,
	dbName string,
	gt int,
	pID *bson.ObjectId,
	guessNum int) (*models.Game, error) {

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

	if er := gc.Session.DB(dbName).C(utils.DBCGames).Insert(game); er != nil {
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
