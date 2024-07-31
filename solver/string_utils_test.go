package solver

import "testing"

func TestContains(t *testing.T) {
	tests := []struct {
		name  string
		lines []string
		runes []rune
		want  []string
	}{
		{
			name:  "all runes present",
			lines: []string{"hello", "world", "hero"},
			runes: []rune{'h', 'e'},
			want:  []string{"hello", "hero"},
		},
		{
			name:  "rune not present",
			lines: []string{"hello", "world", "hero"},
			runes: []rune{'h', 'z'},
			want:  []string{},
		},
		{
			name:  "all lines contain runes",
			lines: []string{"abc", "bac", "cab"},
			runes: []rune{'a', 'b'},
			want:  []string{"abc", "bac", "cab"},
		},
		{
			name:  "empty lines input",
			lines: []string{},
			runes: []rune{'a'},
			want:  []string{},
		},
		{
			name:  "empty runes input",
			lines: []string{"hello", "world"},
			runes: []rune{},
			want:  []string{"hello", "world"},
		},
		{
			name:  "empty lines and runes",
			lines: []string{},
			runes: []rune{},
			want:  []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := contains(tt.lines, tt.runes); !equal(got, tt.want) {
				t.Errorf("contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Вспомогательная функция для сравнения двух слайсов строк
func equal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
