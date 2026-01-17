package util

import (
	"testing"

	"github.com/firefly-zero/firefly-go/firefly"
	"github.com/orsinium-labs/tinymath"
)

func TestAngleDifference(t *testing.T) {
	tests := []struct {
		name    string
		fromDeg float32
		toDeg   float32
		wantDeg float32
	}{
		{fromDeg: 0, toDeg: 0, wantDeg: 0},
		{fromDeg: 0, toDeg: 90, wantDeg: 90},
		{fromDeg: 45, toDeg: 90, wantDeg: 45},
		{fromDeg: 0, toDeg: 179, wantDeg: 179},
		{fromDeg: 179, toDeg: 0, wantDeg: -179},
		{fromDeg: 0, toDeg: -179, wantDeg: -179},
		{fromDeg: -179, toDeg: 0, wantDeg: 179},
		{fromDeg: 720, toDeg: 360, wantDeg: 0},
		{fromDeg: 700, toDeg: 650, wantDeg: -50},
	}

	for _, test := range tests {
		result := AngleDifference(firefly.Degrees(test.fromDeg), firefly.Degrees(test.toDeg))
		resultDeg := tinymath.Round(result.Degrees())
		if resultDeg != test.wantDeg {
			t.Errorf("AngleDifference(%f°, %f°)\nwant: %f°\ngot:  %f°", test.fromDeg, test.toDeg, test.wantDeg, resultDeg)
		}
	}
}
