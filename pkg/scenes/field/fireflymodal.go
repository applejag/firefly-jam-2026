package field

import (
	"firefly-jam-2026/assets"
	"firefly-jam-2026/pkg/state"
	"firefly-jam-2026/pkg/util"
	"strconv"
	"strings"

	"github.com/firefly-zero/firefly-go/firefly"
)

type FireflyModal struct {
	isOpen          bool
	scrollOpenAnim  util.AnimatedSheet
	scrollCloseAnim util.AnimatedSheet
	scrollSprite    firefly.SubImage
	tournamentAnim  util.AnimatedSheet
	firefly         *Firefly
}

func (m *FireflyModal) IsOpen() bool {
	return m.isOpen || !m.scrollCloseAnim.IsPaused()
}

func (m *FireflyModal) IsClosing() bool {
	return !m.isOpen && !m.scrollCloseAnim.IsPaused()
}

func (m *FireflyModal) Open(firefly *Firefly) {
	m.scrollOpenAnim.Play()
	m.isOpen = true
	m.firefly = firefly
}

func (m *FireflyModal) Close() {
	if m.IsClosing() {
		return
	}

	m.scrollCloseAnim.Play()
	m.firefly = nil
	m.isOpen = false
}

func (m *FireflyModal) Boot() {
	m.scrollOpenAnim = assets.ScrollOpen.Animated(12)
	m.scrollOpenAnim.AutoPlay = false
	m.scrollOpenAnim.Stop()
	m.scrollCloseAnim = assets.ScrollClose.Animated(12)
	m.scrollCloseAnim.AutoPlay = false
	m.scrollCloseAnim.Stop()
	m.scrollSprite = assets.ScrollClose[0]
	m.tournamentAnim = assets.TournamentButton.Animated(6)
}

func (m *FireflyModal) Update() {
	m.scrollOpenAnim.Update()
	m.scrollCloseAnim.Update()
	m.tournamentAnim.Update()

	if m.IsClosing() {
		return
	}

	buttons := firefly.ReadButtons(firefly.GetMe())
	if buttons.E {
		m.Close()
	}
}

func (m *FireflyModal) Render() {
	const scrollWidth = 132
	point := firefly.P(firefly.Width/2-scrollWidth/2, 24)
	m.scrollOpenAnim.Draw(point)
	m.scrollCloseAnim.Draw(point)

	if m.isOpen && m.scrollCloseAnim.IsPaused() && m.scrollOpenAnim.IsPaused() {
		m.renderScroll(point)
	}
}

func (m *FireflyModal) renderScroll(point firefly.Point) {
	m.scrollSprite.Draw(point)
	assets.Exit.Draw(point.Add(firefly.P(88, 2)))

	dataIndex := state.Game.FindFireflyByID(m.firefly.id)
	if dataIndex == -1 {
		panic("should never be -1 here")
	}
	data := state.Game.Fireflies[dataIndex]

	const scrollInnerWidth = 92
	const scrollInnerHeight = 92

	innerScrollPoint := point.Add(firefly.P(21, 20))

	text := util.WordWrap(
		data.Name.String(),
		scrollInnerWidth,
		assets.FontEG_6x9.CharWidth(),
	)

	charHeight := assets.FontEG_6x9.CharHeight()

	textPos := innerScrollPoint.Add(firefly.P(0, 10))
	assets.FontEG_6x9.Draw(text, textPos, firefly.ColorDarkGray)
	textHeight := charHeight * (strings.Count(text, "\n") + 1)

	speedPoint := textPos.Add(firefly.P(2, textHeight+4))
	assets.FontEG_6x9.Draw(strconv.Itoa(data.Speed), speedPoint, firefly.ColorBlack)
	assets.FontPico8_4x6.Draw("SPEED", speedPoint.Add(firefly.P(0, charHeight)), firefly.ColorGray)

	nimblenessPoint := textPos.Add(firefly.P(34, textHeight+4))
	assets.FontEG_6x9.Draw(strconv.Itoa(data.Nimbleness), nimblenessPoint, firefly.ColorBlack)
	assets.FontPico8_4x6.Draw("NIMBLE", nimblenessPoint.Add(firefly.P(0, charHeight)), firefly.ColorGray)

	rectPoint := textPos.Add(firefly.P(64, textHeight+4-charHeight))
	rectSize := firefly.S(22, 22)
	firefly.DrawRoundedRect(rectPoint, rectSize, firefly.S(3, 3), firefly.Outlined(firefly.ColorGray, 1))

	assets.FireflySheet[0].Draw(rectPoint.Add(firefly.P(6, 6)))

	changeHatPoint := innerScrollPoint.Add(firefly.P(0, scrollInnerHeight-30))
	assets.FontPico8_4x6.Draw("- CHANGE HAT", changeHatPoint, firefly.ColorGray)

	giveVitaminsPoint := changeHatPoint.Add(firefly.P(0, 9))
	assets.FontPico8_4x6.Draw("> GIVE VITAMINS", giveVitaminsPoint, firefly.ColorBlack)

	tournamentPoint := giveVitaminsPoint.Add(firefly.P(0, 14))
	m.tournamentAnim.Draw(tournamentPoint.Add(firefly.P(8, -9)))
	assets.FontPico8_4x6.Draw("-", tournamentPoint, firefly.ColorGray)
}
