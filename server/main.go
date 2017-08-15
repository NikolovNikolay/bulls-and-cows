package main

import (
	"log"
	"net/http"

	ctrl "github.com/NikolovNikolay/bulls-and-cows/server/controllers"
	r "github.com/NikolovNikolay/bulls-and-cows/server/router"
	"github.com/NikolovNikolay/bulls-and-cows/server/services"
	"github.com/googollee/go-socket.io"
	"github.com/rs/cors"
	"gopkg.in/mgo.v2"
)

func main() {
	session := initMgoSession()
	router := configureRoutes(r.New(), session)

	h := cors.New(
		cors.Options{
			AllowedMethods:   []string{"POST", "GET", "PUT"},
			AllowedOrigins:   []string{"http://localhost:4200"},
			AllowCredentials: true}).Handler(router.R)

	router.R.Handle("/socket.io/", configureSocketIO())

	log.Fatal(http.ListenAndServe(":8081", h))
}

func configureSocketIO() *socketio.Server {
	socket, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}
	socket.On("connection", func(so socketio.Socket) {
		log.Println("on connection")
	})
	socket.On("error", func(so socketio.Socket, err error) {
		log.Println("error:", err)
	})

	return socket
}

func configureRoutes(rt r.BCRouter, s *mgo.Session) r.BCRouter {
	playerController := ctrl.NewPlayerController(s)
	gc := ctrl.NewGameController(s, playerController)

	rt.RegisterService(services.NewInitService(gc))
	rt.RegisterService(services.NewGuessService(gc))
	rt.RegisterService(services.NewGetGameService(gc))

	return rt
}

func initMgoSession() *mgo.Session {
	session, err := mgo.Dial("mongodb://localhost")

	if err != nil {
		panic(err)
	}

	log.Println("Mongo session initialized")

	return session
}
