package utils

import "testing"
import "strconv"

func TestNumgenDistinct(t *testing.T) {
	numGen := GetNumGen()

	checkValid := func(t *testing.T) {
		m := make(map[byte]int)
		gen := numGen.Gen()
		str := strconv.Itoa(gen)
		for i := 0; i < len(str); i++ {
			if m[str[i]] > 0 {
				t.Errorf("Number %d is invalid", gen)
			}
			m[str[i]] = m[str[i]] + 1
		}
	}

	t.Run("Check generated numbers", func(t *testing.T) {
		t.Run("First", checkValid)
		t.Run("Second", checkValid)
		t.Run("Third", checkValid)
		t.Run("Fourth", checkValid)
		t.Run("Fifth", checkValid)
	})
}

func TestNumgenLength(t *testing.T) {
	numGen := GetNumGen()

	checkValid := func(t *testing.T) {
		gen := numGen.Gen()
		s := strconv.Itoa(gen)
		if len(s) != 4 {
			t.Errorf("Number %d is invalid", gen)
		}
	}

	t.Run("Check generated numbers", func(t *testing.T) {
		t.Run("First", checkValid)
		t.Run("Second", checkValid)
		t.Run("Third", checkValid)
		t.Run("Fourth", checkValid)
		t.Run("Fifth", checkValid)
	})
}

func TestNumgenRange(t *testing.T) {
	numGen := GetNumGen()

	checkValid := func(t *testing.T) {
		gen := numGen.Gen()

		if gen < 1234 || gen > 9876 {
			t.Errorf("Number %d is invalid", gen)
		}
	}

	t.Run("Check generated numbers", func(t *testing.T) {
		t.Run("First", checkValid)
		t.Run("Second", checkValid)
		t.Run("Third", checkValid)
		t.Run("Fourth", checkValid)
		t.Run("Fifth", checkValid)
	})
}
