package domain

import (
	"fmt"
	"math/rand/v2"
)

type DefaultEnemyBehaviors struct{}

type ZombieBehaviors struct{ DefaultEnemyBehaviors }
type VampireBehaviors struct{ DefaultEnemyBehaviors }
type GhostBehaviors struct{ DefaultEnemyBehaviors }
type OgreBehaviors struct{ DefaultEnemyBehaviors }
type SnakeMageBehaviors struct{ DefaultEnemyBehaviors }
type MimicBehaviors struct{ DefaultEnemyBehaviors }

func (enemy *Enemy) Behaviors() EnemyBehaviors {
	switch enemy.Type {
	case Zombie:
		return ZombieBehaviors{}
	case Vampire:
		return VampireBehaviors{}
	case Ghost:
		return GhostBehaviors{}
	case Ogre:
		return OgreBehaviors{}
	case SnakeMage:
		return SnakeMageBehaviors{}
	case Mimic:
		return MimicBehaviors{}
	}

	panic(fmt.Sprintf("behaviors not defined for enemy type %v", enemy.Type))
}

func (VampireBehaviors) CalcHitEnemyChance(enemy *Enemy, playerDexterity int) int {
	if enemy.Flags.ContainsAny(FlagDodgy) {
		enemy.Flags.Remove(FlagDodgy)
		return 0
	}

	return DefaultEnemyBehaviors{}.CalcHitEnemyChance(enemy, playerDexterity)
}

func (DefaultEnemyBehaviors) OnPlayerHit(g *GameInfo, enemyIndex int) {}

func (VampireBehaviors) OnPlayerHit(g *GameInfo, enemyIndex int) {
	g.Player.MaxHealth -= VampireMaxHealthReduction
	g.Player.MaxHealth = max(g.Player.MaxHealth, 1)
	g.Player.CurHealth = min(g.Player.CurHealth, g.Player.MaxHealth)
}

func (DefaultEnemyBehaviors) OnEnemyHit(g *GameInfo, enemyIndex int) {}

func (OgreBehaviors) OnEnemyHit(g *GameInfo, enemyIndex int) {
	enemy := &g.Field.Enemies[enemyIndex]
	if enemy.Flags.ContainsAny(FlagAccurate) {
		return
	}

	enemy.Flags.Add(FlagSkipTurn)
	enemy.Flags.Add(FlagAccurate)
}

func (OgreBehaviors) CalcHitPlayerChance(g *GameInfo, enemyIndex int) int {
	enemy := &g.Field.Enemies[enemyIndex]

	if enemy.Flags.ContainsAny(FlagAccurate) {
		return 100
	}

	return DefaultEnemyBehaviors{}.CalcHitPlayerChance(g, enemyIndex)
}

func (SnakeMageBehaviors) OnPlayerHit(g *GameInfo, enemyIndex int) {
	if rand.IntN(4) == 0 {
		g.SkipInput = true
		g.Log = append(g.Log, "you have been put to sleep for one turn")
	}
}
