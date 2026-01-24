package racebattle

import (
	"github.com/applejag/epic-wizard-firefly-gladiators/pkg/state"
	"github.com/applejag/firefly-go-math/ffrand"
)

type Rewards struct {
	Money      int
	Speed      int
	Nimbleness int
}

func CalculateRewards(scene *Scene) Rewards {
	r := Rewards{}
	points := 1

	if len(scene.Players) > 1 {
		// only get money when you challenge other players
		r.Money = 15

		// against opponents  -> 2 extra points
		points += 2
	}

	// starter firefly does map in 40s-50s       -> 1 extra point
	// slightly upgraded firefly does map in 30s -> 2 extra points
	// high stats firefly does map in 20s        -> 4 extra points
	timeSeconds := float32(scene.ticksSinceStart) / FPS
	switch {
	case timeSeconds < 25:
		points += 4
	case timeSeconds < 35:
		points += 2
	case timeSeconds < 50:
		points += 1
	}
	// randomize the distribution
	r.Speed = ffrand.Intn(points + 1)
	r.Nimbleness = points - r.Speed
	return r
}

func (r Rewards) Apply(ff *state.Firefly) {
	state.Game.Money += r.Money
	if ff != nil {
		ff.Speed += r.Speed
		ff.Nimbleness += r.Nimbleness
	}
}
