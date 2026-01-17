package main

import (
	"firefly-jam-2026/assets"
	"firefly-jam-2026/pkg/mainmenu"
	"firefly-jam-2026/pkg/scenes/racebattle"

	"github.com/firefly-zero/firefly-go/firefly"
)

func init() {
	firefly.Boot = boot
	firefly.Update = update
	firefly.Render = render
}

var world = racebattle.World{
	Camera: racebattle.Camera{},
}

var menu = mainmenu.Menu{}

func boot() {
	assets.Load()

	world.Boot()
	menu.Boot()
}

func update() {
	// world.Update()
	menu.Update()
}

func render() {
	// world.Render()
	menu.Render()
}
