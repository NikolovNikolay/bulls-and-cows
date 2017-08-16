package models

import "gopkg.in/mgo.v2/bson"

// Player represent a human player
type Player struct {
	ID       bson.ObjectId  `json:"id" bson:"_id"`
	Name     string         `json:"name" bson:"name"`
	Wins     int            `json:"wins" bson:"wins"`
	Logged   bool           `json:"logged" bson:"logged"`
	LoggedIn *bson.ObjectId `json:"loggedId" bson:"loggedIn"`
}
