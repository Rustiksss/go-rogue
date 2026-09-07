package domain

import "math/rand/v2"

// WeightedItem associates an item with its weight (drop chance)
type WeightedItem struct {
	Item   Item
	Weight int
}

// the higher the weight, the higher the chance of a drop
var lootTable = []WeightedItem{
	{dagger, 9}, {mace, 7}, {longsword, 5}, {twohandedSword, 3}, {divineRapier, 1},
	{ringmail, 7}, {chainmail, 5}, {leatherArmor, 3}, {plateArmor, 1},
	{healthScroll, 50}, {strengthScroll, 25}, {dexterityScroll, 35},
	{healthElixir, 90}, {strengthElixir, 45}, {dexterityElixir, 55},
	{dryRation, 110}, {pie, 70},
}

var enemyList = []Enemy{BaseZombie, BaseVampire, BaseGhost, BaseOgre, BaseSnakeMage, BaseMimic}

func chooseLoot() Item {
	totalLootWeight := 0
	for i := range len(lootTable) {
		totalLootWeight += lootTable[i].Weight
	}

	roll := rand.IntN(totalLootWeight)

	currentWeight := 0
	for _, entry := range lootTable {
		currentWeight += entry.Weight
		if roll < currentWeight {
			return entry.Item
		}
	}

	return Treasure{} // technically unreachable
}

func chooseEnemy() Enemy {
	return enemyList[rand.IntN(len(enemyList))]
}

func (g *GameInfo) TrySpawnEnemyOrLoot(pos Position) {
	lootSpawnChance := BaseLootSpawnChance - int(g.Difficulty*10)
	enemySpawnChance := BaseEnemySpawnChance + int(g.Difficulty*10)
	roll := rand.IntN(1000)
	switch {
	case roll <= lootSpawnChance:
		g.Field.Loot = append(g.Field.Loot, Loot{chooseLoot(), pos})
	case roll <= lootSpawnChance+enemySpawnChance:
		enemy := chooseEnemy()
		enemyHealth := enemy.MaxHealth + int(g.Difficulty*3)
		enemy.CurHealth, enemy.MaxHealth = enemyHealth, enemyHealth
		enemy.Strength += int(g.Difficulty)
		enemy.Dexterity += int(g.Difficulty)
		enemy.Pos = pos
		g.Field.Enemies = append(g.Field.Enemies, enemy)
	}
}
