package main

import (
	"log"
	"net/http"

	r "github.com/NikolovNikolay/bulls-and-cows/server/pkg/router"
	"github.com/NikolovNikolay/bulls-and-cows/server/pkg/services"
	"github.com/NikolovNikolay/bulls-and-cows/server/pkg/socket"
	u "github.com/NikolovNikolay/bulls-and-cows/server/pkg/utils"
	"github.com/googollee/go-socket.io"
	"github.com/rs/cors"
	"gopkg.in/mgo.v2"
)

const (
	servePort    = ":8080"
	clientOrigin = "http://localhost:4200"
)

func main() {
	db := u.GetDBSession()
	router := configureRoutes(r.New(), db)

	h := cors.New(
		cors.Options{
			AllowedMethods:   []string{"POST", "GET", "PUT"},
			AllowedOrigins:   []string{clientOrigin},
			AllowCredentials: true}).Handler(router.R)

	log.Fatal(http.ListenAndServe(servePort, h))
}

func configureRoutes(rt r.BCRouter, s *mgo.Session) r.BCRouter {

	ss, e := socketio.NewServer(nil)
	ws := socket.New(ss, s)
	ws.Init()

	if e != nil {
		panic(e)
	}
	// sc := socket.New() ctrl.NewSocketController(ss, gc, pc)

	rt.RegisterService(services.NewInitService(s))
	rt.RegisterService(services.NewGuessService(s))
	rt.RegisterService(services.NewGetGameService(s))
	rt.R.Handle(socket.SockIoEndpoint, ws.Socket)

	return rt
}
