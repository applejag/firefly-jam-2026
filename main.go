package main

import (
	"github.com/applejag/epic-wizard-firefly-gladiators/assets"
	"github.com/applejag/epic-wizard-firefly-gladiators/pkg/game"
	"github.com/applejag/epic-wizard-firefly-gladiators/pkg/state"

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
	scenemanager.Boot()
}

func update() {
	state.Input.Update()
	scenemanager.Update()
}

func render() {
	scenemanager.Render()
}
