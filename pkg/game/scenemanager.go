package game

import (
	"firefly-jam-2026/assets"
	"firefly-jam-2026/pkg/scenes"
	"firefly-jam-2026/pkg/scenes/insectarium"
	"firefly-jam-2026/pkg/scenes/mainmenu"
	"firefly-jam-2026/pkg/scenes/racebattle"
	"firefly-jam-2026/pkg/scenes/shop"
	"firefly-jam-2026/pkg/util"
	"fmt"

	"github.com/firefly-zero/firefly-go/firefly"
)

type SceneManager struct {
	CurrentScene scenes.Scene
	NextScene    scenes.Scene
	Transition   util.Transition

	MainMenu    mainmenu.Menu
	RaceBattle  racebattle.Scene
	Insectarium insectarium.Scene
	Shop        shop.Scene
}

// SwitchScene implements [scenes.SceneSwitcher].
func (s *SceneManager) SwitchScene(scene scenes.Scene) {
	s.NextScene = scene
	s.Transition.Play()
	firefly.LogDebug(fmt.Sprintf("switching scene to %q", scene))
}

func (s *SceneManager) Boot() {
	// register as global scene switcher
	scenes.SwitchScene = s.SwitchScene

	s.Transition = util.NewTransition(assets.TransitionSheet.Animated(12), firefly.S(8, 8))

	s.MainMenu.Boot()
	s.RaceBattle.Boot()
	s.Insectarium.Boot()
	s.Shop.Boot()
}

func (s *SceneManager) Update() {
	if s.NextScene == s.CurrentScene {
		switch s.CurrentScene {
		case scenes.Insectarium:
			s.Insectarium.Update()
		case scenes.MainMenu:
			s.MainMenu.Update()
		case scenes.RaceBattle:
			s.RaceBattle.Update()
		case scenes.Shop:
			s.Shop.Update()
		}
	}
	s.Transition.Update()
	if s.Transition.IsPaused() {
		s.CurrentScene = s.NextScene
	}
}

func (s *SceneManager) Render() {
	scene := s.CurrentScene
	if s.Transition.IsPastHalf() {
		scene = s.NextScene
	}
	switch scene {
	case scenes.Insectarium:
		s.Insectarium.Render()
	case scenes.MainMenu:
		s.MainMenu.Render()
	case scenes.RaceBattle:
		s.RaceBattle.Render()
	case scenes.Shop:
		s.Shop.Render()
	}
	s.Transition.Draw()
}
