package field

import (
	"strconv"
	"strings"

	"github.com/applejag/epic-wizard-firefly-gladiators/assets"
	"github.com/applejag/epic-wizard-firefly-gladiators/pkg/state"
	"github.com/applejag/epic-wizard-firefly-gladiators/pkg/util"
	"github.com/firefly-zero/firefly-go/firefly"
)

type StatsPage struct {
	changeHatBtn      Button
	giveVitaminsBtn   Button
	playTournamentBtn Button
	focused           StatsButton

	cachedFireflyNameText string
	cachedSpeedText       string
	cachedNimblenessText  string
}

func (p *StatsPage) Boot() {
	p.changeHatBtn = NewButton("CHANGE HAT")
	p.changeHatBtn.Disabled = true
	p.giveVitaminsBtn = NewButton("GIVE VITAMINS")
	p.giveVitaminsBtn.Disabled = true
	p.playTournamentBtn = NewButton("RACING")
}

func (p *StatsPage) OnOpen() {
	p.cachedFireflyNameText = ""
	p.cachedSpeedText = ""
	p.cachedNimblenessText = ""
}

func (p *StatsPage) Update(modal *FireflyModal) {
	p.changeHatBtn.Update()
	p.giveVitaminsBtn.Update()
	p.playTournamentBtn.Update()

	if justPressed := state.Input.JustPressedDPad4(); justPressed != firefly.DPad4None {
		p.handleInputDPad4(justPressed)
	}
	if justPressed := state.Input.JustPressedButtons(); justPressed.Any() {
		p.handleInputButtons(justPressed, modal)
	}
}

func (p *StatsPage) handleInputDPad4(justPressed firefly.DPad4) {
	switch justPressed {
	case firefly.DPad4Up:
		p.focused = p.focused.Up()
	case firefly.DPad4Down:
		p.focused = p.focused.Down()
	}
}

func (p *StatsPage) handleInputButtons(justPressed firefly.Buttons, modal *FireflyModal) {
	switch {
	case justPressed.S:
		switch p.focused {
		case StatsChangeHat:
			// Shake to signify that the button doesn't work
			p.changeHatBtn.Shake()
			// TODO: allow transition to hats page
			// m.state = ModalHats
		case StatsGiveVitamins:
			// Shake to signify that the button doesn't work
			p.giveVitaminsBtn.Shake()
		case StatsRacing:
			modal.OpenPage(ModalRacing)
		}
	}
}

func (p *StatsPage) Render(innerScrollPoint firefly.Point, fireflyID int) {
	dataIndex := state.Game.FindFireflyByID(fireflyID)
	if dataIndex == -1 {
		panic("should never be -1 here")
	}
	data := state.Game.Fireflies[dataIndex]

	text := p.cachedFireflyNameText
	if text == "" {
		var buf [util.LongestPossibleName]byte
		text = util.WordWrap(
			string(buf[:data.Name.WriteInto(buf[:])]),
			scrollInnerWidth,
			assets.FontEG_6x9.CharWidth(),
		)
		p.cachedFireflyNameText = text
	}

	speedText := p.cachedSpeedText
	if speedText == "" {
		speedText = strconv.Itoa(data.Speed)
		p.cachedSpeedText = speedText
	}
	nimblenessText := p.cachedNimblenessText
	if nimblenessText == "" {
		nimblenessText = strconv.Itoa(data.Nimbleness)
		p.cachedNimblenessText = speedText
	}

	charHeight := assets.FontEG_6x9.CharHeight()

	textPos := innerScrollPoint.Add(firefly.P(0, 10))
	assets.FontEG_6x9.Draw(text, textPos, firefly.ColorDarkGray)
	textHeight := charHeight * (strings.Count(text, "\n") + 1)

	speedPoint := textPos.Add(firefly.P(2, textHeight))
	assets.FontEG_6x9.Draw(speedText, speedPoint, firefly.ColorBlack)
	assets.FontPico8_4x6.Draw("SPEED", speedPoint.Add(firefly.P(0, charHeight)), firefly.ColorGray)

	nimblenessPoint := speedPoint.Add(firefly.P(32, 0))
	assets.FontEG_6x9.Draw(nimblenessText, nimblenessPoint, firefly.ColorBlack)
	assets.FontPico8_4x6.Draw("NIMBLE", nimblenessPoint.Add(firefly.P(0, charHeight)), firefly.ColorGray)

	rectPoint := textPos.Add(firefly.P(64, textHeight+4-charHeight))
	rectSize := firefly.S(22, 22)
	firefly.DrawRoundedRect(rectPoint, rectSize, firefly.S(3, 3), firefly.Outlined(firefly.ColorGray, 1))

	assets.FireflySheet[0].Draw(rectPoint.Add(firefly.P(6, 6)))

	changeHatPoint := innerScrollPoint.Add(firefly.P(0, scrollInnerHeight-26))
	p.changeHatBtn.Render(changeHatPoint, p.focused == StatsChangeHat)

	giveVitaminsPoint := changeHatPoint.Add(firefly.P(0, 8))
	p.giveVitaminsBtn.Render(giveVitaminsPoint, p.focused == StatsGiveVitamins)

	tournamentPoint := giveVitaminsPoint.Add(firefly.P(0, 8))
	p.playTournamentBtn.Render(tournamentPoint, p.focused == StatsRacing)

	// m.tournamentAnim.Draw(tournamentPoint.Add(firefly.P(8, -9)))
	// m.playTournamentBtn.Render(tournamentPoint, m.focused)
	// assets.TrainButton.Draw(tournamentPoint.Add(firefly.P(72, -11)))
}

type StatsButton byte

const (
	StatsNone StatsButton = iota
	StatsChangeHat
	StatsGiveVitamins
	StatsRacing
)

func (k StatsButton) Down() StatsButton {
	switch k {
	case StatsChangeHat:
		return StatsGiveVitamins
	case StatsGiveVitamins:
		return StatsRacing
	case StatsNone:
		return StatsChangeHat
	case StatsRacing:
		return StatsChangeHat
	default:
		panic("unexpected field.ButtonKind")
	}
}

func (k StatsButton) Up() StatsButton {
	switch k {
	case StatsChangeHat:
		return StatsRacing
	case StatsGiveVitamins:
		return StatsChangeHat
	case StatsNone:
		return StatsRacing
	case StatsRacing:
		return StatsGiveVitamins
	default:
		panic("unexpected field.ButtonKind")
	}
}
