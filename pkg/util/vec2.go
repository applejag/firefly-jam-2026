package util

import (
	"math"

	"github.com/firefly-zero/firefly-go/firefly"
	"github.com/orsinium-labs/tinymath"
)

type Vec2 struct {
	X float32
	Y float32
}

func V(x, y float32) Vec2 {
	return Vec2{X: x, Y: y}
}

func PointToVec2(point firefly.Point) Vec2 {
	return Vec2{X: float32(point.X), Y: float32(point.Y)}
}

func AngleToVec2(angle firefly.Angle) Vec2 {
	return Vec2{
		X: tinymath.Cos(angle.Radians()),
		Y: -tinymath.Sin(angle.Radians()),
	}
}

func (v Vec2) Point() firefly.Point {
	return firefly.Point{X: int(v.X), Y: int(v.Y)}
}

func (v Vec2) Round() Vec2 {
	return Vec2{X: tinymath.Round(v.X), Y: tinymath.Round(v.Y)}
}

func (v Vec2) Abs() Vec2 {
	return Vec2{X: tinymath.Abs(v.X), Y: tinymath.Abs(v.Y)}
}

func (v Vec2) Add(rhs Vec2) Vec2 {
	return Vec2{X: v.X + rhs.X, Y: v.Y + rhs.Y}
}

func (v Vec2) Sub(rhs Vec2) Vec2 {
	return Vec2{X: v.X - rhs.X, Y: v.Y - rhs.Y}
}

func (v Vec2) Negate() Vec2 {
	return Vec2{X: -v.X, Y: -v.Y}
}

func (v Vec2) Scale(factor float32) Vec2 {
	return Vec2{X: v.X * factor, Y: v.Y * factor}
}

func (v Vec2) Radius() float32 {
	return tinymath.Sqrt(v.X*v.X + v.Y*v.Y)
}

func (v Vec2) RadiusSquared() float32 {
	return v.X*v.X + v.Y*v.Y
}

func (v Vec2) Azimuth() firefly.Angle {
	r := math.Pi / 2. * tinymath.Atan2Norm(-v.Y, v.X)
	return firefly.Radians(r)
}
