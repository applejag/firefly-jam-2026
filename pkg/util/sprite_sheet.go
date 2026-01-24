package util

import (
	"github.com/applejag/firefly-go-math/ffmath"
	"github.com/firefly-zero/firefly-go/firefly"
)

type SpriteSheet []firefly.SubImage

func SplitImageByCount(image firefly.Image, counts firefly.Size) SpriteSheet {
	imageSize := image.Size()
	spriteSize := firefly.S(imageSize.W/counts.W, imageSize.H/counts.H)
	return internalSplitImageBySize(image, spriteSize, imageSize)
}

func SplitImageBySize(image firefly.Image, spriteSize firefly.Size) SpriteSheet {
	return internalSplitImageBySize(image, spriteSize, image.Size())
}

func internalSplitImageBySize(image firefly.Image, spriteSize, imageSize firefly.Size) SpriteSheet {
	countPerAxis := firefly.S(imageSize.W/spriteSize.W, imageSize.H/spriteSize.H)
	count := countPerAxis.W * countPerAxis.H
	arr := make([]firefly.SubImage, count)
	for i := range count {
		point := firefly.P(
			i%countPerAxis.W*spriteSize.W,
			i/countPerAxis.W*spriteSize.H,
		)
		arr[i] = image.Sub(point, spriteSize)
	}
	return SpriteSheet(arr)
}

func (s SpriteSheet) Animated(fps int) AnimatedSheet {
	return NewAnimatedSheet(s, fps)
}

func NewAnimatedSheet(sprites SpriteSheet, fps int) AnimatedSheet {
	return AnimatedSheet{
		sprites:       sprites,
		ticksPerFrame: 60 / fps,
		AutoPlay:      true,
	}
}

type AnimatedSheet struct {
	sprites       SpriteSheet
	index         int
	time          int
	ticksPerFrame int
	AutoPlay      bool
}

func (s *AnimatedSheet) SetFrame(frame int) {
	s.index = ffmath.Clamp(frame, 0, len(s.sprites)-1)
	s.time = 0
}

func (s *AnimatedSheet) Update() {
	if s.IsPaused() {
		return
	}
	s.time++
	if s.time >= s.ticksPerFrame {
		before := s.index
		s.index = (s.index + 1) % len(s.sprites)
		s.time = 0
		if !s.AutoPlay && before > s.index {
			s.time = -1
		}
	}
}

func (s *AnimatedSheet) Draw(point firefly.Point) {
	if s.IsPaused() {
		return
	}
	s.sprites[s.index].Draw(point)
}

func (s *AnimatedSheet) DrawOrLastFrame(point firefly.Point) {
	if s.IsPaused() {
		s.sprites[len(s.sprites)-1].Draw(point)
	} else {
		s.Draw(point)
	}
}

func (s *AnimatedSheet) Play() {
	s.time = 0
	s.index = 0
}

func (s *AnimatedSheet) Stop() {
	s.time = -1
}

func (s *AnimatedSheet) IsPaused() bool {
	return s.time == -1
}
