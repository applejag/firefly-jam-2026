package util

import (
	"testing"

	"github.com/firefly-zero/firefly-go/firefly"
)

// image examples taken from docs: https://docs.fireflyzero.com/internal/formats/image/
var image1BPP = []byte{
	// header
	0x21,       // magic number
	0x01,       // bits per pixel
	0x04, 0x00, // image width
	0xff, // transparency
	0x42, // color palette swap
	// body
	0xc3, // row 1 & row 2, 0b1100_0011
	0x9b, // row 3 & row 4, 0b1001_1011
}

var image2BPP = []byte{
	// header
	0x21,       // magic number
	0x02,       // bits per pixel
	0x04, 0x00, // image width
	0xff,       // transparency
	0x2B, 0x5A, // color palette swap
	// body
	0xec, // row 1
	0xaf, // row 2
	0x50, // row 3
	0x91, // row 4
}

var image4BPP = []byte{
	// header
	0x21,       // magic number
	0x04,       // bits per pixel
	0x04, 0x00, // image width
	0x01,                                           // transparency
	0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef, // color palette swap
	// body
	0x01, 0x23, // row 1
	0x45, 0x67, // row 2
	0x89, 0xab, // row 3
	0xcd, 0xef, // row 4
}

func TestExtImage_GetColorAt(t *testing.T) {
	tests := []struct {
		name  string
		raw   []byte
		pixel firefly.Point
		want  firefly.Color
	}{
		{name: "negative point", raw: image1BPP, pixel: firefly.P(-1, -1), want: firefly.ColorNone},
		{name: "point out of bounds", raw: image1BPP, pixel: firefly.P(100, 100), want: firefly.ColorNone},

		// 1 BPP
		{name: "1 BPP/0x0", raw: image1BPP, pixel: firefly.P(0, 0), want: firefly.ColorRed},
		{name: "1 BPP/1x0", raw: image1BPP, pixel: firefly.P(1, 0), want: firefly.ColorRed},
		{name: "1 BPP/2x0", raw: image1BPP, pixel: firefly.P(2, 0), want: firefly.ColorYellow},
		{name: "1 BPP/3x0", raw: image1BPP, pixel: firefly.P(3, 0), want: firefly.ColorYellow},
		{name: "1 BPP/0x1", raw: image1BPP, pixel: firefly.P(0, 1), want: firefly.ColorYellow},

		// 2 BPP
		{name: "2 BPP/0x0", raw: image2BPP, pixel: firefly.P(0, 0), want: firefly.ColorLightBlue},
		{name: "2 BPP/1x0", raw: image2BPP, pixel: firefly.P(1, 0), want: firefly.ColorLightGreen},
		{name: "2 BPP/2x0", raw: image2BPP, pixel: firefly.P(2, 0), want: firefly.ColorLightBlue},
		{name: "2 BPP/3x0", raw: image2BPP, pixel: firefly.P(3, 0), want: firefly.ColorRed},
		{name: "2 BPP/0x1", raw: image2BPP, pixel: firefly.P(0, 1), want: firefly.ColorLightGreen},
		{name: "2 BPP/1x1", raw: image2BPP, pixel: firefly.P(1, 1), want: firefly.ColorLightGreen},
		{name: "2 BPP/2x1", raw: image2BPP, pixel: firefly.P(2, 1), want: firefly.ColorLightBlue},
		{name: "2 BPP/3x1", raw: image2BPP, pixel: firefly.P(3, 1), want: firefly.ColorLightBlue},
		{name: "2 BPP/0x2", raw: image2BPP, pixel: firefly.P(0, 2), want: firefly.ColorCyan},
		{name: "2 BPP/1x2", raw: image2BPP, pixel: firefly.P(1, 2), want: firefly.ColorCyan},
		{name: "2 BPP/2x2", raw: image2BPP, pixel: firefly.P(2, 2), want: firefly.ColorRed},
		{name: "2 BPP/3x2", raw: image2BPP, pixel: firefly.P(3, 2), want: firefly.ColorRed},
		{name: "2 BPP/2x3", raw: image2BPP, pixel: firefly.P(2, 3), want: firefly.ColorRed},

		// 4 BPP
		{name: "4 BPP/0x0", raw: image4BPP, pixel: firefly.P(0, 0), want: firefly.ColorBlack},
		{name: "4 BPP/1x1", raw: image4BPP, pixel: firefly.P(1, 1), want: firefly.ColorLightGreen},
		{name: "4 BPP/2x3", raw: image4BPP, pixel: firefly.P(2, 3), want: firefly.ColorGray},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			image := NewExtImage(firefly.File{Raw: test.raw})
			got := image.GetColorAt(test.pixel)
			if got != test.want {
				t.Errorf("want %s %#v, but got %s %#v", colorName(test.want), test.want, colorName(got), got)
			}
		})
	}
}

func colorName(color firefly.Color) string {
	switch color {
	case firefly.ColorBlack:
		return "black"
	case firefly.ColorBlue:
		return "blue"
	case firefly.ColorCyan:
		return "cyan"
	case firefly.ColorDarkBlue:
		return "dark blue"
	case firefly.ColorDarkGray:
		return "dark gray"
	case firefly.ColorDarkGreen:
		return "dark green"
	case firefly.ColorGray:
		return "gray"
	case firefly.ColorGreen:
		return "green"
	case firefly.ColorLightBlue:
		return "light blue"
	case firefly.ColorLightGray:
		return "light gray"
	case firefly.ColorLightGreen:
		return "light green"
	case firefly.ColorNone:
		return "none"
	case firefly.ColorOrange:
		return "orange"
	case firefly.ColorPurple:
		return "purple"
	case firefly.ColorRed:
		return "red"
	case firefly.ColorWhite:
		return "white"
	case firefly.ColorYellow:
		return "yellow"
	default:
		panic("unexpected firefly.Color")
	}
}
