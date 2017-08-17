package services

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/NikolovNikolay/bulls-and-cows/server/pkg/game"
	"github.com/NikolovNikolay/bulls-and-cows/server/pkg/player"
	"github.com/NikolovNikolay/bulls-and-cows/server/pkg/utils"
)

var tgame *game.Game

func TestGetGameInit(t *testing.T) {
	gg := NewGetGameService(utils.GetDBSession())
	if gg == nil {
		t.Error("Could not init get game service")
	}
}

func TestGetGameMethod(t *testing.T) {
	gg := NewGetGameService(utils.GetDBSession())
	if gg.Method() != "GET" {
		t.Error("Not correct method for get game service")
	}
}

func TestGetGameEndpoint(t *testing.T) {
	gg := NewGetGameService(utils.GetDBSession())
	if gg.Endpoint() != "/api/game/{gameID}" {
		t.Error("Not correct endpoint for get game service")
	}
}

func TestGetGameHandlerValidID(t *testing.T) {
	gg := NewGetGameService(utils.GetDBSession())
	tgame = game.New(1)
	tgame.AddPlayer(player.New("TestUser", true, utils.DBNameTest, utils.GetDBSession()))
	tgame.Start(time.Now().Unix())
	e := tgame.Add(utils.DBNameTest, utils.GetDBSession())
	if e != nil {
		t.Error("Could not add new game to DB")
	}

	handler := gg.Handle

	rq, _ := http.NewRequest("GET", "/api/game?gameID="+tgame.ID.Hex(), nil)
	addTestHeadInRequest(rq)
	// req := httptest.NewRequest(gg.Method(), "/api/game/", nil)

	w := httptest.NewRecorder()
	handler(w, rq)

	resp := w.Result()

	if resp.StatusCode != 200 {
		t.Error("Invalid response from server")
	}
}

func TestGetGameHandlerInvalidID(t *testing.T) {
	gg := NewGetGameService(utils.GetDBSession())
	handler := gg.Handle

	req := httptest.NewRequest(gg.Method(), "/api/game?gameID=9s8df9ds8f", nil)
	w := httptest.NewRecorder()
	handler(w, req)

	resp := w.Result()

	if resp.StatusCode != 400 {
		t.Error("Invalid response from server")
	}
}

func TestGetGameHandlerNoID(t *testing.T) {
	gg := NewGetGameService(utils.GetDBSession())
	handler := gg.Handle

	req := httptest.NewRequest(gg.Method(), "/api/game/", nil)

	w := httptest.NewRecorder()
	handler(w, req)

	resp := w.Result()

	if resp.StatusCode != 400 {
		t.Error("Invalid response from server")
	}
}
