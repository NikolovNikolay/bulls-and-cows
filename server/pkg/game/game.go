package game

import (
	"errors"

	"github.com/NikolovNikolay/bulls-and-cows/server/pkg/player"

	"gopkg.in/mgo.v2"

	"github.com/NikolovNikolay/bulls-and-cows/server/pkg/utils"
	"gopkg.in/mgo.v2/bson"
)

// Game represents a games session
type Game struct {
	ID        bson.ObjectId    `json:"gameID" bson:"_id"`
	StartTime int64            `josn:"startTime" bson:"startTime"`
	EndTime   int64            `json:"endTime" bson:"endTime"`
	GameType  int              `json:"type" bson:"type"`
	Players   []*player.Player `json:"players" bson:"players"`
	Number    int              `json:"number" bson:"number"`
}

// New returns a new Game instance
func New(
	gameType int,
	dbName string,
	db *mgo.Session) *Game {
	game := &Game{
		ID:       bson.NewObjectId(),
		GameType: gameType}

	return game
}

// Add creates a new game in DB
func (g *Game) Add(dbName string, db *mgo.Session) error {

	if er := db.DB(
		dbName).C(
		utils.DBCGames).Insert(
		g); er != nil {
		return er
	}

	return nil
}

// Start adds a start time to a game
func (g *Game) Start(t int64) {
	g.StartTime = t
}

// Update updates a game in DB
func (g *Game) Update(dbName string, db *mgo.Session) error {

	if bson.IsObjectIdHex(g.ID.Hex()) {
		e := db.DB(
			dbName).C(
			utils.DBCGames).UpdateId(
			g.ID, g)

		if e != nil {
			return e
		}

		return nil
	}
	return errors.New("Invalid gameID")
}

// End marks a game asa ended by adding end timestamp
func (g *Game) End(t int64) {
	g.EndTime = t
}

// AddPlayer adds a new player to the game
func (g *Game) AddPlayer(p *player.Player) {
	g.Players = append(g.Players, []*player.Player{p}...)
}

// GenNumber creates a number to guess
func (g *Game) GenNumber(num int) {
	g.Number = num
}

// FindByID finds a game bt ID in DB
func FindByID(
	gameID string,
	dbName string,
	mon *mgo.Session) (*Game, error) {

	game := Game{}
	var err error
	if bson.IsObjectIdHex(gameID) {
		err = mon.DB(
			dbName).C(
			utils.DBCGames).FindId(
			bson.ObjectIdHex(gameID)).One(&game)
		return &game, err
	}

	err = errors.New("Invalid gameID")
	return &game, err
}
