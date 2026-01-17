package mainmenu

import (
	"firefly-jam-2026/assets"
	"firefly-jam-2026/pkg/util"

	"github.com/firefly-zero/firefly-go/firefly"
)

type Menu struct {
	TitleScreen     util.AnimatedSheet
	ButtonHighlight util.AnimatedSheet
	Transition      Transition
	Button          Button
	lastInput       firefly.Pad
}

type Button byte

const (
	ButtonContinue Button = iota
	ButtonNewGame
)

func (b Button) HighlightPosition() firefly.Point {
	switch b {
	case ButtonContinue:
		return firefly.P(72, 63)
	case ButtonNewGame:
		return firefly.P(72, 81)
	default:
		panic("unexpected mainmenu.Button")
	}
}

func (m *Menu) Boot() {
	m.TitleScreen = assets.TitleScreen.Animated(2)
	m.ButtonHighlight = assets.TitleButtonHighlight.Animated(2)
	m.Transition = NewTransition(assets.TransitionSheet.Animated(10))
}

func (m *Menu) Update() {
	m.TitleScreen.Update()
	m.ButtonHighlight.Update()
	m.Transition.Update()

	pad, ok := firefly.ReadPad(firefly.Combined)
	if ok {
		justPressed := pad.DPad4().JustPressed(m.lastInput.DPad4())
		switch justPressed {
		case firefly.DPad4Up:
			m.Button = ButtonContinue
		case firefly.DPad4Down:
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
	m.ButtonHighlight.Draw(m.Button.HighlightPosition())
	m.Transition.Draw()
}
