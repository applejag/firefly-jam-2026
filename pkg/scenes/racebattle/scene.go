package racebattle

import (
	"cmp"
	"slices"

	"github.com/applejag/firefly-jam-2026/assets"
	"github.com/applejag/firefly-jam-2026/pkg/state"
	"github.com/applejag/firefly-jam-2026/pkg/util"

	"github.com/firefly-zero/firefly-go/firefly"
)

type GameStatus byte

const (
	GameStarting GameStatus = iota
	GamePlaying
	GameOverDefeat
	GameOverVictory
)

type Scene struct {
	AnimatedClouds util.AnimatedSheet
	VictorySplash  util.AnimatedSheet
	DefeatSplash   util.AnimatedSheet

	Players []Firefly
	Camera  Camera
	Status  GameStatus

	countdownNum  int
	countdownTime int

	// Placement among all competitors.
	//
	// - 1 means 1st place
	// - 2 means 2nd place
	// - 3 means 3rd place
	myPlayerPlace byte
}

func (s *Scene) Boot() {
	s.AnimatedClouds = assets.RacingMapClouds.Animated(2)
	s.VictorySplash = assets.VictorySplash.Animated(12)
	s.VictorySplash.AutoPlay = false
	s.DefeatSplash.Stop()
	s.DefeatSplash = assets.DefeatSplash.Animated(12)
	s.DefeatSplash.AutoPlay = false
	s.DefeatSplash.Stop()
	s.Status = GameStarting
}

func (s *Scene) Update() {
	switch s.Status {
	case GamePlaying:
		for i := range s.Players {
			result := s.Players[i].Update()
			if result == PathTrackerLooped {
				isMyPlayer := s.Players[i].IsPlayer && s.Players[i].Peer == state.Input.Me
				if isMyPlayer {
					s.Status = GameOverVictory
					s.VictorySplash.Play()
				} else {
					s.Status = GameOverDefeat
					s.DefeatSplash.Play()
				}
			}
		}
		s.nudgeFirefliesAwayFromEachOther()
		s.updateMyPlayerPlace()
		// Sort by Y-axis so that they're drawn in the right order
		slices.SortFunc(s.Players, func(a, b Firefly) int {
			return cmp.Compare(a.Pos.Y, b.Pos.Y)
		})

		s.Camera.Update(s)
		s.AnimatedClouds.Update()

	case GameStarting:
		// TODO: we should start in "GameStarting", and do a countdown
		s.countdownTime--
		if s.countdownTime <= 0 {
			s.countdownNum--
			if s.countdownNum <= 0 {
				s.Status = GamePlaying
			} else {
				s.countdownTime = 60
			}
		}

	case GameOverVictory:
		s.VictorySplash.Update()
	case GameOverDefeat:
		s.DefeatSplash.Update()
	}
}

func (s *Scene) nudgeFirefliesAwayFromEachOther() {
	for i := 0; i < len(s.Players); i++ {
		for j := i + 1; j < len(s.Players); j++ {
			s.Players[i].MoveAwayFrom(&s.Players[j])
		}
	}
}

func (s *Scene) updateMyPlayerPlace() {
	var myProgress float32
	for _, player := range s.Players {
		if player.IsPlayer && player.Peer == state.Input.Me {
			myProgress = player.Progress() + float32(player.LoopsDone)
			break
		}
	}

	var playersWithHigerProgress byte
	for _, player := range s.Players {
		if player.IsPlayer && player.Peer == state.Input.Me {
			continue
		}
		if player.Progress()+float32(player.LoopsDone) > myProgress {
			playersWithHigerProgress++
		}
	}

	s.myPlayerPlace = playersWithHigerProgress + 1
}

func (s *Scene) Render() {
	// Background
	firefly.ClearScreen(firefly.ColorDarkGray)
	mapPos := s.Camera.WorldPointToCameraSpace(firefly.P(0, 0))
	assets.RacingMap.Draw(mapPos)
	assets.RacingMapTrees.Draw(mapPos)
	// Players
	var me *Firefly
	for i, player := range s.Players {
		if player.IsPlayer && player.Peer == state.Input.Me {
			me = &s.Players[i]
		} else {
			player.Render(s)
		}
	}
	// Draw my player last
	if me != nil {
		me.Render(s)
	}
	// Draw tree tops layer on top
	assets.RacingMapTreetops.Draw(mapPos)
	s.AnimatedClouds.Draw(mapPos)

	switch s.Status {
	case GameStarting:
		center := firefly.P(firefly.Width/2, firefly.Height/2)
		boxSize := firefly.S(34, 14)
		boxPos := firefly.P(center.X-boxSize.W/2, center.Y-boxSize.H/2-2)
		switch s.countdownNum {
		case 3:
			firefly.DrawRoundedRect(boxPos, boxSize, firefly.S(3, 3), firefly.Solid(firefly.ColorBlack))
			util.DrawTextCentered(assets.FontEG_6x9, "[ 3 ]", center.Add(firefly.P(1, 1)), firefly.ColorDarkGray)
			util.DrawTextCentered(assets.FontEG_6x9, "[ 3 ]", center, firefly.ColorLightGray)
		case 2:
			firefly.DrawRoundedRect(boxPos, boxSize, firefly.S(3, 3), firefly.Solid(firefly.ColorBlack))
			util.DrawTextCentered(assets.FontEG_6x9, "[ 2 ]", center.Add(firefly.P(1, 1)), firefly.ColorGray)
			util.DrawTextCentered(assets.FontEG_6x9, "[ 2 ]", center, firefly.ColorWhite)
		case 1:
			firefly.DrawRoundedRect(boxPos, boxSize, firefly.S(3, 3), firefly.Solid(firefly.ColorBlack))
			util.DrawTextCentered(assets.FontEG_6x9, "[ 1 ]", center.Add(firefly.P(1, 1)), firefly.ColorGray)
			util.DrawTextCentered(assets.FontEG_6x9, "[ 1 ]", center, firefly.ColorYellow)
		}

	case GamePlaying:
		if s.myPlayerPlace >= 1 && s.myPlayerPlace <= 3 {
			assets.RacingPlace[s.myPlayerPlace-1].Draw(firefly.P(firefly.Width-28-4, 4))
		}

	case GameOverVictory:
		s.VictorySplash.DrawOrLastFrame(firefly.P(0, 0))
	case GameOverDefeat:
		s.DefeatSplash.DrawOrLastFrame(firefly.P(0, 0))
	}
}

func (s *Scene) OnSceneEnter() {
	clear(s.Players)
	s.Players = s.Players[:0]
	for peer := range state.Game.InRaceBattle {
		s.Players = append(s.Players, NewFireflyPlayer(peer, util.V(41, 390).Add(offsetForPlayer(len(s.Players))), firefly.Degrees(271)))
	}
	if len(s.Players) < 2 {
		s.Players = append(s.Players, NewFireflyBot(util.V(41, 390).Add(offsetForPlayer(len(s.Players))), firefly.Degrees(271)))
	}
	// Update once so it focuses on player when we're transitioning to this scene
	s.Camera.Update(s)

	s.VictorySplash.Stop()
	s.DefeatSplash.Stop()
	s.Status = GameStarting
	s.countdownNum = 4
	s.countdownTime = 20
}

func offsetForPlayer(index int) util.Vec2 {
	if index == 0 {
		return util.V(0, 0)
	}
	angle := firefly.Degrees(60 * float32(index-1))
	return util.AngleToVec2(angle).Scale(12)
}
