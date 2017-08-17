package utils

import (
	"log"

	mgo "gopkg.in/mgo.v2"
)

var session *mgo.Session

func init() {
	session = initMongo()
}

// GetDBSession returns the initialized DB session
func GetDBSession() *mgo.Session {
	return session
}

// initMongo initializes a mongo db connection
func initMongo() *mgo.Session {
	if session != nil {
		return session
	}

	session, err := mgo.Dial("mongodb://localhost")

	if err != nil {
		panic(err)
	}

	log.Println("Mongo session initialized")

	return session
}
