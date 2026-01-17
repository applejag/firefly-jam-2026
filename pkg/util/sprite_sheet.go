package util

import "github.com/firefly-zero/firefly-go/firefly"

type SpriteSheet struct {
	sprites       []firefly.SubImage
	index         int
	time          int
	ticksPerFrame int
}

func SplitImageByCount(image firefly.Image, counts firefly.Size, fps int) SpriteSheet {
	imageSize := image.Size()
	spriteSize := firefly.S(imageSize.W/counts.W, imageSize.H/counts.H)
	return internalSplitImageBySize(image, spriteSize, fps, imageSize)
}

func SplitImageBySize(image firefly.Image, spriteSize firefly.Size, fps int) SpriteSheet {
	return internalSplitImageBySize(image, spriteSize, fps, image.Size())
}

func internalSplitImageBySize(image firefly.Image, spriteSize firefly.Size, fps int, imageSize firefly.Size) SpriteSheet {
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
	return SpriteSheet{
		sprites:       arr,
		ticksPerFrame: 60 / fps,
	}
}

func (s *SpriteSheet) Update() {
	s.time++
	if s.time >= s.ticksPerFrame {
		s.index = (s.index + 1) % len(s.sprites)
		s.time = 0
	}
}

func (s *SpriteSheet) Draw(point firefly.Point) {
	s.sprites[s.index].Draw(point)
}
