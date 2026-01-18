package mainmenu

import (
	"firefly-jam-2026/assets"
	"firefly-jam-2026/pkg/scenes"
	"firefly-jam-2026/pkg/state"
	"firefly-jam-2026/pkg/util"

	"github.com/firefly-zero/firefly-go/firefly"
)

type Menu struct {
	TitleScreen     util.AnimatedSheet
	ButtonHighlight util.AnimatedSheet
	Button          Button
	hasSaveFile     bool
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
	m.hasSaveFile = state.Game.HasSave()
}

func (m *Menu) Update() {
	m.TitleScreen.Update()
	m.ButtonHighlight.Update()

	if justPressed := state.Input.JustPressedDPad4(); justPressed != firefly.DPad4None {
		switch {
		case m.Button == ButtonNone:
			m.Button = ButtonContinue
		case justPressed == firefly.DPad4Up:
			m.Button = ButtonContinue
		case justPressed == firefly.DPad4Down:
			m.Button = ButtonNewGame
		}
	}

	if state.Input.JustPressedButtons().S {
		switch m.Button {
		case ButtonNewGame:
			scenes.SwitchScene(scenes.Shop)
		case ButtonContinue:
			if m.hasSaveFile && state.Game.LoadSave() {
				scenes.SwitchScene(scenes.Field)
			}
		}
	}
}

func (m *Menu) Render() {
	m.TitleScreen.Draw(firefly.P(0, 0))
	if !m.hasSaveFile {
		assets.TitleNoContinue.Draw(firefly.P(91, 61))
	}
	if pos, ok := m.Button.HighlightPosition(); ok {
		m.ButtonHighlight.Draw(pos)
	}
}
