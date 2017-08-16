package player

import (
	"errors"

	"github.com/NikolovNikolay/bulls-and-cows/server/pkg/utils"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Player represent a human player
type Player struct {
	ID       bson.ObjectId  `json:"id" bson:"_id"`
	Name     string         `json:"name" bson:"name"`
	Logged   bool           `json:"logged" bson:"logged"`
	LoggedIn *bson.ObjectId `json:"loggedId" bson:"loggedIn"`
}

// New returns a new player instance
func New(
	id bson.ObjectId,
	name string,
	logged bool) *Player {
	return &Player{
		ID:       id,
		Name:     name,
		Logged:   logged,
		LoggedIn: nil}
}

// FindByName finds a player by name in DB
func FindByName(
	name,
	dbName string,
	mon *mgo.Session) (p *Player, e error) {
	c := mon.DB(dbName).C(utils.DBCPlayers)
	result := Player{}
	err := c.Find(bson.M{"name": name}).One(&result)

	return &result, err
}

// FindByID finds a player by its ID in DB
func FindByID(
	id,
	dbName string,
	mon *mgo.Session) (p *Player, e error) {
	result := Player{}
	if bson.IsObjectIdHex(id) {
		c := mon.DB(dbName).C(utils.DBCPlayers)
		err := c.FindId(bson.ObjectIdHex(id)).One(&result)

		return &result, err
	}

	return &result, errors.New("Invalid player id")
}

// AddToDB creates a player in the DB
func AddToDB(
	p *Player,
	dbName string,
	mon *mgo.Session) (e error) {
	c := mon.DB(dbName).C(utils.DBCPlayers)
	err := c.Insert(p)

	return err
}

// Update updates a player in the DB
func Update(
	p *Player,
	dbName string,
	mon *mgo.Session) (e error) {
	c := mon.DB(dbName).C(utils.DBCPlayers)
	err := c.Update(bson.M{"_id": p.ID}, p)

	return err
}
