package domain

import (
	"slices"
	"strconv"
)

func (g *GameInfo) stepOnLoot() {
	for i := len(g.Field.Loot) - 1; i >= 0; i-- {
		if g.Field.Loot[i].Pos.Y == g.Player.Pos.Y &&
			g.Field.Loot[i].Pos.X == g.Player.Pos.X {
			if backpackTreasureIndex, has := g.hasTreasure(); has &&
				g.Field.Loot[i].Content.GetType() == ItemTreasure {

				newTreasureAmount := g.Field.Loot[i].Content.(Treasure).Amount
				oldTreasureAmount := g.Backpack[backpackTreasureIndex].(Treasure).Amount
				sumTreasure := Treasure{
					Type:   ItemTreasure,
					Amount: oldTreasureAmount + newTreasureAmount,
				}
				g.Backpack[backpackTreasureIndex] = sumTreasure

				g.Log = append(g.Log, "picked up treasure, amount: "+strconv.Itoa(newTreasureAmount))

				g.Field.Loot = slices.Delete(g.Field.Loot, i, i+1)
			} else if len(g.Backpack) < 9 {
				g.Backpack = append(g.Backpack, g.Field.Loot[i].Content)

				itemName := GetItemName(g.Field.Loot[i].Content)
				g.Log = append(g.Log, "picked up the "+itemName)

				g.Field.Loot = slices.Delete(g.Field.Loot, i, i+1)
			} else {
				g.Log = append(g.Log, "not enough space in the backpack")
			}
		}
	}
	if backpackTreasureIndex, has := g.hasTreasure(); has {
		g.Player.TreasAmt = g.Backpack[backpackTreasureIndex].(Treasure).Amount
		g.Statistic.TreasAmt = g.Player.TreasAmt
	}
}

func (g *GameInfo) hasTreasure() (int, bool) {
	for i := range len(g.Backpack) {
		if g.Backpack[i].GetType() == ItemTreasure {
			return i, true
		}
	}
	return -1, false
}

func GetItemName(item Item) (itemName string) {
	switch item.GetType() {
	case ItemWeapon:
		itemName = item.(Weapon).Name
	case ItemArmor:
		itemName = item.(Armor).Name
	case ItemScroll:
		itemName = item.(Scroll).Name
	case ItemElixir:
		itemName = item.(Elixir).Name
	case ItemTreasure:
		itemName = "treasure, amount: " + strconv.Itoa(item.(Treasure).Amount)
	case ItemFood:
		itemName = item.(Food).Name
	case ItemKey:
		itemName = item.(Key).Name
	}
	return
}

func (g *GameInfo) UsingItemModeTurnOn(action Action) {
	g.UsingItem = true
	switch action {
	case EquipWeapon:
		g.UsingItemType = ItemWeapon
	case EquipArmor:
		g.UsingItemType = ItemArmor
	case EatFood:
		g.UsingItemType = ItemFood
	case QuaffElixir:
		g.UsingItemType = ItemElixir
	case ReadScroll:
		g.UsingItemType = ItemScroll
	case DropItem:
		g.UsingItemType = ItemNoType
	}
}

func (g *GameInfo) useItem() bool {
	performedAction := false
	slot := getUsedSlot(g.Action)

	performedAction = g.chooseUsedItem(slot)

	if performedAction || g.Action == Return {
		g.UsingItem = false
	}
	return performedAction
}

func getUsedSlot(action Action) (slot int) {
	switch action {
	case One:
		slot = 0
	case Two:
		slot = 1
	case Three:
		slot = 2
	case Four:
		slot = 3
	case Five:
		slot = 4
	case Six:
		slot = 5
	case Seven:
		slot = 6
	case Eight:
		slot = 7
	case Nine:
		slot = 8
	case Zero:
		slot = -1
	default:
		slot = -2
	}
	return
}

func (g *GameInfo) chooseUsedItem(slot int) bool {
	performedAction := false
	itemAction, acceptableSlot := g.getItemAction(slot)
	if acceptableSlot {
		performedAction = itemAction(slot)
	}
	return performedAction
}

