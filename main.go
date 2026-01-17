package main

import (
	"firefly-jam-2026/assets"
	"firefly-jam-2026/pkg/game"

	"github.com/firefly-zero/firefly-go/firefly"
)

func init() {
	firefly.Boot = boot
	firefly.Update = update
	firefly.Render = render
}

var scenemanager = game.SceneManager{}

func boot() {
	assets.Load()
	scenemanager.Boot()
}

func update() {
	scenemanager.Update()
}

func render() {
	scenemanager.Render()
}
