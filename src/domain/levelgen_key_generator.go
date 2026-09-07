package domain

import (
	"math/rand/v2"
	"slices"
)

type keyGenerator struct {
	layout             *levelLayout
	exploredRoomsSet   [roomCount]bool
	exploredRoomsSlice []roomIndex
	roomsToExplore     []roomIndex
	lockedCorridors    []*corridor
	spawnedKeys        [maxDoorColor]bool
	neededKeys         [maxDoorColor]bool
}

func newkeyGenerator(layout *levelLayout) keyGenerator {
	spawnedKeys := [maxDoorColor]bool{}
	spawnedKeys[NotADoor] = true
	spawnedKeys[NoColor] = true

	return keyGenerator{
		layout:             layout,
		spawnedKeys:        spawnedKeys,
		exploredRoomsSlice: make([]roomIndex, 0, roomCount),
		roomsToExplore:     make([]roomIndex, 0, roomCount),
		lockedCorridors:    make([]*corridor, 0, maxDoorColor-minDoorColor),
	}
}

func (g *keyGenerator) isCorridorLocked(corridor *corridor) bool {
	return !g.spawnedKeys[corridor.door1] || !g.spawnedKeys[corridor.door2]
}

func (g *keyGenerator) addLockedCorridor(corridor *corridor) {
	g.lockedCorridors = append(g.lockedCorridors, corridor)

	if corridor.door1 != NoColor {
		g.neededKeys[corridor.door1] = true
	}

	if corridor.door2 != NoColor {
		g.neededKeys[corridor.door2] = true
	}
}

func (g *keyGenerator) exploreRoom(roomId roomIndex) {
	room := &g.layout.rooms[roomId]

	for _, corridorId := range room.corridors {
		corridor := &g.layout.corridors[corridorId]

		if g.isCorridorLocked(corridor) {
			g.addLockedCorridor(corridor)
			continue
		}

		if corridor.room1 != roomId {
			g.roomsToExplore = append(g.roomsToExplore, corridor.room1)
		} else {
			g.roomsToExplore = append(g.roomsToExplore, corridor.room2)
		}
	}
}

func (g *keyGenerator) unlockCorridors() {
	g.lockedCorridors = slices.DeleteFunc(g.lockedCorridors, func(corridor *corridor) bool {
		if g.isCorridorLocked(corridor) {
			return false
		}

		g.roomsToExplore = append(g.roomsToExplore, corridor.room1)
		g.roomsToExplore = append(g.roomsToExplore, corridor.room2)

		return true
	})
}

func (g *keyGenerator) spawnKey() {
	keyPool := make([]DoorColor, 0, maxDoorColor-minDoorColor)

	for c := minDoorColor; c < maxDoorColor; c++ {
		if g.neededKeys[c] && !g.spawnedKeys[c] {
			keyPool = append(keyPool, c)
		}
	}

	// prevent keys from spawning in entrance room
	viableRooms := make([]roomIndex, 0, len(g.exploredRoomsSlice)-1)
	for _, roomId := range g.exploredRoomsSlice {
		if g.layout.entrance != roomId {
			viableRooms = append(viableRooms, roomId)
		}
	}

	keyColor := keyPool[rand.IntN(len(keyPool))]
	roomId := viableRooms[rand.IntN(len(viableRooms))]
	room := &g.layout.rooms[roomId]
	room.keys = append(room.keys, keyColor)

	g.spawnedKeys[keyColor] = true
}

func (g *keyGenerator) generate() {
	g.roomsToExplore = append(g.roomsToExplore, g.layout.entrance)

	for len(g.exploredRoomsSlice) < roomCount {
		if len(g.roomsToExplore) == 0 {
			g.spawnKey()
			g.unlockCorridors()
			continue
		}

		roomId := slicePop(&g.roomsToExplore)
		if g.exploredRoomsSet[roomId] {
			continue
		}

		g.exploreRoom(roomId)

		g.exploredRoomsSet[roomId] = true
		g.exploredRoomsSlice = append(g.exploredRoomsSlice, roomId)
	}

	for len(g.lockedCorridors) > 0 {
		g.spawnKey()
		g.unlockCorridors()
	}
}

func slicePop[T any](s *[]T) T {
	elem, newS := (*s)[len(*s)-1], (*s)[:len(*s)-1]
	*s = newS

	return elem
}
