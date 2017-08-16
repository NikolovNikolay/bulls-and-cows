package utils

import (
	"fmt"
	"testing"
)

func TestGenPrependZeros(t *testing.T) {
	gStr := fmt.Sprintf("%s%s", genPrependZeros(5), "1234")
	if gStr != "000001234" {
		t.Error("Not prepended digits correctly")
	}

	gStr = fmt.Sprintf("%s%s", genPrependZeros(1), "234")
	if gStr != "0234" {
		t.Error("Not prepended digits correctly")
	}

	gStr = fmt.Sprintf("%s%s", genPrependZeros(0), "234")
	if gStr != "234" {
		t.Error("Not prepended digits correctly")
	}
}

func TestValidateGuess(t *testing.T) {
	bcc := BCCheck{}

	val := bcc.ValidateMadeGuess("1234")
	if val != true {
		t.Error("Guess was not validated correctly")
	}

	val = bcc.ValidateMadeGuess("1234")
	if val != true {
		t.Error("Guess was not validated correctly")
	}

	val = bcc.ValidateMadeGuess("e125")
	if val != false {
		t.Error("Guess was not validated correctly")
	}

	val = bcc.ValidateMadeGuess("1234566")
	if val != false {
		t.Error("Guess was not validated correctly")
	}

	val = bcc.ValidateMadeGuess("0123")
	if val != true {
		t.Error("Guess was not validated correctly")
	}

	val = bcc.ValidateMadeGuess("1231")
	if val != false {
		t.Error("Guess was not validated correctly")
	}

	val = bcc.ValidateMadeGuess("3231")
	if val != false {
		t.Error("Guess was not validated correctly")
	}
}

func TestCheck(t *testing.T) {
	bc := BCCheck{}
	r := bc.Check(1234, 1256)
	if r.Bulls != 2 && r.Cows != 0 {
		t.Error("Incorrect check for bulls and cows")
	}
	r = bc.Check(1234, 1234)
	if r.Bulls != 4 && r.Cows != 0 {
		t.Error("Incorrect check for bulls and cows")
	}
	r = bc.Check(1234, 123)
	if r.Bulls != 0 && r.Cows != 3 {
		t.Error("Incorrect check for bulls and cows")
	}
	r = bc.Check(1234, 12)
	if r.Bulls != 0 && r.Cows != 2 {
		t.Error("Incorrect check for bulls and cows")
	}
	r = bc.Check(7845, 8754)
	if r.Bulls != 0 && r.Cows != 4 {
		t.Error("Incorrect check for bulls and cows")
	}
}
