package render

import (
	"math"
	"rogue/domain"
	"sync"
	"time"

	"github.com/gbin/goncurses"
)

var onceRender sync.Once

func Render(g *domain.GameInfo, mainScr, fieldScr, logScr, statScr, backpackScr, miniMapScr *goncurses.Window) {

	onceRender.Do(func() {
		mainScr.Keypad(true)
		goncurses.Echo(false)

		goncurses.InitPair(1, goncurses.C_YELLOW, goncurses.C_BLACK)
		goncurses.InitPair(2, goncurses.C_BLUE, goncurses.C_BLACK)
		goncurses.InitPair(3, goncurses.C_MAGENTA, goncurses.C_MAGENTA)
		goncurses.InitPair(4, goncurses.C_WHITE, goncurses.C_BLACK)
		goncurses.InitPair(5, goncurses.C_WHITE, goncurses.C_WHITE)
		goncurses.InitPair(6, goncurses.C_RED, goncurses.C_BLACK)
		goncurses.InitPair(7, goncurses.C_BLACK, goncurses.C_WHITE)
		goncurses.InitPair(8, goncurses.C_GREEN, goncurses.C_BLACK)
		goncurses.InitPair(9, goncurses.C_YELLOW, goncurses.C_BLACK)
		goncurses.InitPair(10, goncurses.C_CYAN, goncurses.C_BLACK)
		goncurses.InitPair(11, goncurses.C_MAGENTA, goncurses.C_BLACK)
		goncurses.InitColor(12, 700, 600, 900)
		goncurses.InitPair(12, 12, goncurses.C_BLACK)
	})

	mainScr.Erase()

	mainScr.NoutRefresh()
	renderLog(g, logScr)
	renderStat(g, statScr)
	renderField(g, fieldScr)
	renderFirstPersonView(g, fieldScr, miniMapScr)
	renderBackpack(g, backpackScr)
	renderHelpScreen(g, mainScr)
	renderStartScreen(g, mainScr)
	renderGameOver(g, mainScr)
	renderVictory(g, mainScr)
	if !g.StartScreen {
		renderLeaderboard(g, mainScr, 0)
	}

	goncurses.Update()
}

func renderStartScreen(g *domain.GameInfo, mainScr *goncurses.Window) {
	if !g.StartScreen {
		return
	}
	mainScr.Erase()

	renderLeaderboard(g, mainScr, 4)

	mainScr.AttrOn(goncurses.A_BOLD)
	mainScr.MovePrint(1, 37, "ROGUE")
	mainScr.AttrOff(goncurses.A_BOLD)

	if g.WaitLoadChoice {
		mainScr.MovePrint(3, 3, "Found save file. Load? (y/n)")
	} else if g.WaitNameInput {
		mainScr.MovePrint(3, 3, "Enter name: ")
	}

	mainScr.NoutRefresh()

	if g.WaitLoadChoice {
		goncurses.Cursor(0)
	} else if g.WaitNameInput {
		mainScr.Move(3, 15)
	}

}

func renderLog(g *domain.GameInfo, logScr *goncurses.Window) {
	logScr.Erase()

	logCount := len(g.Log)
	if logCount == 1 {
		logScr.Print(g.Log[0])
	} else if logCount > 1 {
		logScr.Print(g.Log[0], "-more-")
	}

	logScr.NoutRefresh()
}

func renderStat(g *domain.GameInfo, statScr *goncurses.Window) {
	statScr.Erase()

	statScr.ColorOn(1)
	defer statScr.ColorOff(1)

	statScr.AttrOn(goncurses.A_BOLD)
	defer statScr.AttrOff(goncurses.A_BOLD)

	statScr.MovePrint(0, 0, "Level: ", g.Level)
	statScr.MovePrint(0, 10, "Health: ", g.Player.CurHealth, "/", g.Player.MaxHealth)
	statScr.MovePrint(0, 26, "Strength: ", g.Player.Strength)
	statScr.MovePrint(0, 40, "Dexterity: ", g.Player.Dexterity)
	statScr.MovePrint(0, 54, "Armor: ", g.Player.Armor)
	statScr.MovePrint(0, 64, "Treasure: ", g.Player.TreasAmt)

	statScr.NoutRefresh()

}

func renderField(g *domain.GameInfo, fieldScr *goncurses.Window) {
	fieldScr.Erase()

	renderTerrain(g, fieldScr)
	renderLoot(g, fieldScr)
	renderEnemies(g, fieldScr)

	playerPos := domain.Position{Y: g.Player.Pos.Y, X: g.Player.Pos.X}
	printPlayer(fieldScr, playerPos)

	fieldScr.NoutRefresh()
}

func renderTerrain(g *domain.GameInfo, fieldScr *goncurses.Window) {
	for fieldY := range g.Field.Height {
		for fieldX := range g.Field.Width {
			tile := g.Field.Grid[fieldY][fieldX]
			printTerrainTile(tile, fieldScr, fieldY, fieldX)
		}
	}
	stairsPos := domain.Position{Y: g.Field.Staircase.Y, X: g.Field.Staircase.X}
	if g.Field.Grid[stairsPos.Y][stairsPos.X].Fog == domain.Visible {
		printStaircase(fieldScr, stairsPos)
	}
}

