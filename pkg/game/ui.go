package game

import (
	"firefly-jam-2026/assets"
	"firefly-jam-2026/pkg/state"
	"strconv"

	"github.com/firefly-zero/firefly-go/firefly"
)

type UI struct{}

func (u *UI) Render() {
	if len(state.Game.Fireflies) == 0 {
		// wait until player has at least 1 firefly
		return
	}
	assets.CashBanner.Draw(firefly.P(2, 2))
	drawRightAligned(assets.FontEG_6x9, "0000", firefly.P(30, 12), firefly.ColorDarkGray)
	drawRightAligned(assets.FontEG_6x9, strconv.Itoa(len(state.Game.Fireflies)), firefly.P(59, 12), firefly.ColorDarkGray)
}

func drawRightAligned(font firefly.Font, text string, right firefly.Point, color firefly.Color) {
	width := font.CharWidth() * len(text)
	font.Draw(text, right.Add(firefly.P(-width, 0)), color)
}