func (g *GameInfo) getItemAction(slot int) (func(int) bool, bool) {
	acceptableSlot := false
	var itemAction func(int) bool
	itemType := g.UsingItemType
	switch itemType {
	case ItemWeapon:
		if slot >= 0 {
			itemAction = g.equipWeapon
			acceptableSlot = true
		} else if slot == -1 {
			itemAction = g.unequipWeapon
			acceptableSlot = true
		}
	case ItemArmor:
		if slot >= 0 {
			itemAction = g.equipArmor
			acceptableSlot = true
		} else if slot == -1 {
			itemAction = g.unequipArmor
			acceptableSlot = true
		}
	case ItemFood:
		if slot >= 0 {
			itemAction = g.eatFood
			acceptableSlot = true
		}
	case ItemScroll:
		if slot >= 0 {
			itemAction = g.readScroll
			acceptableSlot = true
		}
	case ItemElixir:
		if slot >= 0 {
			itemAction = g.quaffElixir
			acceptableSlot = true
		}
	case ItemNoType:
		if slot >= 0 {
			itemAction = g.dropItem
			acceptableSlot = true
		}
	}
	return itemAction, acceptableSlot
}

func (g *GameInfo) equipWeapon(slot int) bool {
	if len(g.Log) > 0 {
		g.Log = slices.Delete(g.Log, 0, 1)
	}
	if len(g.Backpack) > slot &&
		g.Backpack[slot].GetType() == ItemWeapon {
		if g.Player.HasWeapon {
			currentWeapon := g.Player.WieldedWeapon
			g.Backpack = append(g.Backpack, currentWeapon)
		}
		g.Player.WieldedWeapon = g.Backpack[slot].(Weapon)
		g.Player.HasWeapon = true

		g.Backpack = slices.Delete(g.Backpack, slot, slot+1)

		g.Log = append(g.Log, "equipped the "+g.Player.WieldedWeapon.Name)
		return true
	}
	g.Log = append(g.Log, "this is not a weapon")
	return false
}

func (g *GameInfo) unequipWeapon(slot int) bool {
	if len(g.Log) > 0 {
		g.Log = slices.Delete(g.Log, 0, 1)
	}
	if g.Player.HasWeapon {
		itemName := g.Player.WieldedWeapon.Name
		g.Backpack = append(g.Backpack, g.Player.WieldedWeapon)
		if len(g.Backpack) > 9 {
			g.dropItem(len(g.Backpack) - 1)
		}
		g.Player.HasWeapon = false
		g.Log = slices.Insert(g.Log, 0, "unequipped the "+itemName)
		return true
	}
	g.Log = append(g.Log, "nothing to unequip")
	return false
}

func (g *GameInfo) equipArmor(slot int) bool {
	if len(g.Log) > 0 {
		g.Log = slices.Delete(g.Log, 0, 1)
	}
	if len(g.Backpack) > slot &&
		g.Backpack[slot].GetType() == ItemArmor {
		if g.Player.HasArmor {
			currentArmor := g.Player.WornArmor
			g.Backpack = append(g.Backpack, currentArmor)
		}
		g.Player.WornArmor = g.Backpack[slot].(Armor)
		g.Player.HasArmor = true

		g.Backpack = slices.Delete(g.Backpack, slot, slot+1)
		g.Log = append(g.Log, "equipped the "+g.Player.WornArmor.Name)
		return true
	}
	g.Log = append(g.Log, "this is not an armor")
	return false
}

func (g *GameInfo) unequipArmor(slot int) bool {
	if len(g.Log) > 0 {
		g.Log = slices.Delete(g.Log, 0, 1)
	}
	if g.Player.HasArmor {
		itemName := g.Player.WornArmor.Name
		g.Backpack = append(g.Backpack, g.Player.WornArmor)
		if len(g.Backpack) > 9 {
			g.dropItem(len(g.Backpack) - 1)
		}
		g.Player.HasArmor = false
		g.Log = slices.Insert(g.Log, 0, "unequipped the "+itemName)
		return true
	}
	g.Log = append(g.Log, "nothing to unequip")
	return false
}

