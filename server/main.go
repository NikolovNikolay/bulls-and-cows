package main

import (
	"log"
	"net/http"
	"strconv"

	"gopkg.in/mgo.v2/bson"

	"github.com/NikolovNikolay/bulls-and-cows/server/utils"

	"github.com/NikolovNikolay/bulls-and-cows/server/controllers"
	"github.com/NikolovNikolay/bulls-and-cows/server/models"

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

	log.Fatal(http.ListenAndServe(":8080", h))
}

func configureSocketIO(
	pc controllers.PlayerController,
	gc controllers.GameController) *socketio.Server {
	socket, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}

	roomPMap := make(map[string]int)
	socket.On("connection", func(so socketio.Socket) {
		so.Join("default")
		log.Println("on connection")
		so.On("creater", func(room string) bool {

			if roomPMap[room] > 0 {
				return false
			}

			roomPMap[room] = roomPMap[room] + 1
			so.BroadcastTo("default", "updater", [1]string{room})
			so.Join(room)
			so.Leave("default")

			p := models.Player{
				ID:     bson.NewObjectId(),
				Logged: true,
				Name:   room}
			pc.CreatePlayer(p, utils.DBName)

			g, _ := gc.CreateNewGame(
				utils.DBName,
				utils.P2P,
				&p.ID,
				0)

			p.LoggedIn = &g.GameID
			pc.UpdatePlayer(p, utils.DBName)

			return true
		})
		so.On("joinr", func(room, rival string) bool {

			if roomPMap[room] >= 2 {
				return false
			}

			roomPMap[room] = roomPMap[room] + 1
			so.BroadcastTo(room, "joinmy", rival)
			so.Join(room)

			host, _ := pc.FindPlayerByName(room, utils.DBName)
			gID := host.LoggedIn
			g, _ := gc.GetGameByID(gID.Hex(), utils.DBName)
			p := models.Player{
				ID:       bson.NewObjectId(),
				Logged:   true,
				LoggedIn: gID,
				Name:     rival}
			pc.CreatePlayer(p, utils.DBName)
			g.PlayerTwoID = &p.ID
			gc.UpdateGameByID(g, utils.DBName)
			so.BroadcastTo(room, "confjoin")
			return true
		})
		so.On("getavr", func(p string) []string {
			avrooms := make([]string, 1)

			for k, v := range roomPMap {
				if v < 2 || k != p {
					avrooms = append(avrooms, k)
				}
			}

			return avrooms
		})
		so.On("inputguess", func(guess, name string) bool {
			p, _ := pc.FindPlayerByName(name, utils.DBName)
			g, _ := gc.GetGameByID(p.LoggedIn.Hex(), utils.DBName)

			if g.PlayerOneID == &p.ID {
				g.GuessNum, _ = strconv.Atoi(guess)
			} else {
				g.GuessNumSec, _ = strconv.Atoi(guess)
			}

			gc.UpdateGameByID(g, utils.DBName)

			return true
		})
	})
	socket.On("error", func(so socketio.Socket, err error) {
		log.Println("error:", err)
	})
	socket.On("disconnect", func(so socketio.Socket, err error) {
		log.Println("disconnect:", err)
	})
	return socket
}

func configureRoutes(rt r.BCRouter, s *mgo.Session) r.BCRouter {
	pc := ctrl.NewPlayerController(s)
	gc := ctrl.NewGameController(s, pc)

	rt.RegisterService(services.NewInitService(gc))
	rt.RegisterService(services.NewGuessService(gc))
	rt.RegisterService(services.NewGetGameService(gc))

	rt.R.Handle("/socket.io/", configureSocketIO(pc, gc))

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
