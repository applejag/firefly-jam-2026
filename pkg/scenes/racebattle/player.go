package racebattle

import (
	"firefly-jam-2026/assets"
	"firefly-jam-2026/pkg/util"
	"math"

	"github.com/firefly-zero/firefly-go/firefly"
	"github.com/orsinium-labs/tinymath"
)

const (
	FPS    = 60
	FPSInv = 1.0 / FPS

	FireflyAnimationFPS = 30

	// Rotation speed at top speed. Angle is the rotation speed per second
	RotationSpeedWhenMovingRad = (110.0 / FPS) * util.DegToRad
	// Rotation speed when standing completely still. Angle is the rotation speed per second
	RotationSpeedWhenStillRad = (540.0 / FPS) * util.DegToRad

	/// Maximum movement speed, in pixels/frame
	MoveMaxSpeed = 1.2

	/// Acceleration from 0 to max speed, in %/frame
	MoveAcceleration = FPSInv / 4.5 // 0%-100% in 4.5sec

	// Deacceleration (break force) to go from max speed to 0, in %/frame
	MoveDeacceleration = FPSInv / 1 // 100%-0% in 1sec
)

type Firefly struct {
	IsPlayer    bool
	Peer        firefly.Peer
	PathTracker PathTracker

	SpriteSheet    util.AnimatedSheet
	SpriteSheetRev util.AnimatedSheet

	Pos         util.Vec2
	Angle       firefly.Angle
	SpeedFactor float32
}

func NewFireflyPlayer(peer firefly.Peer, pos util.Vec2, angle firefly.Angle) Firefly {
	return Firefly{
		IsPlayer:       true,
		Peer:           peer,
		PathTracker:    NewPathTracker(path),
		SpriteSheet:    assets.FireflySheet.Animated(FireflyAnimationFPS),
		SpriteSheetRev: assets.FireflySheetRev.Animated(FireflyAnimationFPS),
		Pos:            pos,
		Angle:          angle,
	}
}

func NewFireflyAI(pos util.Vec2, angle firefly.Angle) Firefly {
	return Firefly{
		IsPlayer:       false,
		PathTracker:    NewPathTracker(path),
		SpriteSheet:    assets.FireflySheet.Animated(FireflyAnimationFPS),
		SpriteSheetRev: assets.FireflySheetRev.Animated(FireflyAnimationFPS),
		Pos:            pos,
		Angle:          angle,
	}
}

func (f *Firefly) Update() {
	f.SpriteSheet.Update()
	f.SpriteSheetRev.Update()
	if f.IsPlayer {
		f.UpdatePlayerInput()
	}
	dir := util.AngleToVec2(f.Angle)
	newPos := f.Pos.Add(dir.Scale(MoveMaxSpeed * f.SpeedFactor))
	maskColor := assets.RacingMapMask.GetColorAt(newPos.Point())
	if maskColor == firefly.ColorWhite {
		f.Pos = newPos
		f.PathTracker.Update(f.Pos)
	}
}

func (f *Firefly) UpdatePlayerInput() {
	pad, ok := firefly.ReadPad(f.Peer)
	if !ok {
		f.UpdateSpeedFactor(0)
		return
	}
	radiusPercentage := pad.Radius() / 1000
	// multiply by 1.2 and clamp so that there's some threshold to where
	// top speed is on the input
	targetSpeedFactor := min(radiusPercentage*1.2, 1.0)
	f.UpdateSpeedFactor(targetSpeedFactor)
	if pad.X != 0 && pad.Y != 0 {
		f.UpdateAngle(pad.Azimuth())
	}
}

func (f *Firefly) UpdateSpeedFactor(target float32) {
	if target > f.SpeedFactor {
		f.SpeedFactor = util.MoveTowards(f.SpeedFactor, target, MoveAcceleration)
	} else {
		f.SpeedFactor = util.MoveTowards(f.SpeedFactor, target, MoveDeacceleration)
	}
}

func (f *Firefly) UpdateAngle(target firefly.Angle) {
	rotationSpeed := util.Lerp(RotationSpeedWhenStillRad, RotationSpeedWhenMovingRad, f.SpeedFactor)
	f.Angle = util.RotateTowards(f.Angle, target, firefly.Radians(rotationSpeed))
}

func (f *Firefly) Draw(scene *Scene) {
	point := scene.Camera.WorldVec2ToCameraSpace(f.Pos)
	// Draw shadow
	firefly.DrawCircle(point.Add(firefly.P(-2, 2)), 5, firefly.Solid(firefly.ColorDarkGray))
	if f.IsPlayer && f.Peer == scene.Me {
		// Draw arrow to direction you should move in
		targetPoint := f.PathTracker.PeekNext().Add(f.PathTracker.PeekCurrent()).Scale(0.5)
		targetAngle := targetPoint.Sub(f.Pos).Azimuth()
		targetDir := util.AngleToVec2(targetAngle)
		drawArrow(
			point.Add(targetDir.Scale(10).Point()),
			targetAngle,
			7,
			firefly.L(firefly.ColorDarkGreen, 1))
	}
	// Draw sprite
	isLookingLeft := tinymath.Abs(util.AngleDifference(firefly.Radians(math.Pi), f.Angle).Radians()) < math.Pi/2
	spritePos := point.Sub(firefly.P(4, 5))
	if isLookingLeft {
		f.SpriteSheetRev.Draw(spritePos)
	} else {
		f.SpriteSheet.Draw(spritePos)
	}
}

func drawArrow(from firefly.Point, angle firefly.Angle, length int, lineStyle firefly.LineStyle) {
	fromV := util.PointToVec2(from)
	toV := fromV.Add(util.AngleToVec2(angle).Scale(float32(length)))
	to := toV.Point()
	firefly.DrawLine(
		from,
		to,
		lineStyle)
	firefly.DrawLine(
		to,
		toV.Add(util.AngleToVec2(angle.Add(firefly.Degrees(145))).Scale(3)).Point(),
		lineStyle)
	firefly.DrawLine(
		to,
		toV.Add(util.AngleToVec2(angle.Add(firefly.Degrees(-145))).Scale(3)).Point(),
		lineStyle)
}