func (g *GameInfo) dropItem(slot int) bool {
	if len(g.Log) > 0 {
		g.Log = slices.Delete(g.Log, 0, 1)
	}
	if len(g.Backpack) > slot {
		lootPiece := Loot{
			Content: g.Backpack[slot],
			Pos:     Position{g.Player.Pos.Y, g.Player.Pos.X},
		}
		g.Field.Loot = append(g.Field.Loot, lootPiece)
		g.Backpack = slices.Delete(g.Backpack, slot, slot+1)
		g.Log = append(g.Log, "dropped the "+GetItemName(lootPiece.Content))
		if lootPiece.Content.GetType() == ItemTreasure {
			g.Player.TreasAmt = 0
			g.Statistic.TreasAmt = g.Player.TreasAmt
		}
		return true
	}

	g.Log = append(g.Log, "nothing to drop")
	return false
}

func (g *GameInfo) eatFood(slot int) bool {
	if len(g.Log) > 0 {
		g.Log = slices.Delete(g.Log, 0, 1)
	}
	if len(g.Backpack) > slot &&
		g.Backpack[slot].GetType() == ItemFood {
		food := g.Backpack[slot].(Food)
		g.Player.CurHealth += food.Heal
		if g.Player.CurHealth > g.Player.MaxHealth {
			g.Player.CurHealth = g.Player.MaxHealth
		}
		g.Backpack = slices.Delete(g.Backpack, slot, slot+1)
		g.Statistic.FoodEaten++
		g.Log = append(g.Log, "yum! tasty "+food.Name+", +"+strconv.Itoa(food.Heal)+" health")
		return true
	}
	g.Log = append(g.Log, "this is not a food")
	return false
}

func (g *GameInfo) readScroll(slot int) bool {
	if len(g.Log) > 0 {
		g.Log = slices.Delete(g.Log, 0, 1)
	}
	if len(g.Backpack) > slot &&
		g.Backpack[slot].GetType() == ItemScroll {
		scroll := g.Backpack[slot].(Scroll)
		var text string
		if scroll.Health != 0 {
			g.Player.MaxHealth += scroll.Health
			g.Player.CurHealth += scroll.Health
			text = "you feel sturdier! +" + strconv.Itoa(scroll.Health) + " max health"
		} else if scroll.Strength != 0 {
			g.Player.Strength += scroll.Strength
			text = "you feel stronger! +" + strconv.Itoa(scroll.Strength) + " strength"
		} else if scroll.Dexterity != 0 {
			g.Player.Dexterity += scroll.Dexterity
			text = "you feel more agile! +" + strconv.Itoa(scroll.Dexterity) + " dexterity"
		}
		g.Backpack = slices.Delete(g.Backpack, slot, slot+1)
		g.Statistic.ScrollsRead++
		g.Log = append(g.Log, text)
		return true
	}
	g.Log = append(g.Log, "this is not a scroll")
	return false
}

func (g *GameInfo) quaffElixir(slot int) bool {
	if len(g.Log) > 0 {
		g.Log = slices.Delete(g.Log, 0, 1)
	}
	if len(g.Backpack) > slot &&
		g.Backpack[slot].GetType() == ItemElixir {
		elixir := g.Backpack[slot].(Elixir)
		var text string
		if elixir.Health != 0 {
			g.Player.BuffMaxHealth += elixir.Health
			g.Player.MaxHealth += elixir.Health
			g.Player.CurHealth += elixir.Health
			g.Player.BuffTimers = append(g.Player.BuffTimers, elixir.Buff)
			text = "you feel sturdier! +" + strconv.Itoa(elixir.Health) + " max health for " + strconv.Itoa(elixir.Buff.TicksLeft) + " turns"
		} else if elixir.Strength != 0 {
			g.Player.BuffStrength += elixir.Strength
			g.Player.Strength += elixir.Strength
			g.Player.BuffTimers = append(g.Player.BuffTimers, elixir.Buff)
			text = "you feel stronger! +" + strconv.Itoa(elixir.Strength) + " strength for " + strconv.Itoa(elixir.Buff.TicksLeft) + " turns"
		} else if elixir.Dexterity != 0 {
			g.Player.BuffDexterity += elixir.Dexterity
			g.Player.Dexterity += elixir.Dexterity
			g.Player.BuffTimers = append(g.Player.BuffTimers, elixir.Buff)
			text = "you feel more agile! +" + strconv.Itoa(elixir.Dexterity) + " dexterity for " + strconv.Itoa(elixir.Buff.TicksLeft) + " turns"
		}
		g.Backpack = slices.Delete(g.Backpack, slot, slot+1)
		g.Statistic.ElixirsDrunk++
		g.Log = append(g.Log, text)
		return true
	}
	g.Log = append(g.Log, "this is not an elixir")
	return false
}

