package field

import (
	"cmp"
	"firefly-jam-2026/assets"
	"firefly-jam-2026/pkg/scenes"
	"firefly-jam-2026/pkg/state"
	"firefly-jam-2026/pkg/util"
	"slices"

	"github.com/firefly-zero/firefly-go/firefly"
)

type Scene struct {
	dpad4Old   firefly.DPad4
	buttonsOld firefly.Buttons

	fireflies []Firefly
	modal     FireflyModal

	highlight      util.AnimatedSheet
	shopButtonAnim util.AnimatedSheet
	focusedID      int
}

func (s *Scene) Boot() {
	s.highlight = assets.FireflyHighlight.Animated(2)
	s.shopButtonAnim = assets.ShopButton.Animated(3)
	s.modal.Boot()
}

func (s *Scene) Update() {
	if s.modal.IsOpen() {
		s.modal.Update()
		return
	}

	s.shopButtonAnim.Update()
	s.highlight.Update()
	for i := range s.fireflies {
		s.fireflies[i].Update()
	}
	// Sort by Y-axis so that they're drawn in the right order
	slices.SortFunc(s.fireflies, func(a, b Firefly) int {
		return cmp.Compare(a.pos.Y, b.pos.Y)
	})

	me := firefly.GetMe()
	if pad, ok := firefly.ReadPad(me); ok {
		dpad4 := pad.DPad4()
		justPressed := dpad4.JustPressed(s.dpad4Old)
		if justPressed != firefly.DPad4None {
			s.handleInputDPad4(justPressed)
		}
		s.dpad4Old = dpad4
	} else {
		s.dpad4Old = firefly.DPad4None
	}

	buttons := firefly.ReadButtons(me)
	if justPressed := buttons.JustPressed(s.buttonsOld); justPressed.Any() {
		s.handleInputButtons(justPressed)
	}
	s.buttonsOld = buttons
}

func (s *Scene) handleInputDPad4(justPressed firefly.DPad4) {
	if len(s.fireflies) == 0 {
		return
	}
	if s.focusedID == -1 {
		s.focusedID = s.fireflies[0].id
		return
	}
	ids := make([]int, len(s.fireflies))
	for i, f := range s.fireflies {
		ids[i] = f.id
	}
	slices.Sort(ids)
	idsIndex := slices.Index(ids, s.focusedID)

	switch justPressed {
	case firefly.DPad4Down, firefly.DPad4Right:
		idsIndex = (idsIndex + 1) % len(ids)
	case firefly.DPad4Up, firefly.DPad4Left:
		idsIndex = (idsIndex + len(ids) - 1) % len(ids)
	}

	newID := ids[idsIndex]
	s.focusedID = newID
}

func (s *Scene) handleInputButtons(justPressed firefly.Buttons) {
	switch {
	case justPressed.E:
		s.focusedID = -1

	case justPressed.S && s.focusedID != -1 && len(s.fireflies) > 0:
		idx := s.FindFireflyByID(s.focusedID)
		s.modal.Open(&s.fireflies[idx])

	case justPressed.N:
		scenes.SwitchScene(scenes.Shop)
	}
}

func (s *Scene) Render() {
	firefly.ClearScreen(firefly.ColorBlack)
	assets.Field.Draw(firefly.P(0, 0))

	for i := range s.fireflies {
		s.fireflies[i].Render()
	}

	if s.modal.IsOpen() {
		s.modal.Render()
		return
	}

	if len(s.fireflies) > 0 {
		if idx := s.FindFireflyByID(s.focusedID); idx != -1 {
			s.renderFocused(s.fireflies[idx])
		}
	}

	s.shopButtonAnim.Draw(firefly.P(firefly.Width-46, firefly.Height-16))
}

func (s *Scene) renderFocused(f Firefly) {
	pos := f.pos.Round().Point()
	s.highlight.Draw(pos.Sub(firefly.P(15, 15)))

	dataIndex := state.Game.FindFireflyByID(f.id)
	if dataIndex == -1 {
		panic("should never be -1 here")
	}
	data := state.Game.Fireflies[dataIndex]

	text := util.WordWrap(
		data.Name.String(),
		firefly.Width-75,
		assets.FontEG_6x9.CharWidth(),
	)

	assets.FontEG_6x9.Draw(text, firefly.P(73, 12), firefly.ColorDarkGray)
	assets.FontEG_6x9.Draw(text, firefly.P(73, 11), firefly.ColorWhite)
}

func (s *Scene) FindFireflyByID(id int) int {
	for idx := range s.fireflies {
		if s.fireflies[idx].id == id {
			return idx
		}
	}
	return -1
}

func (s *Scene) OnSceneSwitch() {
	s.focusedID = -1
	for _, f := range state.Game.Fireflies {
		idx := s.FindFireflyByID(f.ID)
		if idx == -1 {
			s.fireflies = append(s.fireflies, NewFirefly(f.ID))
		}
	}
	s.fireflies = slices.DeleteFunc(s.fireflies, func(f Firefly) bool {
		return state.Game.FindFireflyByID(f.id) == -1
	})
}
