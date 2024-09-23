package luna

import (
	"errors"
	"testing"
)

func TestCheck(t *testing.T) {
	testCases := []struct {
		name  string
		input string
		want  bool
		err   error
	}{
		{
			name:  "valid_lunа_number_1",
			input: "7562702519",
			want:  true,
			err:   nil,
		},
		{
			name:  "valid_lunа_number_2",
			input: "33775703563187",
			want:  true,
			err:   nil,
		},
		{
			name:  "valid_lunа_number_3",
			input: "547086702",
			want:  true,
			err:   nil,
		},
		{
			name:  "invalid_lunа_number",
			input: "49927398717",
			want:  false,
			err:   nil,
		},
		{
			name:  "invalid_chars_in_number",
			input: "abc123",
			want:  false,
			err:   ErrorIncorrectNumber,
		},
		{
			name:  "empty_number",
			input: "",
			want:  false,
			err:   ErrorIncorrectNumber,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := Check(tc.input)
			if tc.err == nil {
				if got != tc.want {
					t.Errorf("Check(%q) = %v; want %v", tc.input, got, tc.want)
				}
			} else {
				if err == nil || !errors.Is(err, tc.err) {
					t.Errorf("Check(%q) = %v, %v; want %v, %v", tc.input, got, err, tc.want, tc.err)
				}
			}
		})
	}
}
