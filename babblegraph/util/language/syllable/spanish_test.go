package syllable

import "testing"

func testWord(t *testing.T, word string, expected int64) {
	count, err := countSyllablesForSpanish(word)
	if err != nil {
		t.Errorf("Got error on test for %s: %s", word, err.Error())
	}
	if *count != expected {
		t.Errorf("Error on %s. Expected %d, but got %d", word, expected, *count)
	}
}

func TestCountSyllables(t *testing.T) {
	testWord(t, "cuando", 2)
	testWord(t, "alcanzar", 3)
	testWord(t, "sábana", 3)
	testWord(t, "oro", 2)
	testWord(t, "sombrilla", 3)
}
