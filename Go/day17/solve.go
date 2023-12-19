package day17

import (
	"aoc/utils"
	"math"
	"os"

	"github.com/oleiade/lane/v2"
)

type State struct {
	Position  utils.Point
	Direction utils.Point
	Travel    int
}

type AStarItem struct {
	State State
	F     int
	G     int
}

func (item *AStarItem) Priority() int {
	return item.F
}

var neighbors = map[utils.Point][]utils.Point{
	utils.UP:    {utils.LEFT, utils.UP, utils.RIGHT},
	utils.RIGHT: {utils.UP, utils.RIGHT, utils.DOWN},
	utils.DOWN:  {utils.RIGHT, utils.DOWN, utils.LEFT},
	utils.LEFT:  {utils.DOWN, utils.LEFT, utils.UP},
}

func shortestPath(grid utils.Grid[int], start, end utils.Point, minDist, maxDist int) (int, bool) {
	heap := lane.NewMinPriorityQueue[AStarItem, int]()

	initState := State{Position: start, Direction: utils.RIGHT, Travel: 1}
	initState2 := State{Position: start, Direction: utils.DOWN, Travel: 1}
	dist := map[State]int{initState: 0, initState2: 0}

	itm := AStarItem{State: initState, G: 0, F: utils.ManhattanDistance(start, end)}
	heap.Push(itm, itm.Priority())

	itm = AStarItem{State: initState2, G: 0, F: utils.ManhattanDistance(start, end)}
	heap.Push(itm, itm.Priority())

	for !heap.Empty() {
		// Pop item
		item, _, ok := heap.Pop()
		if !ok {
			panic("Remove from priority queue failed.")
		}

		// We found the shortest path
		if item.State.Position == end {
			return item.G, true
		}

		// If we've already found a better way, we won't visit this node on the current path;
		// this can happen if multiple states with the same value were pushed into the queue
		d, ok := dist[item.State]
		if !ok {
			d = math.MaxInt
		}
		if item.G > d {
			continue
		}

		// Test all neighbors to see if there is a better path to them by going though the current
		// position
		for _, dir := range neighbors[item.State.Direction] {
			neighbor := item.State.Position.Add(dir)
			// Off edge of grid
			if !grid.Contains(neighbor) {
				continue
			}

			// If we are turning, need to be at least min dist
			if item.State.Direction != dir && item.State.Travel < minDist {
				continue
			}

			// Also need to be no more than max dist
			straightLineTravel := item.State.Travel + 1
			if item.State.Direction != dir {
				straightLineTravel = 1
			}

			// Too long in a straight line
			if straightLineTravel > maxDist {
				continue
			}

			node := AStarItem{State: State{Position: neighbor, Direction: dir, Travel: straightLineTravel}, G: item.G + grid.GetCopy(neighbor), F: item.G + grid.GetCopy(neighbor) + utils.ManhattanDistance(neighbor, end)}

			// If cost is lower, add it to the list of nodes to visit and update cost
			d, ok = dist[node.State]
			if !ok {
				d = math.MaxInt
			}
			if node.G < d {
				dist[node.State] = node.G
				heap.Push(node, node.Priority())
			}
		}
	}

	return -1, false
}

// Simply parse the input into our grid structure
func parseInput(path string) (utils.Grid[int], error) {
	f, err := os.Open(path)
	if err != nil {
		return utils.Grid[int]{}, err
	}
	defer f.Close()

	grid, err := utils.ReadGridFromBytes(f, func(b byte, p utils.Point) (int, error) { return int(b - '0'), nil })
	if err != nil {
		return grid, nil
	}

	return grid, nil
}

func PartA(path string) int {
	grid, err := parseInput(path)
	utils.CheckError(err)

	sp, ok := shortestPath(grid, utils.Point{X: 0, Y: 0}, utils.Point{X: grid.Width() - 1, Y: grid.Height() - 1}, 0, 3)
	if !ok {
		panic("No shortest path found")
	}

	return sp
}

func PartB(path string) int {
	grid, err := parseInput(path)
	utils.CheckError(err)

	sp, ok := shortestPath(grid, utils.Point{X: 0, Y: 0}, utils.Point{X: grid.Width() - 1, Y: grid.Height() - 1}, 4, 10)
	if !ok {
		panic("No shortest path found")
	}

	return sp
}
