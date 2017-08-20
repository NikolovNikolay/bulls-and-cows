package utils

import (
	"log"

	mgo "gopkg.in/mgo.v2"
)

var session *mgo.Session

func init() {
	s, e := initMongo()
	if e != nil {
		panic(e)
	}
	session = s
}

// GetDBSession returns the initialized DB session
func GetDBSession() *mgo.Session {
	return session
}

// initMongo initializes a mongo db connection
func initMongo() (*mgo.Session, error) {
	if session != nil {
		return session, nil
	}

	session, e := mgo.Dial("mongodb://localhost")
	if e != nil {
		return nil, e
	}
	log.Println("Mongo session initialized")

	return session, nil
}
