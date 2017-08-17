package services

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/NikolovNikolay/bulls-and-cows/server/pkg/player"
	"github.com/NikolovNikolay/bulls-and-cows/server/pkg/utils"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
)

const (
	testHeader = "x-test"
)

func generateUserName() string {
	rand.Seed(time.Now().UnixNano())
	randPostfix := rand.Intn(10000)
	return fmt.Sprintf("user%d", randPostfix)
}

func generateNewPlayer(name string) player.Player {
	return player.Player{
		ID:       bson.NewObjectId(),
		Name:     name,
		Logged:   false,
		LoggedIn: nil}
}

func addTestHeadInRequest(r *http.Request) {
	r.Header.Set(testHeader, "true")
}

func getVarFromRequest(r *http.Request, key string) string {
	vars := mux.Vars(r)
	paramVal := vars[key]

	if paramVal == "" {
		paramVal = r.FormValue(key)
	}

	if paramVal == "" {
		paramVal = r.PostFormValue(key)
	}

	return paramVal
}

func parseForm(r *http.Request) {
	e := r.ParseForm()
	if e != nil {
		log.Print("Could not parse  form values of response")
	}
}

func getTargetDbName(r *http.Request) string {
	th := r.Header.Get(testHeader)
	isTesting := th != ""
	if isTesting {
		return utils.DBNameTest
	}

	return utils.DBName
}

func validateGameTypeParam(r *http.Request) (int, error) {
	gameType := r.PostFormValue(paramGameTypeKey)
	if gameType == "" {
		return 0, errors.New("Missing parameter for game type")
	}

	gt, e := strconv.Atoi(gameType)
	if e != nil {
		return 0, errors.New("Could not parse game type parameter")
	}
	return gt, nil
}
