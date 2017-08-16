package guess

import "github.com/NikolovNikolay/bulls-and-cows/server/pkg/utils"

// Guess represents a guess, stored in the game object
type Guess struct {
	Guess int                  `json:"g"`
	Bc    *utils.BCCheckResult `json:"bc"`
}
