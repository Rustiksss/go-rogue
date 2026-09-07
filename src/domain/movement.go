package domain

import (
	"slices"
	"strconv"
)

func (g *GameInfo) move(action Action) bool {
	performedAction := false
	newPos := g.Player.Pos
	if g.FirstPersonView {
		action, performedAction = g.threeDeeMove(action)
		if !performedAction {
			return performedAction
		} else {
			performedAction = false
		}
	}
	var lookDir Direction
	switch action {
	case MoveUp:
		newPos.Y--
		lookDir = Up
	case MoveDown:
		newPos.Y++
		lookDir = Down
	case MoveLeft:
		newPos.X--
		lookDir = Left
	case MoveRight:
		newPos.X++
		lookDir = Right
	}
	if !g.FirstPersonView {
		g.LookDirection = lookDir
	}
	if g.checkTerrain(newPos) {
		enemyIndex, hasEnemy := g.checkEnemy(newPos)
		if hasEnemy {
			g.playerAttackEnemy(enemyIndex)
		} else {
			g.Player.Pos = newPos
			g.Statistic.TilesTraveled++
			g.stepOnLoot()
		}
		performedAction = true
	} else if g.checkOpenDoor(newPos) {
		g.Player.Pos = newPos
		g.Statistic.TilesTraveled++
	}
	return performedAction
}

func (g *GameInfo) checkTerrain(Pos Position) bool {
	terrainSteppable := false
	if Pos.Y < 0 || Pos.Y > g.Field.Height-1 ||
		Pos.X < 0 || Pos.X > g.Field.Width-1 {
		return terrainSteppable
	}
	tileType := g.Field.Grid[Pos.Y][Pos.X].Type
	switch tileType {
	case Ground, Corridor:
		terrainSteppable = true
	case Door:
		doorColor := g.Field.Grid[Pos.Y][Pos.X].Door
		if doorColor == NoColor {
			terrainSteppable = true
		}
	}
	return terrainSteppable
}

func (g *GameInfo) checkOpenDoor(Pos Position) bool {
	openedDoor := false
	if Pos.Y < 0 || Pos.Y > g.Field.Height-1 ||
		Pos.X < 0 || Pos.X > g.Field.Width-1 {
		return openedDoor
	}
	tileType := g.Field.Grid[Pos.Y][Pos.X].Type
	if tileType == Door {
		doorColor := g.Field.Grid[Pos.Y][Pos.X].Door
		var colorName string
		switch doorColor {
		case Red:
			colorName = "red"
		case Green:
			colorName = "green"
		case Cyan:
			colorName = "cyan"
		}
		for i := range len(g.Backpack) {
			if g.Backpack[i].GetType() == ItemKey && g.Backpack[i].(Key).Color == doorColor {
				g.Log = append(g.Log, "opened the "+colorName+" door!")
				g.Statistic.DoorsOpened++
				g.Field.Grid[Pos.Y][Pos.X].Door = NoColor
				g.Backpack = slices.Delete(g.Backpack, i, i+1)
				openedDoor = true
				break
			}
		}
		if !openedDoor {
			g.Log = append(g.Log, "the "+colorName+" door is closed")
		}
	}
	return openedDoor
}

func (g *GameInfo) checkEnemy(Pos Position) (enemyIndex int, hasEnemy bool) {
	for i := range len(g.Field.Enemies) {
		if g.Field.Enemies[i].Pos == Pos {
			enemyIndex = i
			hasEnemy = true
		}
	}
	return
}

func (g *GameInfo) nextLevel() (err error) {
	if g.Field.Staircase.Y != g.Player.Pos.Y ||
		g.Field.Staircase.X != g.Player.Pos.X {
		g.Log = append(g.Log, "I see no way down")
	} else {
		if g.Level == MaxLevel {
			g.Statistic.IsWon = true
			g.Victory = true

			g.Leaderboard = append(g.Leaderboard, g.Statistic)
			g.sortLeaderboard()
			err = g.SaveLoader.SaveLeaderboard(g.Leaderboard)

			g.calcCurGameTop()

			err = g.SaveLoader.DeleteSavedGame()
			g.SkipInput = false
			return
		}
		g.Level++
		g.Statistic.LevelReached = g.Level
		g.newGameMap()
		g.deleteKeys()
		g.Log = append(g.Log, "welcome to level "+strconv.Itoa(g.Level))

		g.SaveLoader.SaveGame(*g)
	}
	return
}

func (g *GameInfo) threeDeeMove(action Action) (correctedAction Action, performedAction bool) {
	correctedAction = action

	var moveTransformations = map[Direction]map[Action]Action{
		Up: {
			MoveUp:   MoveUp,
			MoveDown: MoveDown,
		},
		Down: {
			MoveUp:   MoveDown,
			MoveDown: MoveUp,
		},
		Left: {
			MoveUp:   MoveLeft,
			MoveDown: MoveRight,
		},
		Right: {
			MoveUp:   MoveRight,
			MoveDown: MoveLeft,
		},
	}
	if transforms, ok := moveTransformations[g.LookDirection]; ok {
		if correctedAction, found := transforms[action]; found {
			return correctedAction, true
		}
	}

	var moveRotations = map[Direction]map[Action]Direction{
		Up: {
			MoveLeft:  Left,
			MoveRight: Right,
		},
		Down: {
			MoveLeft:  Right,
			MoveRight: Left,
		},
		Left: {
			MoveLeft:  Down,
			MoveRight: Up,
		},
		Right: {
			MoveLeft:  Up,
			MoveRight: Down,
		},
	}

	if rotations, ok := moveRotations[g.LookDirection]; ok {
		if lookDirection, found := rotations[action]; found {
			g.LookDirection = lookDirection
			return
		}
	}

	return
}
