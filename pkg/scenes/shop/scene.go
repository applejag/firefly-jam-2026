package shop

import "github.com/firefly-zero/firefly-go/firefly"

type Scene struct{}

func (s *Scene) Boot() {
}

func (s *Scene) Update() {
}

func (s *Scene) Render() {
	firefly.ClearScreen(firefly.ColorBlue)
}
