package assets

import (
	"firefly-jam-2026/pkg/util"
	"slices"

	"github.com/firefly-zero/firefly-go/firefly"
)

var (
	fieldBuf            [19213]byte
	Field               firefly.Image
	fireflyHighlightBuf [262]byte
	FireflyHighlight    util.SpriteSheet
	scrollBuf           [13895]byte
	ScrollClose         util.SpriteSheet
	ScrollOpen          util.SpriteSheet

	racingMapBuf         [192013]byte
	RacingMap            firefly.Image
	racingMapTreesBuf    [192013]byte
	RacingMapTrees       firefly.Image
	racingMapTreetopsBuf [192013]byte
	RacingMapTreetops    firefly.Image
	racingMapCloudsBuf   [32207]byte
	RacingMapClouds      util.SpriteSheet
	racingMapMaskBuf     [96007]byte
	RacingMapMask        util.ExtImage
	fireflySheetBuf      [333]byte
	FireflySheet         util.SpriteSheet
	fireflySheetRevBuf   [333]byte
	FireflySheetRev      util.SpriteSheet

	titleScreenBuf          [38413]byte
	TitleScreen             util.SpriteSheet
	titleButtonHighlightBuf [702]byte
	TitleButtonHighlight    util.SpriteSheet
	titleNoContinueBuf      [423]byte
	TitleNoContinue         firefly.Image

	shopBGBuf      [9607]byte
	ShopBG         firefly.Image
	shopFrogBuf    [31693]byte
	ShopFrog       util.SpriteSheet
	shopPropsBuf   [23533]byte
	ShopProps      util.SpriteSheet
	shopChatboxBuf [3073]byte
	ShopChatbox    firefly.Image
	shopItemBuf    [6061]byte
	ShopItem       util.SpriteSheet

	exitBuf [95]byte
	Exit    firefly.Image

	transitionSheetBuf [134]byte
	TransitionSheet    util.SpriteSheet

	cashBannerBuf [589]byte
	CashBanner    firefly.Image

	fontEG_6x9Buf    [655]byte
	FontEG_6x9       firefly.Font
	fontPico8_4x6Buf [295]byte
	FontPico8_4x6    firefly.Font
)

func Load() {
	// firefly.LogDebug(strconv.Itoa(firefly.GetFileSize("scroll")))
	Field = firefly.LoadImage("field", fieldBuf[:])
	FireflyHighlight = util.SplitImageBySize(firefly.LoadImage("firefly-hi", fireflyHighlightBuf[:]), firefly.S(32, 32))
	ScrollClose = util.SplitImageByCount(firefly.LoadImage("scroll", scrollBuf[:]), firefly.S(4, 1))
	ScrollOpen = slices.Clone(ScrollClose)
	slices.Reverse(ScrollOpen)
	RacingMap = firefly.LoadImage("racing-map", racingMapBuf[:])
	RacingMapTrees = firefly.LoadImage("racing-map-trees", racingMapTreesBuf[:])
	RacingMapTreetops = firefly.LoadImage("racing-map-treetops", racingMapTreetopsBuf[:])
	RacingMapClouds = util.SplitImageByCount(firefly.LoadImage("racing-map-clouds", racingMapCloudsBuf[:]), firefly.S(2, 1))
	RacingMapMask = util.NewExtImage(firefly.LoadFile("racing-map-mask", racingMapMaskBuf[:]))
	FireflySheet = util.SplitImageByCount(firefly.LoadImage("firefly", fireflySheetBuf[:]), firefly.S(7, 1))
	FireflySheetRev = util.SplitImageByCount(firefly.LoadImage("firefly-rev", fireflySheetRevBuf[:]), firefly.S(7, 1))
	TitleScreen = util.SplitImageByCount(firefly.LoadImage("title-screen", titleScreenBuf[:]), firefly.S(2, 1))
	TitleButtonHighlight = util.SplitImageByCount(firefly.LoadImage("title-button-hi", titleButtonHighlightBuf[:]), firefly.S(2, 1))
	TitleNoContinue = firefly.LoadImage("title-no-continue", titleNoContinueBuf[:])
	ShopBG = firefly.LoadImage("shop-bg", shopBGBuf[:])
	ShopFrog = util.SplitImageByCount(firefly.LoadImage("shop-frog", shopFrogBuf[:]), firefly.S(3, 2))
	ShopProps = util.SplitImageByCount(firefly.LoadImage("shop-props", shopPropsBuf[:]), firefly.S(3, 2))
	ShopChatbox = firefly.LoadImage("shop-chatbox", shopChatboxBuf[:])
	ShopItem = util.SplitImageByCount(firefly.LoadImage("shop-item", shopItemBuf[:]), firefly.S(4, 3))
	TransitionSheet = util.SplitImageByCount(firefly.LoadImage("transition", transitionSheetBuf[:]), firefly.S(4, 4))
	Exit = firefly.LoadImage("exit", exitBuf[:])
	CashBanner = firefly.LoadImage("cash-banner", cashBannerBuf[:])
	FontEG_6x9 = firefly.LoadFont("eg_6x9", fontEG_6x9Buf[:])
	FontPico8_4x6 = firefly.LoadFont("pico8_4x6", fontPico8_4x6Buf[:])
}
