package controllers

// import (
// 	"io"
// 	"net/http"
// 	"net/http/httptest"
// 	"strings"
// 	"testing"
// )

// func TestGenerateUserName(t *testing.T) {

// 	testFunc := func(t *testing.T) {
// 		if un := generateUserName(); un == "" || !strings.Contains(un, "user") {
// 			t.Error("Username is not generated correctly")
// 		} else {
// 			t.Log("Username generated correctly")
// 		}
// 	}

// 	t.Run("genUserName", func(t *testing.T) {
// 		t.Run("FirstGen", testFunc)
// 		t.Run("SecondGen", testFunc)
// 		t.Run("ThirdGen", testFunc)
// 		t.Run("FourthGen", testFunc)
// 	})
// }

// func TestInitHandler(t *testing.T) {

// 	req := httptest.NewRequest("GET", "http://localhost:8080/api/init", nil)
// 	req.Header.Add("x-test", "true")

// 	w := httptest.NewRecorder()
// 	handler := http.HandlerFunc(HealthCheckHandler)

// 	handler.ServeHTTP(w, req)

// 	// Check the status code is what we expect.
// 	if status := w.Code; status != http.StatusOK {
// 		t.Errorf("handler returned wrong status code: got %v want %v",
// 			status, http.StatusOK)
// 	}

// 	// Check the response body is what we expect.
// 	expected := `{"alive": true}`
// 	if w.Body.String() != expected {
// 		t.Errorf("handler returned unexpected body: got %v want %v",
// 			w.Body.String(), expected)
// 	}
// }

// func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
// 	w.WriteHeader(http.StatusOK)
// 	w.Header().Set("Content-Type", "application/json")

// 	io.WriteString(w, `{"alive": true}`)
// }
