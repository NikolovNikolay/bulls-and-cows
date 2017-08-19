package utils

import "strconv"
import "fmt"

// BCCheck represents a tool that checks a
// number for bulls and cows occurances
type BCCheck struct{}

// BCCheckResult represents the result of the bulls and cows check
type BCCheckResult struct {
	Bulls int `json:"b"`
	Cows  int `json:"c"`
}

// Check takes the original number and the guess
// and returns the number of bulls and cows
func (bc BCCheck) Check(origin, guess int) *BCCheckResult {

	oStr := strconv.Itoa(origin)
	oMap := make(map[byte]bool)
	for i := 0; i < len(oStr); i++ {
		oMap[oStr[i]] = true
	}
	gStr := strconv.Itoa(guess)
	if len(gStr) <= 3 {

		gStr = fmt.Sprintf("%s%s", genPrependZeros(len(oStr)-len(gStr)), gStr)
	}
	bulls := 0
	cows := 0
	for i := 0; i < len(oStr); i++ {

		// check for bull
		if oStr[i] == gStr[i] {
			bulls++
			continue
		}
		// check for cow
		if oMap[gStr[i]] == true {
			cows++
		}
	}
	return &BCCheckResult{
		Bulls: bulls,
		Cows:  cows}
}

// ValidateMadeGuess check for a valid user gess input
func (bc BCCheck) ValidateMadeGuess(gs string) bool {
	dMap := make(map[byte]int)

	if len(gs) != 4 {
		return false
	}
	for i := 0; i < len(gs); i++ {
		if gs[i] < byte('0') || gs[i] > byte('9') {
			return false
		}
		if dMap[gs[i]] > 0 {
			return false
		}
		dMap[gs[i]] = dMap[gs[i]] + 1
	}

	return true
}

func genPrependZeros(len int) string {
	p := ""
	for i := 0; i < len; i++ {
		p = p + "0"
	}

	return p
}