func printTerrainTile(tile domain.Terrain, scr *goncurses.Window, y, x int) {
	var ch goncurses.Char
	var color int16
	switch tile.Type {
	case domain.Nothing:
		ch = ' '
	case domain.Ground:
		ch, color = getGroundSymbolAndColor(tile.Fog)
	case domain.Wall:
		ch, color = getWallSymbolAndColor(tile.Fog)
	case domain.Corridor:
		ch, color = getCorridorSymbolAndColor(tile.Fog)
	case domain.Door:
		ch, color = getDoorSymbolAndColor(tile.Fog, tile.Door)
	}
	if tile.Type == domain.Door && tile.Fog != domain.Unknown {
		scr.AttrOn(goncurses.A_BOLD)
	}
	scr.ColorOn(color)
	scr.MoveAddChar(y, x, ch)
	scr.ColorOff(color)
	if tile.Type == domain.Door {
		scr.AttrOff(goncurses.A_BOLD)
	}
}

func getGroundSymbolAndColor(fog domain.FogOfWar) (ch goncurses.Char, color int16) {
	switch fog {
	case domain.Unknown:
		ch = ' '
		color = 4
	case domain.NotVisible:
		ch = ' '
		color = 4
	case domain.Visible:
		ch = '.'
		color = 2
	}
	return
}

func getWallSymbolAndColor(fog domain.FogOfWar) (ch goncurses.Char, color int16) {
	ch = ' '
	switch fog {
	case domain.Unknown:
		color = 4
	case domain.NotVisible:
		color = 7
	case domain.Visible:
		color = 3
	}
	return
}

func getCorridorSymbolAndColor(fog domain.FogOfWar) (ch goncurses.Char, color int16) {
	switch fog {
	case domain.Unknown:
		ch = ' '
		color = 4
	case domain.NotVisible:
		ch = '#'
		color = 4
	case domain.Visible:
		ch = '#'
		color = 2
	}
	return
}

func getDoorSymbolAndColor(fog domain.FogOfWar, doorColor domain.DoorColor) (ch goncurses.Char, color int16) {
	switch fog {
	case domain.Unknown:
		ch = ' '
		color = 4
	case domain.NotVisible:
		ch = '+'
		color = 4
	case domain.Visible:
		ch = '+'
		color = getDoorColorPair(doorColor)
	}
	return
}

func getDoorColorPair(doorColor domain.DoorColor) (color int16) {
	switch doorColor {
	case domain.NoColor:
		color = 11
	case domain.Red:
		color = 6
	case domain.Green:
		color = 8
	case domain.Cyan:
		color = 10
	}
	return
}

func printStaircase(scr *goncurses.Window, stairsPos domain.Position) {
	scr.ColorOn(4)
	scr.MoveAddChar(stairsPos.Y, stairsPos.X, '%')
	scr.ColorOff(4)
}

func renderLoot(g *domain.GameInfo, fieldScr *goncurses.Window) {
	for i := range len(g.Field.Loot) {
		loot := &g.Field.Loot[i]
		if g.Field.Grid[loot.Pos.Y][loot.Pos.X].Fog == domain.Visible {
			ch, color := getLootSymbolAndColor(loot)
			fieldScr.ColorOn(color)
			fieldScr.MoveAddChar(g.Field.Loot[i].Pos.Y, g.Field.Loot[i].Pos.X, ch)
			fieldScr.ColorOff(color)
		}
	}
}

func getLootSymbolAndColor(loot *domain.Loot) (ch goncurses.Char, color int16) {
	color = 4
	switch loot.Content.GetType() {
	case domain.ItemWeapon:
		ch = ')'
	case domain.ItemArmor:
		ch = '&'
	case domain.ItemScroll:
		ch = '?'
	case domain.ItemElixir:
		ch = '!'
	case domain.ItemTreasure:
		ch = '*'
	case domain.ItemFood:
		ch = ';'
	case domain.ItemKey:
		ch = '$'
		color = getDoorColorPair(loot.Content.(domain.Key).Color)
	}
	return
}

func renderEnemies(g *domain.GameInfo, fieldScr *goncurses.Window) {
	for i := range len(g.Field.Enemies) {
		enemy := &g.Field.Enemies[i]
		if g.Field.Grid[enemy.Pos.Y][enemy.Pos.X].Fog == domain.Visible {
			ch, color := getEnemySymbolAndColor(enemy)
			fieldScr.ColorOn(color)
			if !(enemy.Type == domain.Ghost && enemy.Flags.ContainsAny(domain.FlagInvisible)) {
				fieldScr.MoveAddChar(enemy.Pos.Y, enemy.Pos.X, ch)
			}
			fieldScr.ColorOff(color)
		}
	}
}

