package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/NikolovNikolay/bulls-and-cows/server/models"
	"github.com/NikolovNikolay/bulls-and-cows/server/response"
	"github.com/julienschmidt/httprouter"
	uuid "github.com/nu7hatch/gouuid"
)

// GameController holds methods for managing a game session
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

// Init initializes a game process
func (gc GameController) Init(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	r.ParseForm()

	// Gettings the player's username
	userName := r.Form.Get("userName")
	if userName == "" {
		// if no username is set then we autogenerate one
		rand.Seed(time.Now().UnixNano())
		randPostfix := rand.Intn(10000)
		userName = fmt.Sprintf("user%d", randPostfix)
	}

	// Searching if a player exists
	p, err := gc.pc.findPlayerByName(userName)
	if err != nil {
		// if he dows not exist we get an error, then we create it
		p = models.Player{
			ID:     bson.NewObjectId(),
			Name:   userName,
			Wins:   0,
			Logged: false}
		if insErr := gc.pc.createPlayer(p); insErr != nil {
			log.Fatal(insErr)
		}
	}

	response := response.New(200, "", nil)

	if p.Logged == true {
		// if the player exists and he is logged in we return forbidden - the new guy can not continue with this name
		response.Status = http.StatusForbidden
		response.Error = errors.New("Player with that name has already logged").Error()
	} else {
		// if the player exists and he is not logged in then its ok
		p.Logged = true

		// we update the player in the DB as now logged in and continue
		if e := gc.pc.updatePlayer(p); e != nil {
			log.Fatal(e)
		}

		u, err := uuid.NewV4()
		if err != nil {
			response.Status = http.StatusInternalServerError
			response.Error = errors.New("Could not initialize game session").Error()
		} else {
			response.Payload = initResponse{GameSessionID: u.String()}
		}
	}

	w.WriteHeader(response.Status)
	json.NewEncoder(w).Encode(response)
}