func (g *GameInfo) decreaseBuffs() {
	for i := len(g.Player.BuffTimers) - 1; i >= 0; i-- {
		g.Player.BuffTimers[i].TicksLeft--
		if g.Player.BuffTimers[i].TicksLeft == 0 {
			switch g.Player.BuffTimers[i].Attribute {
			case maxHealth:
				g.Player.BuffMaxHealth -= g.Player.BuffTimers[i].Amount
				g.Player.MaxHealth -= g.Player.BuffTimers[i].Amount
				g.Player.CurHealth -= g.Player.BuffTimers[i].Amount
				if g.Player.CurHealth < 1 {
					g.Player.CurHealth = 1
				}
				g.Log = append(g.Log, "the max health buff expires!")
			case strength:
				g.Player.BuffStrength -= g.Player.BuffTimers[i].Amount
				g.Player.Strength -= g.Player.BuffTimers[i].Amount
				g.Log = append(g.Log, "the strength buff expires!")
			case dexterity:
				g.Player.BuffDexterity -= g.Player.BuffTimers[i].Amount
				g.Player.Dexterity -= g.Player.BuffTimers[i].Amount
				g.Log = append(g.Log, "the dexterity buff expires!")
			}
			g.Player.BuffTimers = slices.Delete(g.Player.BuffTimers, i, i+1)
		}
	}
}

func (g *GameInfo) deleteKeys() {
	for i := len(g.Backpack) - 1; i >= 0; i-- {
		if g.Backpack[i].GetType() == ItemKey {
			g.Backpack = slices.Delete(g.Backpack, i, i+1)
		}
	}
}

func (g *GameInfo) recalculatePlayerDamage() {
	g.Player.DmgMin = (3 * g.Player.Strength) / 4
	g.Player.DmgMax = g.Player.Strength + (g.Player.Strength / 4)

	if g.Player.HasWeapon {
		strInsufficiency := g.Player.WieldedWeapon.ReqStr - g.Player.Strength
		if strInsufficiency < 0 {
			strInsufficiency = 0
		}

		g.Player.EquipDmgMin = g.Player.WieldedWeapon.DmgMin - (2 * strInsufficiency)
		g.Player.EquipDmgMax = g.Player.WieldedWeapon.DmgMax - (3 * strInsufficiency)
	} else {
		g.Player.EquipDmgMin = 0
		g.Player.EquipDmgMax = 0
	}

	g.Player.DmgMin += g.Player.EquipDmgMin
	g.Player.DmgMax += g.Player.EquipDmgMax
}

func (g *GameInfo) recalculatePlayerArmorAndDexterity() {
	g.Player.Armor -= g.Player.EquipArmor
	g.Player.Dexterity -= g.Player.EquipDexterity

	if g.Player.HasArmor {
		strInsufficiency := g.Player.WornArmor.ReqStr - g.Player.Strength
		if strInsufficiency < 0 {
			strInsufficiency = 0
		}

		g.Player.EquipArmor = g.Player.WornArmor.Armor - strInsufficiency
		g.Player.EquipDexterity = g.Player.WornArmor.Dexterity - (2 * strInsufficiency)
	} else {
		g.Player.EquipArmor = 0
		g.Player.EquipDexterity = 0
	}

	g.Player.Armor += g.Player.EquipArmor
	g.Player.Dexterity += g.Player.EquipDexterity
}
