package racebattle

import "firefly-jam-2026/pkg/util"

const PathToCurrentThresholdSquared = 60 * 60

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
