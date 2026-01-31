package game

import "testing"

func TestFormatPaddedInt(t *testing.T) {
	tests := []struct {
		name  string
		num   int
		width int
		want  string
	}{
		{
			name:  "zero width",
			num:   1,
			width: 0,
			want:  "",
		},
		{
			name:  "width shorter than value",
			num:   99,
			width: 1,
			want:  "9",
		},
		{
			name:  "just the correct width",
			num:   99,
			width: 2,
			want:  "99",
		},
		{
			name:  "shorter than width",
			num:   99,
			width: 4,
			want:  "0099",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := formatPaddedInt(test.num, test.width)
			if got != test.want {
				t.Errorf("wrong result from formatPaddedInt(%d, %d)\nwant: %q\ngot:  %q",
					test.num, test.width, test.want, got)
			}
		})
	}
}
