package domain

import (
	"slices"
	"sort"
	"time"
)

func (g *GameInfo) Update(action Action) (err error) {
	g.Action = action
	if g.SkipInput {
		action = NoAction
	}
	if g.Action == Quit {
		err = g.processExit()
		return
	}

	screenHandled, err := g.handleScreens()
	if screenHandled || err != nil {
		return
	}

	skipTurn, err := g.handlePlayerAction()
	if skipTurn || err != nil {
		return
	}

	err = g.processGameWorld()
	return
}

// handleScreens controls the logic of windows on top of the game world
// returns true if the action has been processed by the screen and the turn should not continue
func (g *GameInfo) handleScreens() (isScreenHandled bool, err error) {
	isScreenHandled = true
	if g.LeaderboardScr && !g.StartScreen {
		if g.Action == Confirm || g.Action == Return || g.Action == Leaderboard {
			g.LeaderboardScr = false
		}
		return
	}
	if g.Action == Leaderboard && !g.HelpScreen {
		g.calcCurGameTop()
		g.LeaderboardScr = true
		return
	}

	if g.GameOver || g.Victory {
		return
	}

	if g.StartScreen {
		if g.WaitLoadChoice {
			err = g.loadChoice()
		}
		return
	}

	if g.HelpScreen {
		if g.Action == Confirm || g.Action == Return || g.Action == Help {
			g.HelpScreen = false
		}
		return
	}
	if g.Action == Help {
		g.HelpScreen = true
		return
	}

	if g.OpenBackpack {
		if g.Action == Confirm || g.Action == Return || g.Action == Inventory {
			g.OpenBackpack = false
		}
		return
	}
	if g.Action == Inventory && !g.UsingItem {
		g.OpenBackpack = true
		return
	}

	if g.Action == FirstPerson {
		if g.FirstPersonView == false {
			g.FirstPersonView = true
		} else {
			g.FirstPersonView = false
		}
		return
	}
	if g.FirstPersonView && g.Action == MiniMap {
		if g.FirstPersonMiniMap {
			g.FirstPersonMiniMap = false
		} else {
			g.FirstPersonMiniMap = true
		}
	}

	if g.logSkipTurn() {
		return
	}

	if g.UsingItem {
		if !g.useItem() {
			return
		} else {
			g.Action = Idle
		}
	}

	isScreenHandled = false
	return
}

func (g *GameInfo) handlePlayerAction() (skipTurn bool, err error) {
	skipTurn = true
	switch g.Action {
	case MoveUp, MoveDown, MoveLeft, MoveRight:
		if !g.move(g.Action) {
			return
		}
	case EquipWeapon, EquipArmor, EatFood, QuaffElixir, ReadScroll, DropItem:
		g.UsingItemModeTurnOn(g.Action)
		return
	case GoDownTheStairs:
		err = g.nextLevel()
		return
	case Idle:
	default:
		if !g.SkipInput {
			return
		}
	}
	skipTurn = false
	return
}

func (g *GameInfo) processGameWorld() (err error) {
	if g.SkipInput {
		time.Sleep(time.Duration(500) * time.Millisecond)
		g.SkipInput = false
	}

	g.decreaseBuffs()
	g.recalculatePlayerDamage()
	g.recalculatePlayerArmorAndDexterity()

	err = g.enemiesAction()

	g.calcFogOfWar()
	g.calcDifficulty()
	return
}

func NewGameInfo(saveLoader SaveLoad) (g *GameInfo, err error) {
	g = &GameInfo{
		Level:          1,
		Player:         basePlayer,
		StartScreen:    true,
		WaitNameInput:  true,
		LeaderboardScr: true,
		SkipInput:      true,
	}
	g.newGameMap()
	g.SaveLoader = saveLoader
	g.Leaderboard, err = g.SaveLoader.LoadLeaderboard()
	if err != nil {
		return
	}
	has, err := g.SaveLoader.HasSavedGame()
	if err != nil {
		return
	}
	if has {
		g.loadGame()
	} else {
		g.statisticInit()
	}

	g.calcDifficulty()

	return
}

func (g *GameInfo) SetPlayerName(name string) {
	g.Statistic.PlayerName = name
	g.StartScreen = false
	g.WaitNameInput = false
	g.SkipInput = false
	g.LeaderboardScr = false
}

