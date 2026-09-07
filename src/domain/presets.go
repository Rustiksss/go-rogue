package domain

import (
	"math"
	"strconv"
)

const (
	FieldHeight = 22
	FieldWidth  = 80
)

const VisionRange = 8

const MaxLevel = 21

const (
	BaseLootSpawnChance  = 30
	BaseEnemySpawnChance = 45
)

const VampireMaxHealthReduction = 3

// 1 rad = 180 deg = Pi
const FOV = math.Pi / 2

var (
	TerrainVoid = Terrain{
		Type: Nothing,
		Fog:  Unknown,
		Door: NotADoor,
	}
	TerrainGround = Terrain{
		Type: Ground,
		Fog:  Unknown,
		Door: NotADoor,
	}
	TerrainWall = Terrain{
		Type: Wall,
		Fog:  Unknown,
		Door: NotADoor,
	}
	TerrainCorridor = Terrain{
		Type: Corridor,
		Fog:  Unknown,
		Door: NotADoor,
	}
	TerrainDoor = Terrain{
		Type: Door,
		Fog:  Unknown,
		Door: NoColor,
	}
	TerrainDoorRed = Terrain{
		Type: Door,
		Fog:  Unknown,
		Door: Red,
	}
	TerrainDoorGreen = Terrain{
		Type: Door,
		Fog:  Unknown,
		Door: Green,
	}
	TerrainDoorCyan = Terrain{
		Type: Door,
		Fog:  Unknown,
		Door: Cyan,
	}
)

var basePlayer = Player{
	MaxHealth:     100,
	CurHealth:     100,
	Strength:      16,
	Dexterity:     5,
	Armor:         0,
	BuffTimers:    []Buff{},
	WieldedWeapon: Weapon{},
	WornArmor:     Armor{},
}

var (
	dagger = Weapon{
		Type:   ItemWeapon,
		Name:   "dagger(" + strconv.Itoa(16) + ")",
		ReqStr: 16,
		DmgMin: 1,
		DmgMax: 4,
	}
	mace = Weapon{
		Type:   ItemWeapon,
		Name:   "mace(" + strconv.Itoa(17) + ")",
		ReqStr: 17,
		DmgMin: 2,
		DmgMax: 6,
	}
	longsword = Weapon{
		Type:   ItemWeapon,
		Name:   "longsword(" + strconv.Itoa(18) + ")",
		ReqStr: 18,
		DmgMin: 3,
		DmgMax: 8,
	}
	twohandedSword = Weapon{
		Type:   ItemWeapon,
		Name:   "two-handed sword(" + strconv.Itoa(19) + ")",
		ReqStr: 19,
		DmgMin: 4,
		DmgMax: 10,
	}
	divineRapier = Weapon{
		Type:   ItemWeapon,
		Name:   "divine rapier(" + strconv.Itoa(20) + ")",
		ReqStr: 20,
		DmgMin: 5,
		DmgMax: 13,
	}
)

var (
	ringmail = Armor{
		Type:      ItemArmor,
		Name:      "ringmail(" + strconv.Itoa(17) + ")",
		ReqStr:    17,
		Armor:     1,
		Dexterity: 2,
	}
	chainmail = Armor{
		Type:      ItemArmor,
		Name:      "chainmail(" + strconv.Itoa(18) + ")",
		ReqStr:    18,
		Armor:     3,
		Dexterity: 2,
	}
	leatherArmor = Armor{
		Type:      ItemArmor,
		Name:      "leather armor(" + strconv.Itoa(19) + ")",
		ReqStr:    19,
		Armor:     2,
		Dexterity: 4,
	}
	plateArmor = Armor{
		Type:      ItemArmor,
		Name:      "plate armor(" + strconv.Itoa(20) + ")",
		ReqStr:    20,
		Armor:     5,
		Dexterity: 1,
	}
)

