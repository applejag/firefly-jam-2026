package racebattle

import (
	"firefly-jam-2026/assets"
	"firefly-jam-2026/pkg/util"

	"github.com/firefly-zero/firefly-go/firefly"
)

const (
	ScreenWidthHalf  = 120
	ScreenHeightHalf = 80
)

type Camera struct {
	pos util.Vec2
}

func (c Camera) WorldVec2ToCameraSpace(pos util.Vec2) firefly.Point {
	return (pos.Sub(c.pos)).Round().Point()
}

func (c Camera) WorldPointToCameraSpace(pos firefly.Point) firefly.Point {
	rhs := c.pos.Round().Point()
	return firefly.P(pos.X-rhs.X, pos.Y-rhs.Y)
}

func (c *Camera) Update(world *World) {
	for _, player := range world.Players {
		if world.Me == player.Peer {
			c.pos = player.Pos.Sub(util.V(ScreenWidthHalf, ScreenHeightHalf))

			c.pos.X = util.Clamp(c.pos.X, 0, float32(assets.RacingMap.Width()-firefly.Width))
			c.pos.Y = util.Clamp(c.pos.Y, 0, float32(assets.RacingMap.Height()-firefly.Height))
			break
		}
	}
}
