package racebattle

import (
	"github.com/applejag/firefly-jam-2026/pkg/util"
)

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
	util.V(40, 405),
}

type Path []util.Vec2

type PathTracker struct {
	path     Path
	index    int
	previous util.Vec2
	current  util.Vec2
	next     util.Vec2
}

func NewPathTracker(path Path) PathTracker {
	return PathTracker{
		path:     path,
		previous: path[len(path)-1],
		current:  path[0],
		next:     path[1],
	}
}

func (p PathTracker) PeekPrevious() util.Vec2 {
	return p.previous
}

func (p PathTracker) PeekCurrent() util.Vec2 {
	return p.current
}

func (p PathTracker) PeekNext() util.Vec2 {
	return p.next
}

func (p *PathTracker) goNext() {
	p.previous = p.current
	p.index = (p.index + 1) % len(p.path)
	p.current = p.path[p.index]
	p.next = p.path[(p.index+1)%len(p.path)]
}

func (p *PathTracker) PeekSoftNext(currentPos util.Vec2) util.Vec2 {
	currentTarget := p.PeekCurrent()
	nextTarget := p.PeekNext()
	prevTarget := p.PeekPrevious()

	// TODO: this implementation is a little buggy
	// it's not smooth at all when it switches checkpoints.
	// Maybe if we checked the distance to a point that's projected on a line
	// that's perpendicular with the prev->current line and that crosses the current target.
	distSqToCurrent := currentTarget.Sub(currentPos).RadiusSquared()
	distSqFromPrev := currentTarget.Sub(prevTarget).RadiusSquared()
	distWeight := 1 - min(distSqToCurrent/distSqFromPrev, 1)

	return util.V(
		util.Lerp(currentTarget.X, nextTarget.X, distWeight),
		util.Lerp(currentTarget.Y, nextTarget.Y, distWeight),
	)
}

// Progress returns the percentage (0.0-1.0) of progress made throughout the
// path. The "pos" is used to calculate fractional progress between checkpoints.
func (p *PathTracker) Progress(pos util.Vec2) float32 {
	prev := p.PeekPrevious()
	current := p.PeekCurrent()

	distSqToCurrent := current.Sub(pos).RadiusSquared()
	distSqFromPrev := current.Sub(prev).RadiusSquared()
	distWeight := 1 - min(distSqToCurrent/distSqFromPrev, 1)

	return float32(p.index+1)/float32(len(p.path)) + distWeight/float32(len(p.path))
}

func (p *PathTracker) Update(pos util.Vec2) PathTrackerResult {
	curr := p.PeekCurrent()
	prev := p.PeekPrevious()
	distSqToCurr := curr.Sub(pos).RadiusSquared()
	distSqToPrev := pos.Sub(prev).RadiusSquared()
	distSqBetweenPoints := curr.Sub(prev).RadiusSquared()
	switch {
	case distSqToPrev < distSqBetweenPoints:
		// haven't gotten far enough away from previous point
		return PathTrackerKeepCurrent
	case distSqToPrev < distSqToCurr:
		// moving backwards
		return PathTrackerMovingBackwards
	}
	p.goNext()
	if p.index == 0 {
		return PathTrackerLooped
	}
	return PathTrackerNextCheckpoint
}

type PathTrackerResult byte

const (
	PathTrackerKeepCurrent PathTrackerResult = iota
	PathTrackerMovingBackwards
	PathTrackerNextCheckpoint
	PathTrackerLooped
)
