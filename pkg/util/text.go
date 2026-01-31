package util

import (
	"bytes"
	"strings"
	"unsafe"

	"github.com/firefly-zero/firefly-go/firefly"
)

func WordWrap(s string, maxWidth, charWidth int) string {
	if len(s)*charWidth <= maxWidth {
		return s
	}
	var buf [128]byte
	sb := bytes.NewBuffer(buf[:])
	sb.Reset()
	lineWidth := 0

	for word := range strings.FieldsSeq(s) {
		wordWidth := len(word) * charWidth
		if lineWidth == 0 {
			lineWidth += wordWidth
			sb.WriteString(word)
			continue
		}
		wordWidthPlusSpace := wordWidth + charWidth
		if lineWidth+wordWidthPlusSpace > maxWidth {
			sb.WriteByte('\n')
			lineWidth = wordWidth
		} else {
			sb.WriteByte(' ')
			lineWidth += wordWidthPlusSpace
		}
		sb.WriteString(word)
	}
	return sb.String()
}

func DrawTextRightAligned(font firefly.Font, text string, right firefly.Point, color firefly.Color) {
	width := font.CharWidth() * len(text)
	font.Draw(text, right.Add(firefly.P(-width, 0)), color)
}

func DrawTextCentered(font firefly.Font, text string, center firefly.Point, color firefly.Color) {
	width := font.CharWidth() * len(text)
	font.Draw(text, center.Add(firefly.P(-width/2, 0)), color)
}

// ConcatInto writes the given strings into a buffer with zero allocations.
//
// Returns the number of bytes written. Intended usage:
//
//	var buf [10]
//	written := ConcatInto(buf[:], "hello", "world")
//	foobar(string(buf[:written]))
//
// Panics if "buf" is too small.
func ConcatInto(buf []byte, parts ...string) int {
	written := 0
	for i := range parts {
		partBytes := unsafe.Slice(unsafe.StringData(parts[i]), len(parts[i]))
		if len(partBytes) > len(buf) {
			panic("ConcatBuf: buffer is not big enough")
		}
		written += copy(buf, partBytes)
		buf = buf[len(partBytes):]
	}
	return written
}

func FormatIntInto(buf []byte, num int) int {
	if num < 0 {
		buf[0] = '-'
		return FormatIntInto(buf[1:], -num) + 1
	}
	if num < 10 {
		buf[0] = '0' + byte(num)
		return 1
	}
	size := numberOfDigits(num)
	index := size - 1
	for num > 0 && index >= 0 {
		buf[index] = '0' + byte(num%10)
		num /= 10
		index--
	}
	return size
}

func numberOfDigits(num int) int {
	switch {
	case num < 0:
		return numberOfDigits(-num) + 1
	case num < 1e1:
		return 1
	case num < 1e2:
		return 2
	case num < 1e3:
		return 3
	case num < 1e4:
		return 4
	case num < 1e5:
		return 5
	default:
		panic("number is too big")
	}
}
