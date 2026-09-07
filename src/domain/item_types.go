package domain

type ItemType int

const (
	ItemNoType ItemType = iota
	ItemWeapon
	ItemArmor
	ItemScroll
	ItemElixir
	ItemTreasure
	ItemFood
	ItemKey
)

type Weapon struct {
	Type   ItemType `json:"type,omitempty"`
	Name   string   `json:"name,omitempty"`
	ReqStr int      `json:"reqStr,omitempty"` // Determines the strength required to wield the weapon
	DmgMin int      `json:"dmgMin,omitempty"`
	DmgMax int      `json:"dmgMax,omitempty"`
}

func (item Weapon) GetType() ItemType {
	return ItemWeapon
}

type Armor struct {
	Type      ItemType `json:"type,omitempty"`
	Name      string   `json:"name,omitempty"`
	ReqStr    int      `json:"reqStr,omitempty"` // Determines the strength required to wear armor.
	Armor     int      `json:"armor,omitempty"`
	Dexterity int      `json:"dexterity,omitempty"`
}

func (item Armor) GetType() ItemType {
	return ItemArmor
}

type Scroll struct {
	Type      ItemType `json:"type,omitempty"`
	Name      string   `json:"name,omitempty"`
	Health    int      `json:"health,omitempty"`
	Strength  int      `json:"strength,omitempty"`
	Dexterity int      `json:"dexterity,omitempty"`
}

func (item Scroll) GetType() ItemType {
	return ItemScroll
}

type Elixir struct {
	Type      ItemType `json:"type,omitempty"`
	Name      string   `json:"name,omitempty"`
	Health    int      `json:"health,omitempty"`
	Strength  int      `json:"strength,omitempty"`
	Dexterity int      `json:"dexterity,omitempty"`
	Buff      Buff     `json:"buff"`
}

func (item Elixir) GetType() ItemType {
	return ItemElixir
}

// Treasures are stored in one slot
type Treasure struct {
	Type   ItemType `json:"type,omitempty"`
	Amount int      `json:"amount,omitempty"`
}

func (item Treasure) GetType() ItemType {
	return ItemTreasure
}

type Food struct {
	Type ItemType `json:"type,omitempty"`
	Name string   `json:"name,omitempty"`
	Heal int      `json:"heal,omitempty"`
}

func (item Food) GetType() ItemType {
	return ItemFood
}

type Key struct {
	Type  ItemType  `json:"type,omitempty"`
	Name  string    `json:"name,omitempty"`
	Color DoorColor `json:"color,omitempty"`
}

func (item Key) GetType() ItemType {
	return ItemKey
}
