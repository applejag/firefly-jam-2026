package racebattle

import (
	"math"

	"github.com/applejag/epic-wizard-firefly-gladiators/assets"
	"github.com/applejag/epic-wizard-firefly-gladiators/pkg/state"
	"github.com/applejag/epic-wizard-firefly-gladiators/pkg/util"
	"github.com/applejag/firefly-go-math/ffmath"
	"github.com/applejag/firefly-go-math/ffrand"

	"github.com/firefly-zero/firefly-go/firefly"
	"github.com/orsinium-labs/tinymath"
)

const (
	FPS    = 60
	FPSInv = 1.0 / FPS

	FireflyAnimationFPS = 10.0

	RotationSpeedFactorWhenStill = 5 // 5x faster rotation when still

	// Fireflies have stats like "SPEED: 12"
	// but top speed isn't 12 pixels per frame.
	// A decent movement speed is ~70pixels/s
	// Lower means slower, higher means faster
	StatsMoveSpeedFactor = 4.0 / FPS // div by FPS to make it "per second"

	// Fireflies have stats like "NIMBLE: 12"
	// but top rotation isn't 12 degrees per frame.
	// A decent rotation speed is ~120deg/s
	// This is later passed into firefly.Degrees(), so we're not talking radians here.
	// Lower means wider turns, higher is snappier turns.
	StatsRotationSpeedFactor = 8.0 / FPS // div by FPS to make it "per second"

	MoveAwayFromEachOtherSpeed     = 0.3
	MoveAwayFromEachOtherThreshold = 5.0

	/// Acceleration from 0 to max speed, in %/frame
	MoveAcceleration = FPSInv / 2.5 // 0%-100% in 2.5sec

	// Deacceleration (break force) to go from max speed to 0, in %/frame
	MoveDeacceleration = FPSInv / 1.0 // 100%-0% in 1sec
)

type Firefly struct {
	IsPlayer    bool
	Peer        firefly.Peer
	PathTracker PathTracker
	LoopsDone   int

	SpriteSheet    util.AnimatedSheet
	SpriteSheetRev util.AnimatedSheet

	Pos         ffmath.Vec
	Angle       firefly.Angle
	SpeedFactor float32

	MoveSpeed  float32
	Nimbleness float32
}

func NewFireflyPlayer(peer firefly.Peer, stats state.Firefly, pos ffmath.Vec, angle firefly.Angle) Firefly {
	return Firefly{
		IsPlayer:       true,
		Peer:           peer,
		PathTracker:    NewPathTracker(path),
		SpriteSheet:    assets.FireflySheet.Animated(FireflyAnimationFPS),
		SpriteSheetRev: assets.FireflySheetRev.Animated(FireflyAnimationFPS),
		Pos:            pos,
		Angle:          angle,
		MoveSpeed:      float32(stats.Speed),
		Nimbleness:     float32(stats.Nimbleness),
	}
}

func NewFireflyBot(pos ffmath.Vec, angle firefly.Angle) Firefly {
	// Randomize skills
	// These skills must be better than the starting score when buying a basic firefly
	speed := ffrand.IntRange(12, 18)
	nimbleness := 12 + (18 - speed)
	return Firefly{
		IsPlayer:       false,
		PathTracker:    NewPathTracker(path),
		SpriteSheet:    assets.FireflySheet.Animated(FireflyAnimationFPS),
		SpriteSheetRev: assets.FireflySheetRev.Animated(FireflyAnimationFPS),
		Pos:            pos,
		Angle:          angle,
		MoveSpeed:      float32(speed),
		Nimbleness:     float32(nimbleness),
	}
}

func (f *Firefly) Update() PathTrackerResult {
	f.SpriteSheet.Update()
	f.SpriteSheetRev.Update()
	if f.IsPlayer {
		f.updatePlayerInput()
	} else {
		f.updateBotInput()
	}
	dir := ffmath.VAngle(f.Angle)
	newPos := f.Pos.Add(dir.Scale(f.MoveSpeed * f.SpeedFactor * StatsMoveSpeedFactor))
	f.Move(newPos)
	trackerResult := f.PathTracker.Update(f.Pos)
	if trackerResult == PathTrackerLooped {
		// loop loop!
		f.LoopsDone++
	}
	return trackerResult
}

func (f *Firefly) Move(to ffmath.Vec) {
	delta := to.Sub(f.Pos)
	if delta.RadiusSquared() <= 0.01 {
		// no need to move
		return
	}
	for delta.RadiusSquared() > 0.01 {
		to = f.Pos.Add(delta)
		switch assets.RacingMapMask.GetColorAt(to.Point()) {
		case firefly.ColorWhite:
			f.Pos = to
			return
		default:
			// reduce distance, and try move again
			delta = delta.Scale(0.7)
		}
	}
	f.SpeedFactor = 0 // reset momentum when colliding into the wall
}

