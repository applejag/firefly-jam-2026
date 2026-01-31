package field

import (
	"cmp"
	"slices"

	"github.com/applejag/epic-wizard-firefly-gladiators/assets"
	"github.com/applejag/epic-wizard-firefly-gladiators/pkg/scenes"
	"github.com/applejag/epic-wizard-firefly-gladiators/pkg/state"
	"github.com/applejag/epic-wizard-firefly-gladiators/pkg/util"

	"github.com/firefly-zero/firefly-go/firefly"
)

type Scene struct {
	fireflies []Firefly
	modal     FireflyModal

	highlight      util.AnimatedSheet
	shopButtonAnim util.AnimatedSheet
	focusedID      int

	cachedFireflyNameText string
	fireflyIDs            []int
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

	if justPressed := state.Input.JustPressedDPad4(); justPressed != firefly.DPad4None {
		s.handleInputDPad4(justPressed)
	}
	if justPressed := state.Input.JustPressedButtons(); justPressed.Any() {
		s.handleInputButtons(justPressed)
	}
}

func (s *Scene) handleInputDPad4(justPressed firefly.DPad4) {
	if len(s.fireflies) == 0 {
		return
	}
	if s.focusedID == -1 {
		s.focusedID = s.fireflies[0].id
		return
	}
	idsIndex := slices.Index(s.fireflyIDs, s.focusedID)

	switch justPressed {
	case firefly.DPad4Down, firefly.DPad4Right:
		idsIndex = (idsIndex + 1) % len(s.fireflies)
	case firefly.DPad4Up, firefly.DPad4Left:
		idsIndex = (idsIndex + len(s.fireflies) - 1) % len(s.fireflies)
	}

	newID := s.fireflyIDs[idsIndex]
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
			s.renderFocused(&s.fireflies[idx])
		}
	}

	s.shopButtonAnim.Draw(firefly.P(firefly.Width-46, firefly.Height-16))
}

func (s *Scene) renderFocused(f *Firefly) {
	pos := f.pos.Round().Point()
	s.highlight.Draw(pos.Sub(firefly.P(15, 15)))

	dataIndex := state.Game.FindFireflyByID(f.id)
	if dataIndex == -1 {
		panic("should never be -1 here")
	}
	data := state.Game.Fireflies[dataIndex]

	text := s.cachedFireflyNameText
	if text == "" {
		var buf [util.LongestPossibleName]byte
		text = util.WordWrap(
			string(buf[:data.Name.WriteInto(buf[:])]),
			firefly.Width-75,
			assets.FontEG_6x9.CharWidth(),
		)
		s.cachedFireflyNameText = text
	}

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

func (s *Scene) OnSceneEnter() {
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
	s.cachedFireflyNameText = ""
	s.fireflyIDs = make([]int, len(s.fireflies))
	for i := range s.fireflies {
		s.fireflyIDs[i] = s.fireflies[i].id
	}
	slices.Sort(s.fireflyIDs)
}
