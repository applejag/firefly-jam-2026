package game

import (
	"github.com/applejag/epic-wizard-firefly-gladiators/assets"
	"github.com/applejag/epic-wizard-firefly-gladiators/pkg/scenes"
	"github.com/applejag/epic-wizard-firefly-gladiators/pkg/scenes/field"
	"github.com/applejag/epic-wizard-firefly-gladiators/pkg/scenes/insectarium"
	"github.com/applejag/epic-wizard-firefly-gladiators/pkg/scenes/mainmenu"
	"github.com/applejag/epic-wizard-firefly-gladiators/pkg/scenes/racebattle"
	"github.com/applejag/epic-wizard-firefly-gladiators/pkg/scenes/shop"
	"github.com/applejag/epic-wizard-firefly-gladiators/pkg/util"

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
	prev := s.currentScene
	s.nextScene = scene
	s.Transition.Play()
	s.onSceneSwitch(scene)

	var buf [22 + scenes.LongestSceneName + 6 + scenes.LongestSceneName + 1]byte
	written := util.ConcatInto(buf[:],
		"switching scene from '", prev.String(), "' to '", scene.String(), "'",
	)
	util.LogDebugBytes(buf[:written])
}

func (s *SceneManager) SwitchSceneNoTransition(scene scenes.Scene) {
	prev := s.currentScene
	s.nextScene = scene
	s.currentScene = scene
	s.Transition.Stop()
	s.onSceneSwitch(scene)

	var buf [22 + scenes.LongestSceneName + 6 + scenes.LongestSceneName + 22]byte
	written := util.ConcatInto(buf[:],
		"switching scene from '", prev.String(), "' to '", scene.String(), "' (without transition)",
	)
	util.LogDebugBytes(buf[:written])
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
}

func (s *SceneManager) Update() {
	if s.nextScene == s.currentScene && s.Transition.IsPaused() {
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
		case scenes.RacingBattle, scenes.RacingTraining:
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
	case scenes.RacingBattle, scenes.RacingTraining:
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
	case scenes.RacingBattle:
		s.RaceBattle.OnSceneEnter(2)
	case scenes.RacingTraining:
		s.RaceBattle.OnSceneEnter(1)
	case scenes.Shop:
		s.Shop.OnSceneEnter()
	}
}
