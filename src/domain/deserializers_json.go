package domain

import (
	"encoding/json"
	"fmt"
)

func (g *GameInfo) UnmarshalJSON(serializedData []byte) error {
	var err error
	type Alias GameInfo
	var gameInfoTemplate = struct {
		SaveLoader json.RawMessage   `json:"saveLoader,omitempty"`
		Backpack   []json.RawMessage `json:"backpack,omitempty"`
		*Alias
	}{
		Alias: (*Alias)(g),
	}

	if err = json.Unmarshal(serializedData, &gameInfoTemplate); err != nil {
		return err
	}

	if len(gameInfoTemplate.Backpack) == 0 || string(gameInfoTemplate.Backpack[0]) == "null" {
		return nil
	}

	g.Backpack = make([]Item, len(gameInfoTemplate.Backpack))

	for i := range len(g.Backpack) {
		var peek struct {
			Type ItemType `json:"type,omitempty"`
		}
		if err = json.Unmarshal(gameInfoTemplate.Backpack[i], &peek); err != nil {
			return fmt.Errorf("failed to read item type: %w, raw data: %s", err, gameInfoTemplate.Backpack[i])
		}
		switch peek.Type {
		case ItemWeapon:
			var w Weapon
			err = json.Unmarshal(gameInfoTemplate.Backpack[i], &w)
			g.Backpack[i] = w
		case ItemArmor:
			var a Armor
			err = json.Unmarshal(gameInfoTemplate.Backpack[i], &a)
			g.Backpack[i] = a
		case ItemScroll:
			var s Scroll
			err = json.Unmarshal(gameInfoTemplate.Backpack[i], &s)
			g.Backpack[i] = s
		case ItemElixir:
			var e Elixir
			err = json.Unmarshal(gameInfoTemplate.Backpack[i], &e)
			g.Backpack[i] = e
		case ItemTreasure:
			var t Treasure
			err = json.Unmarshal(gameInfoTemplate.Backpack[i], &t)
			g.Backpack[i] = t
		case ItemFood:
			var f Food
			err = json.Unmarshal(gameInfoTemplate.Backpack[i], &f)
			g.Backpack[i] = f
		case ItemKey:
			var k Key
			err = json.Unmarshal(gameInfoTemplate.Backpack[i], &k)
			g.Backpack[i] = k
		default:
			err = fmt.Errorf("unknown item type: %w, raw data: %s", err, gameInfoTemplate.Backpack[i])
		}
		if err != nil {
			return err
		}
	}

	return err
}

func (l *Loot) UnmarshalJSON(serializedData []byte) error {
	var err error
	type Alias Loot
	var lootTemplate = struct {
		Content json.RawMessage `json:"content,omitempty"`
		*Alias
	}{
		Alias: (*Alias)(l),
	}

	if err = json.Unmarshal(serializedData, &lootTemplate); err != nil {
		return err
	}

	if len(lootTemplate.Content) == 0 || string(lootTemplate.Content) == "null" {
		return nil
	}

	var peek struct {
		Type ItemType `json:"type,omitempty"`
	}
	if err = json.Unmarshal(lootTemplate.Content, &peek); err != nil {
		return fmt.Errorf("failed to read item type: %w", err)
	}
	switch peek.Type {
	case ItemWeapon:
		var w Weapon
		err = json.Unmarshal(lootTemplate.Content, &w)
		l.Content = w
	case ItemArmor:
		var a Armor
		err = json.Unmarshal(lootTemplate.Content, &a)
		l.Content = a
	case ItemScroll:
		var s Scroll
		err = json.Unmarshal(lootTemplate.Content, &s)
		l.Content = s
	case ItemElixir:
		var e Elixir
		err = json.Unmarshal(lootTemplate.Content, &e)
		l.Content = e
	case ItemTreasure:
		var t Treasure
		err = json.Unmarshal(lootTemplate.Content, &t)
		l.Content = t
	case ItemFood:
		var f Food
		err = json.Unmarshal(lootTemplate.Content, &f)
		l.Content = f
	case ItemKey:
		var k Key
		err = json.Unmarshal(lootTemplate.Content, &k)
		l.Content = k
	default:
		err = fmt.Errorf("unknown item type: %w", err)
	}
	return err
}
