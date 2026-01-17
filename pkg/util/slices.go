package util

import "github.com/firefly-zero/firefly-go/firefly"

func RandomSliceElem[E any](slice []E) E {
	return slice[int(firefly.GetRandom()%uint32(len(slice)))]
}
