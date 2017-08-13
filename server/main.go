package main

import (
	"log"
	"net/http"

	ctrl "github.com/NikolovNikolay/bulls-and-cows/server/controllers"
	r "github.com/NikolovNikolay/bulls-and-cows/server/router"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2"
)

func main() {
	session := initMgoSession()
	router := r.New()
	router = configureRoutes(r.New(), session)
	log.Fatal(http.ListenAndServe(":8080", router))
}

func configureRoutes(router *httprouter.Router, s *mgo.Session) *httprouter.Router {
	playerController := ctrl.NewPlayerController(s)
	gameController := ctrl.NewGameController(s, playerController)

	router.POST("/api/init", gameController.InitHandler)

	return router
}

func initMgoSession() *mgo.Session {
	session, err := mgo.Dial("mongodb://localhost")

	if err != nil {
		panic(err)
	}

	log.Println("Mongo session initialized")

	return session
}
