package domain

import "math/rand/v2"

func randomRoam(g *GameInfo, enemyIndex int) error {
	enemy := &g.Field.Enemies[enemyIndex]

	pos := enemy.Pos
	switch rand.IntN(4) {
	case 0:
		pos.X--
	case 1:
		pos.X++
	case 2:
		pos.Y--
	case 3:
		pos.Y++
	}

	return g.enemyMove(enemyIndex, pos)
}

func (ZombieBehaviors) Roam(g *GameInfo, enemyIndex int) error {
	enemy := &g.Field.Enemies[enemyIndex]

	pos := enemy.Pos
	if rand.IntN(2) == 0 {
		pos.X--
	} else {
		pos.X++
	}

	return g.enemyMove(enemyIndex, pos)
}

func (VampireBehaviors) Roam(g *GameInfo, enemyIndex int) error {
	return randomRoam(g, enemyIndex)
}

func (GhostBehaviors) Roam(g *GameInfo, enemyIndex int) error {
	enemy := &g.Field.Enemies[enemyIndex]

	enemy.Flags.Remove(FlagInvisible)
	if (enemy.RoamingCounter/3)%2 == 0 {
		enemy.Flags.Add(FlagInvisible)
	}

	dx := rand.IntN(15) - 7
	dy := rand.IntN(7) - 3

	randPos := Position{enemy.Pos.Y + dy, enemy.Pos.X + dx}
	newPos := enemy.Pos

	g.rayCast(enemy.Pos.Y, enemy.Pos.X, randPos.Y, randPos.X, func(y, x int) {
		pos := Position{y, x}
		// don't teleport onto character
		if g.Player.Pos == pos || !g.checkTerrain(pos) {
			return
		}

		newPos = pos
	})

	return g.enemyMove(enemyIndex, newPos)
}

func (GhostBehaviors) Chase(g *GameInfo, enemyIndex int) error {
	enemy := &g.Field.Enemies[enemyIndex]
	enemy.Flags.Remove(FlagInvisible)

	return DefaultEnemyBehaviors{}.Chase(g, enemyIndex)
}

// move twice, but if the first movement ends up attacking the player,
// cancel the second movement, so as to not attack twice
func moveTwice(g *GameInfo, enemyIndex int, move func(*GameInfo, int) error) (err error) {
	enemy := &g.Field.Enemies[enemyIndex]

	err = move(g, enemyIndex)
	if err == nil && !enemy.Flags.ContainsAny(FlagJustAttacked) {
		err = move(g, enemyIndex)
	}

	return
}

func (OgreBehaviors) Roam(g *GameInfo, enemyIndex int) error {
	enemy := &g.Field.Enemies[enemyIndex]
	defer enemy.Flags.Remove(FlagAccurate)

	return moveTwice(g, enemyIndex, randomRoam)
}

func (OgreBehaviors) Chase(g *GameInfo, enemyIndex int) error {
	enemy := &g.Field.Enemies[enemyIndex]
	defer enemy.Flags.Remove(FlagAccurate)

	return moveTwice(g, enemyIndex, DefaultEnemyBehaviors{}.Chase)
}

func (SnakeMageBehaviors) Roam(g *GameInfo, enemyIndex int) error {
	enemy := &g.Field.Enemies[enemyIndex]
	diagonalRoam := func(g *GameInfo, enemyIndex int) error {
		pos := enemy.Pos

		coord := &pos.X
		if (pos.X+pos.Y)%2 == 0 {
			coord = &pos.Y
		}

		if rand.IntN(2) == 0 {
			*coord--
		} else {
			*coord++
		}

		if g.Field.Grid[pos.Y][pos.X].Type != Ground {
			return nil
		}

		return g.enemyMove(enemyIndex, pos)
	}

	return moveTwice(g, enemyIndex, diagonalRoam)
}

func (SnakeMageBehaviors) Chase(g *GameInfo, enemyIndex int) error {
	return moveTwice(g, enemyIndex, DefaultEnemyBehaviors{}.Chase)
}

func (SnakeMageBehaviors) GetTileNeighbors(pos Position) []Position {
	if (pos.X+pos.Y)%2 == 0 {
		return []Position{
			{pos.Y - 1, pos.X},
			{pos.Y + 1, pos.X},
		}
	}

	return []Position{
		{pos.Y, pos.X - 1},
		{pos.Y, pos.X + 1},
	}
}

func (MimicBehaviors) Roam(g *GameInfo, enemyIndex int) error {
	enemy := &g.Field.Enemies[enemyIndex]
	num := rand.IntN(6)
	if enemy.Flags.ContainsAny(FlagMimicking) {
		return nil
	}
	switch num {
	case 0:
		enemy.Flags.Add(FlagMimickingWeapon)
	case 1:
		enemy.Flags.Add(FlagMimickingArmor)
	case 2:
		enemy.Flags.Add(FlagMimickingScroll)
	case 3:
		enemy.Flags.Add(FlagMimickingElixir)
	case 4:
		enemy.Flags.Add(FlagMimickingFood)
	case 5:
		enemy.Flags.Add(FlagMimickingTreasure)
	}
	return nil
}

func (MimicBehaviors) Chase(g *GameInfo, enemyIndex int) error {
	enemy := &g.Field.Enemies[enemyIndex]
	dY, dX := abs(g.Player.Pos.Y-enemy.Pos.Y), abs(g.Player.Pos.X-enemy.Pos.X)
	if dY+dX == 1 || !enemy.Flags.ContainsAny(FlagMimicking) {
		enemy.Flags.Remove(FlagMimickingWeapon | FlagMimickingArmor |
			FlagMimickingScroll | FlagMimickingElixir |
			FlagMimickingFood | FlagMimickingTreasure)

		return DefaultEnemyBehaviors{}.Chase(g, enemyIndex)
	}
	return nil
}
