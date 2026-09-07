package domain

var colorToKey = map[DoorColor]Key{
	Red:   redKey,
	Green: greenKey,
	Cyan:  cyanKey,
}

func (g *GameInfo) buildRoomWalls(rect roomRect) {
	for j := range rect.width + 2 {
		x := rect.left + j - 1
		g.Field.Grid[rect.top-1][x] = TerrainWall
		g.Field.Grid[rect.top+rect.height][x] = TerrainWall
	}

	for i := range rect.height {
		y := rect.top + i
		g.Field.Grid[y][rect.left-1] = TerrainWall
		g.Field.Grid[y][rect.left+rect.width] = TerrainWall
	}
}

func (g *GameInfo) buildRoom(layout *levelLayout, roomId roomIndex) {
	room := &layout.rooms[roomId]
	rect := room.rect

	iterateRectPostions(rect, func(y, x int) {
		g.Field.Grid[y][x] = TerrainGround
	})
	g.buildRoomWalls(rect)

	positions := generateRandomPositions(rect)

	if roomId == layout.entrance {
		g.Player.Pos, positions = positions[0], positions[1:]
	}

	if roomId == layout.exit {
		g.Field.Staircase, positions = positions[0], positions[1:]
	}

	for _, keyColor := range room.keys {
		var keyPos Position
		keyPos, positions = positions[0], positions[1:]

		g.Field.Loot = append(g.Field.Loot, Loot{colorToKey[keyColor], keyPos})
	}

	if layout.entrance == roomId {
		return
	}

	for _, pos := range positions {
		g.TrySpawnEnemyOrLoot(pos)
	}
}

func (g *GameInfo) buildCorridor(corridor *corridor) {
	// access field grid in by main and cross axis instead of x and y
	var setCell func(main int, cross int, terrain Terrain)

	if corridor.isHorizontal() {
		setCell = func(main, cross int, terrain Terrain) { g.Field.Grid[cross][main] = terrain }
	} else {
		setCell = func(main, cross int, terrain Terrain) { g.Field.Grid[main][cross] = terrain }
	}

	bends := corridor.bends
	cross := corridor.crossAxisStart

	startDoor := TerrainDoor
	startDoor.Door = corridor.door1
	setCell(corridor.mainAxisStart, cross, startDoor)

	for main := corridor.mainAxisStart + 1; main < corridor.mainAxisEnd-1; main++ {
		if len(bends) > 0 && main == bends[0].mainAxisPos {
			crossTarget := bends[0].crossAxisTarget
			bends = bends[1:]

			var step int
			if crossTarget < cross {
				step = -1
			} else {
				step = 1
			}

			for ; cross != crossTarget; cross += step {
				setCell(main, cross, TerrainCorridor)
			}
		}
		setCell(main, cross, TerrainCorridor)
	}

	endDoor := TerrainDoor
	endDoor.Door = corridor.door2
	setCell(corridor.mainAxisEnd-1, cross, endDoor)
}

func (g *GameInfo) newGameMap() {
	g.Field = GameMap{
		Grid:      [FieldHeight][FieldWidth]Terrain{},
		Staircase: Position{},
		Loot:      []Loot{},
		Height:    FieldHeight,
		Width:     FieldWidth,
	}

	layout := generateLayout()
	for index := range layout.rooms {
		g.buildRoom(&layout, roomIndex(index))
	}

	for index := range layout.corridors {
		g.buildCorridor(&layout.corridors[index])
	}

	g.calcFogOfWar()
}
