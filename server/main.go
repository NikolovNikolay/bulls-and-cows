package main

import (
	"log"
	"net/http"

	ctrl "github.com/NikolovNikolay/bulls-and-cows/server/controllers"
	r "github.com/NikolovNikolay/bulls-and-cows/server/router"
	"github.com/googollee/go-socket.io"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"gopkg.in/mgo.v2"
)

func main() {
	session := initMgoSession()
	router := configureRoutes(r.New(), session)
	h := cors.Default().Handler(router)

	server, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}
	server.On("connection", func(so socketio.Socket) {
		log.Println("on connection")
	})
	server.On("error", func(so socketio.Socket, err error) {
		log.Println("error:", err)
	})
	router.Handle("/socket.io/", server)

	log.Fatal(http.ListenAndServe(":8081", h))
}

func configureRoutes(router *mux.Router, s *mgo.Session) *mux.Router {
	playerController := ctrl.NewPlayerController(s)
	gameController := ctrl.NewGameController(s, playerController)

	router.HandleFunc("/api/init", gameController.InitHandler).Methods("POST")
	router.HandleFunc("/api/init", gameController.InitHandler).Methods("POST")
	router.HandleFunc("/api/guess/{guessNum:[0-9]+}", gameController.GuessHandler).Methods("PUT")
	router.HandleFunc("/api/game/{gameID}", gameController.GetGameDataHandler).Methods("GET")

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