func getEnemySymbolAndColor(enemy *domain.Enemy) (ch goncurses.Char, color int16) {
	switch enemy.Type {
	case domain.Zombie:
		ch = 'Z'
		color = 8
	case domain.Vampire:
		ch = 'V'
		color = 6
	case domain.Ghost:
		ch = 'G'
		color = 4
	case domain.Ogre:
		ch = 'O'
		color = 9
	case domain.SnakeMage:
		ch = 'S'
		color = 4
	case domain.Mimic:
		if enemy.Flags.ContainsAll(domain.FlagMimickingWeapon) {
			ch = ')'
		} else if enemy.Flags.ContainsAll(domain.FlagMimickingArmor) {
			ch = '&'
		} else if enemy.Flags.ContainsAll(domain.FlagMimickingScroll) {
			ch = '?'
		} else if enemy.Flags.ContainsAll(domain.FlagMimickingElixir) {
			ch = '!'
		} else if enemy.Flags.ContainsAll(domain.FlagMimickingFood) {
			ch = ';'
		} else if enemy.Flags.ContainsAll(domain.FlagMimickingTreasure) {
			ch = '*'
		} else {
			ch = 'M'
		}
		color = 4
	}
	return
}

func printPlayer(scr *goncurses.Window, playerPos domain.Position) {
	goncurses.Cursor(1)
	scr.ColorOn(4)
	scr.MoveAddChar(playerPos.Y, playerPos.X, '@')
	scr.ColorOff(4)
	scr.Move(playerPos.Y, playerPos.X)
}

func renderBackpack(g *domain.GameInfo, backpackScr *goncurses.Window) {
	if !g.OpenBackpack && !g.UsingItem {
		return
	}
	backpackScr.Erase()

	maxY, maxX := backpackScr.MaxYX()
	backpackScr.ColorOn(5)
	for row := range maxY {
		for col := range maxX {
			if row == 0 || row == maxY-1 || col == 0 || col == maxX-1 {
				backpackScr.MoveAddChar(row, col, ' ')
			}
		}
	}
	backpackScr.ColorOff(5)
	backpackScr.ColorOn(4)
	for i := range len(g.Backpack) {
		itemName := domain.GetItemName(g.Backpack[i])
		backpackScr.MovePrint(i+1, 1, i+1, ") ", itemName)
	}
	if g.Player.HasWeapon {
		backpackScr.MovePrint(maxY-3, 1, "wielded: ", g.Player.WieldedWeapon.Name)
	}
	if g.Player.HasArmor {
		backpackScr.MovePrint(maxY-4, 1, "worn: ", g.Player.WornArmor.Name)
	}

	lineToPrint := ""

	if g.OpenBackpack {
		lineToPrint = "--press space to continue--"
	} else if g.UsingItem {
		switch g.UsingItemType {
		case domain.ItemWeapon:
			lineToPrint = "--select weapon (or 0 for unequip)--"
		case domain.ItemArmor:
			lineToPrint = "--select armor (or 0 for unequip)--"
		case domain.ItemFood:
			lineToPrint = "--select what to eat--"
		case domain.ItemScroll:
			lineToPrint = "--select what to read--"
		case domain.ItemElixir:
			lineToPrint = "--select what to quaff--"
		case domain.ItemNoType:
			lineToPrint = "--select what to drop--"
		}
	}
	backpackScr.MovePrint(maxY-2, 1, lineToPrint)
	backpackScr.ColorOff(4)

	goncurses.Cursor(0)

	backpackScr.NoutRefresh()
}

func renderHelpScreen(g *domain.GameInfo, mainScr *goncurses.Window) {
	if !g.HelpScreen {
		return
	}
	mainScr.Erase()

	goncurses.Cursor(0)

	mainScr.ColorOn(4)

	_, maxX := mainScr.MaxYX()
	firstColumn, secondColumn := 7, 42

	mainScr.AttrOn(goncurses.A_BOLD)
	mainScr.MovePrintf(1, (maxX/2)-5, "CONTROLS")
	mainScr.AttrOff(goncurses.A_BOLD)

	mainScr.MovePrintf(3, firstColumn, "'.' - Idle")
	mainScr.MovePrintf(4, firstColumn, "' ' - Confirm (Space Bar)")
	mainScr.MovePrintf(5, firstColumn, "'W' - Move Up")
	mainScr.MovePrintf(6, firstColumn, "'S' - Move Down")
	mainScr.MovePrintf(7, firstColumn, "'A' - Move Left")
	mainScr.MovePrintf(8, firstColumn, "'D' - Move Right")
	mainScr.MovePrintf(9, firstColumn, "'V' - Go Down The Stairs")
	mainScr.MovePrintf(10, firstColumn, "'I' - Open Backpack")
	mainScr.MovePrintf(11, firstColumn, "'H' - Equip Weapon (Or Unequip)")
	mainScr.MovePrintf(12, firstColumn, "'G' - Equip Armor (Or Unequip)")
	mainScr.MovePrintf(13, firstColumn, "'J' - Eat Food")
	mainScr.MovePrintf(14, firstColumn, "'K' - Quaff Elixir")
	mainScr.MovePrintf(15, firstColumn, "'E' - Read Scroll")
	mainScr.MovePrintf(16, firstColumn, "'T' - Drop Item")
	mainScr.MovePrintf(17, firstColumn, "'L' - Open Leaderboard")
	mainScr.MovePrintf(18, firstColumn, "'`' - Return (Under 'Esc')")
	mainScr.MovePrintf(3, secondColumn, "'Q' - Exit (With Save)")
	mainScr.MovePrintf(4, secondColumn, "'0' - Equipped Slot")
	mainScr.MovePrintf(5, secondColumn, "'1' - First Backpack Slot")
	mainScr.MovePrintf(6, secondColumn, "'2' - Second Backpack Slot")
	mainScr.MovePrintf(7, secondColumn, "'3' - Third Backpack Slot")
	mainScr.MovePrintf(8, secondColumn, "'4' - Fourth Backpack Slot")
	mainScr.MovePrintf(9, secondColumn, "'5' - Fifth Backpack Slot")
	mainScr.MovePrintf(10, secondColumn, "'6' - Sixth Backpack Slot")
	mainScr.MovePrintf(11, secondColumn, "'7' - Seventh Backpack Slot")
	mainScr.MovePrintf(12, secondColumn, "'8' - Eighth Backpack Slot")
	mainScr.MovePrintf(13, secondColumn, "'9' - Nineth Backpack Slot")
	mainScr.MovePrintf(14, secondColumn, "'Y' - Yes")
	mainScr.MovePrintf(15, secondColumn, "'N' - No")
	mainScr.MovePrintf(16, secondColumn, "'O' - First Person Mode Toggle")
	mainScr.MovePrintf(17, secondColumn, "'M' - First Person Mode Mini Map Toggle")
	mainScr.MovePrintf(18, secondColumn, "'?' - Open Controls Screen")

	mainScr.AttrOn(goncurses.A_BOLD)
	mainScr.MovePrintf(21, (maxX/2)-14, "--press space to continue--")
	mainScr.AttrOff(goncurses.A_BOLD)

	mainScr.ColorOff(4)

	mainScr.NoutRefresh()
}

