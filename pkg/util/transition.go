package util

import (
	"github.com/firefly-zero/firefly-go/firefly"
)

type Transition struct {
	AnimatedSheet
	size firefly.Size
}

func NewTransition(sprites AnimatedSheet, size firefly.Size) Transition {
	sprites.AutoPlay = false
	sprites.Stop()
	return Transition{AnimatedSheet: sprites, size: size}
}

func (t *Transition) Draw() {
	// tile the sprite
	for x := 0; x < firefly.Width; x += t.size.W {
		for y := 0; y < firefly.Height; y += t.size.H {
			t.AnimatedSheet.Draw(firefly.P(x, y))
		}
	}
}

func (t *Transition) IsPastHalf() bool {
	return t.index >= len(t.sprites)/2
}
