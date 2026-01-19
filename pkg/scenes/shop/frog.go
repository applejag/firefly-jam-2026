package shop

import (
	"github.com/applejag/epic-wizard-firefly-gladiators/assets"
	"github.com/applejag/epic-wizard-firefly-gladiators/pkg/util"

	"github.com/firefly-zero/firefly-go/firefly"
)

const FPS = 60

const (
	chatWobbleTicks   = 40
	frogAnimRandomMin = FPS * 2
	frogAnimRandomMax = FPS * 8
)

type Frog struct {
	animatedFrog util.AnimatedSheet
	frogAnimWait int

	animatedProps util.AnimatedSheet

	chatWobbleTime   int
	chatWobbleOffset int
}

func (f *Frog) Boot() {
	f.animatedFrog = assets.ShopFrog.Animated(3)
	f.animatedFrog.AutoPlay = false
	f.animatedFrog.Stop()
	f.resetRandomWait()

	f.animatedProps = assets.ShopProps.Animated(8)
}

func (f *Frog) Update() {
	if f.animatedFrog.IsPaused() && f.frogAnimWait > 0 {
		f.frogAnimWait--
		if f.frogAnimWait <= 0 {
			f.animatedFrog.Play()
			f.resetRandomWait()
		}
	}
	f.animatedFrog.Update()

	f.animatedProps.Update()

	f.chatWobbleTime++
	if f.chatWobbleTime >= chatWobbleTicks {
		f.chatWobbleTime -= chatWobbleTicks
		f.chatWobbleOffset = 1 - f.chatWobbleOffset
	}
}

func (f *Frog) resetRandomWait() {
	f.frogAnimWait = frogAnimRandomMin + int(firefly.GetRandom()%(frogAnimRandomMax-frogAnimRandomMin))
}

func (f *Frog) Render() {
	if f.animatedFrog.IsPaused() {
		assets.ShopFrog[0].Draw(firefly.P(0, 22))
	} else {
		f.animatedFrog.Draw(firefly.P(0, 22))
	}
	f.animatedProps.Draw(firefly.P(0, 79))

	assets.ShopChatbox.Draw(firefly.P(4, 22+f.chatWobbleOffset))
	// fits ~17 chars per line, and max 2 lines
	assets.FontEG_6x9.Draw("oy m8, u here 4\nsum foirefloies??", firefly.P(12, 38), firefly.ColorBlack)
}