var onceGameOver sync.Once

func renderGameOver(g *domain.GameInfo, mainScr *goncurses.Window) {
	if !g.GameOver {
		return
	}

	mainScr.Erase()
	mainScr.ColorOn(4)

	onceGameOver.Do(func() {
		mainScr.Timeout(0)
	})

	goncurses.Cursor(0)

	tick := time.Now().Second() % 2

	skull := []string{
		"          _,.-------.,_          ",
		"      ,;~'             '~;,      ",
		"    ,;                     ;,    ",
		"   ;                         ;   ",
		"  ,'                         ',  ",
		" ,;                           ;, ",
		" ; ;      .           .      ; ; ",
		" | ;   ______       ______   ; | ",
		" |  `/~\"     ~\" . \"~     \"~\\'  | ",
		" |  ~  ,-^^^^-,   ,-^^^^-,  ~  | ",
		"  |   |        }:{        |   |  ",
		"  |   !   *   / | \\   *   !   |  ",
		"  .~  (__,.--\" .^. \"--.,__)  ~.  ",
		"  |     ---;' / | \\ `;---     |  ",
		"   \\__.       \\/^\\/       .__/   ",
		"    V| \\                 / |V    ",
		"     | |T~\\___!___!___/~T| |     ",
		"     | |`_I_I_I_I_I_I_I_'| |     ",
		"     |  \\,I_I I I I I_I,/  |     ",
		"      \\   `~~~~~~~~~~~'   /      ",
		"        \\   .   ^   .   /        ",
		"          \\._________./          ",
	}

	startY := 1
	startX := 23

	for y, line := range skull {
		for x, char := range line {
			if goncurses.Char(char) == '*' && tick == 0 {
				mainScr.ColorOn(6)
				mainScr.MoveAddChar(startY+y, startX+x, goncurses.Char(char))
			} else if goncurses.Char(char) != '*' {
				mainScr.MoveAddChar(startY+y, startX+x, goncurses.Char(char))
			}
			if goncurses.Char(char) == '*' && tick == 0 {
				mainScr.ColorOff(6)
			}
		}
	}

	mainScr.AttrOn(goncurses.A_BOLD)
	mainScr.MovePrint(1, 6, "GAME OVER")
	mainScr.MovePrint(1, 63, "GAME OVER")
	mainScr.AttrOff(goncurses.A_BOLD)

	topStatOffset := 0
	if g.CurGameTop == 1 {
		topStatOffset = 1
		mainScr.AttrOn(goncurses.A_BOLD)
		mainScr.MovePrint(3, 2, "new high-score!")
		mainScr.MovePrint(3, 59, "new high-score!")
		mainScr.AttrOff(goncurses.A_BOLD)
	}

	renderTopStatistics(g.Statistic, mainScr, 3+topStatOffset, 1, g.CurGameTop, true)
	renderTopStatistics(g.Statistic, mainScr, 3+topStatOffset, 58, g.CurGameTop, true)

	mainScr.ColorOff(4)

	mainScr.NoutRefresh()
}

