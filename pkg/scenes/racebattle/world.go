package racebattle

import (
	"firefly-jam-2026/assets"
	"firefly-jam-2026/pkg/util"

	"github.com/firefly-zero/firefly-go/firefly"
)

var path = []util.Vec2{
	util.V(41, 406),
	util.V(56, 451),
	util.V(125, 504),
	util.V(232, 576),
	util.V(420, 556),
	util.V(431, 516),
	util.V(418, 440),
	util.V(550, 327),
	util.V(542, 257),
	util.V(512, 212),
	util.V(470, 195),
	util.V(426, 112),
	util.V(361, 123),
	util.V(242, 252),
	util.V(213, 233),
	util.V(235, 135),
	util.V(226, 82),
	util.V(154, 84),
	util.V(98, 119),
	util.V(76, 172),
	util.V(78, 250),
	util.V(72, 288),
	util.V(79, 319),
	util.V(47, 372),
}

type World struct {
	AnimatedClouds util.AnimatedSheet
	Players        []Firefly
	Camera         Camera
	Me             firefly.Peer
	// path Path
}

func (w *World) Boot() {
	w.AnimatedClouds = assets.RacingMapClouds.Animated(2)
	w.Me = firefly.GetMe()
	w.Players = []Firefly{
		NewFireflyPlayer(w.Me, util.V(41, 390), firefly.Degrees(270)),
	}
}

func (w *World) Update() {
	for i := range w.Players {
		w.Players[i].Update()
	}
	w.Camera.Update(w)
	w.AnimatedClouds.Update()
}

func (w *World) Render() {
	// Background
	firefly.ClearScreen(firefly.ColorDarkGray)
	mapPos := w.Camera.WorldPointToCameraSpace(firefly.P(0, 0))
	assets.RacingMap.Draw(mapPos)
	assets.RacingMapTrees.Draw(mapPos)
	// Players
	var me *Firefly
	for i, player := range w.Players {
		if player.Peer == w.Me {
			me = &w.Players[i]
		} else {
			player.Draw(w)
		}
	}
	// Draw my player last
	if me != nil {
		me.Draw(w)
	}
	// Draw tree tops layer on top
	assets.RacingMapTreetops.Draw(mapPos)
	w.AnimatedClouds.Draw(mapPos)
}
