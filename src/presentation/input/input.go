package input

import (
	"rogue/domain"

	"github.com/gbin/goncurses"
)

func Input(scr *goncurses.Window) (action domain.Action) {
	key := scr.GetChar()
	switch key {
	case '.':
		action = domain.Idle
	case ' ':
		action = domain.Confirm
	case 'w', 'W':
		action = domain.MoveUp
	case 's', 'S':
		action = domain.MoveDown
	case 'a', 'A':
		action = domain.MoveLeft
	case 'd', 'D':
		action = domain.MoveRight
	case 'v', 'V':
		action = domain.GoDownTheStairs
	case 'i', 'I':
		action = domain.Inventory
	case 'h', 'H':
		action = domain.EquipWeapon
	case 'g', 'G':
		action = domain.EquipArmor
	case 'j', 'J':
		action = domain.EatFood
	case 'k', 'K':
		action = domain.QuaffElixir
	case 'e', 'E':
		action = domain.ReadScroll
	case 't', 'T':
		action = domain.DropItem
	case '`':
		action = domain.Return
	case 'q', 'Q':
		action = domain.Quit
	case '0':
		action = domain.Zero
	case '1':
		action = domain.One
	case '2':
		action = domain.Two
	case '3':
		action = domain.Three
	case '4':
		action = domain.Four
	case '5':
		action = domain.Five
	case '6':
		action = domain.Six
	case '7':
		action = domain.Seven
	case '8':
		action = domain.Eight
	case '9':
		action = domain.Nine
	case 'y', 'Y':
		action = domain.Yes
	case 'n', 'N':
		action = domain.No
	case 'l', 'L':
		action = domain.Leaderboard
	case 'o', 'O':
		action = domain.FirstPerson
	case 'm', 'M':
		action = domain.MiniMap
	case '?':
		action = domain.Help
	}
	return
}

func GetPlayerName(scr *goncurses.Window) string {
	goncurses.Echo(true)
	scr.Timeout(-1)
	name, _ := scr.GetString(13)
	goncurses.Echo(false)
	scr.Timeout(1000)
	return name
}
