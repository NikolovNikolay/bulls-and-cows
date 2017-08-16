package controllers

import (
	"errors"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/NikolovNikolay/bulls-and-cows/server/models"
	"github.com/NikolovNikolay/bulls-and-cows/server/utils"
)

// GameController holds methods for managing a
// game session including http handlers
type GameController struct {
	Session *mgo.Session
	Pc      PlayerController
}

// NewGameController returns a new instance of GameController
func NewGameController(
	s *mgo.Session,
	pc PlayerController) GameController {
	return GameController{s, pc}
}

// CreateNewGame creates a new game in DB
func (gc GameController) CreateNewGame(
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
		guessNum,
		0)

	if e != nil {
		return nil, e
	}

	if er := gc.Session.DB(dbName).C(utils.DBCGames).Insert(game); er != nil {
		return nil, er
	}

	return game, nil
}

// GetGameByID finds a game bt ID in DB
func (gc GameController) GetGameByID(
	gameID string,
	dbName string) (models.Game, error) {
	game := models.Game{}
	var err error
	if bson.IsObjectIdHex(gameID) {
		err = gc.Session.DB(
			dbName).C(
			utils.DBCGames).FindId(
			bson.ObjectIdHex(gameID)).One(&game)
	} else {
		err = errors.New("Invalid gameID")
	}
	return game, err
}

// UpdateGameByID finds a game by ID in DB
// then updates it
func (gc GameController) UpdateGameByID(
	g models.Game,
	dbName string) error {
	if bson.IsObjectIdHex(g.GameID.Hex()) {
		e := gc.Session.DB(
			dbName).C(
			utils.DBCGames).UpdateId(
			g.GameID, g)
		if e != nil {
			return e
		}

		return nil
	}
	return errors.New("Invalid gameID")
}
