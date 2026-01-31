package racebattle

import (
	"github.com/applejag/epic-wizard-firefly-gladiators/assets"
	"github.com/applejag/epic-wizard-firefly-gladiators/pkg/state"
	"github.com/applejag/firefly-go-math/ffmath"

	"github.com/firefly-zero/firefly-go/firefly"
)

const (
	ScreenWidthHalf  = 120
	ScreenHeightHalf = 80
)

type Camera struct {
	pos ffmath.Vec
}

func (c Camera) WorldVec2ToCameraSpace(pos ffmath.Vec) firefly.Point {
	return (pos.Sub(c.pos)).Round().Point()
}

func (c Camera) WorldPointToCameraSpace(pos firefly.Point) firefly.Point {
	rhs := c.pos.Round().Point()
	return firefly.P(pos.X-rhs.X, pos.Y-rhs.Y)
}

func (c *Camera) Update(scene *Scene) {
	for i := range scene.Players {
		player := &scene.Players[i]
		if player.IsPlayer && player.Peer == state.Input.Me {
			c.pos = player.Pos.Sub(ffmath.V(ScreenWidthHalf, ScreenHeightHalf))

			c.pos.X = ffmath.Clamp(c.pos.X, 0, float32(assets.RacingMap.Width()-firefly.Width))
			c.pos.Y = ffmath.Clamp(c.pos.Y, 0, float32(assets.RacingMap.Height()-firefly.Height))
			break
		}
	}
}
