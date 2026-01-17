package assets

import (
	"firefly-jam-2026/pkg/util"

	"github.com/firefly-zero/firefly-go/firefly"
)

var (
	fieldBuf             [19213]byte
	Field                firefly.Image
	racingMapBuf         [192013]byte
	RacingMap            firefly.Image
	racingMapTreesBuf    [192013]byte
	RacingMapTrees       firefly.Image
	racingMapTreetopsBuf [192013]byte
	RacingMapTreetops    firefly.Image
	fireflySheetBuf      [333]byte
	FireflySheet         util.SpriteSheet
	fireflySheetRevBuf   [333]byte
	FireflySheetRev      util.SpriteSheet
)

func Load() {
	Field = firefly.LoadImage("field", fieldBuf[:])
	RacingMap = firefly.LoadImage("racing-map", racingMapBuf[:])
	RacingMapTrees = firefly.LoadImage("racing-map-trees", racingMapTreesBuf[:])
	RacingMapTreetops = firefly.LoadImage("racing-map-treetops", racingMapTreetopsBuf[:])
	FireflySheet = util.SplitImageByCount(firefly.LoadImage("firefly", fireflySheetBuf[:]), firefly.S(7, 1))
	FireflySheetRev = util.SplitImageByCount(firefly.LoadImage("firefly-rev", fireflySheetRevBuf[:]), firefly.S(7, 1))
}
