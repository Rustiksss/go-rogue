package domain

import (
	"container/heap"
	"math/rand"
)

// https://en.wikipedia.org/wiki/A*_search_algorithm

func pathfindHeuristic(a Position, b Position) int {
	dx := a.X - b.X
	dy := a.Y - b.Y
	return abs(dx) + abs(dy)
}

type pathNode struct {
	pos        Position
	previous   *pathNode
	nodeCost   int
	pathLength int
}

func (n *pathNode) isOrigin() bool {
	return n.previous == nil
}

func (n *pathNode) getFirstPathPosition() Position {
	for {
		if n.isOrigin() || n.previous.isOrigin() {
			return n.pos
		}

		n = n.previous
	}
}

type nodeHeap []*pathNode

func (h nodeHeap) Len() int               { return len(h) }
func (h nodeHeap) Less(i int, j int) bool { return h[i].nodeCost < h[j].nodeCost }
func (h nodeHeap) Swap(i int, j int)      { h[i], h[j] = h[j], h[i] }

func (h *nodeHeap) Push(item any) {
	*h = append(*h, item.(*pathNode))
}

func (h *nodeHeap) Pop() any {
	item, newH := (*h)[h.Len()-1], (*h)[:h.Len()-1]
	*h = newH

	return item
}

type pathfinder struct {
	gameinfo          *GameInfo
	enemy             *Enemy
	exploredPositions map[Position]struct{}
	pathQueue         nodeHeap
	origin            Position
	target            Position
}

func newPathFinder(gameinfo *GameInfo, enemy *Enemy, origin Position, target Position) pathfinder {
	return pathfinder{
		gameinfo:          gameinfo,
		enemy:             enemy,
		exploredPositions: make(map[Position]struct{}),
		pathQueue:         nodeHeap{},
		origin:            origin,
		target:            target,
	}
}

func (DefaultEnemyBehaviors) GetTileNeighbors(pos Position) []Position {
	return []Position{
		{pos.Y - 1, pos.X},
		{pos.Y + 1, pos.X},
		{pos.Y, pos.X - 1},
		{pos.Y, pos.X + 1},
	}
}

func (p *pathfinder) iterateNeighbors(pos Position, f func(Position)) {
	neighbors := p.enemy.Behaviors().GetTileNeighbors(pos)

	rand.Shuffle(len(neighbors), func(i, j int) { neighbors[i], neighbors[j] = neighbors[j], neighbors[i] })

	for _, neighbor := range neighbors {
		_, hasEnemy := p.gameinfo.checkEnemy(neighbor)
		if !p.gameinfo.checkTerrain(neighbor) || hasEnemy {
			continue
		}

		f(neighbor)
	}
}

func (p *pathfinder) addPositionToQueue(pos Position, previous *pathNode) {
	pathLength := 0
	if previous != nil {
		pathLength = previous.pathLength + 1
	}

	nodeCost := pathLength + pathfindHeuristic(pos, p.target)

	node := &pathNode{
		pos,
		previous,
		nodeCost,
		pathLength,
	}

	heap.Push(&p.pathQueue, node)
	p.exploredPositions[node.pos] = struct{}{}
}

func (p *pathfinder) pathfind() Position {
	p.addPositionToQueue(p.origin, nil)
	i := 0

	for p.pathQueue.Len() > 0 {
		i++
		node := heap.Pop(&p.pathQueue).(*pathNode)

		if node.pos == p.target {
			return node.getFirstPathPosition()
		}

		p.iterateNeighbors(node.pos, func(pos Position) {
			if _, positionExplored := p.exploredPositions[pos]; positionExplored {
				return
			}

			p.addPositionToQueue(pos, node)
		})
	}

	return p.origin
}

// returns the position of the next step from origin towards target
//
// if there is no path from origin to target, returns the origin position
func (g *GameInfo) pathfind(enemy *Enemy, origin Position, target Position) Position {
	p := newPathFinder(g, enemy, origin, target)
	return p.pathfind()
}
