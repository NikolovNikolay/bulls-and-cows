package game

import (
	"errors"
	"time"

	"gopkg.in/mgo.v2"

	"github.com/NikolovNikolay/bulls-and-cows/server/pkg/guess"
	"github.com/NikolovNikolay/bulls-and-cows/server/pkg/utils"
	"gopkg.in/mgo.v2/bson"
)

// Game represents a games session
type Game struct {
	GameID           bson.ObjectId  `json:"gameID" bson:"_id"`
	StartTime        int64          `josn:"startTime" bson:"startTime"`
	EndTime          int64          `json:"endTime" bson:"endTime"`
	GameType         int            `json:"type" bson:"type"`
	PlayerOneID      *bson.ObjectId `json:"playerOne" bson:"playerOne"`
	PlayerTwoID      *bson.ObjectId `json:"playerTwo" bson:"playerTwo"`
	PlayerOneGuesses []guess.Guess  `json:"playerOneGuesses" bson:"playerOneGuesses"`
	PlayerTwoGuesses []guess.Guess  `json:"playerTwoGuesses" bson:"playerTwoGuesses"`
	GuessNum         int            `json:"guess" bson:"guess"`
	GuessNumSec      int            `json:"guessSec" bson:"guessSec"`
}

// New returns a new Game instance
func New(
	gameID bson.ObjectId,
	gameType int,
	pOneID *bson.ObjectId,
	pTwoID *bson.ObjectId,
	guessNum int,
	guessNumSec int) (*Game, error) {

	if pOneID == nil {
		return nil, errors.New("Player one is required")
	}

	game := &Game{
		GameID:           gameID,
		GameType:         gameType,
		PlayerOneID:      pOneID,
		PlayerTwoID:      pTwoID,
		PlayerOneGuesses: nil,
		PlayerTwoGuesses: nil,
		GuessNum:         guessNum,
		GuessNumSec:      guessNumSec,
		StartTime:        time.Now().Unix()}

	return game, nil
}

// AddToDb creates a new game in DB
func AddToDb(
	dbName string,
	gameType int,
	pID *bson.ObjectId,
	guessNum int,
	mon *mgo.Session) (*Game, error) {

	gameID := bson.NewObjectId()
	game, e := New(
		gameID,
		gameType,
		pID,
		nil,
		guessNum,
		0)

	if e != nil {
		return nil, e
	}
	if er := mon.DB(
		dbName).C(
		utils.DBCGames).Insert(
		game); er != nil {
		return nil, er
	}

	return game, nil
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

// UpdateByID finds a game by ID in DB
// then updates it
func UpdateByID(
	g *Game,
	dbName string,
	mon *mgo.Session) error {

	if bson.IsObjectIdHex(g.GameID.Hex()) {
		e := mon.DB(
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
