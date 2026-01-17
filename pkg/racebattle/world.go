package racebattle

import (
	"firefly-jam-2026/assets"
	"firefly-jam-2026/pkg/util"

	"github.com/firefly-zero/firefly-go/firefly"
)

var path = []util.Vec2{
	util.V(68, 331),
	// util.V(),
	// util.V(),
	// util.V(),
	// util.V(),
	// util.V(),
	// util.V(),
	// util.V(),
	// util.V(),
	// util.V(),
	// util.V(),
	// util.V(),
	// util.V(),
	// util.V(),
	// util.V(),
	// util.V(),
	// util.V(),
	// util.V(),
	// util.V(),
	// util.V(),
	// util.V(),
	// util.V(),
	// util.V(),
	// util.V(),
	// util.V(),
	// util.V(),
	// util.V(),
	// util.V(),
}

type World struct {
	Players []Firefly
	Camera  Camera
	Me      firefly.Peer
	// path Path
}

func (w *World) Update() {
	for i := range w.Players {
		w.Players[i].Update()
	}
	w.Camera.Update(w)
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
}
