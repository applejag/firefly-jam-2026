package racebattle

import (
	"cmp"
	"slices"

	"github.com/applejag/firefly-jam-2026/assets"
	"github.com/applejag/firefly-jam-2026/pkg/state"
	"github.com/applejag/firefly-jam-2026/pkg/util"

	"github.com/firefly-zero/firefly-go/firefly"
)

type Scene struct {
	AnimatedClouds util.AnimatedSheet
	Players        []Firefly
	Camera         Camera
}

func (s *Scene) Boot() {
	s.AnimatedClouds = assets.RacingMapClouds.Animated(2)
}

func (s *Scene) Update() {
	for i := range s.Players {
		s.Players[i].Update()
	}
	s.nudgeFirefliesAwayFromEachOther()
	s.Camera.Update(s)
	s.AnimatedClouds.Update()

	// Sort by Y-axis so that they're drawn in the right order
	slices.SortFunc(s.Players, func(a, b Firefly) int {
		return cmp.Compare(a.Pos.Y, b.Pos.Y)
	})
}

func (s *Scene) nudgeFirefliesAwayFromEachOther() {
	for i := 0; i < len(s.Players); i++ {
		for j := i + 1; j < len(s.Players); j++ {
			s.Players[i].MoveAwayFrom(&s.Players[j])
		}
	}
}

func (s *Scene) Render() {
	// Background
	firefly.ClearScreen(firefly.ColorDarkGray)
	mapPos := s.Camera.WorldPointToCameraSpace(firefly.P(0, 0))
	assets.RacingMap.Draw(mapPos)
	assets.RacingMapTrees.Draw(mapPos)
	// Players
	var me *Firefly
	for i, player := range s.Players {
		if player.IsPlayer && player.Peer == state.Input.Me {
			me = &s.Players[i]
		} else {
			player.Render(s)
		}
	}
	// Draw my player last
	if me != nil {
		me.Render(s)
	}
	// Draw tree tops layer on top
	assets.RacingMapTreetops.Draw(mapPos)
	s.AnimatedClouds.Draw(mapPos)
}

func (s *Scene) OnSceneEnter() {
	clear(s.Players)
	s.Players = s.Players[:0]
	for peer := range state.Game.InRaceBattle {
		s.Players = append(s.Players, NewFireflyPlayer(peer, util.V(41, 390).Add(offsetForPlayer(len(s.Players))), firefly.Degrees(271)))
	}
	for len(s.Players) < 2 {
		s.Players = append(s.Players, NewFireflyAI(util.V(41, 390).Add(offsetForPlayer(len(s.Players))), firefly.Degrees(271)))
	}
	s.Camera.Update(s)
}

func offsetForPlayer(index int) util.Vec2 {
	if index == 0 {
		return util.V(0, 0)
	}
	angle := firefly.Degrees(60 * float32(index-1))
	return util.AngleToVec2(angle).Scale(12)
}
