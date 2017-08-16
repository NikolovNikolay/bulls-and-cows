package controllers

// import (
// 	"log"
// 	"testing"

// 	"github.com/NikolovNikolay/bulls-and-cows/server/utils"

// 	mgo "gopkg.in/mgo.v2"
// )

// var db *mgo.Session
// var c PlayerController

// func init() {
// 	db = initMgoSession()
// }

// func initMgoSession() *mgo.Session {
// 	session, err := mgo.Dial("mongodb://localhost")

// 	if err != nil {
// 		panic(err)
// 	}

// 	log.Println("Mongo session initialized")

// 	return session
// }

// func TestCreatePlayer(t *testing.T) {

// 	c = NewPlayerController(db)
// 	p := generateNewPlayer("TestUser")

// 	e := c.createPlayer(p, utils.DBNameTest)
// 	if e != nil {
// 		t.Error("Could not insert user to db")
// 	} else {
// 		t.Log("Saved user into database")
// 	}
// }

// func TestFindPlayer(t *testing.T) {
// 	if p, e := c.findPlayerByName("TestUser", utils.DBNameTest); e != nil || p.Name != "TestUser" {
// 		t.Error("Could not get user from database")
// 	} else {
// 		t.Log("Found user in database")
// 	}
// }

// func TestUpdatePlayer(t *testing.T) {
// 	if p, e := c.findPlayerByName("TestUser", utils.DBNameTest); e != nil || p.Name != "TestUser" {
// 		t.Error("Could not get user from database")
// 	} else {
// 		p.Wins = p.Wins + 1
// 		if e = c.updatePlayer(p, utils.DBNameTest); e != nil {
// 			t.Error("Error while trying to update player in DB")
// 		} else {
// 			t.Log("Updated user in database")
// 		}
// 	}
// }

// func TestUpdatedPlayer(t *testing.T) {
// 	if p, e := c.findPlayerByName("TestUser", utils.DBNameTest); e != nil || p.Name != "TestUser" {
// 		t.Error("Could not get user from database")
// 	} else {
// 		if p.Wins != 1 {
// 			t.Error("Player was not updated correctly")
// 		} else {
// 			db.DB(utils.DBNameTest).C(utils.DBCPlayers).RemoveAll(nil)
// 			t.Log("Player updated correctly")
// 		}
// 	}
// }
