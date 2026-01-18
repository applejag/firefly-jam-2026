package util

import (
	"bytes"
	"strings"

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
