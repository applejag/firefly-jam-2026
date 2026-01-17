package util

import (
	"cmp"
	"math"

	"github.com/firefly-zero/firefly-go/firefly"
	"github.com/orsinium-labs/tinymath"
)

const (
	Tau         = math.Pi * 2
	RadToDeg    = 360 / Tau
	DegToRad    = Tau / 360
	cmp_epsilon = 0.00001
)

// Returns the difference between the two angles, in the range of `[-PI, +PI]`.
// When `self` and `to` are opposite,
// returns `-PI` if `self` is smaller than `to`, or `PI` otherwise.
//
// Input angles do not need to be normalized.
// Based on the Godot `angle_difference` (licensed under MIT):
// https://github.com/godotengine/godot/blob/50277787eacaf4bc4d8683a706fe54dc65762020/core/math/math_funcs.h#L482-L489
func AngleDifference(from, to firefly.Angle) firefly.Angle {
	// tinymath has "RemEuclid" https://github.com/orsinium-labs/tinymath/blob/de812093edff2384fd94a922c5b255b0e39139a6/tinymath.go#L320-L328
	// but it has some math bugs, so we have to resort to the big math functions.
	diff := math.Mod(float64(to.Radians()-from.Radians()), Tau)
	return firefly.Radians(float32(math.Mod(2*diff, Tau) - diff))
}

func RotateTowards(from, to, delta firefly.Angle) firefly.Angle {
	diff := AngleDifference(from, to).Radians()
	abs_diff := tinymath.Abs(diff)
	return firefly.Radians(
		from.Radians() + Clamp(delta.Radians(), abs_diff-math.Pi, abs_diff)*tinymath.Sign(diff),
	)
}

func Clamp[T cmp.Ordered](val, minimum, maximum T) T {
	switch {
	case val < minimum:
		return minimum
	case val > maximum:
		return maximum
	default:
		return val
	}
}

func MoveTowards(start, end, delta float32) float32 {
	if tinymath.Abs(end-start) <= delta {
		return end
	} else {
		return start + tinymath.Sign(end-start)*delta
	}
}

func Lerp(from, to, weight float32) float32 {
	return from + (to-from)*weight
}

func RandomRange(start, end int) int {
	return start + int(firefly.GetRandom()%uint32(end-start))
}
