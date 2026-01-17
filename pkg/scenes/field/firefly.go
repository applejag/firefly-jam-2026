package field

import (
	"firefly-jam-2026/assets"
	"firefly-jam-2026/pkg/util"

	"github.com/firefly-zero/firefly-go/firefly"
)

const (
	FPS                       = 60
	FireflyAnimationFPS       = 15
	FireflySpeedPixelsPerTick = 0.2
)

type Firefly struct {
	id         int
	sprites    util.AnimatedSheet
	spritesRev util.AnimatedSheet

	pos        util.Vec2
	nextPos    util.Vec2
	sleepTicks int
}

func NewFirefly(id int) Firefly {
	sprites := assets.FireflySheet.Animated(FireflyAnimationFPS)
	spritesRev := assets.FireflySheetRev.Animated(FireflyAnimationFPS)
	frame := util.RandomRange(0, len(assets.FireflySheet))
	sprites.SetFrame(frame)
	spritesRev.SetFrame(frame)
	return Firefly{
		id: id,
		pos: util.V(
			float32(util.RandomRange(40, firefly.Width-40)),
			float32(util.RandomRange(30, firefly.Height-30)),
		),
		sleepTicks: util.RandomRange(FPS*1, FPS*3),
		sprites:    sprites,
		spritesRev: spritesRev,
	}
}

func (f *Firefly) Update() {
	f.sprites.Update()
	f.spritesRev.Update()

	if f.sleepTicks > 0 {
		// stay idle for a bit
		f.sleepTicks--
		if f.sleepTicks <= 0 {
			f.nextPos = util.V(
				float32(util.RandomRange(40, firefly.Width-40)),
				float32(util.RandomRange(30, firefly.Height-30)),
			)
		}
	} else {
		// move to nextPos
		f.pos = f.pos.MoveTowards(f.nextPos, FireflySpeedPixelsPerTick)

		if f.nextPos.Sub(f.pos).RadiusSquared() < 10 {
			f.sleepTicks = util.RandomRange(FPS*1, FPS*3)
		}
	}
}

func (f *Firefly) Render() {
	f.sprites.Draw(f.pos.Round().Point())
}
