package domain

func (g *GameInfo) calcFogOfWar() {
	g.resetVisibleTiles()

	pos := g.Player.Pos

	// semi-axes of an ellipse
	a := VisionRange * 2
	b := VisionRange

	y0 := g.Player.Pos.Y
	x0 := g.Player.Pos.X

	for dy := -b; dy <= b; dy++ {
		y := pos.Y + dy
		for dx := -a; dx <= a; dx++ {
			x := pos.X + dx
			if dx*dx*b*b+dy*dy*a*a < a*a*b*b {
				g.rayCast(y0, x0, y, x, func(y, x int) {
					g.Field.Grid[y][x].Fog = Visible
				})
			}
		}
	}
}

func (g *GameInfo) resetVisibleTiles() {
	for y := range g.Field.Height {
		for x := range g.Field.Width {
			if g.Field.Grid[y][x].Fog == Visible {
				g.Field.Grid[y][x].Fog = NotVisible
			}
		}
	}
}
