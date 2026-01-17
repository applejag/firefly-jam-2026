package mainmenu

import (
	"firefly-jam-2026/assets"
	"firefly-jam-2026/pkg/scenes"
	"firefly-jam-2026/pkg/util"

	"github.com/firefly-zero/firefly-go/firefly"
)

type Menu struct {
	TitleScreen     util.AnimatedSheet
	ButtonHighlight util.AnimatedSheet
	Button          Button
	padOld          firefly.Pad
	buttonsOld      firefly.Buttons
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
}

func (m *Menu) Update() {
	m.TitleScreen.Update()
	m.ButtonHighlight.Update()

	pad, ok := firefly.ReadPad(firefly.GetMe())
	if ok {
		justPressed := pad.DPad4().JustPressed(m.padOld.DPad4())
		switch {
		case m.Button == ButtonNone:
			m.Button = ButtonContinue
		case justPressed == firefly.DPad4Up:
			m.Button = ButtonContinue
		case justPressed == firefly.DPad4Down:
			m.Button = ButtonNewGame
		}
	}
	m.padOld = pad

	buttons := firefly.ReadButtons(firefly.GetMe())
	if buttons.JustPressed(m.buttonsOld).S {
		switch m.Button {
		case ButtonNewGame:
			scenes.SwitchScene(scenes.Shop)
		}
	}
	m.buttonsOld = buttons
}

func (m *Menu) Render() {
	m.TitleScreen.Draw(firefly.P(0, 0))
	assets.TitleNoContinue.Draw(firefly.P(91, 61))
	if pos, ok := m.Button.HighlightPosition(); ok {
		m.ButtonHighlight.Draw(pos)
	}
}
