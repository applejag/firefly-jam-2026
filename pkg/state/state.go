package state

import (
	"fmt"

	"github.com/applejag/firefly-jam-2026/pkg/util"
	gamev1 "github.com/applejag/firefly-jam-2026/proto/game/v1"

	"github.com/firefly-zero/firefly-go/firefly"
)

var (
	nextID int
	Game   = GameState{
		InRaceBattle: map[firefly.Peer]Firefly{},
	}
)

type Firefly struct {
	ID            int
	Name          util.Name
	Speed         int
	Nimbleness    int
	BattlesPlayed int
	BattlesWon    int
}

type GameState struct {
	Fireflies          []Firefly
	BattlesPlayedTotal int
	BattlesWonTotal    int
	Money              int

	// Not saved, it's only ephemeral data
	InRaceBattle map[firefly.Peer]Firefly
}

func (g *GameState) AddFirefly() int {
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
	return nextID
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
	state := gamev1.Save{
		Fireflies:          make([]*gamev1.Firefly, len(g.Fireflies)),
		BattlesWonTotal:    int32(g.BattlesWonTotal),
		BattlesPlayedTotal: int32(g.BattlesPlayedTotal),
		Money:              int32(g.Money),
	}
	for i, f := range g.Fireflies {
		state.Fireflies[i] = &gamev1.Firefly{
			Id:            int32(f.ID),
			Name:          int32(f.Name),
			Speed:         int32(f.Speed),
			Nimbleness:    int32(f.Nimbleness),
			BattlesWon:    int32(f.BattlesWon),
			BattlesPlayed: int32(f.BattlesPlayed),
		}
	}
	b, err := state.MarshalVT()
	if err != nil {
		panic("failed to marshal save file")
	}
	firefly.DumpFile("save", b)
	firefly.LogDebug(fmt.Sprintf("saved game, size: %d B", len(b)))
}

func (g *GameState) HasSave() bool {
	return firefly.FileExists("save")
}

func (g *GameState) LoadSave() bool {
	file := firefly.LoadFile("save", nil)
	if !file.Exists() {
		return false
	}
	var state gamev1.Save
	if err := state.UnmarshalVT(file.Raw); err != nil {
		firefly.LogError(fmt.Sprintf("failed to load save: %s", err))
		return false
	}

	g.Reset()
	g.BattlesPlayedTotal = int(state.BattlesPlayedTotal)
	g.BattlesWonTotal = int(state.BattlesWonTotal)
	g.Money = int(state.Money)
	g.Fireflies = make([]Firefly, len(state.Fireflies))
	for i, f := range state.Fireflies {
		g.Fireflies[i] = Firefly{
			ID:            int(f.Id),
			Name:          util.Name(f.Name),
			Speed:         int(f.Speed),
			Nimbleness:    int(f.Nimbleness),
			BattlesWon:    int(f.BattlesWon),
			BattlesPlayed: int(f.BattlesPlayed),
		}
	}

	firefly.LogDebug(fmt.Sprintf("loaded saved game, size: %d B", len(file.Raw)))
	return true
}

func (g *GameState) Reset() {
	*g = GameState{}
}
