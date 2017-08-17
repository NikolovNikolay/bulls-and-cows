package player

import (
	"errors"

	"github.com/NikolovNikolay/bulls-and-cows/server/pkg/guess"

	"github.com/NikolovNikolay/bulls-and-cows/server/pkg/utils"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Player represent a human player
type Player struct {
	ID       bson.ObjectId  `json:"id" bson:"_id"`
	Name     string         `json:"name" bson:"name"`
	Guesses  []*guess.Guess `json:"guesses" bson:"guesses"`
	Logged   bool           `json:"logged" bson:"logged"`
	LoggedIn *bson.ObjectId `json:"loggedId" bson:"loggedIn"`
	Number   int            `json:"number" bson:"number"`
}

// New returns a new player instance
func New(
	name string,
	logged bool,
	dbName string,
	db *mgo.Session) *Player {
	return &Player{
		ID:       bson.NewObjectId(),
		Name:     name,
		Logged:   logged,
		LoggedIn: nil}
}

// Add creates a player in the DB
func (p *Player) Add(dbName string, db *mgo.Session) (e error) {
	c := db.DB(dbName).C(utils.DBCPlayers)
	err := c.Insert(p)

	return err
}

// Update updates a player in the DB
func (p *Player) Update(dbName string, db *mgo.Session) (e error) {
	c := db.DB(dbName).C(utils.DBCPlayers)
	err := c.Update(bson.M{"_id": p.ID}, p)

	return err
}

// LogIn marks the player as logged into a game
func (p *Player) LogIn(gID *bson.ObjectId) {
	p.Logged = true
	p.LoggedIn = gID
}

// LogOut marks the player as logged out of a game
func (p *Player) LogOut() {
	p.Logged = false
	p.LoggedIn = nil
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
