package day10

import (
	"aoc/utils"
	"errors"
	"os"
	"slices"
)

// type Node struct {
// 	Position utils.Point
// 	Edges    []utils.Point
// }

type QueueNode struct {
	Position utils.Point
	Distance int
}

var charToNeighbors map[byte][]utils.Point = map[byte][]utils.Point{
	'|': {utils.UP, utils.DOWN},
	'-': {utils.RIGHT, utils.LEFT},
	'L': {utils.UP, utils.RIGHT},
	'J': {utils.UP, utils.LEFT},
	'7': {utils.DOWN, utils.LEFT},
	'F': {utils.DOWN, utils.RIGHT},
}

func parseInput(path string) (utils.Grid[byte], utils.Point, error) {

	f, err := os.Open(path)
	if err != nil {
		return utils.Grid[byte]{}, utils.Point{}, err
	}
	defer f.Close()

	//scanner := bufio.NewScanner(f)
	//scanner.Split(bufio.ScanLines)

	grid, err := utils.ReadGridFromBytes(f, func(b byte, p utils.Point) (byte, error) { return b, nil })
	if err != nil {
		return grid, utils.Point{}, err
	}

	// Find starting node
	startIdx := slices.Index(grid.Slice, 'S')
	if startIdx == -1 {
		return grid, utils.Point{}, errors.New("bad data")
	}
	start := grid.PosFromIndex(startIdx)

	// Find nodes connected to the start.  Bidirectional edges are guaranteed at the start (and
	// _ONLY_ the start) per the instructions
	var startNeighbors []utils.Point
	for _, dir := range utils.CARDINAL_DIRS_CLOCKWISE {
		neighbor := start.Add(dir)
		if !grid.Contains(neighbor) {
			continue
		}

		neighborPerspectiveDir := utils.GetOppositeDir(dir)
		slice := charToNeighbors[grid.GetCopy(neighbor)]
		if slices.Contains(slice, neighborPerspectiveDir) {
			startNeighbors = append(startNeighbors, dir)
		}
	}
	if len(startNeighbors) != 2 {
		return grid, start, errors.New("incorrect starting neighbors")
	}

	var startChar byte
	for k, v := range charToNeighbors {
		if v[0] == startNeighbors[0] && v[1] == startNeighbors[1] || v[1] == startNeighbors[0] && v[0] == startNeighbors[1] {
			startChar = k
			break
		}
	}

	if startChar == 0 {
		return grid, start, errors.New("couldn't find pipe for start")
	}

	grid.Set(start, startChar)

	return grid, start, nil

	// graph := map[utils.Point]Node{}
	// y := 0
	// start := utils.Point{}

	// for scanner.Scan() {
	// 	text := scanner.Bytes()

	// 	for x, b := range text {
	// 		pos := utils.Point{X: x, Y: y}
	// 		node := Node{Position: pos}
	// 		for _, n := range charToNeighbors[b] {
	// 			node.Edges = append(node.Edges, pos.Add(n))
	// 		}

	// 		if b == 'S' {
	// 			start = pos
	// 		}

	// 		graph[pos] = node
	// 	}
	// 	y++
	// }

	// // Find connections for start
	// for _, d := range utils.CARDINAL_DIRS_CLOCKWISE {
	// 	pos += i
	// }

}

func PartA(path string) int {
	grid, start, err := parseInput(path)
	utils.CheckError(err)

	visited := map[utils.Point]bool{}
	distances := map[utils.Point]int{}

	queue := []QueueNode{{Position: start, Distance: 0}}

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]

		if visited[cur.Position] {
			continue
		}

		visited[cur.Position] = true

		distances[cur.Position] = cur.Distance

		for _, dir := range charToNeighbors[grid.GetCopy(cur.Position)] {
			neighbor := cur.Position.Add(dir)
			if grid.Contains(neighbor) && !visited[neighbor] {
				queue = append(queue, QueueNode{Position: neighbor, Distance: cur.Distance + 1})
			}
		}
	}

	max := 0
	for _, v := range distances {
		if v > max {
			max = v
		}
	}

	return max
}

func PartB(path string) string {
	_, _, err := parseInput(path)
	utils.CheckError(err)

	return "Not implemented"
}
