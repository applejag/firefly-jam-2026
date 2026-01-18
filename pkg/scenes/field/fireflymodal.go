package field

import (
	"strconv"
	"strings"

	"github.com/applejag/firefly-jam-2026/assets"
	"github.com/applejag/firefly-jam-2026/pkg/state"
	"github.com/applejag/firefly-jam-2026/pkg/util"

	"github.com/firefly-zero/firefly-go/firefly"
	"github.com/orsinium-labs/tinymath"
)

type ModalState byte

const (
	ModalClosed ModalState = iota
	ModalStats
	ModalHats
	ModalTournament
)

type FireflyModal struct {
	state           ModalState
	scrollOpenAnim  util.AnimatedSheet
	scrollCloseAnim util.AnimatedSheet
	scrollSprite    firefly.SubImage
	tournamentAnim  util.AnimatedSheet
	firefly         *Firefly

	changeHatBtn      Button
	giveVitaminsBtn   Button
	playTournamentBtn Button
	focused           ButtonKind
}

func (m *FireflyModal) IsOpen() bool {
	return m.state != ModalClosed || !m.scrollCloseAnim.IsPaused()
}

func (m *FireflyModal) IsClosing() bool {
	return m.state == ModalClosed && !m.scrollCloseAnim.IsPaused()
}

func (m *FireflyModal) Open(firefly *Firefly) {
	m.scrollOpenAnim.Play()
	m.state = ModalStats
	m.focused = ButtonNone
	m.firefly = firefly
}

func (m *FireflyModal) Close() {
	if m.IsOpen() && m.IsClosing() {
		return
	}

	m.scrollCloseAnim.Play()
	m.firefly = nil
	m.state = ModalClosed
}

func (m *FireflyModal) CloseWithoutTransition() {
	if m.IsOpen() && m.IsClosing() {
		return
	}

	m.scrollCloseAnim.Stop()
	m.firefly = nil
	m.state = ModalClosed
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
	m.changeHatBtn = NewButton(ButtonChangeHat, "CHANGE HAT")
	m.giveVitaminsBtn = NewButton(ButtonGiveVitamins, "GIVE VITAMINS")
	m.giveVitaminsBtn.Disabled = true
	m.playTournamentBtn = NewButton(ButtonTournament, "")
}

func (m *FireflyModal) Update() {
	m.scrollOpenAnim.Update()
	m.scrollCloseAnim.Update()
	m.tournamentAnim.Update()

	m.changeHatBtn.Update()
	m.giveVitaminsBtn.Update()
	m.playTournamentBtn.Update()

	if m.IsClosing() {
		return
	}

	if justPressed := state.Input.JustPressedDPad4(); justPressed != firefly.DPad4None {
		m.handleInputDPad4(justPressed)
	}
	if justPressed := state.Input.JustPressedButtons(); justPressed.Any() {
		m.handleInputButtons(justPressed)
	}
}

func (m *FireflyModal) handleInputDPad4(justPressed firefly.DPad4) {
	switch justPressed {
	case firefly.DPad4Up:
		m.focused = m.focused.Up()
	case firefly.DPad4Down:
		m.focused = m.focused.Down()
	}
}

func (m *FireflyModal) handleInputButtons(justPressed firefly.Buttons) {
	switch {
	case justPressed.E:
		if m.state == ModalStats {
			m.Close()
		} else {
			m.state = ModalStats
		}

	case justPressed.S:
		switch m.state {
		case ModalStats:
			// yo i heard you like nested switch statements?
			switch m.focused {
			case ButtonChangeHat:
				m.state = ModalHats
			case ButtonGiveVitamins:
				// Shake to signify that the button doesn't work
				m.giveVitaminsBtn.Shake()
			case ButtonTournament:
				m.state = ModalTournament

				// state.Game.AddMyFireflyToRaceBattle(m.firefly.id)
				// m.CloseWithoutTransition()
				// scenes.SwitchScene(scenes.RacingTraining)
			}
		}
	}
}

