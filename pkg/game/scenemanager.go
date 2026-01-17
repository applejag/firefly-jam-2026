package game

import (
	"firefly-jam-2026/pkg/scenes/insectarium"
	"firefly-jam-2026/pkg/scenes/mainmenu"
	"firefly-jam-2026/pkg/scenes/racebattle"
	"firefly-jam-2026/pkg/scenes/shop"
)

type SceneManager struct {
	CurrentScene Scene

	MainMenu    mainmenu.Menu
	RaceBattle  racebattle.Scene
	Insectarium insectarium.Scene
	Shop        shop.Scene
}

type Scene byte

const (
	SceneMainMenu Scene = iota
	SceneInsectarium
	SceneShop
	SceneRaceBattle
)

func (s *SceneManager) Boot() {
	s.MainMenu.Boot()
	s.RaceBattle.Boot()
	s.Insectarium.Boot()
	s.Shop.Boot()
}

func (s *SceneManager) Update() {
	switch s.CurrentScene {
	case SceneInsectarium:
		s.Insectarium.Update()
	case SceneMainMenu:
		s.MainMenu.Update()
	case SceneRaceBattle:
		s.RaceBattle.Update()
	case SceneShop:
		s.Shop.Update()
	}
}

func (s *SceneManager) Render() {
	switch s.CurrentScene {
	case SceneInsectarium:
		s.Insectarium.Render()
	case SceneMainMenu:
		s.MainMenu.Render()
	case SceneRaceBattle:
		s.RaceBattle.Render()
	case SceneShop:
		s.Shop.Render()
	}
}
