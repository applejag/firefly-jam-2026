package shop

import (
	"firefly-jam-2026/assets"

	"github.com/firefly-zero/firefly-go/firefly"
)

type Scene struct{}

func (s *Scene) Boot() {
}

func (s *Scene) Update() {
}

func (s *Scene) Render() {
	assets.Shop[1].Draw(firefly.P(0, 0))

	// fits ~17 chars per line, and max 2 lines
	assets.FontEG_6x9.Draw("oy m8, u here 4\nsum foirefloies??", firefly.P(12, 38), firefly.ColorBlack)
}
