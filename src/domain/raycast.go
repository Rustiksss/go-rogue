package domain

func (g *GameInfo) rayCast(y0, x0, y1, x1 int, visitCell func(y, x int)) {
	if y1 < 0 || y1 >= g.Field.Height ||
		x1 < 0 || x1 >= g.Field.Width {
		return
	}

	dirY := 1
	if y1 < y0 {
		dirY = -1
	}
	dirX := 1
	if x1 < x0 {
		dirX = -1
	}
	dy := (y1 - y0) * dirY
	dx := (x1 - x0) * dirX
	if dx > dy {
		g.rayCastHoriz(y0, x0, dy, dx, dirY, dirX, visitCell)
	} else {
		g.rayCastVert(y0, x0, dy, dx, dirY, dirX, visitCell)
	}
}

func (g *GameInfo) rayCastHoriz(y0, x0, dy, dx, dirY, dirX int, visitCell func(y, x int)) {
	y := y0
	x := x0
	d := 2*dy - dx
	for i := range dx + 1 {
		visitCell(y, x)
		if i != 0 && !g.CheckTileTransparency(y, x) {
			break
		}
		if d > 0 {
			y += dirY
			d -= 2 * dx
		}
		d += 2 * dy
		x += dirX
	}
}

func (g *GameInfo) rayCastVert(y0, x0, dy, dx, dirY, dirX int, visitCell func(y, x int)) {
	y := y0
	x := x0
	d := 2*dx - dy
	for i := range dy + 1 {
		visitCell(y, x)
		if i != 0 && !g.CheckTileTransparency(y, x) {
			break
		}
		if d > 0 {
			x += dirX
			d -= 2 * dy
		}
		d += 2 * dx
		y += dirY
	}
}

func (g *GameInfo) CheckTileTransparency(y, x int) (transparent bool) {
	tile := g.Field.Grid[y][x].Type
	transparent = true
	enemy, hasEnemyOnTile := g.HasEnemyOnTile(y, x)
	_, hasLootOnTile := g.HasLootOnTile(y, x)
	if tile == Wall || tile == Nothing ||
		(tile == Door && !hasLootOnTile) ||
		hasEnemyOnTile {

		transparent = false
		if hasEnemyOnTile && (enemy.Flags.ContainsAny(FlagInvisible) || enemy.Flags.ContainsAny(FlagMimicking)) {
			transparent = true
		}
	}
	return
}

func (g *GameInfo) HasEnemyOnTile(y, x int) (enemy *Enemy, has bool) {
	pos := Position{y, x}
	for i := range len(g.Field.Enemies) {
		if g.Field.Enemies[i].Pos == pos {
			has = true
			enemy = &g.Field.Enemies[i]
			return
		}
	}
	return
}

func (g *GameInfo) HasLootOnTile(y, x int) (loot *Loot, has bool) {
	pos := Position{y, x}
	for i := range len(g.Field.Loot) {
		if g.Field.Loot[i].Pos == pos {
			has = true
			loot = &g.Field.Loot[i]
			return
		}
	}
	return
}

func abs(num int) int {
	if num < 0 {
		return num * -1
	}
	return num
}
