package main

import (
	"firefly-jam-2026/assets"
	"firefly-jam-2026/pkg/racebattle"
	"firefly-jam-2026/pkg/util"

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

func boot() {
	assets.Load()

	world.AnimatedClouds = assets.RacingMapClouds.Animated(2)
	world.Me = firefly.GetMe()
	world.Players = []racebattle.Firefly{
		racebattle.NewFireflyPlayer(world.Me, util.V(41, 390), firefly.Degrees(270)),
	}
}

func update() {
	world.Update()
}

func render() {
	world.Render()
}
