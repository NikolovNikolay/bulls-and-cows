package services

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/NikolovNikolay/bulls-and-cows/server/pkg/game"
	"github.com/NikolovNikolay/bulls-and-cows/server/pkg/player"
	"github.com/NikolovNikolay/bulls-and-cows/server/pkg/utils"
)

func TestGuessInit(t *testing.T) {
	gs := NewGuessService(utils.GetDBSession())
	if gs == nil {
		t.Error("Could not init guess service")
	}
}

func TestGuessMethod(t *testing.T) {
	gs := NewGuessService(utils.GetDBSession())
	if gs.Method() != "PUT" {
		t.Error("Not correct method for get game service")
	}
}

func TestGuessEndpoint(t *testing.T) {
	gg := NewGuessService(utils.GetDBSession())
	if gg.Endpoint() != "/api/guess/{guessNum:[0-9]+}" {
		t.Error("Not correct endpoint for guess  service")
	}
}

func TestGuessValidation(t *testing.T) {
	gg := NewGuessService(utils.GetDBSession())
	t.Run("Validate guess", func(t *testing.T) {
		t.Run("valid guess", func(t *testing.T) {
			g, e := gg.validateGuessNumberParam("9567")
			if g != 9567 || e != nil {
				t.Error("Could not validate a valid guess")
			}
		})
		t.Run("invalid guess number", func(t *testing.T) {
			_, e := gg.validateGuessNumberParam("95674")
			if e == nil || e.Error() != "Invalid guess number" {
				t.Error("Could not validate a valid guess")
			}
		})
		t.Run("invalid guess mix string", func(t *testing.T) {
			_, e := gg.validateGuessNumberParam("95f6")
			if e == nil || e.Error() != "Invalid guess number" {
				t.Error("Could not validate a valid guess")
			}
		})
		t.Run("blank", func(t *testing.T) {
			_, e := gg.validateGuessNumberParam("")
			if e == nil || e.Error() != "Missing parameter guess" {
				t.Error("Could not validate a valid guess")
			}
		})
	})
}

func TestGuessHandler(t *testing.T) {
	// Init the service
	gg := NewGuessService(utils.GetDBSession())

	// Create a test game and initialize it
	tgame = game.New(1)
	tgame.AddPlayer(player.New("TestUser", true, utils.DBNameTest, utils.GetDBSession()))
	tgame.Start(time.Now().Unix())

	// Save the game to the DB
	e := tgame.Add(utils.DBNameTest, utils.GetDBSession())
	if e != nil {
		t.Error("Could not add new game to DB")
	}

	t.Run("Test  gues handler", func(t *testing.T) {
		t.Run("Correct request", func(t *testing.T) {
			form := url.Values{}
			form.Add(paramGameIDKey, tgame.ID.Hex())
			form.Add(paramGuessKey, "1234")

			// Create a test request and add custom headers
			rq, _ := http.NewRequest("PUT", "/api/guess/", strings.NewReader(form.Encode()))
			rq.Header.Add("Content-Type", "application/x-www-form-urlencoded")
			addTestHeadInRequest(rq)

			w := httptest.NewRecorder()
			gg.Handle(w, rq)

			resp := w.Result()

			if resp.StatusCode != 200 {
				t.Error("Invalid response from server")
			}
		})
		t.Run("Missing game id", func(t *testing.T) {
			form := url.Values{}
			form.Add(paramGuessKey, "1234")

			// Create a test request and add custom headers
			rq, _ := http.NewRequest("PUT", "/api/guess/", strings.NewReader(form.Encode()))
			rq.Header.Add("Content-Type", "application/x-www-form-urlencoded")
			addTestHeadInRequest(rq)

			w := httptest.NewRecorder()
			gg.Handle(w, rq)

			resp := w.Result()

			if resp.StatusCode != 400 {
				t.Errorf("Expected bad request. Received %d", resp.StatusCode)
			}
		})
		t.Run("Missing guess number", func(t *testing.T) {
			form := url.Values{}
			form.Add(paramGameIDKey, tgame.ID.Hex())

			// Create a test request and add custom headers
			rq, _ := http.NewRequest("PUT", "/api/guess/", strings.NewReader(form.Encode()))
			rq.Header.Add("Content-Type", "application/x-www-form-urlencoded")
			addTestHeadInRequest(rq)

			w := httptest.NewRecorder()
			gg.Handle(w, rq)

			resp := w.Result()

			if resp.StatusCode != 400 {
				t.Errorf("Expected bad request. Received %d", resp.StatusCode)
			}
		})
		t.Run("Invalid game id", func(t *testing.T) {
			form := url.Values{}
			form.Add(paramGameIDKey, "as98d09as8d")
			form.Add(paramGuessKey, "1234")

			// Create a test request and add custom headers
			rq, _ := http.NewRequest("PUT", "/api/guess/", strings.NewReader(form.Encode()))
			rq.Header.Add("Content-Type", "application/x-www-form-urlencoded")
			addTestHeadInRequest(rq)

			w := httptest.NewRecorder()
			gg.Handle(w, rq)

			resp := w.Result()

			if resp.StatusCode != 400 {
				t.Errorf("Expected bad request. Received %d", resp.StatusCode)
			}
		})
	})

}
