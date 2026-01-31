package game

import (
	"strings"

	"github.com/applejag/epic-wizard-firefly-gladiators/assets"
	"github.com/applejag/epic-wizard-firefly-gladiators/pkg/state"

	"github.com/firefly-zero/firefly-go/firefly"
)

type UI struct {
	cachedMoneyValue     int
	cachedMoneyText      string
	cachedFirefliesValue int
	cachedFirefliesText  string
}

func (u *UI) Render() {
	if len(state.Game.Fireflies) == 0 {
		// wait until player has at least 1 firefly
		return
	}
	assets.CashBanner.Draw(firefly.P(2, 2))

	moneyText := u.cachedMoneyText
	money := min(state.Game.Money, 9999)
	if moneyText == "" || money != u.cachedMoneyValue {
		moneyText = formatPaddedInt(money, 4)
		u.cachedMoneyText = moneyText
		u.cachedMoneyValue = money
	}
	firefliesText := u.cachedFirefliesText
	fireflies := min(len(state.Game.Fireflies), 99)
	if firefliesText == "" || fireflies != u.cachedFirefliesValue {
		firefliesText = formatPaddedInt(fireflies, 2)
		u.cachedFirefliesText = firefliesText
		u.cachedFirefliesValue = fireflies
	}

	drawRightAlignedWithColoredZeros(
		assets.FontEG_6x9,
		moneyText,
		firefly.P(29, 12),
		firefly.ColorLightGray,
		firefly.ColorDarkGray)
	drawRightAlignedWithColoredZeros(
		assets.FontEG_6x9,
		firefliesText,
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

func formatPaddedInt(value, width int) string {
	buf := make([]byte, width)
	for i := range buf {
		buf[i] = '0'
	}
	index := len(buf) - 1
	for value > 0 && index >= 0 {
		buf[index] = '0' + byte(value%10)
		value /= 10
		index -= 1
	}
	return string(buf)
}
