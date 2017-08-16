package utils

import (
	"log"

	mgo "gopkg.in/mgo.v2"
)

// InitMongo initializes a mongo db connection
func InitMongo() *mgo.Session {
	session, err := mgo.Dial("mongodb://localhost")

	if err != nil {
		panic(err)
	}

	log.Println("Mongo session initialized")

	return session
}
