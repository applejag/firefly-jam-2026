package game

import (
	"fmt"
	"strings"

	"github.com/applejag/epic-wizard-firefly-gladiators/assets"
	"github.com/applejag/epic-wizard-firefly-gladiators/pkg/state"

	"github.com/firefly-zero/firefly-go/firefly"
)

type UI struct{}

func (u *UI) Render() {
	if len(state.Game.Fireflies) == 0 {
		// wait until player has at least 1 firefly
		return
	}
	assets.CashBanner.Draw(firefly.P(2, 2))
	drawRightAlignedWithColoredZeros(
		assets.FontEG_6x9,
		fmt.Sprintf("%04d", min(state.Game.Money, 9999)),
		firefly.P(29, 12),
		firefly.ColorLightGray,
		firefly.ColorDarkGray)
	drawRightAlignedWithColoredZeros(
		assets.FontEG_6x9,
		fmt.Sprintf("%02d", min(len(state.Game.Fireflies), 99)),
		firefly.P(59, 12),
		firefly.ColorLightGray,
		firefly.ColorDarkGray)
}

func drawRightAlignedWithColoredZeros(font firefly.Font, text string, right firefly.Point, zeroColor, textColor firefly.Color) {
	width := font.LineWidth(text)
	left := right.Add(firefly.P(-width, 0))
	withoutZeros := strings.TrimLeft(text, "0")
	zeros := text[:len(text)-len(withoutZeros)]
	font.Draw(zeros, left, zeroColor)
	font.Draw(withoutZeros, left.Add(firefly.P(font.LineWidth(zeros), 0)), textColor)
}
