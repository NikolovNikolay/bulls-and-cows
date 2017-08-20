package guess

import "github.com/NikolovNikolay/bulls-and-cows/server/pkg/utils"

// Guess represents a guess, stored in the game object
type Guess struct {
	Number int                  `json:"g"`
	Bc     *utils.BCCheckResult `json:"bc"`
}

// New returns a new Guess instance
func New(num, bulls, cows int) Guess {
	bc := utils.BCCheckResult{
		Bulls: bulls,
		Cows:  cows}

	return Guess{
		Number: num,
		Bc:     &bc}
}
