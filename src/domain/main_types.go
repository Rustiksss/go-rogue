package domain

type GameInfo struct {
	Field              GameMap      `json:"field"`
	Level              int          `json:"level,omitempty"`
	Player             Player       `json:"player"`
	Action             Action       `json:"action"`
	Difficulty         float64      `json:"difficulty"`         // the difficulty level of the game, the lower the easier
	Backpack           []Item       `json:"backpack,omitempty"` // the backpack holds up to 9 items
	OpenBackpack       bool         `json:"openBackpack,omitempty"`
	UsingItem          bool         `json:"usingItem,omitempty"` // select item mode for use from backpack
	UsingItemType      ItemType     `json:"usingItemType,omitempty"`
	Log                []string     `json:"log,omitempty"`
	Statistic          Statistics   `json:"statistic"`
	CurGameTop         int          `json:"curGameTop,omitempty"`  // current game number in the top
	Leaderboard        []Statistics `json:"leaderboard,omitempty"` // statistics for all games (not including this one)
	SaveLoader         SaveLoad     `json:"saveLoader,omitempty"`
	StartScreen        bool         `json:"startScreen,omitempty"`
	WaitLoadChoice     bool         `json:"waitLoadChoice,omitempty"`
	WaitNameInput      bool         `json:"waitNameInput,omitempty"`
	GameOver           bool         `json:"gameOver,omitempty"`
	Victory            bool         `json:"victory,omitempty"`
	LeaderboardScr     bool         `json:"leaderboardScr,omitempty"`
	HelpScreen         bool         `json:"helpScreen,omitempty"`
	FirstPersonView    bool         `json:"firstPersonView,omitempty"`
	FirstPersonMiniMap bool         `json:"firstPersonMiniMap,omitempty"`
	LookDirection      Direction    `json:"lookDirection,omitempty"` // first-person view direction
	SkipInput          bool         `json:"skipInput,omitempty"`
	Quit               bool         `json:"quit,omitempty"`
}

type GameMap struct {
	Grid      [FieldHeight][FieldWidth]Terrain `json:"grid,omitempty"`
	Staircase Position                         `pos:"x,omitempty"`
	Loot      []Loot                           `json:"loot,omitempty"`
	Enemies   []Enemy                          `json:"enemies,omitempty"`
	Height    int                              `json:"height,omitempty"`
	Width     int                              `json:"width,omitempty"`
}

type Terrain struct {
	Type TerrainType `json:"type,omitempty"`
	Fog  FogOfWar    `json:"fog,omitempty"`
	Door DoorColor   `json:"door,omitempty"`
}

type TerrainType int

const (
	Nothing TerrainType = iota
	Ground
	Wall
	Corridor
	Door
)

type DoorColor int

const (
	NotADoor DoorColor = iota
	NoColor
	Red
	Green
	Cyan
)

type FogOfWar int

const (
	Unknown FogOfWar = iota
	NotVisible
	Visible
)

type AiState int

const (
	AiRoaming AiState = iota
	AiChasing
)

type StatusFlags uint64

const (
	FlagNone StatusFlags = 0
	// dodges first hit
	FlagDodgy StatusFlags = 1 << 0
	// skips a turn
	FlagSkipTurn StatusFlags = 1 << 1
	// guaranteed first hit against player
	FlagAccurate StatusFlags = 1 << 2
	// Set after hitting the enemy
	FlagJustAttacked StatusFlags = 1 << 3
	// Is enemy invisible
	FlagInvisible StatusFlags = 1 << 4
	// Is enemy mimicking an item
	FlagMimicking         StatusFlags = 1 << 5
	FlagMimickingWeapon   StatusFlags = 1<<6 | FlagMimicking
	FlagMimickingArmor    StatusFlags = 1<<7 | FlagMimicking
	FlagMimickingScroll   StatusFlags = 1<<8 | FlagMimicking
	FlagMimickingElixir   StatusFlags = 1<<9 | FlagMimicking
	FlagMimickingFood     StatusFlags = 1<<10 | FlagMimicking
	FlagMimickingTreasure StatusFlags = 1<<11 | FlagMimicking
)

func (a *StatusFlags) Add(attrs StatusFlags) {
	*a |= attrs
}

func (a *StatusFlags) Remove(attrs StatusFlags) {
	*a &= ^attrs
}

func (a StatusFlags) ContainsAny(attrs StatusFlags) bool {
	return a&attrs != 0
}

func (a StatusFlags) ContainsAll(attrs StatusFlags) bool {
	return a&attrs == attrs
}

type Position struct {
	Y int `json:"y,omitempty"`
	X int `json:"x,omitempty"`
}

