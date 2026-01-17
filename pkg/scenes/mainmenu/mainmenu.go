package mainmenu

import (
	"firefly-jam-2026/assets"
	"firefly-jam-2026/pkg/util"

	"github.com/firefly-zero/firefly-go/firefly"
)

type Menu struct {
	TitleScreen     util.AnimatedSheet
	ButtonHighlight util.AnimatedSheet
	Transition      util.Transition
	Button          Button
	lastInput       firefly.Pad
}

type Button byte

const (
	ButtonNone Button = iota
	ButtonContinue
	ButtonNewGame
)

func (b Button) HighlightPosition() (firefly.Point, bool) {
	switch b {
	case ButtonContinue:
		return firefly.P(72, 63), true
	case ButtonNewGame:
		return firefly.P(72, 81), true
	default:
		return firefly.Point{}, false
	}
}

func (m *Menu) Boot() {
	m.TitleScreen = assets.TitleScreen.Animated(2)
	m.ButtonHighlight = assets.TitleButtonHighlight.Animated(2)
	m.Transition = util.NewTransition(assets.TransitionSheet.Animated(10), firefly.S(8, 8))
}

func (m *Menu) Update() {
	m.TitleScreen.Update()
	m.ButtonHighlight.Update()
	m.Transition.Update()

	pad, ok := firefly.ReadPad(firefly.Combined)
	if ok {
		justPressed := pad.DPad4().JustPressed(m.lastInput.DPad4())
		switch {
		case m.Button == ButtonNone:
			m.Button = ButtonContinue
		case justPressed == firefly.DPad4Up:
			m.Button = ButtonContinue
		case justPressed == firefly.DPad4Down:
			m.Button = ButtonNewGame
		}
	}
	m.lastInput = pad

	if firefly.ReadButtons(firefly.Combined).S {
		m.Transition.Play()
	}
}

func (m *Menu) Render() {
	m.TitleScreen.Draw(firefly.P(0, 0))
	assets.TitleNoContinue.Draw(firefly.P(91, 61))
	if pos, ok := m.Button.HighlightPosition(); ok {
		m.ButtonHighlight.Draw(pos)
	}
	m.Transition.Draw()
}
