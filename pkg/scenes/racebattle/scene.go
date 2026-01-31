package racebattle

import (
	"cmp"
	"slices"

	"github.com/applejag/epic-wizard-firefly-gladiators/assets"
	"github.com/applejag/epic-wizard-firefly-gladiators/pkg/scenes"
	"github.com/applejag/epic-wizard-firefly-gladiators/pkg/state"
	"github.com/applejag/epic-wizard-firefly-gladiators/pkg/util"
	"github.com/applejag/firefly-go-math/ffmath"

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
	AnimatedClouds  util.AnimatedSheet
	VictorySplash   util.AnimatedSheet
	DefeatSplash    util.AnimatedSheet
	ButtonHighlight util.AnimatedSheet

	Players      []Firefly
	Camera       Camera
	status       GameStatus
	defeatButton DefeatButton

	countdownNum  int
	countdownTime int

	// Placement among all competitors.
	//
	// - 1 means 1st place
	// - 2 means 2nd place
	// - 3 means 3rd place
	myPlayerPlace   byte
	rewards         Rewards
	rewardsText     string
	ticksSinceStart int
}

func (s *Scene) Boot() {
	s.AnimatedClouds = assets.RacingMapClouds.Animated(2)
	s.VictorySplash = assets.VictorySplash.Animated(12)
	s.VictorySplash.AutoPlay = false
	s.DefeatSplash.Stop()
	s.DefeatSplash = assets.DefeatSplash.Animated(12)
	s.DefeatSplash.AutoPlay = false
	s.DefeatSplash.Stop()
	s.ButtonHighlight = assets.TitleButtonHighlight.Animated(2)
}

