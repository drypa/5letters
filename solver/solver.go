package solver

type Solver struct {
	allWords            []string
	containRune         []rune
	notContainRune      []rune
	correctRunePlaces   []rune
	incorrectRunePlaces []RunePlace
}

type RunePlace struct {
	Rune rune
	Pos  int
}

func NewSolver(words []string, len int) *Solver {
	return &Solver{
		allWords:            words,
		correctRunePlaces:   make([]rune, len),
		incorrectRunePlaces: make([]RunePlace, 0),
	}
}

// Contains Appends existing runes in Solver.
func (s *Solver) Contains(runes []rune) {
	for _, r := range runes {
		contain := false
		for _, cr := range s.containRune {
			if r == cr {
				contain = true
			}
		}
		if !contain {
			s.containRune = append(s.containRune, r)
		}
	}
}

// NotContain Appends not existing runes in Solver.
func (s *Solver) NotContain(runes []rune) {
	for _, r := range runes {
		contain := false
		for _, cr := range s.notContainRune {
			if r == cr {
				contain = true
			}
		}
		if !contain {
			s.notContainRune = append(s.notContainRune, r)
		}
	}
}

// CorrectRunePlaces appends runes in right places.
func (s *Solver) CorrectRunePlaces(runes []RunePlace) {
	for _, r := range runes {
		s.correctRunePlaces[r.Pos] = r.Rune
	}
}

// IncorrectRunePlaces append runes in incorrect places.
func (s *Solver) IncorrectRunePlaces(runes []RunePlace) {
	for _, r := range runes {
		s.incorrectRunePlaces = append(s.incorrectRunePlaces, r)
	}
}

func (s *Solver) GetSuitable() []string {
	result := s.allWords
	if s.containRune != nil && len(s.containRune) > 0 {
		result = contains(result, s.containRune)
	}
	if s.notContainRune != nil && len(s.notContainRune) > 0 {
		result = notContainsAny(result, s.notContainRune)
	}
	if s.correctRunePlaces != nil && len(s.correctRunePlaces) > 0 {
		result = correctRunePlaces(result, s.correctRunePlaces)
	}
	if s.incorrectRunePlaces != nil && len(s.incorrectRunePlaces) > 0 {
		result = incorrectRunePlaces(result, s.incorrectRunePlaces)
	}

	return result
}
