package guess

import (
	"testing"

	"github.com/NikolovNikolay/bulls-and-cows/server/pkg/utils"
)

func TestBson(t *testing.T) {
	bc := utils.BCCheckResult{Bulls: 0, Cows: 3}
	g := New(
		1234,
		bc.Bulls,
		bc.Cows)

	if g.Bc.Bulls != 0 && g.Bc.Cows != 3 && g.Guess != 1234 {
		t.Error("Guess was not initialied as expected")
	}

	bc = utils.BCCheckResult{Bulls: 1, Cows: 3}
	g = New(
		4577,
		bc.Bulls,
		bc.Cows)

	if g.Bc.Bulls != 1 && g.Bc.Cows != 3 && g.Guess != 4577 {
		t.Error("Guess was not initialied as expected")
	}

	bc = utils.BCCheckResult{Bulls: 2, Cows: 2}
	g = New(
		1278,
		bc.Bulls,
		bc.Cows)

	if g.Bc.Bulls != 2 && g.Bc.Cows != 2 && g.Guess != 1278 {
		t.Error("Guess was not initialied as expected")
	}
}
