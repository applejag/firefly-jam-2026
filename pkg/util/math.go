package util

import (
	"cmp"
	"math"

	"github.com/firefly-zero/firefly-go/firefly"
	"github.com/orsinium-labs/tinymath"
)

const Tau = math.Pi * 2

// Returns the difference between the two angles, in the range of `[-PI, +PI]`.
// When `self` and `to` are opposite,
// returns `-PI` if `self` is smaller than `to`, or `PI` otherwise.
//
// Input angles do not need to be normalized.
// Based on the Godot `angle_difference` (licensed under MIT):
// https://github.com/godotengine/godot/blob/50277787eacaf4bc4d8683a706fe54dc65762020/core/math/math_funcs.h#L482-L489
func AngleDifference(from, to firefly.Angle) firefly.Angle {
	diff := tinymath.RemEuclid(to.Radians()-from.Radians(), Tau)
	return firefly.Radians(tinymath.RemEuclid(2*diff, Tau) - diff)
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
