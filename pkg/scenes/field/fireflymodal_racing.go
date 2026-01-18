package field

import (
	"github.com/applejag/firefly-jam-2026/assets"
	"github.com/applejag/firefly-jam-2026/pkg/scenes"
	"github.com/applejag/firefly-jam-2026/pkg/state"
	"github.com/applejag/firefly-jam-2026/pkg/util"
	"github.com/firefly-zero/firefly-go/firefly"
)

type RacingPage struct {
	trainingAnim   util.AnimatedSheet
	tournamentAnim util.AnimatedSheet
	trainingBtn    Button
	tournamentBtn  Button
	focused        RacingButton
}

func (p *RacingPage) Boot() {
	p.trainingAnim = assets.TrainButton.Animated(2)
	p.tournamentAnim = assets.TournamentButton.Animated(6)
	p.trainingBtn = NewButton("")
	p.tournamentBtn = NewButton("")
}

func (p *RacingPage) Update() {
	p.tournamentAnim.Update()
	p.trainingAnim.Update()

	if justPressed := state.Input.JustPressedDPad4(); justPressed != firefly.DPad4None {
		p.handleInputDPad4(justPressed)
	}
	if justPressed := state.Input.JustPressedButtons(); justPressed.Any() {
		p.handleInputButtons(justPressed)
	}
}

func (p *RacingPage) handleInputDPad4(justPressed firefly.DPad4) {
	switch justPressed {
	case firefly.DPad4Up:
		p.focused = p.focused.Up()
	case firefly.DPad4Down:
		p.focused = p.focused.Down()
	}
}

func (p *RacingPage) handleInputButtons(justPressed firefly.Buttons) {
	switch {
	case justPressed.S:
		switch p.focused {
		case RacingTraining:
			scenes.SwitchScene(scenes.RacingTraining)
		case RacingTournament:
			scenes.SwitchScene(scenes.RacingBattle)
		}
	}
}

func (p *RacingPage) Render(innerScrollPoint firefly.Point) {
	p.trainingAnim.Draw(innerScrollPoint.Add(firefly.P(14, 2)))
	p.trainingBtn.Render(innerScrollPoint.Add(firefly.P(4, 10)), p.focused == RacingTraining)
	assets.FontEG_6x9.Draw(
		"Train firefly\nNo opponents",
		innerScrollPoint.Add(firefly.P(0, 22)),
		firefly.ColorDarkGray,
	)
	assets.FontPico8_4x6.Draw(
		"+SPEED  +NIMBLE",
		innerScrollPoint.Add(firefly.P(0, 39)),
		firefly.ColorYellow,
	)

	p.tournamentAnim.Draw(innerScrollPoint.Add(firefly.P(14, 43)))
	p.tournamentBtn.Render(innerScrollPoint.Add(firefly.P(4, 51)), p.focused == RacingTournament)
	assets.FontEG_6x9.Draw(
		"Race others\nFastest wins",
		innerScrollPoint.Add(firefly.P(0, 63)),
		firefly.ColorDarkGray,
	)
	assets.FontPico8_4x6.Draw(
		"+SPEED  +NIMBLE  +MONEY",
		innerScrollPoint.Add(firefly.P(0, 80)),
		firefly.ColorYellow,
	)
}

type RacingButton byte

const (
	RacingNone RacingButton = iota
	RacingTraining
	RacingTournament
)

func (k RacingButton) Down() RacingButton {
	switch k {
	case RacingNone:
		return RacingTraining
	case RacingTraining:
		return RacingTournament
	case RacingTournament:
		return RacingTraining
	default:
		panic("unexpected field.ButtonKind")
	}
}

func (k RacingButton) Up() RacingButton {
	switch k {
	case RacingNone:
		return RacingTournament
	case RacingTraining:
		return RacingTournament
	case RacingTournament:
		return RacingTraining
	default:
		panic("unexpected field.ButtonKind")
	}
}
