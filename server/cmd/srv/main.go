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
	r, e := configureRoutes(r.New(), db)
	if e != nil {
		panic(e)
	}

	h := cors.New(
		cors.Options{
			AllowedMethods:   []string{"POST", "GET", "PUT"},
			AllowedOrigins:   []string{clientOrigin},
			AllowCredentials: true}).Handler(r.R)

	log.Fatal(http.ListenAndServe(servePort, h))
}

func configureRoutes(rt *r.BCRouter, s *mgo.Session) (*r.BCRouter, error) {

	sio, e := socketio.NewServer(nil)
	if e != nil {
		return nil, e
	}

	ws := socket.New(sio, s, false)
	e = ws.Init()
	if e != nil {
		return nil, e
	}

	rt.RegisterService(services.NewInitService(s))
	rt.RegisterService(services.NewGuessService(s))
	rt.RegisterService(services.NewGetGameService(s))
	rt.R.Handle(socket.SockIoEndpoint, ws.Socket)

	return rt, nil
}
