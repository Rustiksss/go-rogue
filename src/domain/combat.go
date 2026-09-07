package domain

import (
	"math"
	"math/rand"
	"slices"
)

func (g *GameInfo) playerAttackEnemy(enemyIndex int) {
	damage := calcPlayerDmg(g.Player.DmgMin, g.Player.DmgMax)
	enemy := &g.Field.Enemies[enemyIndex]
	isHit := enemy.rollToHitEnemy(g.Player.Dexterity)
	if isHit {
		enemy.CurHealth -= damage - enemy.Armor
		g.Log = append(g.Log, "hit the "+enemy.Name)
		g.Statistic.HitsLanded++

		enemy.Behaviors().OnEnemyHit(g, enemyIndex)
	} else {
		g.Log = append(g.Log, "miss the "+enemy.Name)
		g.Statistic.HitsMissed++
	}

	// after this function, enemyIndex is invalidated
	g.checkEnemyDeath(enemyIndex)
}

func calcPlayerDmg(dmgMin, dmgMax int) int {
	dmgRange := dmgMax - dmgMin + 1
	return dmgMin + rand.Intn(dmgRange)
}

func (enemy *Enemy) rollToHitEnemy(playerDexterity int) bool {
	roll := rand.Intn(100)
	return enemy.Behaviors().CalcHitEnemyChance(enemy, playerDexterity) > roll
}

func (DefaultEnemyBehaviors) CalcHitEnemyChance(enemy *Enemy, playerDexterity int) int {
	evasCoeff := 0.95
	return int(math.Round(100.0 * math.Pow(evasCoeff, float64(enemy.Dexterity-playerDexterity))))
}

func (g *GameInfo) checkEnemyDeath(enemyIndex int) {
	enemy := &g.Field.Enemies[enemyIndex]
	if enemy.CurHealth <= 0 {
		for i := range len(g.Statistic.Kills) {
			if g.Statistic.Kills[i].Type == enemy.Type {
				g.Log = slices.Delete(g.Log, len(g.Log)-1, len(g.Log))
				g.Log = append(g.Log, "killed the "+enemy.Name+"!")
				g.Statistic.Kills[i].Count++
				break
			}
		}
		g.spawnTresureOnDeath(enemyIndex)
		g.Field.Enemies = slices.Delete(g.Field.Enemies, enemyIndex, enemyIndex+1)
	}
}

func (g *GameInfo) spawnTresureOnDeath(enemyIndex int) {
	enemy := &g.Field.Enemies[enemyIndex]
	treasAmt := enemy.MaxHealth / 4
	treasAmt += enemy.Strength / 2
	treasAmt += enemy.Dexterity / 2
	treasAmt += enemy.Hostility / 2
	g.Field.Loot = append(g.Field.Loot, Loot{Treasure{ItemTreasure, treasAmt}, enemy.Pos})
}

func (g *GameInfo) enemyAttackPlayer(enemyIndex int) (err error) {
	damage := g.calcEnemyDmg(enemyIndex)
	enemy := &g.Field.Enemies[enemyIndex]
	isHit := g.rollToHitPlayer(enemyIndex)
	if isHit {
		g.Player.CurHealth -= damage - g.Player.Armor
		g.Log = append(g.Log, "take a hit from the "+enemy.Name)
		g.Statistic.HitsTaken++
		enemy.Behaviors().OnPlayerHit(g, enemyIndex)
	} else {
		g.Log = append(g.Log, "dodge a hit from the "+enemy.Name)
		g.Statistic.HitsDodged++
	}
	enemy.Flags.Add(FlagJustAttacked)

	err = g.checkCharacterDeath(enemyIndex)
	return
}

func (g *GameInfo) calcEnemyDmg(enemyIndex int) int {
	enemy := &g.Field.Enemies[enemyIndex]
	dmgMin := (3 * enemy.Strength) / 4
	dmgMax := enemy.Strength + (enemy.Strength / 4)
	dmgRange := dmgMax - dmgMin + 1
	return dmgMin + rand.Intn(dmgRange)
}

func (g *GameInfo) rollToHitPlayer(enemyIndex int) bool {
	enemy := &g.Field.Enemies[enemyIndex]
	roll := rand.Intn(100)

	return enemy.Behaviors().CalcHitPlayerChance(g, enemyIndex) > roll
}

func (DefaultEnemyBehaviors) CalcHitPlayerChance(g *GameInfo, enemyIndex int) int {
	enemy := &g.Field.Enemies[enemyIndex]
	evasCoeff := 0.95
	return int(math.Round(100.0 * math.Pow(evasCoeff, float64(g.Player.Dexterity-enemy.Dexterity))))
}

func (g *GameInfo) checkCharacterDeath(enemyIndex int) (err error) {
	if g.Player.CurHealth <= 0 {
		g.Statistic.IsDead = true
		g.Statistic.EnemyWhoKilled = g.Field.Enemies[enemyIndex]
		g.GameOver = true

		g.Leaderboard = append(g.Leaderboard, g.Statistic)
		g.sortLeaderboard()
		err = g.SaveLoader.SaveLeaderboard(g.Leaderboard)

		g.calcCurGameTop()

		err = g.SaveLoader.DeleteSavedGame()
		g.SkipInput = false
	}
	return
}
