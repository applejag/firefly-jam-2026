package assets

import (
	"firefly-jam-2026/pkg/util"

	"github.com/firefly-zero/firefly-go/firefly"
)

var (
	fieldBuf                [19213]byte
	Field                   firefly.Image
	racingMapBuf            [192013]byte
	RacingMap               firefly.Image
	racingMapTreesBuf       [192013]byte
	RacingMapTrees          firefly.Image
	racingMapTreetopsBuf    [192013]byte
	RacingMapTreetops       firefly.Image
	racingMapCloudsBuf      [32207]byte
	RacingMapClouds         util.SpriteSheet
	racingMapMaskBuf        [96007]byte
	RacingMapMask           util.ExtImage
	fireflySheetBuf         [333]byte
	FireflySheet            util.SpriteSheet
	fireflySheetRevBuf      [333]byte
	FireflySheetRev         util.SpriteSheet
	titleScreenBuf          [38413]byte
	TitleScreen             util.SpriteSheet
	titleButtonHighlightBuf [702]byte
	TitleButtonHighlight    util.SpriteSheet
	transitionSheetBuf      [134]byte
	TransitionSheet         util.SpriteSheet
)

func Load() {
	// firefly.LogDebug(strconv.Itoa(firefly.GetFileSize("title-button-hi")))
	Field = firefly.LoadImage("field", fieldBuf[:])
	RacingMap = firefly.LoadImage("racing-map", racingMapBuf[:])
	RacingMapTrees = firefly.LoadImage("racing-map-trees", racingMapTreesBuf[:])
	RacingMapTreetops = firefly.LoadImage("racing-map-treetops", racingMapTreetopsBuf[:])
	RacingMapClouds = util.SplitImageByCount(firefly.LoadImage("racing-map-clouds", racingMapCloudsBuf[:]), firefly.S(2, 1))
	RacingMapMask = util.NewExtImage(firefly.LoadFile("racing-map-mask", racingMapMaskBuf[:]))
	FireflySheet = util.SplitImageByCount(firefly.LoadImage("firefly", fireflySheetBuf[:]), firefly.S(7, 1))
	FireflySheetRev = util.SplitImageByCount(firefly.LoadImage("firefly-rev", fireflySheetRevBuf[:]), firefly.S(7, 1))
	TitleScreen = util.SplitImageByCount(firefly.LoadImage("title-screen", titleScreenBuf[:]), firefly.S(2, 1))
	TitleButtonHighlight = util.SplitImageByCount(firefly.LoadImage("title-button-hi", titleButtonHighlightBuf[:]), firefly.S(2, 1))
	TransitionSheet = util.SplitImageByCount(firefly.LoadImage("transition", transitionSheetBuf[:]), firefly.S(4, 4))
}
