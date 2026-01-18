package game

import (
	"firefly-jam-2026/assets"
	"firefly-jam-2026/pkg/scenes"
	"firefly-jam-2026/pkg/scenes/field"
	"firefly-jam-2026/pkg/scenes/insectarium"
	"firefly-jam-2026/pkg/scenes/mainmenu"
	"firefly-jam-2026/pkg/scenes/racebattle"
	"firefly-jam-2026/pkg/scenes/shop"
	"firefly-jam-2026/pkg/util"
	"fmt"

	"github.com/firefly-zero/firefly-go/firefly"
)

type SceneManager struct {
	currentScene scenes.Scene
	nextScene    scenes.Scene
	Transition   util.Transition

	MainMenu    mainmenu.Menu
	RaceBattle  racebattle.Scene
	Insectarium insectarium.Scene
	Shop        shop.Scene
	Field       field.Scene

	UI UI
}

// SwitchScene implements [scenes.SceneSwitcher].
func (s *SceneManager) SwitchScene(scene scenes.Scene) {
	if s.currentScene == scene || s.nextScene == scene {
		return
	}
	s.nextScene = scene
	s.Transition.Play()
	firefly.LogDebug(fmt.Sprintf("switching scene from %q to %q", s.currentScene, scene))
	s.onSceneSwitch(scene)
}

func (s *SceneManager) SwitchSceneNoTransition(scene scenes.Scene) {
	if s.currentScene == scene {
		return
	}
	prev := s.currentScene
	s.nextScene = scene
	s.currentScene = scene
	s.Transition.Stop()
	firefly.LogDebug(fmt.Sprintf("switching scene from %q to %q (without transition)", prev, scene))
	s.onSceneSwitch(scene)
}

func (s *SceneManager) Boot() {
	s.nextScene = s.currentScene
	// register as global scene switcher
	scenes.SwitchScene = s.SwitchScene

	s.Transition = util.NewTransition(assets.TransitionSheet.Animated(12), firefly.S(8, 8))

	s.Insectarium.Boot()
	s.Field.Boot()
	s.MainMenu.Boot()
	s.RaceBattle.Boot()
	s.Shop.Boot()

	s.Shop.Shop.AddFireflyItem(0, 10, assets.ShopItem[6])
}

func (s *SceneManager) Update() {
	if s.nextScene == s.currentScene {
		// intentionally coded this way with enums because I want to avoid using
		// interfaces or function pointers for the heaviest functions like
		// Update and Render
		switch s.currentScene {
		case scenes.Insectarium:
			s.Insectarium.Update()
		case scenes.Field:
			s.Field.Update()
		case scenes.MainMenu:
			s.MainMenu.Update()
		case scenes.RaceBattle:
			s.RaceBattle.Update()
		case scenes.Shop:
			s.Shop.Update()
		}
	}
	s.Transition.Update()
	if s.currentScene != s.nextScene && s.Transition.IsPaused() {
		s.currentScene = s.nextScene
	}
}

func (s *SceneManager) Render() {
	scene := s.currentScene
	if s.Transition.IsPastHalf() {
		scene = s.nextScene
	}
	switch scene {
	case scenes.Insectarium:
		s.Insectarium.Render()
	case scenes.Field:
		s.Field.Render()
	case scenes.MainMenu:
		s.MainMenu.Render()
	case scenes.RaceBattle:
		s.RaceBattle.Render()
	case scenes.Shop:
		s.Shop.Render()
	}
	s.Transition.Draw()

	if scene != scenes.MainMenu {
		s.UI.Render()
	}
}

func (s *SceneManager) onSceneSwitch(scene scenes.Scene) {
	switch scene {
	case scenes.Field:
		s.Field.OnSceneEnter()
	case scenes.RaceBattle:
		s.RaceBattle.OnSceneEnter()
	}
}
