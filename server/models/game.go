package models

import "gopkg.in/mgo.v2/bson"
import "errors"
import "time"

// Game represents a games session
type Game struct {
	GameID           bson.ObjectId  `json:"gameID" bson:"_id"`
	StartTime        int64          `josn:"startTime" bson:"startTime"`
	EndTime          int64          `json:"endTime" bson:"endTime"`
	GameType         *int           `json:"type" bson:"type"`
	PlayerOneID      *bson.ObjectId `json:"playerOne" bson:"playerOne"`
	PlayerTwoID      *bson.ObjectId `json:"playerTwo" bson:"playerTwo"`
	PlayerOneGuesses []int16        `json:"playerOneGuesses" bson:"playerOneGuesses"`
	PlayerTwoGuesses []int16        `json:"playerTwoGuesses" bson:"playerTwoGuesses"`
	GuessNum         int            `json:"guess" bson:"guess"`
}

// NewGame returns a new Game instance
func NewGame(
	gameID bson.ObjectId,
	gameType *int,
	pOneID *bson.ObjectId,
	pTwoID *bson.ObjectId,
	guessNum int) (*Game, error) {

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
		StartTime:        time.Now().Unix()}

	return game, nil
}
