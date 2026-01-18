package state

import (
	"firefly-jam-2026/pkg/util"
	gamev1 "firefly-jam-2026/proto/game/v1"
	"fmt"

	"github.com/firefly-zero/firefly-go/firefly"
)

var (
	nextID int
	Game   = GameState{
		InRaceBattle: map[firefly.Peer]Firefly{},
	}
)

type Firefly struct {
	ID         int
	Name       util.Name
	Speed      int
	Nimbleness int
}

type GameState struct {
	Fireflies    []Firefly
	InRaceBattle map[firefly.Peer]Firefly
}

func (g *GameState) AddFirefly() {
	nextID++
	name := util.RandomName()
	randomness := util.RandomRange(8, 14)
	g.Fireflies = append(g.Fireflies, Firefly{
		ID:         nextID,
		Name:       name,
		Speed:      randomness,
		Nimbleness: 8 + (14 - randomness),
	})
	g.Save()
}

func (g *GameState) FindFireflyByID(id int) int {
	for idx := range g.Fireflies {
		if g.Fireflies[idx].ID == id {
			return idx
		}
	}
	return -1
}

func (g *GameState) AddMyFireflyToRaceBattle(id int) {
	dataIndex := g.FindFireflyByID(id)
	if dataIndex == -1 {
		panic("should never be -1 here")
	}
	g.InRaceBattle[Input.Me] = g.Fireflies[dataIndex]
}

func (g *GameState) RemoveMyFireflyFromRaceBattle() {
	delete(g.InRaceBattle, Input.Me)
}

func (g *GameState) Save() {
	state := gamev1.GameState{
		Fireflies: make([]*gamev1.Firefly, len(g.Fireflies)),
	}
	for i, f := range g.Fireflies {
		state.Fireflies[i] = &gamev1.Firefly{
			Id:         int32(f.ID),
			Name:       int32(f.Name),
			Speed:      int32(f.Speed),
			Nimbleness: int32(f.Nimbleness),
		}
	}
	var buf [1000]byte
	written, err := state.MarshalToVT(buf[:])
	if err != nil {
		panic("failed to marshal save file")
	}
	firefly.DumpFile("save", buf[:written])
	firefly.LogDebug(fmt.Sprint("saved game, size: ", written))
}
