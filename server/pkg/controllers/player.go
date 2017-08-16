package controllers

import (
	"errors"

	"github.com/NikolovNikolay/bulls-and-cows/server/pkg/models"
	"github.com/NikolovNikolay/bulls-and-cows/server/pkg/utils"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// PlayerController holds methods for managing a players
type PlayerController struct {
	session *mgo.Session
}

// NewPlayerController returns a new instance of GameController
func NewPlayerController(s *mgo.Session) PlayerController {
	return PlayerController{s}
}

// FindPlayerByName finds a player by name in DB
func (pc *PlayerController) FindPlayerByName(name, dbName string) (p models.Player, e error) {
	c := pc.session.DB(dbName).C(utils.DBCPlayers)
	result := models.Player{}
	err := c.Find(bson.M{"name": name}).One(&result)

	return result, err
}

// FindPlayerByID finds a player by its ID in DB
func (pc *PlayerController) FindPlayerByID(id, dbName string) (p models.Player, e error) {
	result := models.Player{}
	if bson.IsObjectIdHex(id) {
		c := pc.session.DB(dbName).C(utils.DBCPlayers)
		err := c.FindId(bson.ObjectIdHex(id)).One(&result)

		return result, err
	}

	return result, errors.New("Invalid player id")
}

// CreatePlayer creates a player in the DB
func (pc *PlayerController) CreatePlayer(p models.Player, dbName string) (e error) {
	c := pc.session.DB(dbName).C(utils.DBCPlayers)
	err := c.Insert(p)

	return err
}

// UpdatePlayer updates a player in the DB
func (pc *PlayerController) UpdatePlayer(p models.Player, dbName string) (e error) {
	c := pc.session.DB(dbName).C(utils.DBCPlayers)
	err := c.Update(bson.M{"_id": p.ID}, p)

	return err
}
