package main

import (
	"log"
	"net/http"

	ctrl "github.com/NikolovNikolay/bulls-and-cows/server/pkg/controllers"
	r "github.com/NikolovNikolay/bulls-and-cows/server/pkg/router"
	"github.com/NikolovNikolay/bulls-and-cows/server/pkg/services"
	u "github.com/NikolovNikolay/bulls-and-cows/server/pkg/utils"
	"github.com/googollee/go-socket.io"
	"github.com/rs/cors"
	"gopkg.in/mgo.v2"
)

const (
	servePort    = "8080"
	clientOrigin = "http://localhost:4200"
)

func main() {
	db := u.InitMongo()
	router := configureRoutes(r.New(), db)

	h := cors.New(
		cors.Options{
			AllowedMethods:   []string{"POST", "GET", "PUT"},
			AllowedOrigins:   []string{clientOrigin},
			AllowCredentials: true}).Handler(router.R)

	log.Fatal(http.ListenAndServe(servePort, h))
}

func configureRoutes(rt r.BCRouter, s *mgo.Session) r.BCRouter {
	pc := ctrl.NewPlayerController(s)
	gc := ctrl.NewGameController(s, pc)
	ss, e := socketio.NewServer(nil)
	if e != nil {
		panic(e)
	}
	sc := ctrl.NewSocketController(ss, gc, pc)

	rt.RegisterService(services.NewInitService(gc))
	rt.RegisterService(services.NewGuessService(gc))
	rt.RegisterService(services.NewGetGameService(gc))
	rt.R.Handle(ctrl.SockIoEndpoint, sc.Socket)

	return rt
}