func renderVictory(g *domain.GameInfo, mainScr *goncurses.Window) {
	if !g.Victory {
		return
	}

	mainScr.Erase()
	goncurses.Cursor(0)
	mainScr.ColorOn(4)

	trophy := []string{
		"  ___________  ",
		" '._==_==_=_.' ",
		" .-\\:      /-. ",
		"| (|:.     |) |",
		" '-|:.     |-' ",
		"   \\::.    /   ",
		"    '::. .'    ",
		"      ) (      ",
		"    _.' '._    ",
		"   `\"\"\"\"\"\"\"`   ",
	}

	mainScr.ColorOn(9)
	startY := 1
	startX := 5
	for i := range 4 {
		switch i {
		case 1:
			startY = 11
			startX = 5
		case 2:
			startY = 1
			startX = 59
		case 3:
			startY = 11
			startX = 59
		}
		for y, line := range trophy {
			for x, char := range line {
				mainScr.MoveAddChar(startY+y, startX+x, goncurses.Char(char))
			}
		}
	}

	mainScr.ColorOff(9)

	mainScr.AttrOn(goncurses.A_BOLD)
	mainScr.MovePrint(1, 35, "VICTORY!")
	mainScr.AttrOff(goncurses.A_BOLD)

	topStatOffset := 0
	if g.CurGameTop == 1 {
		topStatOffset = 1
		mainScr.AttrOn(goncurses.A_BOLD)
		mainScr.MovePrint(3, 31, "new high-score!")
		mainScr.AttrOff(goncurses.A_BOLD)
	}
	renderTopStatistics(g.Statistic, mainScr, 3+topStatOffset, 30, g.CurGameTop, true)

	mainScr.ColorOff(4)
	mainScr.NoutRefresh()
}

func renderLeaderboard(g *domain.GameInfo, mainScr *goncurses.Window, offset int) {
	if !g.LeaderboardScr && !g.StartScreen {
		return
	}
	mainScr.Erase()
	goncurses.Cursor(0)
	mainScr.ColorOn(4)

	mainScr.AttrOn(goncurses.A_BOLD)
	mainScr.MovePrint(1+offset, 34, "LEADERBOARD")
	mainScr.AttrOff(goncurses.A_BOLD)

	leadLength := len(g.Leaderboard) + 1
	if leadLength > 3 {
		leadLength = 3
	} else {
		if g.WaitNameInput || g.GameOver || g.Victory {
			leadLength--
		}
	}

	for i, numOfInsertedRecords := 0, 0; i < leadLength; i++ {
		var record domain.Statistics
		if g.CurGameTop == i+1 && !g.GameOver && !g.Victory {
			numOfInsertedRecords = 1
			record = g.Statistic
		} else {
			record = g.Leaderboard[i-numOfInsertedRecords]
		}
		renderTopStatistics(record, mainScr, 3+offset, (i*26 + 1), i+1, true)
	}

	mainScr.ColorOff(4)
	mainScr.NoutRefresh()
}

func renderTopStatistics(stat domain.Statistics, scr *goncurses.Window, firstRow, firstColumn, recordNumber int, needNum bool) {
	secondColumn := firstColumn + 18
	if needNum {
		scr.MovePrint(firstRow, firstColumn+1, "top ", recordNumber)
	}
	scr.MovePrint(firstRow+1, firstColumn, "player: "+stat.PlayerName)

	if stat.IsWon {
		scr.MovePrint(firstRow+2, firstColumn, "victory")
	} else if stat.IsDead {
		scr.MovePrint(firstRow+2, firstColumn, "killed by: ", stat.EnemyWhoKilled.Name)
	} else if stat.IsAbandoned {
		scr.MovePrint(firstRow+2, firstColumn, "game is adandoned")
	} else {
		scr.AttrOn(goncurses.A_BOLD)
		scr.MovePrint(firstRow+2, firstColumn, "current game")
		scr.AttrOff(goncurses.A_BOLD)
	}
	scr.MovePrint(firstRow+3, firstColumn, "treasure:")
	scr.MovePrint(firstRow+3, secondColumn, stat.TreasAmt)

	scr.MovePrint(firstRow+4, firstColumn, "level reached:")
	scr.MovePrint(firstRow+4, secondColumn, stat.LevelReached)

	scr.MovePrint(firstRow+5, firstColumn, "food eaten:")
	scr.MovePrint(firstRow+5, secondColumn, stat.FoodEaten)

	scr.MovePrint(firstRow+6, firstColumn, "scrolls read:")
	scr.MovePrint(firstRow+6, secondColumn, stat.ScrollsRead)

	scr.MovePrint(firstRow+7, firstColumn, "elixirs drunk:")
	scr.MovePrint(firstRow+7, secondColumn, stat.ElixirsDrunk)

	enemiesKilled := 0
	for enemyTypeIndex := range len(stat.Kills) {
		enemiesKilled += stat.Kills[enemyTypeIndex].Count
	}
	scr.MovePrint(firstRow+8, firstColumn, "enemies killed:")
	scr.MovePrint(firstRow+8, secondColumn, enemiesKilled)

	scr.MovePrint(firstRow+9, firstColumn, "hits landed:")
	scr.MovePrint(firstRow+9, secondColumn, stat.HitsLanded)

	scr.MovePrint(firstRow+10, firstColumn, "hits missed:")
	scr.MovePrint(firstRow+10, secondColumn, stat.HitsMissed)

	scr.MovePrint(firstRow+11, firstColumn, "hits taken:")
	scr.MovePrint(firstRow+11, secondColumn, stat.HitsTaken)

	scr.MovePrint(firstRow+12, firstColumn, "hits dodged:")
	scr.MovePrint(firstRow+12, secondColumn, stat.HitsDodged)

	scr.MovePrint(firstRow+13, firstColumn, "tiles traveled:")
	scr.MovePrint(firstRow+13, secondColumn, stat.TilesTraveled)

	scr.MovePrint(firstRow+14, firstColumn, "doors opened:")
	scr.MovePrint(firstRow+14, secondColumn, stat.DoorsOpened)
}