func (s *Scene) Update() {
	switch s.status {
	case GamePlaying:
		s.ticksSinceStart++

		for i := range s.Players {
			result := s.Players[i].Update()
			if result == PathTrackerLooped {
				isMyPlayer := s.Players[i].IsPlayer && s.Players[i].Peer == state.Input.Me
				if isMyPlayer {
					s.changeStatus(GameOverVictory)
				} else {
					s.changeStatus(GameOverDefeat)
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
		s.countdownTime--
		if s.countdownTime <= 0 {
			s.countdownNum--
			if s.countdownNum <= 0 {
				s.changeStatus(GamePlaying)
			} else {
				s.countdownTime = 60
			}
		}

	case GameOverVictory:
		s.VictorySplash.Update()
		s.ButtonHighlight.Update()

		if !s.VictorySplash.IsPaused() {
			return
		}
		justPressedButtons := state.Input.JustPressedButtons()
		if justPressedButtons.S || justPressedButtons.E {
			// only 1 button so "S" always goes back to field
			scenes.SwitchScene(scenes.Field)
		}
	case GameOverDefeat:
		s.DefeatSplash.Update()
		s.ButtonHighlight.Update()

		switch state.Input.JustPressedDPad4() {
		case firefly.DPad4Up:
			s.defeatButton = s.defeatButton.Next()
		case firefly.DPad4Down:
			s.defeatButton = s.defeatButton.Next()
		}
		justPressed := state.Input.JustPressedButtons()
		switch {
		case justPressed.S:
			switch s.defeatButton {
			case DefeatButtonBackToField:
				scenes.SwitchScene(scenes.Field)
			case DefeatButtonTryAgain:
				scenes.SwitchScene(scenes.RacingBattle)
			}

		case justPressed.E:
			scenes.SwitchScene(scenes.Field)
		}
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
	for i := range s.Players {
		player := &s.Players[i]
		if player.IsPlayer && player.Peer == state.Input.Me {
			myProgress = player.Progress() + float32(player.LoopsDone)
			break
		}
	}

	var playersWithHigerProgress byte
	for i := range s.Players {
		player := &s.Players[i]
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
	for i := range s.Players {
		player := &s.Players[i]
		if player.IsPlayer && player.Peer == state.Input.Me {
			me = player
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

	switch s.status {
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
		if len(s.Players) > 1 && s.myPlayerPlace >= 1 && s.myPlayerPlace <= 3 {
			assets.RacingPlace[s.myPlayerPlace-1].Draw(firefly.P(firefly.Width-28-4, 4))
		}

	case GameOverVictory:
		s.VictorySplash.DrawOrLastFrame(firefly.P(0, 0))
		if s.VictorySplash.IsPaused() {
			s.ButtonHighlight.Draw(firefly.P(137, 69))
		}

		point := firefly.P(143, 117)
		assets.FontEG_6x9.Draw(s.rewardsText, point.Add(firefly.P(1, 1)), firefly.ColorDarkGray)
		assets.FontEG_6x9.Draw(s.rewardsText, point, firefly.ColorWhite)

	case GameOverDefeat:
		s.DefeatSplash.DrawOrLastFrame(firefly.P(0, 0))
		if s.DefeatSplash.IsPaused() {
			switch s.defeatButton {
			case DefeatButtonBackToField:
				s.ButtonHighlight.Draw(firefly.P(5, 71))
			case DefeatButtonTryAgain:
				s.ButtonHighlight.Draw(firefly.P(5, 90))
			}
		}
	}
}

func (s *Scene) changeStatus(newStatus GameStatus) {
	s.status = newStatus
	switch newStatus {
	case GameOverDefeat:
		s.DefeatSplash.Play()

		if ff, ok := state.Game.InRaceBattle[state.Input.Me]; ok {
			idx := state.Game.FindFireflyByID(ff.ID)
			if idx == -1 {
				panic("this should never happen")
			}
			state.Game.Fireflies[idx].BattlesPlayed++
			state.Game.BattlesPlayedTotal++
			state.Game.Save()
		}
	case GameOverVictory:
		s.VictorySplash.Play()

		if ff, ok := state.Game.InRaceBattle[state.Input.Me]; ok {
			idx := state.Game.FindFireflyByID(ff.ID)
			if idx == -1 {
				panic("this should never happen")
			}
			state.Game.Fireflies[idx].BattlesPlayed++
			state.Game.Fireflies[idx].BattlesWon++
			state.Game.BattlesPlayedTotal++
			state.Game.BattlesWonTotal++
			s.rewards = CalculateRewards(s)

			var out [len("+") + 2 + len(" speed\n+") + 2 + len(" nimble\n+") + 3 + len(" money")]byte
			index := copy(out[0:], "+")
			index += util.FormatIntInto(out[index:], s.rewards.Speed)
			index += copy(out[index:], " speed\n+")
			index += util.FormatIntInto(out[index:], s.rewards.Nimbleness)
			index += copy(out[index:], " nimble\n+")
			index += util.FormatIntInto(out[index:], s.rewards.Money)
			s.rewardsText = string(out[:index])

			s.rewards.Apply(&state.Game.Fireflies[idx])
			state.Game.Save()
		}
	case GamePlaying:
		s.ticksSinceStart = 0
	case GameStarting:
		s.countdownNum = 4
		s.countdownTime = 20
	default:
		panic("unexpected racebattle.GameStatus")
	}
}

func (s *Scene) OnSceneEnter(players int) {
	// clear the slices instead of setting to nil, just to avoid extra allocations
	clear(s.Players)
	s.Players = s.Players[:0]

	for peer, stats := range state.Game.InRaceBattle {
		s.Players = append(s.Players, NewFireflyPlayer(peer, stats, ffmath.V(41, 390).Add(offsetForPlayer(len(s.Players))), firefly.Degrees(271)))
	}
	for len(s.Players) < players {
		s.Players = append(s.Players, NewFireflyBot(ffmath.V(41, 390).Add(offsetForPlayer(len(s.Players))), firefly.Degrees(271)))
	}
	// Update once so it focuses on player when we're transitioning to this scene
	s.Camera.Update(s)

	s.VictorySplash.Stop()
	s.DefeatSplash.Stop()
	s.defeatButton = DefeatButtonBackToField
	s.rewards = Rewards{}
	s.rewardsText = ""
	s.changeStatus(GameStarting)
}

func offsetForPlayer(index int) ffmath.Vec {
	if index == 0 {
		return ffmath.V(0, 0)
	}
	angle := firefly.Degrees(60 * float32(index-1))
	return ffmath.VAngle(angle).Scale(12)
}

type DefeatButton byte

const (
	DefeatButtonBackToField DefeatButton = iota
	DefeatButtonTryAgain

	defeatButtonCount = 2
)

func (b DefeatButton) Next() DefeatButton {
	return DefeatButton((byte(b) + 1) % defeatButtonCount)
}

func (b DefeatButton) Previous() DefeatButton {
	return DefeatButton((byte(b) + defeatButtonCount - 1) % defeatButtonCount)
}
