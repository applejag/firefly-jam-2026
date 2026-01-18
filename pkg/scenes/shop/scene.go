package shop

import (
	"github.com/applejag/firefly-jam-2026/assets"
	"github.com/applejag/firefly-jam-2026/pkg/state"

	"github.com/firefly-zero/firefly-go/firefly"
)

type Scene struct {
	Shop Shop
	Frog Frog
}

func (s *Scene) Boot() {
	s.Shop.Boot()
	s.Frog.Boot()
}

func (s *Scene) Update() {
	s.Shop.Update()
	s.Frog.Update()
}

func (s *Scene) Render() {
	assets.ShopBG.Draw(firefly.P(0, 0))

	s.Shop.Render()
	s.Frog.Render()
}

func (s Scene) OnSceneEnter() {
	if state.Game.BattlesPlayedTotal == 0 && len(state.Game.Fireflies) == 0 && len(s.Shop.Items) == 0 {
		// brand new player, they get a free firefly
		s.Shop.AddFireflyItem(0, 1, assets.ShopItem[6])
		return
	}

	// sync shop items
	if len(s.Shop.Items) == 0 {
		s.Shop.AddFireflyItem(100, -1, assets.ShopItem[6])
	}
}
