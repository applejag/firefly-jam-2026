package field

import (
	"github.com/applejag/epic-wizard-firefly-gladiators/assets"
	"github.com/applejag/epic-wizard-firefly-gladiators/pkg/state"
	"github.com/applejag/epic-wizard-firefly-gladiators/pkg/util"

	"github.com/firefly-zero/firefly-go/firefly"
)

const (
	scrollInnerWidth  = 92
	scrollInnerHeight = 92
)

type ModalState byte

const (
	ModalClosed ModalState = iota
	ModalStats
	ModalHats
	ModalRacing
)

type FireflyModal struct {
	state           ModalState
	scrollOpenAnim  util.AnimatedSheet
	scrollCloseAnim util.AnimatedSheet
	scrollSprite    firefly.SubImage
	firefly         *Firefly

	racingPage RacingPage
	statsPage  StatsPage
}

func (m *FireflyModal) IsOpen() bool {
	return m.state != ModalClosed || !m.scrollCloseAnim.IsPaused()
}

func (m *FireflyModal) IsClosing() bool {
	return m.state == ModalClosed && !m.scrollCloseAnim.IsPaused()
}

func (m *FireflyModal) Open(firefly *Firefly) {
	m.scrollOpenAnim.Play()
	m.firefly = firefly
	m.OpenPage(ModalStats)
}

func (m *FireflyModal) OpenPage(target ModalState) {
	m.state = target
	switch target {
	case ModalStats:
		m.statsPage.OnOpen()
		m.statsPage.focused = StatsNone
	case ModalRacing:
		m.racingPage.OnOpen()
		m.racingPage.focused = RacingNone
	}
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
	m.statsPage.Boot()
	m.racingPage.Boot()
}

func (m *FireflyModal) Update() {
	m.scrollOpenAnim.Update()
	m.scrollCloseAnim.Update()

	if m.IsClosing() {
		return
	}

	switch m.state {
	case ModalStats:
		m.statsPage.Update(m)
	case ModalRacing:
		m.racingPage.Update(m)
	}

	if justPressed := state.Input.JustPressedButtons(); justPressed.Any() {
		switch {
		case justPressed.E:
			if m.state == ModalStats {
				m.Close()
			} else {
				m.OpenPage(ModalStats)
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

	innerScrollPoint := point.Add(firefly.P(21, 20))

	switch m.state {
	case ModalStats:
		m.statsPage.Render(innerScrollPoint, m.firefly.id)
	case ModalRacing:
		m.racingPage.Render(innerScrollPoint)
	}
}