func (g *GameInfo) processExit() (err error) {
	g.Quit = true

	if !g.Statistic.IsDead && !g.Statistic.IsWon && !g.Statistic.IsAbandoned {
		err = g.SaveLoader.SaveGame(*g)
	}

	return
}

func (g *GameInfo) statisticInit() {
	g.Statistic.Kills = append(g.Statistic.Kills, ZombieStat)
	g.Statistic.Kills = append(g.Statistic.Kills, VampireStat)
	g.Statistic.Kills = append(g.Statistic.Kills, GhostStat)
	g.Statistic.Kills = append(g.Statistic.Kills, OgreStat)
	g.Statistic.Kills = append(g.Statistic.Kills, SnakeMageStat)

	g.Statistic.LevelReached = 1
}

func (g *GameInfo) loadChoice() (err error) {
	g.SkipInput = false
	switch g.Action {
	case Quit:
		g.Quit = true
		return
	case Yes:
		g.StartScreen = false
		return
	case No:
		g.SaveLoader.DeleteSavedGame()

		g.Statistic.IsAbandoned = true

		g.Leaderboard = append(g.Leaderboard, g.Statistic)
		g.sortLeaderboard()
		err = g.SaveLoader.SaveLeaderboard(g.Leaderboard)

		var gNewPtr *GameInfo
		gNewPtr, err = NewGameInfo(g.SaveLoader)
		*g = *gNewPtr
	default:
		return
	}
	g.WaitLoadChoice = false
	return
}

func (g *GameInfo) loadGame() (err error) {
	copySaveLoader := g.SaveLoader
	*g, err = g.SaveLoader.LoadGame()
	g.Quit = false
	g.SaveLoader = copySaveLoader

	g.Leaderboard, err = g.SaveLoader.LoadLeaderboard()

	g.calcCurGameTop()
	g.StartScreen = true
	g.SkipInput = false
	g.WaitLoadChoice = true
	g.WaitNameInput = false
	return
}

