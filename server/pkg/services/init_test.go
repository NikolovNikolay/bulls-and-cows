package services

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/NikolovNikolay/bulls-and-cows/server/pkg/utils"
)

func TestInitOfInit(t *testing.T) {
	s := NewInitService(utils.GetDBSession())
	if s == nil {
		t.Error("Could not init the init service")
	}
}

func TestInitMethod(t *testing.T) {
	s := NewInitService(utils.GetDBSession())
	if s.Method() != "POST" {
		t.Error("Not correct method for init service")
	}
}

func TestInitEndpoint(t *testing.T) {
	s := NewInitService(utils.GetDBSession())
	if s.Endpoint() != "/api/init" {
		t.Error("Not correct endpoint for init service")
	}
}

func TestInitHandler(t *testing.T) {
	s := NewInitService(utils.GetDBSession())
	t.Run("Test with user name", func(t *testing.T) {
		t.Run("Correct request", func(t *testing.T) {
			form := url.Values{}
			form.Add(paramUserNameKey, "Joe Black")
			form.Add(paramGameTypeKey, "2")

			// Create a test request and add custom headers
			rq, _ := http.NewRequest("POST", "/api/init/", strings.NewReader(form.Encode()))
			rq.Header.Add("Content-Type", "application/x-www-form-urlencoded")
			addTestHeadInRequest(rq)

			w := httptest.NewRecorder()
			s.Handle(w, rq)

			resp := w.Result()

			if resp.StatusCode != 200 {
				t.Error("Invalid response from server")
			}
		})
	})
	t.Run("Test with no user name", func(t *testing.T) {
		t.Run("Correct request", func(t *testing.T) {
			form := url.Values{}
			form.Add(paramUserNameKey, "Joe Black")
			form.Add(paramGameTypeKey, "1")

			// Create a test request and add custom headers
			rq, _ := http.NewRequest("POST", "/api/init/", strings.NewReader(form.Encode()))
			rq.Header.Add("Content-Type", "application/x-www-form-urlencoded")
			addTestHeadInRequest(rq)

			w := httptest.NewRecorder()
			s.Handle(w, rq)

			resp := w.Result()

			if resp.StatusCode != 200 {
				t.Error("Invalid response from server")
			}
		})
	})
	t.Run("Test with no game type", func(t *testing.T) {
		t.Run("Correct request", func(t *testing.T) {
			form := url.Values{}
			form.Add(paramUserNameKey, "Joe Black")
			// form.Add(paramGameTypeKey, "1")

			// Create a test request and add custom headers
			rq, _ := http.NewRequest("POST", "/api/init/", strings.NewReader(form.Encode()))
			rq.Header.Add("Content-Type", "application/x-www-form-urlencoded")
			addTestHeadInRequest(rq)

			w := httptest.NewRecorder()
			s.Handle(w, rq)

			resp := w.Result()

			if resp.StatusCode != 400 {
				t.Error("Invalid response from server")
			}
		})
	})
	t.Run("Test with no params", func(t *testing.T) {
		t.Run("Correct request", func(t *testing.T) {

			// Create a test request and add custom headers
			rq, _ := http.NewRequest("POST", "/api/init/", nil)
			addTestHeadInRequest(rq)

			w := httptest.NewRecorder()
			s.Handle(w, rq)

			resp := w.Result()

			if resp.StatusCode != 400 {
				t.Error("Invalid response from server")
			}
		})
	})

}
