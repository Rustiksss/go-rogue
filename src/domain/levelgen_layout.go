package domain

import (
	"math/rand/v2"
)

const (
	roomRows  = 3
	roomCols  = 3
	roomCount = roomRows * roomCols

	roomSpaceHeight = (FieldHeight - (roomRows - 1)) / roomRows
	roomSpaceWidth  = (FieldWidth - (roomCols - 1)) / roomCols

	// max height and width without walls
	maxRoomHeight = roomSpaceHeight - 2
	maxRoomWidth  = roomSpaceWidth - 2
	minRoomHeight = 2
	minRoomWidth  = 4

	corridorStraightness = 8

	minDoorColor = Red
	maxDoorColor = Cyan + 1
)

type roomIndex uint8

func (i roomIndex) row() uint8 {
	return uint8(i) / roomCols
}

func (i roomIndex) col() uint8 {
	return uint8(i) % roomCols
}

// room rect without walls
type roomRect struct {
	top    int
	left   int
	height int
	width  int
}

func generateRoomRect(room roomIndex) roomRect {
	var r roomRect
	r.height = rand.IntN(maxRoomHeight-minRoomHeight+1) + minRoomHeight
	r.width = rand.IntN(maxRoomWidth-minRoomWidth+1) + minRoomWidth

	maxTop := maxRoomHeight - r.height
	maxLeft := maxRoomWidth - r.width

	r.top = rand.IntN(maxTop + 1)
	r.left = rand.IntN(maxLeft + 1)

	// add margin for walls
	r.top += 1
	r.left += 1

	r.top += (roomSpaceHeight+1)*int(room.row()) + 1
	r.left += (roomSpaceWidth + 1) * int(room.col())

	return r
}

func iterateRectPostions(rect roomRect, f func(y int, x int)) {
	for i := range rect.height {
		for j := range rect.width {
			y := rect.top + i
			x := rect.left + j

			f(y, x)
		}
	}
}

func generateRandomPositions(rect roomRect) []Position {
	positions := make([]Position, 0, rect.width*rect.height)

	iterateRectPostions(rect, func(y, x int) {
		positions = append(positions, Position{Y: y, X: x})
	})

	rand.Shuffle(len(positions), func(i, j int) {
		positions[i], positions[j] = positions[j], positions[i]
	})

	return positions
}

type corridorIndex uint8

type room struct {
	corridors []corridorIndex
	rect      roomRect
	keys      []DoorColor
}

type corridor struct {
	room1          roomIndex
	room2          roomIndex
	door1          DoorColor
	door2          DoorColor
	crossAxisStart int
	mainAxisStart  int
	mainAxisEnd    int
	bends          []corridorBend
}

type corridorBend struct {
	mainAxisPos     int
	crossAxisTarget int
}

func newCorridor(room1 roomIndex, room2 roomIndex) corridor {
	// room2 should either be to the right or to the bottom of room1
	if room1 > room2 {
		room1, room2 = room2, room1
	}
	return corridor{room1: room1, room2: room2, door1: NoColor, door2: NoColor}
}

func (c *corridor) isHorizontal() bool {
	return c.room1.row() == c.room2.row()
}

func (c *corridor) isVertical() bool {
	return !c.isHorizontal()
}

type levelLayout struct {
	rooms     [roomCount]room
	corridors []corridor
	entrance  roomIndex
	exit      roomIndex
}

func (l *levelLayout) connectRoomsFromAdjacency(adjacency *roomAdjacency) {
	for i := range roomIndex(roomCount) {
		for j := i + 1; j < roomCount; j++ {
			if !adjacency.matrix[i][j] {
				continue
			}

			corrIndex := corridorIndex(len(l.corridors))
			l.corridors = append(l.corridors, newCorridor(i, j))

			l.rooms[i].corridors = append(l.rooms[i].corridors, corrIndex)
			l.rooms[j].corridors = append(l.rooms[j].corridors, corrIndex)
		}
	}
}

