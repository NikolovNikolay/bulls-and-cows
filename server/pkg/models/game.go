package models

import (
	"errors"
	"time"

	"github.com/NikolovNikolay/bulls-and-cows/server/pkg/utils"
	"gopkg.in/mgo.v2/bson"
)

// Game represents a games session
type Game struct {
	GameID           bson.ObjectId    `json:"gameID" bson:"_id"`
	StartTime        int64            `josn:"startTime" bson:"startTime"`
	EndTime          int64            `json:"endTime" bson:"endTime"`
	GameType         int              `json:"type" bson:"type"`
	PlayerOneID      *bson.ObjectId   `json:"playerOne" bson:"playerOne"`
	PlayerTwoID      *bson.ObjectId   `json:"playerTwo" bson:"playerTwo"`
	PlayerOneGuesses []GuessDBContent `json:"playerOneGuesses" bson:"playerOneGuesses"`
	PlayerTwoGuesses []GuessDBContent `json:"playerTwoGuesses" bson:"playerTwoGuesses"`
	GuessNum         int              `json:"guess" bson:"guess"`
	GuessNumSec      int              `json:"guessSec" bson:"guessSec"`
}

// GuessDBContent represents a guess, stored in the game object
type GuessDBContent struct {
	Guess int                  `json:"g"`
	Bc    *utils.BCCheckResult `json:"bc"`
}

// NewGame returns a new Game instance
func NewGame(
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
