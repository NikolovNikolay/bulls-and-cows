package services

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"gopkg.in/mgo.v2"

	"github.com/NikolovNikolay/bulls-and-cows/server/pkg/game"
	"github.com/NikolovNikolay/bulls-and-cows/server/pkg/player"
	"github.com/NikolovNikolay/bulls-and-cows/server/pkg/response"
	"github.com/NikolovNikolay/bulls-and-cows/server/pkg/utils"
)

const paramUserNameKey = "userName"
const paramGameTypeKey = "gameType"

// InitService initiates a new game and player in DB
type InitService struct {
	numGen utils.NumGen
	db     *mgo.Session
}

type initPayload struct {
	GameSessionID string `json:"gameID"`
	Guess         string `json:"guess"`
	PlayerName    string `json:"name"`
}

// NewInitService returns a new instance of InitService
func NewInitService(db *mgo.Session) *InitService {
	return &InitService{
		numGen: utils.GetNumGen(),
		db:     db}
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
	response := response.New(200, "", nil)
	if e := r.ParseForm(); e != nil {
		response.Error = e.Error()
		response.Status = http.StatusBadRequest
		DefSendResponseBeh(w, response)
	}

	dbName := getTargetDbName(r)
	db := utils.GetDBSession()

	// Gettings the player's username
	userName := getVarFromRequest(r, paramUserNameKey)
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
	p, err := player.FindByName(userName, utils.DBName, is.db)
	if err != nil {
		// if he dows not exist we get an error, then we create it
		p = player.New(userName, false, dbName, db)

		if insErr := p.Add(dbName, db); insErr != nil {
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

	g := game.New(gt)
	g.AddPlayer(p)
	g.SetNumber(is.numGen.Gen())
	g.Start(time.Now().Unix())
	e = g.Add(dbName, db)
	if e != nil {
		response.Error = e.Error()
		response.Status = http.StatusBadRequest
		DefSendResponseBeh(w, response)
		return
	}

	// if the player exists and he is not logged
	// in a game then we mark him as now logged
	p.LogIn(&g.ID)

	// we update the player in DB
	if e := p.Update(dbName, utils.GetDBSession()); e != nil {
		response.Error = e.Error()
		response.Status = http.StatusInternalServerError
		log.Fatal(e)
		DefSendResponseBeh(w, response)
		return
	}

	payload.GameSessionID = g.ID.Hex()
	payload.PlayerName = p.Name

	// If Comp. vs. Comp mode is selected
	// we return the generated guess also
	if gt == utils.CVC {
		payload.Guess = strconv.Itoa(g.Number)
	}

	response.Payload = payload
	DefSendResponseBeh(w, response)
}
