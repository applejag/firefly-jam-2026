package field

import (
	"github.com/applejag/firefly-jam-2026/assets"
	"github.com/applejag/firefly-jam-2026/pkg/util"
	"github.com/firefly-zero/firefly-go/firefly"
)

type TournamentPage struct {
	tournamentAnim util.AnimatedSheet
	trainingAnim   util.AnimatedSheet
}

func (p *TournamentPage) Boot() {
	p.tournamentAnim = assets.TournamentButton.Animated(6)
	p.trainingAnim = assets.TrainButton.Animated(2)
}

func (p *TournamentPage) Update() {
	p.tournamentAnim.Update()
	p.trainingAnim.Update()
}

func (p *TournamentPage) Render(innerScrollPoint firefly.Point) {
	p.trainingAnim.Draw(innerScrollPoint.Add(firefly.P(14, 20)))
	p.tournamentAnim.Draw(innerScrollPoint.Add(firefly.P(14, 40)))
}
