package racebattle

import (
	"github.com/applejag/firefly-jam-2026/pkg/util"
)

const PathToCurrentThresholdSquared = 60 * 60

var path = Path{
	util.V(40, 414),
	util.V(42, 434),
	util.V(56, 451),
	util.V(125, 504),
	util.V(194, 573),
	util.V(257, 588),
	util.V(337, 583),
	util.V(428, 575),
	util.V(438, 541),
	util.V(434, 513),
	util.V(431, 516),
	util.V(405, 456),
	util.V(420, 415),
	util.V(465, 410),
	util.V(532, 365),
	util.V(566, 294),
	util.V(547, 247),
	util.V(512, 212),
	util.V(470, 195),
	util.V(438, 110),
	util.V(395, 90),
	util.V(349, 114),
	util.V(321, 192),
	util.V(288, 195),
	util.V(284, 241),
	util.V(245, 268),
	util.V(207, 261),
	util.V(198, 230),
	util.V(215, 183),
	util.V(243, 150),
	util.V(245, 91),
	util.V(234, 69),
	util.V(188, 62),
	util.V(143, 68),
	util.V(115, 84),
	util.V(101, 122),
	util.V(69, 138),
	util.V(57, 165),
	util.V(56, 199),
	util.V(85, 232),
	util.V(82, 271),
	util.V(67, 292),
	util.V(74, 312),
	util.V(82, 331),
	util.V(75, 350),
	util.V(41, 361),
	util.V(36, 380),
}

type Path []util.Vec2

type PathTracker struct {
	path  Path
	index int
}

func NewPathTracker(path Path) PathTracker {
	return PathTracker{path: path}
}

func (p PathTracker) PeekPrevious() util.Vec2 {
	return p.path[(p.index-1+len(p.path))%len(p.path)]
}

func (p PathTracker) PeekCurrent() util.Vec2 {
	return p.path[p.index]
}

func (p PathTracker) PeekNext() util.Vec2 {
	return p.path[(p.index+1)%len(p.path)]
}

func (p *PathTracker) GoNext() {
	p.index = (p.index + 1) % len(p.path)
}

func (p *PathTracker) PeekSoftNext(from util.Vec2) util.Vec2 {
	current := p.PeekCurrent()
	next := p.PeekNext()
	prev := p.PeekPrevious()

	// TODO: this implementation is a little buggy
	// it's not smooth at all when it switches checkpoints.
	// Maybe if we checked the distance to a point that's projected on a line
	// that's perpendicular with the prev->current line and that crosses the current target.
	distSqToCurrent := current.Sub(from).RadiusSquared()
	distSqFromPrev := current.Sub(prev).RadiusSquared()
	distWeight := 1 - min(distSqToCurrent/distSqFromPrev, 1)

	return util.V(
		util.Lerp(current.X, next.X, distWeight),
		util.Lerp(current.Y, next.Y, distWeight),
	)
}

func (p *PathTracker) Update(pos util.Vec2) {
	curr := p.PeekCurrent()
	distSquaredToCurr := curr.Sub(pos).RadiusSquared()
	if distSquaredToCurr >= PathToCurrentThresholdSquared {
		return // keep current
	}
	prev := p.PeekPrevious()
	distSquaredToPrev := pos.Sub(prev).RadiusSquared()
	distSquaredBetweenPoints := curr.Sub(prev).RadiusSquared()
	if distSquaredToPrev < distSquaredBetweenPoints {
		return // keep current
	}
	p.GoNext()
}
