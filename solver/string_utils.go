package solver

import (
	"strings"
)

func contains(lines []string, runes []rune) []string {
	result := make([]string, 0)
	for _, line := range lines {
		containsAll := true
		for _, char := range runes {
			if !strings.ContainsRune(line, char) {
				containsAll = false
				break
			}
		}
		if containsAll {
			result = append(result, line)
		}
	}
	return result
}

func notContainsAny(lines []string, runes []rune) []string {
	result := make([]string, 0)
	for _, line := range lines {
		containsAny := false
		for _, r := range runes {
			if strings.ContainsRune(line, r) {
				containsAny = true
				break
			}
		}
		if !containsAny {
			result = append(result, line)
		}
	}
	return result
}

func correctRunePlaces(lines []string, runes []rune) []string {
	var result []string
	for _, line := range lines {
		containsAll := true
		for i, r := range runes {
			if r != 0 && runeAt(line, i) != r {
				containsAll = false
				break
			}
		}
		if containsAll {
			result = append(result, line)
		}
	}
	return result
}

func incorrectRunePlaces(lines []string, runes []rune) []string {
	var result []string
	for _, line := range lines {
		containAny := false
		for i, r := range runes {
			if r != 0 && runeAt(line, i) == r {
				containAny = true
				break
			}
		}
		if !containAny {
			result = append(result, line)
		}
	}
	return result
}

func runeAt(s string, pos int) rune {
	if pos < 0 || pos >= len(s) {
		return -1 // Неверная позиция
	}
	for _, r := range s {
		if pos == 0 {
			return r
		}
		pos--
	}
	return -1
}