var (
	lastPlayerAngle = -1.0
	lastPlayerPosX  = -1.0
	lastPlayerPosY  = -1.0
)

func lerp(a, b, t float64) float64 {
	return b*t + a*(1.0-t)
}

func findClosetAngle(from float64, to float64) float64 {
	more := to + math.Pi*2
	less := to - math.Pi*2
	closest := to

	if math.Abs(more-from) < math.Abs(closest-from) {
		closest = more
	} else if math.Abs(less-from) < math.Abs(closest-from) {
		closest = less
	}

	return closest
}

func renderFirstPersonView(g *domain.GameInfo, fieldScr, miniMapScr *goncurses.Window) {
	if !g.FirstPersonView {
		return
	}

	playerAngle := getPlayerAngle(g)
	playerPosY, playerPosX := float64(g.Player.Pos.Y)+0.5, float64(g.Player.Pos.X)+0.5

	if lastPlayerAngle >= 0.0 && lastPlayerPosX >= 0.0 && lastPlayerPosY >= 0.0 {
		dy := math.Abs(lastPlayerPosY - playerPosY)
		dx := math.Abs(lastPlayerPosX - playerPosX)

		if (dx+dy > 0.0 && dx+dy < 1.1) || lastPlayerAngle != playerAngle {
			closestAngle := findClosetAngle(lastPlayerAngle, playerAngle)
			for i := range 10 {
				t := float64(i) / 10
				angle := lerp(lastPlayerAngle, closestAngle, t)
				x := lerp(lastPlayerPosX, playerPosX, t)
				y := lerp(lastPlayerPosY, playerPosY, t)

				renderFirstPersonFromCoords(g, fieldScr, y, x, angle)
				renderFirstPersonMiniMap(g, miniMapScr)
				time.Sleep(time.Millisecond * 16)
				fieldScr.Refresh()
			}
		}
	}

	lastPlayerAngle = playerAngle
	lastPlayerPosY = playerPosY
	lastPlayerPosX = playerPosX

	renderFirstPersonFromCoords(g, fieldScr, playerPosY, playerPosX, playerAngle)

	renderFirstPersonMiniMap(g, miniMapScr)
}

func renderFirstPersonFromCoords(g *domain.GameInfo, fieldScr *goncurses.Window, playerPosY, playerPosX, playerAngle float64) {
	fieldScr.Erase()
	goncurses.Cursor(0)

	scrMaxY, scrMaxX := fieldScr.MaxYX()

	for x := range scrMaxX {
		rayAngle := playerAngle + (float64(x)/float64(scrMaxX-1)-0.5)*domain.FOV
		vectorY := math.Sin(rayAngle)
		vectorX := math.Cos(rayAngle)
		distance, isHitTile, hitTile, hitTileOffset, hitTilePos, isHitLoot, hitLoot, isHitEnemy, hitEnemy, isHitStairs :=
			DDA(g, playerPosY, playerPosX, vectorY, vectorX)

		// fisheye correction
		distance *= math.Cos(rayAngle - playerAngle)
		if distance < 0.001 {
			distance = 0.001
		}

		wallHeight := int(float64(scrMaxY)/distance + 0.01)
		if wallHeight > scrMaxY {
			wallHeight = scrMaxY
		}
		wallTop := (scrMaxY - wallHeight) / 2
		wallBottom := wallTop + wallHeight

		var color int16
		for y := 0; y < scrMaxY; y++ {
			var ch goncurses.Char
			ch, color = getFirstPersonSymbolAndColor(
				y,
				wallTop,
				wallBottom,
				isHitTile,
				hitTileOffset,
				hitTilePos,
				hitTile,
				isHitLoot,
				hitLoot,
				isHitEnemy,
				hitEnemy,
				isHitStairs,
			)
			fieldScr.ColorOn(color)
			if hitEnemy == nil || !(hitEnemy.Type == domain.Ghost && hitEnemy.Flags.ContainsAny(domain.FlagInvisible)) {
				fieldScr.MoveAddChar(y, x, ch)
			}
			fieldScr.ColorOff(color)
		}
	}

	fieldScr.NoutRefresh()
}

func getPlayerAngle(g *domain.GameInfo) (playerAngle float64) {
	switch g.LookDirection {
	case domain.Up:
		playerAngle = math.Pi + math.Pi/2
	case domain.Down:
		playerAngle = math.Pi / 2
	case domain.Left:
		playerAngle = math.Pi
	case domain.Right:
		playerAngle = 0
	}
	return
}

