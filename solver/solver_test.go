package solver

import (
	"testing"
)

var allWords = []string{
	"abcde",
	"aaaaa",
	"bcdef",
	"cdefg",
	"defgh",
}

type testCase struct {
	incorrectPlaces []RunePlace
	resultsCount    int
	name            string
}

func TestSolver_IncorrectRunePlaces(t *testing.T) {

	tests := []testCase{
		{
			incorrectPlaces: make([]RunePlace, 0),
			resultsCount:    len(allWords),
			name:            "Empty incorrect places collection",
		},
		{
			incorrectPlaces: []RunePlace{{Rune: 'z', Pos: 0}},
			resultsCount:    len(allWords),
			name:            "Not contained rune",
		},
		{
			incorrectPlaces: []RunePlace{{Rune: 'z', Pos: 0}, {Rune: 'z', Pos: 1}},
			resultsCount:    len(allWords),
			name:            "Many not contained runes",
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
				t.Errorf("GetSuitable() len = %d, want %d", len(suitable), tt.resultsCount)
			}
		})
	}

}

func TestSolver_IncorrectRunePlaces_ConsistentlyAdd(t *testing.T) {
	tests := []testCase{
		{
			incorrectPlaces: []RunePlace{{Rune: 'a', Pos: 0}, {Rune: 'b', Pos: 0}},
			resultsCount:    2,
			name:            "Incorrect places on same position",
		},
		{
			incorrectPlaces: []RunePlace{{Rune: 'a', Pos: 0}, {Rune: 'b', Pos: 1}},
			resultsCount:    3,
			name:            "Incorrect places on different positions",
		},
		{
			incorrectPlaces: []RunePlace{{Rune: 'f', Pos: 0}, {Rune: 'z', Pos: 1}},
			resultsCount:    len(allWords),
			name:            "Many not contained runes",
		},
		{
			incorrectPlaces: []RunePlace{{Rune: 'h', Pos: 4}, {Rune: 'f', Pos: 3}},
			resultsCount:    3,
			name:            "Many not contained runes",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSolver(allWords, 5)
			for _, r := range tt.incorrectPlaces {
				s.IncorrectRunePlaces([]RunePlace{r})
			}
			if suitable := s.GetSuitable(); len(suitable) != tt.resultsCount {
				t.Errorf("GetSuitable() len = %d, want %d", len(suitable), tt.resultsCount)
			}
		})
	}
}
