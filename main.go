package main

import (
	"firefly-jam-2026/assets"
	"firefly-jam-2026/pkg/game"
	"firefly-jam-2026/pkg/state"

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

	state.Input.Boot()
	// scenemanager.SwitchSceneNoTransition(scenes.Shop)
	scenemanager.Boot()
}

func update() {
	state.Input.Update()
	scenemanager.Update()
}

func render() {
	scenemanager.Render()
}
