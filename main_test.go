package main

import (
	"github.com/stretchr/testify/require"
	"testing"
	"unicode"
)

func Test_GetIncorrectCharPositionFromRequest(t *testing.T) {
	tests := []struct {
		message      string
		expectedChar rune
		expectedPos  int
		expectedErr  bool
	}{
		{
			message:      "2-а",
			expectedChar: 'а',
			expectedPos:  2,
			expectedErr:  false,
		},
		{
			message:      "3-ы",
			expectedChar: 'ы',
			expectedPos:  3,
			expectedErr:  false,
		},
		{
			message:      "2a",
			expectedChar: 0,
			expectedPos:  0,
			expectedErr:  true,
		},
		{
			message:      "а-2",
			expectedChar: 0,
			expectedPos:  0,
			expectedErr:  true,
		},
		{
			message:      "а2",
			expectedChar: 0,
			expectedPos:  0,
			expectedErr:  true,
		},
		{
			message:      "a-а",
			expectedChar: 0,
			expectedPos:  0,
			expectedErr:  true,
		},
		{
			message:      "2-9",
			expectedChar: 0,
			expectedPos:  0,
			expectedErr:  true,
		},
		{
			message:      "12-а",
			expectedChar: 0,
			expectedPos:  0,
			expectedErr:  true,
		},
		{
			message:      "2-абв",
			expectedChar: 0,
			expectedPos:  0,
			expectedErr:  true,
		},
		{
			message:      "",
			expectedChar: 0,
			expectedPos:  0,
			expectedErr:  true,
		},
		{
			message:      "  ",
			expectedChar: 0,
			expectedPos:  0,
			expectedErr:  true,
		},
		{
			message:      "1/а",
			expectedChar: 0,
			expectedPos:  0,
			expectedErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.message, func(t *testing.T) {
			char, pos, err := getIncorrectCharPositionFromRequest(tt.message)

			if tt.expectedErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedChar, char)
				require.Equal(t, tt.expectedPos, pos)
				require.True(t, unicode.IsLetter(char), "Expected character should be a letter")
			}
		})
	}
}
