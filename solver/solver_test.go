package solver

import (
	"testing"
)

func TestSolver_IncorrectRunePlaces(t *testing.T) {
	type TestCase struct {
		incorrectPlaces []RunePlace
		resultsCount    int
		name            string
	}

	allWords := []string{
		"abcde",
		"aaaaa",
		"bcdef",
		"cdefg",
		"defgh",
	}

	tests := []TestCase{
		{
			incorrectPlaces: make([]RunePlace, 0),
			resultsCount:    len(allWords),
			name:            "Empty incorrect places collection",
		},
		{
			incorrectPlaces: []RunePlace{{Rune: 'a', Pos: 0}},
			resultsCount:    len(allWords) - 2,
			name:            "Single incorrect place(0 pos)",
		},
		{
			incorrectPlaces: []RunePlace{{Rune: 'a', Pos: 1}},
			resultsCount:    len(allWords) - 1,
			name:            "Single incorrect place(1 pos)",
		},
		{
			incorrectPlaces: []RunePlace{{Rune: 'a', Pos: 0}, {Rune: 'b', Pos: 1}},
			resultsCount:    len(allWords) - 2,
			name:            "Many incorrect places",
		},
		{
			incorrectPlaces: []RunePlace{{Rune: 'a', Pos: 0}, {Rune: 'b', Pos: 0}},
			resultsCount:    len(allWords) - 3,
			name:            "Many incorrect places on same position",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSolver(allWords, 5)
			s.IncorrectRunePlaces(tt.incorrectPlaces)
			if suitable := s.GetSuitable(); len(suitable) != tt.resultsCount {
				t.Errorf("GetSuitable() = %d, want %d", len(suitable), tt.resultsCount)
			}
		})
	}

}
