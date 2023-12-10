package day10

import (
	"aoc/utils"
	"errors"
	"os"
	"slices"
)

type QueueNode struct {
	Position utils.Point
	Distance int
}

// Maps each character to the directions it leads to (exactly 2 per character)
var charToNeighbors map[byte][]utils.Point = map[byte][]utils.Point{
	'|': {utils.UP, utils.DOWN},
	'-': {utils.RIGHT, utils.LEFT},
	'L': {utils.UP, utils.RIGHT},
	'J': {utils.UP, utils.LEFT},
	'7': {utils.DOWN, utils.LEFT},
	'F': {utils.DOWN, utils.RIGHT},
}

func parseInput(path string) (utils.Grid[byte], utils.Point, error) {
	// Open file
	f, err := os.Open(path)
	if err != nil {
		return utils.Grid[byte]{}, utils.Point{}, err
	}
	defer f.Close()

	// Read data into grid
	grid, err := utils.ReadGridFromBytes(f, func(b byte, p utils.Point) (byte, error) { return b, nil })
	if err != nil {
		return grid, utils.Point{}, err
	}

	// Find starting node (grid value of 'S')
	startIdx := slices.Index(grid.Slice, 'S')
	if startIdx == -1 {
		return grid, utils.Point{}, errors.New("bad data")
	}
	start := grid.PosFromIndex(startIdx)

	// Find nodes connected to the start.  In the general case, both edges of a pipe need to meet
	// for it to be connected (so we'd need to check bidirectional connectivity).  However, for the
	// start and ONLY the start, we were told we can assume it has exactly 2 things connected to it.
	var startNeighbors []utils.Point
	for _, dir := range utils.CARDINAL_DIRS_CLOCKWISE {
		neighbor := start.Add(dir)
		if !grid.Contains(neighbor) {
			continue
		}

		// From the perspective of the neighbor, it needs to connect the OPPOSITE way from the
		// direction we came from in order to connect.  So, check if the opposite way is in
		// the neighbor's connectivity list.
		neighborPerspectiveDir := utils.GetOppositeDir(dir)
		slice := charToNeighbors[grid.GetCopy(neighbor)]
		if slices.Contains(slice, neighborPerspectiveDir) {
			startNeighbors = append(startNeighbors, dir)
		}
	}

	// We were told there are exactly 2 starting neighbors
	if len(startNeighbors) != 2 {
		return grid, start, errors.New("bad data: incorrect number of starting neighbors")
	}

	// Go through our character mappings, and find the one whose edge list matches the list of
	// neighbors we calculated for our start.  They don't have to be in order; we check both ways
	// around
	var startChar byte
	for k, v := range charToNeighbors {
		if v[0] == startNeighbors[0] && v[1] == startNeighbors[1] || v[1] == startNeighbors[0] && v[0] == startNeighbors[1] {
			startChar = k
			break
		}
	}

	// Bad graph or we messed up
	if startChar == 0 {
		return grid, start, errors.New("couldn't find pipe for start")
	}

	// Replace the starting position with the actual pipe character that belongs here
	grid.Set(start, startChar)

	// Return the final grid, and the start location
	return grid, start, nil
}

// Given a starting position and grid, returns the distances from start of all positions on the main loop
// (including the start) using BFS.  Any value NOT in the resulting map can be assumed to NOT be a part
// of the main loop
func findLoop(grid *utils.Grid[byte], start utils.Point) map[utils.Point]int {
	// Nodes visited and distances found to each node
	distances := map[utils.Point]int{}

	// Queue starts with just the start position
	queue := []QueueNode{{Position: start, Distance: 0}}

	for len(queue) > 0 {
		// Pop element off the front of the queue
		cur := queue[0]
		queue = queue[1:]

		// If we've already found a path, ignore it
		_, ok := distances[cur.Position]
		if ok {
			continue
		}

		// Mark node as visited, and the distance we've found from wherever we came from
		distances[cur.Position] = cur.Distance

		// For each node connected to this one (according to our pipe diagram), add it to the queue
		// if it's within the grid and not visited.  The directions which are part of the pipe
		// connectivity map are ALWAYS a subset of the 4 cardinal directions, so we can always add
		// one to indicate the distance along the path to get to the neighbor (so the "unweighted")
		// invariant necessary to use BFS will still apply
		for _, dir := range charToNeighbors[grid.GetCopy(cur.Position)] {
			neighbor := cur.Position.Add(dir)
			_, ok := distances[neighbor]
			if grid.Contains(neighbor) && !ok {
				queue = append(queue, QueueNode{Position: neighbor, Distance: cur.Distance + 1})
			}
		}
	}

	return distances
}

func PartA(path string) int {
	grid, start, err := parseInput(path)
	utils.CheckError(err)

	mainLoop := findLoop(&grid, start)

	// Find the max distance of any position
	maxVal := 0
	for _, v := range mainLoop {
		if v > maxVal {
			maxVal = v
		}
	}

	return maxVal
}

func PartB(path string) int {
	grid, start, err := parseInput(path)
	utils.CheckError(err)

	mainLoop := findLoop(&grid, start)

	// The trick here is to realize that the loop is always a closed polygon by definition;
	// this opens up a number of methods (shoelace algorithm for area, etc.  Here, I decided to use
	// a typical polygon scanning algorithms from 2D computer graphics.  To figure out if areas are
	// inside or outside of a polygon, we can count the number of times a horizontal line at a given
	// row intersects edges: odd number == inside, even number == outside.
	//
	// There is a tricky corner case, though; cases where you end up intersecting a "vertex"; aka,
	// right on a bend where the 2 ends go different ways.  There are only 2 ways this can occur:
	//     - FJ
	//     - L7
	// This could occur either as above (optionally with a horizontal section of pipes in between
	// like L---7).  In either cases, the horizontal pieces don't count as intersects regardless,
	// so it is irrelevant which case we consider.
	//
	// In either case, moving from one side to the other counts as crossing from outer to inner; so
	// we'll handle it by only considering the 90-degree bends of 'F' and '7' "crossing" the polygon.
	//
	// Choosing 'F' and '7' as the two values to consider also ensures that we handle U-bends;
	// those are:
	//     - LJ (aka a U-shape)
	//     - F7 (upside down u-shape)
	// Again, these can occur with horizontal bars in the middle and the same logic will apply.
	//
	// In these cases, we should _not_ go from outer to inner after traversing these.  If we count
	// only 'F' and '7', then either we count NEITHER, or BOTH bends; in either case it does not
	// change the parity of the number of times we cross (because we either add 0 or 2 times), so
	// this doesn't break either.
	//
	// So, we will count instances of '|', 'F', and '7' _that occur on the main loop_ to be
	// "crossing" the boundary in our count.
	numInner := 0
	for y := 0; y < grid.Height(); y++ {
		isInner := false
		for x := 0; x < grid.Width(); x++ {
			pos := utils.Point{X: x, Y: y}
			val := grid.GetCopy(pos)

			// If value is not on the main loop, it can count as an "inner" point per the directions;
			// so we'll just interpret it as ground since tunnels outside the main loop are irrelevant
			_, ok := mainLoop[pos]
			if !ok {
				val = '.'
			}

			// '.' == "ground", eg NOT a boundary and eligible for inner/outer determination.
			// '|', 'F', or '7' == boundaries
			switch val {
			case '.':
				if isInner {
					numInner++
				}
			case 'F':
				fallthrough
			case '7':
				fallthrough
			case '|':
				isInner = !isInner
			}
		}
	}

	return numInner
}