func (m *FireflyModal) Render() {
	const scrollWidth = 132
	point := firefly.P(firefly.Width/2-scrollWidth/2, 24)
	m.scrollOpenAnim.Draw(point)
	m.scrollCloseAnim.Draw(point)

	if m.state != ModalClosed && m.scrollCloseAnim.IsPaused() && m.scrollOpenAnim.IsPaused() {
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

	switch m.state {
	case ModalStats:
		text := util.WordWrap(
			data.Name.String(),
			scrollInnerWidth,
			assets.FontEG_6x9.CharWidth(),
		)

		charHeight := assets.FontEG_6x9.CharHeight()

		textPos := innerScrollPoint.Add(firefly.P(0, 10))
		assets.FontEG_6x9.Draw(text, textPos, firefly.ColorDarkGray)
		textHeight := charHeight * (strings.Count(text, "\n") + 1)

		speedPoint := textPos.Add(firefly.P(2, textHeight))
		assets.FontEG_6x9.Draw(strconv.Itoa(data.Speed), speedPoint, firefly.ColorBlack)
		assets.FontPico8_4x6.Draw("SPEED", speedPoint.Add(firefly.P(0, charHeight)), firefly.ColorGray)

		nimblenessPoint := speedPoint.Add(firefly.P(32, 0))
		assets.FontEG_6x9.Draw(strconv.Itoa(data.Nimbleness), nimblenessPoint, firefly.ColorBlack)
		assets.FontPico8_4x6.Draw("NIMBLE", nimblenessPoint.Add(firefly.P(0, charHeight)), firefly.ColorGray)

		rectPoint := textPos.Add(firefly.P(64, textHeight+4-charHeight))
		rectSize := firefly.S(22, 22)
		firefly.DrawRoundedRect(rectPoint, rectSize, firefly.S(3, 3), firefly.Outlined(firefly.ColorGray, 1))

		assets.FireflySheet[0].Draw(rectPoint.Add(firefly.P(6, 6)))

		changeHatPoint := innerScrollPoint.Add(firefly.P(0, scrollInnerHeight-26))
		m.changeHatBtn.Render(changeHatPoint, m.focused)

		giveVitaminsPoint := changeHatPoint.Add(firefly.P(0, 8))
		m.giveVitaminsBtn.Render(giveVitaminsPoint, m.focused)

		tournamentPoint := giveVitaminsPoint.Add(firefly.P(0, 13))
		m.tournamentAnim.Draw(tournamentPoint.Add(firefly.P(8, -9)))
		m.playTournamentBtn.Render(tournamentPoint, m.focused)

		assets.TrainButton.Draw(tournamentPoint.Add(firefly.P(72, -11)))
	}
}

type ButtonKind byte

const (
	ButtonNone ButtonKind = iota
	ButtonChangeHat
	ButtonGiveVitamins
	ButtonTournament

	buttonCount = 4
)

func (k ButtonKind) Down() ButtonKind {
	switch k {
	case ButtonChangeHat:
		return ButtonGiveVitamins
	case ButtonGiveVitamins:
		return ButtonTournament
	case ButtonNone:
		return ButtonChangeHat
	case ButtonTournament:
		return ButtonChangeHat
	default:
		panic("unexpected field.ButtonKind")
	}
}

func (k ButtonKind) Up() ButtonKind {
	switch k {
	case ButtonChangeHat:
		return ButtonTournament
	case ButtonGiveVitamins:
		return ButtonChangeHat
	case ButtonNone:
		return ButtonTournament
	case ButtonTournament:
		return ButtonGiveVitamins
	default:
		panic("unexpected field.ButtonKind")
	}
}

const ButtonShakeDuration = 45

type Button struct {
	kind     ButtonKind
	text     string
	Disabled bool
	shake    int
}

func NewButton(kind ButtonKind, text string) Button {
	return Button{
		kind: kind,
		text: text,
	}
}

func (b *Button) Update() {
	if b.shake > 0 {
		b.shake--
	}
}

func (b *Button) Render(point firefly.Point, focused ButtonKind) {
	prefix := "- "
	color := firefly.ColorGray
	if focused == b.kind {
		prefix = "> "
		color = firefly.ColorBlack
	}
	if b.Disabled {
		color = firefly.ColorLightGray
	}
	assets.FontPico8_4x6.Draw(prefix, point, color)
	if b.text != "" {
		if b.shake > 0 {
			t := float32(b.shake) / ButtonShakeDuration
			point = point.Add(firefly.P(int(tinymath.Sin(t*45)*t*4), 0))
		}
		textPoint := point.Add(firefly.P(assets.FontPico8_4x6.LineWidth(prefix), 0))
		assets.FontPico8_4x6.Draw(b.text, textPoint, color)
		if b.Disabled {
			// Draw strikethrough
			firefly.DrawLine(textPoint, textPoint.Add(firefly.P(
				assets.FontPico8_4x6.LineWidth(b.text),
				-assets.FontEG_6x9.CharHeight()/2,
			)), firefly.L(firefly.ColorGray, 1))
		}
	}
}

func (b *Button) Shake() {
	b.shake = ButtonShakeDuration
}