func DDA(g *domain.GameInfo, startY, startX, vectorY, vectorX float64) (
	distance float64,
	isHit bool,
	hitTile domain.Terrain,
	hitTileOffset float64,
	hitTilePos domain.Position,
	isHitLoot bool,
	hitLoot *domain.Loot,
	isHitEnemy bool,
	hitEnemy *domain.Enemy,
	hitStairs bool,
) {
	curY, curX := int(startY), int(startX)

	// distance to travel 1 square on a named axis
	deltaDistY := math.Abs(1 / vectorY)
	deltaDistX := math.Abs(1 / vectorX)

	var sideDistY, sideDistX float64
	var stepY, stepX int
	if vectorY < 0 {
		sideDistY = (startY - float64(curY)) * deltaDistY
		stepY = -1
	} else {
		sideDistY = (float64(curY+1) - startY) * deltaDistY
		stepY = 1

	}
	if vectorX < 0 {
		sideDistX = (startX - float64(curX)) * deltaDistX
		stepX = -1
	} else {
		sideDistX = (float64(curX+1) - startX) * deltaDistX
		stepX = 1
	}

	// 0 = side, 1 = top/bottom
	side := 0

	for range domain.VisionRange {
		if sideDistY > sideDistX {
			sideDistX += deltaDistX
			curX += stepX
			side = 0
		} else {
			sideDistY += deltaDistY
			curY += stepY
			side = 1
		}
		if curY < 0 || curY >= g.Field.Height {
			break
		}
		if curX < 0 || curX >= g.Field.Width {
			break
		}
		_, hasLootOnTile := g.HasLootOnTile(curY, curX)
		enemyOnTile, hasEnemyOnTile := g.HasEnemyOnTile(curY, curX)
		if hasLootOnTile ||
			(hasEnemyOnTile && enemyOnTile.Type == domain.Mimic) ||
			!g.CheckTileTransparency(curY, curX) ||
			(g.Field.Staircase.Y == curY && g.Field.Staircase.X == curX) {
			hitTile = g.Field.Grid[curY][curX]
			hitTilePos = domain.Position{Y: curY, X: curX}
			isHit = true
			hitLoot, isHitLoot = g.HasLootOnTile(curY, curX)
			hitEnemy, isHitEnemy = g.HasEnemyOnTile(curY, curX)
			hitStairs = (g.Field.Staircase.Y == curY && g.Field.Staircase.X == curX)
			break
		}
	}
	if !isHit {
		distance = float64(domain.VisionRange)
		return
	}

	if side == 0 {
		distance = (float64(curX) - float64(stepX)/2 - startX + 0.5) / vectorX
		hitTileOffset = startY + distance*vectorY - float64(curY)
	} else {
		distance = (float64(curY) - float64(stepY)/2 - startY + 0.5) / vectorY
		hitTileOffset = startX + distance*vectorX - float64(curX)
	}

	return
}

func getFirstPersonSymbolAndColor(
	y int,
	wallTop int,
	wallBottom int,
	isHitTile bool,
	hitTileOffset float64,
	hitTilePos domain.Position,
	hitTile domain.Terrain,
	isHitLoot bool,
	loot *domain.Loot,
	isHitEnemy bool,
	enemy *domain.Enemy,
	isHitStairs bool,
) (ch goncurses.Char, color int16) {
	if y < wallTop {
		ch = '-'
	} else if y < wallBottom {
		if isHitTile {
			switch hitTile.Type {
			case domain.Wall:
				ch, color = getFirstPersonWallTexture(hitTileOffset, hitTilePos)
			case domain.Door:
				ch, color = getFirstPersonDoorTexture(hitTileOffset, hitTile)
			case domain.Nothing:
				ch, color = getFirstPersonVoidTexture(hitTileOffset, hitTilePos)
			}

			if isHitEnemy {
				ch, color = getEnemySymbolAndColor(enemy)
			} else if isHitLoot {
				ch, color = getLootSymbolAndColor(loot)
			} else if isHitStairs {
				ch = '%'
				color = 4
			}
		} else {
			ch = ' '
		}
	} else {
		color = 2
		ch = '='
	}
	return
}

func getFirstPersonWallTexture(hitTileOffset float64, hitTilePos domain.Position) (ch goncurses.Char, color int16) {
	color = 11
	if hitTileOffset < 0.33 {
		ch = '#'
		if (hitTilePos.Y+hitTilePos.X)%2 == 0 {
			ch = '%'
		}
	} else if hitTileOffset < 0.66 {
		ch = '%'
		if (hitTilePos.Y+hitTilePos.X)%2 == 0 {
			ch = '#'
		}
	} else {
		ch = '#'
		if (hitTilePos.Y+hitTilePos.X)%2 == 0 {
			ch = '%'
		}
	}
	return
}

func getFirstPersonDoorTexture(hitTileOffset float64, hitTile domain.Terrain) (ch goncurses.Char, color int16) {
	if hitTileOffset < 0.1 {
		ch = '|'
	} else {
		ch = '+'
	}
	switch hitTile.Door {
	case domain.NoColor:
		color = 12
	case domain.Red:
		color = 6
	case domain.Green:
		color = 8
	case domain.Cyan:
		color = 10
	}
	return
}

