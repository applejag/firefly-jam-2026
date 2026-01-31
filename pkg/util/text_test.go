package util

import "testing"

func TestConcatInto(t *testing.T) {
	tests := []struct {
		name    string
		buf     []byte
		strings []string
		want    string
	}{
		{
			name:    "no strings nil buffer",
			buf:     nil,
			strings: nil,
			want:    "",
		},
		{
			name:    "no strings empty buffer",
			buf:     make([]byte, 0),
			strings: nil,
			want:    "",
		},
		{
			name:    "just big enough buffer",
			buf:     make([]byte, 10),
			strings: []string{"hello", "world"},
			want:    "helloworld",
		},
		{
			name:    "bigger buffer",
			buf:     make([]byte, 14),
			strings: []string{"hello", "world"},
			want:    "helloworld",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			written := ConcatInto(test.buf, test.strings...)
			got := string(test.buf[:written])
			if got != test.want {
				t.Errorf("wrong result\nwant: %q (len=%d)\ngot:  %q (len=%d)", test.want, len(test.want), got, len(got))
			}
			t.Logf("buf: %q", test.buf)
		})
	}
}

func TestFormatIntInto(t *testing.T) {
	tests := []struct {
		name string
		buf  []byte
		num  int
		want string
	}{
		{
			name: "just big enough/5",
			buf:  make([]byte, 1),
			num:  5,
			want: "5",
		},
		{
			name: "just big enough/12",
			buf:  make([]byte, 2),
			num:  12,
			want: "12",
		},
		{
			name: "just big enough/333",
			buf:  make([]byte, 3),
			num:  333,
			want: "333",
		},
		{
			name: "bigger buffer/5",
			buf:  make([]byte, 10),
			num:  5,
			want: "5",
		},
		{
			name: "bigger buffer/12",
			buf:  make([]byte, 10),
			num:  12,
			want: "12",
		},
		{
			name: "bigger buffer/333",
			buf:  make([]byte, 10),
			num:  333,
			want: "333",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			written := FormatIntInto(test.buf, test.num)
			got := string(test.buf[:written])
			if got != test.want {
				t.Errorf("wrong result\nwant: %q (len=%d)\ngot:  %q (len=%d)", test.want, len(test.want), got, len(got))
			}
			t.Logf("buf: %q", test.buf)
		})
	}
}