func (f *Firefly) MoveAwayFrom(other *Firefly) {
	delta := other.Pos.Sub(f.Pos)
	if delta.RadiusSquared() > (MoveAwayFromEachOtherThreshold * MoveAwayFromEachOtherThreshold) {
		// far enough away from each other
		return
	}

	f.Move(f.Pos.MoveTowards(f.Pos.Sub(delta), MoveAwayFromEachOtherSpeed))
	other.Move(other.Pos.MoveTowards(other.Pos.Add(delta), MoveAwayFromEachOtherSpeed))
}

func (f *Firefly) updatePlayerInput() {
	pad, ok := firefly.ReadPad(f.Peer)
	if !ok {
		f.updateSpeedFactor(0)
		return
	}
	radiusPercentage := pad.Radius() / 1000
	// multiply by 1.2 and clamp so that there's some threshold to where
	// top speed is on the input
	targetSpeedFactor := min(radiusPercentage*1.2, 1.0)
	f.updateSpeedFactor(targetSpeedFactor)
	if pad.X != 0 && pad.Y != 0 {
		f.updateAngle(pad.Azimuth())
	}
}

func (f *Firefly) updateBotInput() {
	// current := f.PathTracker.PeekCurrent()
	// next := f.PathTracker.PeekNext()
	// firefly.LogDebug(fmt.Sprintf("current: %s, next: %s", current, next))

	target := f.PathTracker.PeekSoftNext(f.Pos)
	delta := target.Sub(f.Pos)
	// TODO: slow down if it's a tight curve
	f.updateSpeedFactor(1)
	f.updateAngle(delta.Azimuth())
}

func (f *Firefly) updateSpeedFactor(target float32) {
	if target > f.SpeedFactor {
		f.SpeedFactor = ffmath.MoveTowards(f.SpeedFactor, target, MoveAcceleration)
	} else {
		f.SpeedFactor = ffmath.MoveTowards(f.SpeedFactor, target, MoveDeacceleration)
	}
}

func (f *Firefly) updateAngle(target firefly.Angle) {
	rotationSpeedDeg := ffmath.Lerp(
		f.Nimbleness*StatsRotationSpeedFactor*RotationSpeedFactorWhenStill,
		f.Nimbleness*StatsRotationSpeedFactor,
		f.SpeedFactor)
	f.Angle = ffmath.RotateTowards(f.Angle, target, firefly.Degrees(rotationSpeedDeg))
}

func (f *Firefly) Render(scene *Scene) {
	point := scene.Camera.WorldVec2ToCameraSpace(f.Pos)
	// Draw shadow
	firefly.DrawCircle(point.Add(firefly.P(-2, 2)), 5, firefly.Solid(firefly.ColorDarkGray))
	if f.IsPlayer && f.Peer == state.Input.Me {
		// Draw arrow to direction you should move in
		targetPoint := f.PathTracker.PeekSoftNext(f.Pos)
		targetAngle := targetPoint.Sub(f.Pos).Azimuth()
		targetDir := ffmath.VAngle(targetAngle)
		drawArrow(
			point.Add(targetDir.Scale(10).Point()),
			targetAngle,
			7,
			firefly.L(firefly.ColorDarkGreen, 1))
	}
	// Draw sprite
	isLookingLeft := tinymath.Abs(ffmath.AngleDifference(firefly.Radians(math.Pi), f.Angle).Radians()) < math.Pi/2
	spritePos := point.Sub(firefly.P(4, 5))
	if isLookingLeft {
		f.SpriteSheetRev.Draw(spritePos)
	} else {
		f.SpriteSheet.Draw(spritePos)
	}
}

func drawArrow(from firefly.Point, angle firefly.Angle, length int, lineStyle firefly.LineStyle) {
	fromV := ffmath.VPoint(from)
	toV := fromV.Add(ffmath.VAngle(angle).Scale(float32(length)))
	to := toV.Point()
	firefly.DrawLine(
		from,
		to,
		lineStyle)
	firefly.DrawLine(
		to,
		toV.Add(ffmath.VAngle(angle.Add(firefly.Degrees(145))).Scale(3)).Point(),
		lineStyle)
	firefly.DrawLine(
		to,
		toV.Add(ffmath.VAngle(angle.Add(firefly.Degrees(-145))).Scale(3)).Point(),
		lineStyle)
}

func (f *Firefly) Progress() float32 {
	return f.PathTracker.Progress(f.Pos)
}
