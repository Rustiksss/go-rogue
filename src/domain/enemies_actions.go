package domain

func (g *GameInfo) enemiesAction() (err error) {
	for i := range len(g.Field.Enemies) {
		err = g.enemyAction(i)
		if err != nil {
			return
		}
	}
	return
}

func (g *GameInfo) enemyAction(enemyIndex int) (err error) {
	enemy := &g.Field.Enemies[enemyIndex]

	enemy.Flags.Remove(FlagJustAttacked)

	if enemy.Flags.ContainsAny(FlagSkipTurn) {
		enemy.Flags.Remove(FlagSkipTurn)
		return
	}

	switch enemy.AiState {
	case AiRoaming:
		err = enemy.Behaviors().Roam(g, enemyIndex)
		enemy.RoamingCounter++
		if g.checkEnemySeesPlayer(enemy) {
			enemy.AiState = AiChasing
		}
	case AiChasing:
		err = enemy.Behaviors().Chase(g, enemyIndex)
	}

	return
}

func (g *GameInfo) checkEnemySeesPlayer(enemy *Enemy) bool {
	y0 := enemy.Pos.Y
	x0 := enemy.Pos.X
	y := g.Player.Pos.Y
	x := g.Player.Pos.X

	dy := (y - y0) * 2
	dx := x - x0
	distanceSq := dx*dx + dy*dy

	maxDist := enemy.Hostility * 2
	maxDistSq := maxDist * maxDist

	if distanceSq >= maxDistSq {
		return false
	}

	seesPlayer := false
	g.rayCast(y0, x0, y, x, func(y, x int) {
		pos := Position{y, x}
		if g.Player.Pos == pos {
			seesPlayer = true
		}
	})

	return seesPlayer
}

func (DefaultEnemyBehaviors) Chase(g *GameInfo, enemyIndex int) (err error) {
	enemy := &g.Field.Enemies[enemyIndex]
	pathPosition := g.pathfind(enemy, enemy.Pos, g.Player.Pos)

	if pathPosition == enemy.Pos {
		enemy.AiState = AiRoaming
		return
	}

	err = g.enemyMove(enemyIndex, pathPosition)

	return
}

func (g *GameInfo) enemyMove(enemyIndex int, newPos Position) (err error) {
	enemy := &g.Field.Enemies[enemyIndex]

	if !g.checkTerrain(newPos) {
		return
	}

	if _, hasOtherEnemy := g.checkEnemy(newPos); hasOtherEnemy {
		return
	}

	if g.Player.Pos == newPos {
		err = g.enemyAttackPlayer(enemyIndex)
	} else {
		enemy.Pos = newPos
	}
	return
}