func getFirstPersonVoidTexture(hitTileOffset float64, hitTilePos domain.Position) (ch goncurses.Char, color int16) {
	color = 11
	if hitTileOffset < 0.33 {
		ch = '!'
		if (hitTilePos.Y+hitTilePos.X)%2 == 0 {
			ch = '|'
		}
	} else if hitTileOffset < 0.66 {
		ch = '|'
		if (hitTilePos.Y+hitTilePos.X)%2 == 0 {
			ch = '!'
		}
	} else {
		ch = '!'
		if (hitTilePos.Y+hitTilePos.X)%2 == 0 {
			ch = '|'
		}
	}
	return
}

func renderFirstPersonMiniMap(g *domain.GameInfo, miniMapScr *goncurses.Window) {
	if !g.FirstPersonView || !g.FirstPersonMiniMap {
		return
	}

	miniMapScr.Erase()

	scrMaxY, scrMaxX := miniMapScr.MaxYX()
	mapCenter := domain.Position{Y: scrMaxY / 2, X: scrMaxX / 2}
	playerPos := domain.Position{Y: g.Player.Pos.Y, X: g.Player.Pos.X}

	for mapY := range scrMaxY {
		for mapX := range scrMaxX {
			if mapY == 0 || mapY == scrMaxY-1 ||
				mapX == 0 || mapX == scrMaxX-1 {
				miniMapScr.ColorOn(5)
				miniMapScr.MovePrint(mapY, mapX, " ")
				miniMapScr.ColorOff(5)
				continue
			} else {
				fieldPos := getMapRealativeFieldPos(g, playerPos, mapCenter, mapY, mapX)
				if fieldPos.Y < 0 || fieldPos.Y >= g.Field.Height ||
					fieldPos.X < 0 || fieldPos.X >= g.Field.Width {
					continue
				}

				tile := g.Field.Grid[fieldPos.Y][fieldPos.X]

				printTerrainTile(tile, miniMapScr, mapY, mapX)
				stairsPos := domain.Position{Y: g.Field.Staircase.Y, X: g.Field.Staircase.X}
				if tile.Fog == domain.Visible &&
					stairsPos.Y == fieldPos.Y &&
					stairsPos.X == fieldPos.X {
					printStaircase(miniMapScr, domain.Position{Y: mapY, X: mapX})
				}
				printLoot(g, miniMapScr, fieldPos.Y, fieldPos.X, mapY, mapX)
				printEnemy(g, miniMapScr, fieldPos.Y, fieldPos.X, mapY, mapX)
			}
		}
	}
	printPlayer(miniMapScr, mapCenter)

	miniMapScr.NoutRefresh()
}

func getMapRealativeFieldPos(g *domain.GameInfo, playerPos, mapCenter domain.Position, mapY, mapX int) (fieldPos domain.Position) {
	switch g.LookDirection {
	case domain.Up:
		fieldPos.Y = playerPos.Y + (mapY - mapCenter.Y)
		fieldPos.X = playerPos.X + (mapX - mapCenter.X)
	case domain.Down:
		fieldPos.Y = playerPos.Y - (mapY - mapCenter.Y)
		fieldPos.X = playerPos.X - (mapX - mapCenter.X)
	case domain.Left:
		fieldPos.Y = playerPos.Y - (mapX - mapCenter.X)
		fieldPos.X = playerPos.X + (mapY - mapCenter.Y)
	case domain.Right:
		fieldPos.Y = playerPos.Y + (mapX - mapCenter.X)
		fieldPos.X = playerPos.X - (mapY - mapCenter.Y)
	}
	return
}

func printEnemy(g *domain.GameInfo, scr *goncurses.Window, fieldY, fieldX, scrY, scrX int) {
	for i := range len(g.Field.Enemies) {
		enemy := &g.Field.Enemies[i]
		if enemy.Pos.Y != fieldY ||
			enemy.Pos.X != fieldX {
			continue
		}
		if g.Field.Grid[fieldY][fieldX].Fog == domain.Visible {
			ch, color := getEnemySymbolAndColor(enemy)
			scr.ColorOn(color)
			if !(enemy.Type == domain.Ghost && enemy.Flags.ContainsAny(domain.FlagInvisible)) {
				scr.MoveAddChar(scrY, scrX, ch)
			}
			scr.ColorOff(color)
		}
	}
}

func printLoot(g *domain.GameInfo, scr *goncurses.Window, fieldY, fieldX, scrY, scrX int) {
	for i := range len(g.Field.Loot) {
		loot := &g.Field.Loot[i]
		if loot.Pos.Y != fieldY ||
			loot.Pos.X != fieldX {
			continue
		}
		if g.Field.Grid[fieldY][fieldX].Fog == domain.Visible {
			ch, color := getLootSymbolAndColor(loot)
			scr.ColorOn(color)
			scr.MoveAddChar(scrY, scrX, ch)
			scr.ColorOff(color)
		}
	}
}