func (c *corridor) generateBends(minCrossAxis int, maxCrossAxis int, minCrossAxisEnd int, maxCrossAxisEnd int) {
	for i := c.mainAxisStart + 1; i < c.mainAxisEnd-1; i++ {
		// 1/straightness chance to generate a bend at any given point
		if rand.IntN(corridorStraightness) == 0 {
			// crossAxisTarget will be filled in later
			c.bends = append(c.bends, corridorBend{mainAxisPos: i})
			// skip next position, so we don't have 2 bends in a row
			i++
		}
	}

	// force at least one bend
	if len(c.bends) == 0 {
		minPos := c.mainAxisStart + 1
		maxPos := c.mainAxisEnd - 1
		mainAxisPos := rand.IntN(maxPos-minPos) + minPos
		c.bends = append(c.bends, corridorBend{mainAxisPos: mainAxisPos})
	}

	for i := range c.bends {
		if i+1 >= len(c.bends) {
			c.bends[i].crossAxisTarget = rand.IntN(maxCrossAxisEnd-minCrossAxisEnd) + minCrossAxisEnd
		} else {
			c.bends[i].crossAxisTarget = rand.IntN(maxCrossAxis-minCrossAxis) + minCrossAxis
		}
	}
}

func (l *levelLayout) generateCorridorShape(corridor *corridor) {
	var minCrossAxis int
	var maxCrossAxis int

	var minCrossAxisEnd int
	var maxCrossAxisEnd int

	rect1 := l.rooms[corridor.room1].rect
	rect2 := l.rooms[corridor.room2].rect

	if corridor.isHorizontal() {
		// main axis = x
		// cross axis = y
		corridor.crossAxisStart = rand.IntN(rect1.height) + rect1.top
		minCrossAxisEnd = rect2.top
		maxCrossAxisEnd = rect2.top + rect2.height
		corridor.mainAxisStart = rect1.left + rect1.width
		corridor.mainAxisEnd = rect2.left
		minCrossAxis = 1 + (roomSpaceHeight+1)*int(corridor.room1.row())
		maxCrossAxis = minCrossAxis + roomSpaceHeight
	} else {
		// main axis = y
		// cross axis = x
		corridor.crossAxisStart = rand.IntN(rect1.width) + rect1.left
		minCrossAxisEnd = rect2.left
		maxCrossAxisEnd = rect2.left + rect2.width
		corridor.mainAxisStart = rect1.top + rect1.height
		corridor.mainAxisEnd = rect2.top
		minCrossAxis = (roomSpaceWidth + 1) * int(corridor.room1.col())
		maxCrossAxis = minCrossAxis + roomSpaceWidth
	}

	corridor.generateBends(minCrossAxis, maxCrossAxis, minCrossAxisEnd, maxCrossAxisEnd)
}

func generateDoorColor() DoorColor {
	switch rand.IntN(15) {
	case int(Red):
		return Red
	case int(Green):
		return Green
	case int(Cyan):
		return Cyan
	default:
		return NoColor
	}
}

func (l *levelLayout) generateColoredDoors() {
	builtDoors := make(map[DoorColor]bool)

	corridorPtrs := make([]*corridor, len(l.corridors))
	for i := range l.corridors {
		corridorPtrs[i] = &l.corridors[i]
	}

	rand.Shuffle(len(corridorPtrs), func(i, j int) {
		corridorPtrs[i], corridorPtrs[j] = corridorPtrs[j], corridorPtrs[i]
	})

	for _, corridor := range corridorPtrs {
		if len(builtDoors) >= int(maxDoorColor-minDoorColor) {
			break
		}

		// Don't add locked doors in entrance room
		if corridor.room1 == l.entrance || corridor.room2 == l.entrance {
			continue
		}

		color1 := generateDoorColor()
		if color1 != NoColor && !builtDoors[color1] {
			builtDoors[color1] = true
			corridor.door1 = color1
		}

		color2 := generateDoorColor()
		if color2 != NoColor && !builtDoors[color2] {
			builtDoors[color2] = true
			corridor.door2 = color2
		}
	}
}

func (l *levelLayout) generateKeys() {
	generator := newkeyGenerator(l)
	generator.generate()
}

func generateLayout() levelLayout {
	layout := levelLayout{}
	layout.entrance = roomIndex(rand.UintN(roomCount))
	layout.exit = roomIndex(rand.UintN(roomCount - 1))
	if layout.exit >= layout.entrance {
		layout.exit++
	}

	adjacency := generateRoomAdjacency()
	layout.connectRoomsFromAdjacency(&adjacency)

	layout.generateColoredDoors()
	layout.generateKeys()

	for i := range layout.rooms {
		layout.rooms[i].rect = generateRoomRect(roomIndex(i))
	}

	for i := range layout.corridors {
		layout.generateCorridorShape(&layout.corridors[i])
	}

	return layout
}