// test map
func (g *GameInfo) newTestMap() {
	g.Field = GameMap{
		Grid:      [FieldHeight][FieldWidth]Terrain{},
		Staircase: Position{},
		Loot:      []Loot{},
		Height:    FieldHeight,
		Width:     FieldWidth,
	}

	for row := range g.Field.Height {
		for col := range g.Field.Width {
			if row == 10 && col == 10 {
				g.Field.Grid[row][col] = TerrainDoor
			} else if row == 0 || row == FieldHeight-1 || col == 0 || col == FieldWidth-1 ||
				(row == 6 && col >= 10 && col <= 50) || (row == 14 && col >= 10 && col <= 50) ||
				(col == 10 && row >= 6 && row <= 14) || (col == 50 && row >= 6 && row <= 14) {
				g.Field.Grid[row][col] = TerrainWall
			} else {
				g.Field.Grid[row][col] = TerrainGround
			}

		}
	}

	g.Player.Pos.Y = 10
	g.Player.Pos.X = 20

	g.Field.Grid[19][40] = TerrainDoorRed
	g.Field.Grid[19][42] = TerrainDoorGreen
	g.Field.Grid[19][44] = TerrainDoorCyan
	g.Field.Loot = append(g.Field.Loot, Loot{redKey, Position{19, 38}})
	g.Field.Loot = append(g.Field.Loot, Loot{greenKey, Position{19, 36}})
	g.Field.Loot = append(g.Field.Loot, Loot{cyanKey, Position{19, 34}})

	g.Field.Staircase.Y = 10
	g.Field.Staircase.X = 22

	g.Field.Loot = append(g.Field.Loot, Loot{dagger, Position{10, 24}})
	g.Field.Loot = append(g.Field.Loot, Loot{chainmail, Position{10, 26}})
	g.Field.Loot = append(g.Field.Loot, Loot{healthScroll, Position{10, 28}})
	g.Field.Loot = append(g.Field.Loot, Loot{healthElixir, Position{10, 30}})
	g.Field.Loot = append(g.Field.Loot, Loot{strengthScroll, Position{9, 28}})
	g.Field.Loot = append(g.Field.Loot, Loot{strengthElixir, Position{9, 30}})
	g.Field.Loot = append(g.Field.Loot, Loot{dexterityScroll, Position{8, 28}})
	g.Field.Loot = append(g.Field.Loot, Loot{dexterityElixir, Position{8, 30}})
	g.Field.Loot = append(g.Field.Loot, Loot{Treasure{ItemTreasure, 100}, Position{10, 32}})
	g.Field.Loot = append(g.Field.Loot, Loot{dryRation, Position{10, 34}})
	g.Field.Loot = append(g.Field.Loot, Loot{Treasure{ItemTreasure, 150}, Position{10, 36}})
	g.Field.Loot = append(g.Field.Loot, Loot{Treasure{ItemTreasure, 70}, Position{10, 38}})
	g.Field.Loot = append(g.Field.Loot, Loot{dagger, Position{11, 24}})
	g.Field.Loot = append(g.Field.Loot, Loot{dagger, Position{12, 24}})
	g.Field.Loot = append(g.Field.Loot, Loot{dagger, Position{13, 24}})
	g.Field.Loot = append(g.Field.Loot, Loot{dagger, Position{11, 26}})
	g.Field.Loot = append(g.Field.Loot, Loot{dagger, Position{12, 26}})
	g.Field.Loot = append(g.Field.Loot, Loot{dagger, Position{13, 26}})
	g.Field.Loot = append(g.Field.Loot, Loot{dagger, Position{11, 28}})
	g.Field.Loot = append(g.Field.Loot, Loot{dagger, Position{12, 28}})
	g.Field.Loot = append(g.Field.Loot, Loot{dagger, Position{13, 28}})
	for range 10 {
		g.Field.Loot = append(g.Field.Loot, Loot{dagger, Position{13, 28}})
	}

	g.Field.Enemies = append(g.Field.Enemies, BaseZombie)
	g.Field.Enemies[len(g.Field.Enemies)-1].Pos.Y = 3
	g.Field.Enemies[len(g.Field.Enemies)-1].Pos.X = 20

	g.Field.Enemies = append(g.Field.Enemies, BaseVampire)
	g.Field.Enemies[len(g.Field.Enemies)-1].Pos.Y = 3
	g.Field.Enemies[len(g.Field.Enemies)-1].Pos.X = 22

	g.Field.Enemies = append(g.Field.Enemies, BaseGhost)
	g.Field.Enemies[len(g.Field.Enemies)-1].Pos.Y = 3
	g.Field.Enemies[len(g.Field.Enemies)-1].Pos.X = 24

	g.Field.Enemies = append(g.Field.Enemies, BaseOgre)
	g.Field.Enemies[len(g.Field.Enemies)-1].Pos.Y = 3
	g.Field.Enemies[len(g.Field.Enemies)-1].Pos.X = 26

	g.Field.Enemies = append(g.Field.Enemies, BaseSnakeMage)
	g.Field.Enemies[len(g.Field.Enemies)-1].Pos.Y = 3
	g.Field.Enemies[len(g.Field.Enemies)-1].Pos.X = 28

	g.Field.Enemies = append(g.Field.Enemies, BaseMimic)
	g.Field.Enemies[len(g.Field.Enemies)-1].Pos.Y = 3
	g.Field.Enemies[len(g.Field.Enemies)-1].Pos.X = 30

	g.calcFogOfWar()
}

func (g *GameInfo) logSkipTurn() (skip bool) {
	if g.Action == NoAction && !g.SkipInput {
		skip = true
		return
	} else if g.SkipInput {
		skip = false
		return
	}
	logCount := len(g.Log)
	if logCount > 0 {
		if g.Action == Confirm || logCount == 1 {
			g.Log = slices.Delete(g.Log, 0, 1)
		}
		if len(g.Log) > 0 {
			skip = true
		}
	}
	return
}

func (g *GameInfo) sortLeaderboard() {
	sort.Slice(g.Leaderboard, func(i, j int) bool {
		return g.Leaderboard[i].TreasAmt > g.Leaderboard[j].TreasAmt
	})
}

func (g *GameInfo) calcCurGameTop() {
	g.CurGameTop = 1
	for i := range g.Leaderboard {
		if g.Statistic.TreasAmt < g.Leaderboard[i].TreasAmt {
			g.CurGameTop += 1
		} else {
			break
		}
	}
}

func (g *GameInfo) calcDifficulty() {
	g.Difficulty = ((float64(g.Level) - 10) / 3) +
		((float64(g.Player.MaxHealth) - 200) / 100) +
		((float64(g.Player.Strength) - 17) / 10) +
		((float64(g.Statistic.HitsLanded) - float64(g.Statistic.HitsMissed)) / 100) +
		((float64(g.Statistic.HitsDodged) - float64(g.Statistic.HitsTaken)) / 100)
}