var (
	healthScroll = Scroll{
		Type:      ItemScroll,
		Name:      "scroll of health",
		Health:    20,
		Strength:  0,
		Dexterity: 0,
	}
	strengthScroll = Scroll{
		Type:      ItemScroll,
		Name:      "scroll of strength",
		Health:    0,
		Strength:  1,
		Dexterity: 0,
	}
	dexterityScroll = Scroll{
		Type:      ItemScroll,
		Name:      "scroll of dexterity",
		Health:    0,
		Strength:  0,
		Dexterity: 1,
	}
)

var (
	healthElixir = Elixir{
		Type:      ItemElixir,
		Name:      "elixir of health",
		Health:    20,
		Strength:  0,
		Dexterity: 0,
		Buff:      buffInfoHealth,
	}
	strengthElixir = Elixir{
		Type:      ItemElixir,
		Name:      "elixir of strength",
		Health:    0,
		Strength:  1,
		Dexterity: 0,
		Buff:      buffInfoStrength,
	}
	dexterityElixir = Elixir{
		Type:      ItemElixir,
		Name:      "elixir of dexterity",
		Health:    0,
		Strength:  0,
		Dexterity: 1,
		Buff:      buffInfoDexterity,
	}
)

var (
	buffInfoHealth = Buff{
		Attribute: maxHealth,
		Amount:    10,
		TicksLeft: 200,
	}
	buffInfoStrength = Buff{
		Attribute: strength,
		Amount:    1,
		TicksLeft: 200,
	}
	buffInfoDexterity = Buff{
		Attribute: dexterity,
		Amount:    1,
		TicksLeft: 200,
	}
)

var (
	dryRation = Food{
		Type: ItemFood,
		Name: "dry ration",
		Heal: 30,
	}
	pie = Food{
		Type: ItemFood,
		Name: "pie",
		Heal: 60,
	}
)

var (
	redKey = Key{
		Type:  ItemKey,
		Name:  "red key",
		Color: Red,
	}
	greenKey = Key{
		Type:  ItemKey,
		Name:  "green key",
		Color: Green,
	}
	cyanKey = Key{
		Type:  ItemKey,
		Name:  "cyan key",
		Color: Cyan,
	}
)

var (
	ZombieStat = KillsStat{
		Type:  Zombie,
		Count: 0,
	}
	VampireStat = KillsStat{
		Type:  Vampire,
		Count: 0,
	}
	GhostStat = KillsStat{
		Type:  Ghost,
		Count: 0,
	}
	OgreStat = KillsStat{
		Type:  Ogre,
		Count: 0,
	}
	SnakeMageStat = KillsStat{
		Type:  SnakeMage,
		Count: 0,
	}
)

var (
	BaseZombie = Enemy{
		Type:      Zombie,
		Name:      "zombie",
		MaxHealth: 30,
		CurHealth: 30,
		Strength:  16,
		Dexterity: 1,
		Armor:     0,
		Hostility: 7,
	}
	BaseVampire = Enemy{
		Type:      Vampire,
		Name:      "vampire",
		MaxHealth: 30,
		CurHealth: 30,
		Strength:  16,
		Dexterity: 7,
		Armor:     0,
		Hostility: 10,
		Flags:     FlagDodgy,
	}
	BaseGhost = Enemy{
		Type:      Ghost,
		Name:      "ghost",
		MaxHealth: 15,
		CurHealth: 15,
		Strength:  10,
		Dexterity: 10,
		Armor:     0,
		Hostility: 2,
	}
	BaseOgre = Enemy{
		Type:      Ogre,
		Name:      "ogre",
		MaxHealth: 35,
		CurHealth: 35,
		Strength:  21,
		Dexterity: 1,
		Armor:     1,
		Hostility: 7,
	}
	BaseSnakeMage = Enemy{
		Type:      SnakeMage,
		Name:      "snake mage",
		MaxHealth: 15,
		CurHealth: 15,
		Strength:  13,
		Dexterity: 13,
		Armor:     0,
		Hostility: 10,
	}
	BaseMimic = Enemy{
		Type:      Mimic,
		Name:      "mimic",
		MaxHealth: 30,
		CurHealth: 30,
		Strength:  11,
		Dexterity: 8,
		Armor:     0,
		Hostility: 2,
	}
)
