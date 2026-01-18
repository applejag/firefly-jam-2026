package game

import (
	"strconv"

	"github.com/applejag/firefly-jam-2026/assets"
	"github.com/applejag/firefly-jam-2026/pkg/state"
	"github.com/applejag/firefly-jam-2026/pkg/util"

	"github.com/firefly-zero/firefly-go/firefly"
)

type UI struct{}

func (u *UI) Render() {
	if len(state.Game.Fireflies) == 0 {
		// wait until player has at least 1 firefly
		return
	}
	assets.CashBanner.Draw(firefly.P(2, 2))
	util.DrawTextRightAligned(assets.FontEG_6x9, "0000", firefly.P(30, 12), firefly.ColorDarkGray)
	util.DrawTextRightAligned(assets.FontEG_6x9, strconv.Itoa(len(state.Game.Fireflies)), firefly.P(59, 12), firefly.ColorDarkGray)
}
