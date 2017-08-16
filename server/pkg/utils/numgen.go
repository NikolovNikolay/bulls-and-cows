package utils

import (
	"math/rand"
	"strconv"
	"time"
)

const (
	maxNumber = 9876
	minNumber = 1234
)

// NumGen represents a tool for generating 4 distinct digit numbers
// for the bulls and cows game
type NumGen struct{}

// GetNumGen returns a new instance of NumGen
func GetNumGen() NumGen {
	return NumGen{}
}

// Gen generates the guess number
func (ng NumGen) Gen() int {
	rand.Seed(time.Now().UTC().UnixNano())
	r := rand.Intn(maxNumber) + minNumber
	if r > maxNumber {
		r = maxNumber
	}
	if r < minNumber {
		r = minNumber
	}

	rStr := strconv.Itoa(r)

	dgb := make(map[byte]bool)
	dgn := [4]byte{}

	// check and set the first char of the number
	dgn[0] = rStr[0]
	dgb[rStr[0]] = true

	for i := 1; i < len(rStr); i++ {
		if dgb[rStr[i]] == false {
			dgn[i] = rStr[i]
			dgb[rStr[i]] = true
		} else {
			for {
				nd := strconv.Itoa(rand.Intn(9))
				if dgb[nd[0]] == false {
					dgn[i] = nd[0]
					dgb[nd[0]] = true
					break
				}
			}
		}
	}

	s := string(dgn[:])
	i, _ := strconv.Atoi(s)
	return i
}
