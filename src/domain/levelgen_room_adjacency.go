package domain

import (
	"math/rand/v2"
)

// Here we generate an adjacency matrix of rooms, telling us which rooms should be connected with corridors

// How likely we are to connect two rooms with corridors.
//
// The algorithm will keep adding random corridors until all rooms are connected.
// After that there is a 1/corridorFrequency chance to place the next corridor
const corridorFrequency = 8

type roomAdjacency struct {
	matrix [roomCount][roomCount]bool
}

func (a *roomAdjacency) connectRooms(room1 roomIndex, room2 roomIndex) {
	a.matrix[room1][room2] = true
	a.matrix[room2][room1] = true
}

// Helper data structure letting us quickly (more or less) check if two rooms are reachable from each other,
// i.e. they belong to the same "island" of connected rooms
//
// It works by assingning each roomIndex a parent room, the room with no parent is a set representative.
// Two rooms are reachable from each otehr if their set representaive (i.e. their topmost parent) is the same.
//
// https://en.wikipedia.org/wiki/Disjoint-set_data_structure
type roomDisjointSet struct {
	parent [roomCount]roomIndex
}

const noRoom roomIndex = 255

func newRoomDisjointSet() (s roomDisjointSet) {
	for i := range roomCount {
		s.parent[i] = noRoom
	}

	return
}

func (s *roomDisjointSet) getSetRepresentative(room roomIndex) roomIndex {
	repr := room

	for s.parent[repr] != noRoom {
		repr = s.parent[repr]
	}

	return repr
}

func (s *roomDisjointSet) connectRooms(room1 roomIndex, room2 roomIndex) {
	repr1 := s.getSetRepresentative(room1)
	repr2 := s.getSetRepresentative(room2)

	if repr1 != repr2 {
		s.parent[repr1] = repr2
	}
}

// All rooms are reachable from each other when their set representative is the same
func (s *roomDisjointSet) areAllRoomsReachable() bool {
	repr := s.getSetRepresentative(0)

	for i := roomIndex(1); i < roomCount; i++ {
		if s.getSetRepresentative(i) != repr {
			return false
		}
	}

	return true
}

type connectionPair struct {
	a roomIndex
	b roomIndex
}

func allConnectionPairs() []connectionPair {
	pairs := []connectionPair{}

	for i := range roomIndex(roomCount) {

		if i.col()+1 < roomCols {
			right := i + 1
			pairs = append(pairs, connectionPair{i, right})
		}
		if i.row()+1 < roomRows {
			down := i + roomCols
			pairs = append(pairs, connectionPair{i, down})
		}
	}

	return pairs
}

func generateRoomAdjacency() roomAdjacency {
	adjacency := roomAdjacency{}
	disjointSet := newRoomDisjointSet()
	connections := allConnectionPairs()

	rand.Shuffle(len(connections), func(i, j int) {
		connections[i], connections[j] = connections[j], connections[i]
	})

	for len(connections) > 0 {
		if disjointSet.areAllRoomsReachable() {
			if rand.IntN(corridorFrequency) != 0 {
				break
			}
		}

		var connection connectionPair
		connection, connections = connections[0], connections[1:]

		adjacency.connectRooms(connection.a, connection.b)
		disjointSet.connectRooms(connection.a, connection.b)
	}

	return adjacency
}
