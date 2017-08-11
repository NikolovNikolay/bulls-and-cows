package controllers

import (
	"github.com/NikolovNikolay/bulls-and-cows/server/models"
	"github.com/NikolovNikolay/bulls-and-cows/server/utils"
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

func (pc *PlayerController) findPlayerByName(name string) (p models.Player, e error) {
	c := pc.session.DB(utils.DBName).C(utils.DBCPlayers)
	result := models.Player{}
	err := c.Find(bson.M{"name": name}).One(&result)

	return result, err
}

func (pc *PlayerController) createPlayer(p models.Player) (e error) {
	c := pc.session.DB(utils.DBName).C(utils.DBCPlayers)
	err := c.Insert(p)

	return err
}

func (pc *PlayerController) updatePlayer(p models.Player) (e error) {
	c := pc.session.DB(utils.DBName).C(utils.DBCPlayers)
	err := c.Update(bson.M{"_id": p.ID}, p)

	return err
}
