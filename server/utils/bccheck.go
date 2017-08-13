package utils

import "strconv"
import "fmt"

// BCCheck represents a tool that checks a
// number for bulls and cows occurances
type BCCheck struct{}

// BCCheckResult represents the result of the bulls and cows check
type BCCheckResult struct {
	Bulls int `json:"c"`
	Cows  int `json:"c"`
}

// Check takes the original number and the guess
// and returns the number of bulls and cows
func (bc BCCheck) Check(origin, guess int) BCCheckResult {

	oStr := strconv.Itoa(origin)
	oMap := make(map[byte]bool)
	for i := 0; i < len(oStr); i++ {
		oMap[oStr[i]] = true
	}

	gStr := strconv.Itoa(guess)
	bulls := 0
	cows := 0

	for i := 0; i < len(oStr); i++ {

		// check for bull
		fmt.Println(oStr[i])
		fmt.Println(gStr[i])
		if oStr[i] == gStr[i] {
			bulls++
			continue
		}

		// check for cow
		if oMap[gStr[i]] == true {
			cows++
		}
	}

	return BCCheckResult{
		Bulls: bulls,
		Cows:  cows}
}