type Player struct {
	MaxHealth      int      `json:"maxHealth,omitempty"`
	BuffMaxHealth  int      `json:"buffMaxHealth,omitempty"`
	CurHealth      int      `json:"curHealth,omitempty"`
	Strength       int      `json:"strength,omitempty"` // determines damage
	BuffStrength   int      `json:"buffStrength,omitempty"`
	Dexterity      int      `json:"dexterity,omitempty"` // determines evasion and hit chance
	BuffDexterity  int      `json:"buffDexterity,omitempty"`
	EquipDexterity int      `json:"equipDexterity,omitempty"` // change in agility depending on the worn armor
	Armor          int      `json:"armor,omitempty"`
	EquipArmor     int      `json:"equipArmor,omitempty"` // change in armor value depending on the worn armor
	DmgMin         int      `json:"dmgMin,omitempty"`
	DmgMax         int      `json:"dmgMax,omitempty"`
	EquipDmgMin    int      `json:"equipDmgMin,omitempty"` // change in the minimum damage from the wielded weapon
	EquipDmgMax    int      `json:"equipDmgMax,omitempty"` // change in the maximum damage from the wielded weapon
	BuffTimers     []Buff   `json:"buffTimers,omitempty"`
	HasWeapon      bool     `json:"hasWeapon,omitempty"`
	WieldedWeapon  Weapon   `json:"wieldedWeapon"`
	HasArmor       bool     `json:"hasArmor,omitempty"`
	WornArmor      Armor    `json:"wornArmor"`
	TreasAmt       int      `json:"treasAmt,omitempty"` // for ease of rendering, the amount of treasure from the Backpack is duplicated here
	Pos            Position `json:"pos"`
}

type Buff struct {
	Attribute BuffType `json:"attribute,omitempty"`
	Amount    int      `json:"amount,omitempty"`
	TicksLeft int      `json:"ticksLeft,omitempty"`
}

type BuffType int

const (
	maxHealth BuffType = iota
	strength
	dexterity
)

type Item interface {
	GetType() ItemType
}

type Loot struct {
	Content Item     `json:"content,omitempty"`
	Pos     Position `json:"pos"`
}

type Statistics struct {
	PlayerName     string      `json:"playerName,omitempty"`
	IsDead         bool        `json:"isDead,omitempty"`
	IsWon          bool        `json:"isWon,omitempty"`
	IsAbandoned    bool        `json:"isAbandoned,omitempty"`
	EnemyWhoKilled Enemy       `json:"enemyWhoKilled"`
	TreasAmt       int         `json:"treasAmt,omitempty"`
	LevelReached   int         `json:"levelReached,omitempty"`
	Kills          []KillsStat `json:"kills,omitempty"`
	FoodEaten      int         `json:"foodEaten,omitempty"`
	ScrollsRead    int         `json:"scrollsRead,omitempty"`
	ElixirsDrunk   int         `json:"elixirsDrunk,omitempty"`
	HitsLanded     int         `json:"hitsLanded,omitempty"`
	HitsMissed     int         `json:"hitsMissed,omitempty"`
	HitsTaken      int         `json:"hitsTaken,omitempty"`
	HitsDodged     int         `json:"hitsDodged,omitempty"`
	TilesTraveled  int         `json:"tilesTraveled,omitempty"`
	DoorsOpened    int         `json:"doorsOpened,omitempty"`
}

type KillsStat struct {
	Type  EnemyType `json:"type,omitempty"`
	Count int       `json:"count,omitempty"`
}

type Enemy struct {
	Type           EnemyType   `json:"type,omitempty"`
	Name           string      `json:"name,omitempty"`
	MaxHealth      int         `json:"maxHealth,omitempty"`
	CurHealth      int         `json:"curHealth,omitempty"`
	Strength       int         `json:"strength,omitempty"`
	Dexterity      int         `json:"dexterity,omitempty"`
	Armor          int         `json:"armor,omitempty"`
	Hostility      int         `json:"hostility,omitempty"` // determines the radius from which the enemy begins to pursue the player
	Pos            Position    `json:"pos"`
	AiState        AiState     `json:"aiState"`
	RoamingCounter int         `json:"roamingCounter"`
	Flags          StatusFlags `json:"enemyFlags"`
}

type EnemyType int

type EnemyBehaviors interface {
	Roam(g *GameInfo, enemyIndex int) error
	Chase(g *GameInfo, enemyIndex int) error
	OnPlayerHit(g *GameInfo, enemyIndex int)
	OnEnemyHit(g *GameInfo, enemyIndex int)
	CalcHitEnemyChance(enemy *Enemy, playerDexterity int) int
	CalcHitPlayerChance(g *GameInfo, enemyIndex int) int
	GetTileNeighbors(pos Position) []Position
}

const (
	Zombie EnemyType = iota
	Vampire
	Ghost
	Ogre
	SnakeMage
	Mimic
)

type Action int

const (
	NoAction Action = iota
	Idle
	Confirm
	MoveUp
	MoveDown
	MoveLeft
	MoveRight
	GoDownTheStairs
	Inventory
	EquipWeapon
	EquipArmor
	EatFood
	QuaffElixir
	ReadScroll
	DropItem
	Return
	Quit
	Zero
	One
	Two
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Yes
	No
	Leaderboard
	FirstPerson
	MiniMap
	Help
)

type SaveLoad interface {
	SaveLeaderboard([]Statistics) error
	LoadLeaderboard() ([]Statistics, error)
	SaveGame(GameInfo) error
	LoadGame() (GameInfo, error)
	HasSavedGame() (bool, error)
	DeleteSavedGame() error
}

type Direction int

const (
	Up Direction = iota
	Down
	Left
	Right
)
