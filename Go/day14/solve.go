package day14

import (
	"aoc/utils"
	"fmt"
	"os"
	"strings"
)

// Simply parse the input into our grid structure
func parseInput(path string) (utils.Grid[byte], error) {
	f, err := os.Open(path)
	if err != nil {
		return utils.Grid[byte]{}, err
	}
	defer f.Close()

	grid, err := utils.ReadGridFromBytes(f, func(b byte, p utils.Point) (byte, error) { return b, nil })
	if err != nil {
		return grid, nil
	}

	return grid, nil
}

// Rolls a rock at the given position in the given position as far in that direction as it can
// roll
func rollRock(dish utils.Grid[byte], pos, dir utils.Point) {
	for {
		neighbor := pos.Add(dir)
		if !dish.Contains(neighbor) {
			break
		}

		if dish.GetCopy(neighbor) != '.' {
			break
		}

		dish.Set(pos, '.')
		pos = neighbor
		dish.Set(pos, 'O')
	}
}

// Rolls all rocks north until they come to a stop
func rollNorth(dish utils.Grid[byte]) {
	var posIt = dish.Positions()
	for posIt.Next() {
		cur := posIt.Current()

		if dish.GetCopy(cur) != 'O' {
			continue
		}

		rollRock(dish, cur, utils.UP)
	}
}

// Rolls all rocks south until they come to a stop
func rollSouth(dish utils.Grid[byte]) {
	for y := dish.Height() - 1; y >= 0; y-- {
		for x := 0; x < dish.Width(); x++ {
			cur := utils.Point{X: x, Y: y}

			if dish.GetCopy(cur) != 'O' {
				continue
			}

			rollRock(dish, cur, utils.DOWN)
		}
	}
}

// Rolls all rocks west until they come to a stop
func rollWest(dish utils.Grid[byte]) {
	for x := 0; x < dish.Width(); x++ {
		for y := 0; y < dish.Height(); y++ {
			cur := utils.Point{X: x, Y: y}

			if dish.GetCopy(cur) != 'O' {
				continue
			}

			rollRock(dish, cur, utils.LEFT)
		}
	}
}

// Rolls all rocks east until they come to a stop
func rollEast(dish utils.Grid[byte]) {
	for x := dish.Width() - 1; x >= 0; x-- {
		for y := 0; y < dish.Height(); y++ {
			cur := utils.Point{X: x, Y: y}

			if dish.GetCopy(cur) != 'O' {
				continue
			}

			rollRock(dish, cur, utils.RIGHT)
		}
	}
}

// Creates a string representation of the positions of movable rocks in the grid.
// This encodes the whole of the state of a grid which can change.  We need to encode that state as
// something that can be used as a map key, so we use a string.
func rocksToString(grid utils.Grid[byte]) string {
	var s strings.Builder

	posIt := grid.Positions()

	for posIt.Next() {
		cur := posIt.Current()
		if grid.GetCopy(cur) == 'O' {
			s.WriteString(fmt.Sprintf("%d", cur.ToIndex(grid.Width())) + " ")
		}
	}

	return s.String()
}

// Calculate the load on the northern support beams, per instructions
func calculateNorthernSupportLoad(dish utils.Grid[byte]) int {
	sum := 0
	posIt := dish.Positions()
	for posIt.Next() {
		cur := posIt.Current()
		if dish.GetCopy(cur) == 'O' {
			sum += dish.Height() - cur.Y
		}
	}

	return sum
}

// Performs a cycle as defined by part 2
func performCycle(dish utils.Grid[byte]) {
	rollNorth(dish)
	rollWest(dish)
	rollSouth(dish)
	rollEast(dish)
}

// Performs "n" cycles (as defined by part 2) on the given grid.  The trick is to cut down on the
// number of iterations by detecting cycles where the grid repeatedly moves between the same states,
// and to use those cycles to skip a bunch of iterations.
func performNCycles(dish utils.Grid[byte], iterations int) {
	stateMap := map[string]int{}
	stateMap[rocksToString(dish)] = 0

	// Perform cycles, and look for a duplicate state
	curIter := 0
	period := 0
	for i := 1; i <= iterations; i++ {
		performCycle(dish)

		state := rocksToString(dish)

		// We identified a cycle.  Figure out the period, eg how long it takes to reach the cycle
		// again
		if v, ok := stateMap[state]; ok {
			curIter = i
			period = i - v
			break

		} else {
			stateMap[state] = i
		}
	}

	// Found a cycle; so we know that, from the current grid state, every "period" iterations
	// we'll end up at this very same grid state.  So we can use this to skip a bunch of iterations;
	// just find the multiple of period closest to (while not exceeding) the number of iterations
	// we have to do, then perform the last few iterations manually.
	//
	// If we didn't find a period then we will have completed the entire set of iterations looking,
	// so we're already done
	if period != 0 {
		// Calculate how many iterations are remaining
		remainingIters := iterations - curIter

		// Since we know period iters from this point is a cycle, figure out the ending point of the
		// cycle that is nearest to the number of iterations we have left to do
		remainingIters %= period

		// Perform remaining cycles
		for i := 0; i < remainingIters; i++ {
			performCycle(dish)
		}
	}
}

func PartA(path string) int {
	dish, err := parseInput(path)
	utils.CheckError(err)

	// Tilt all rocks north
	rollNorth(dish)

	// Calculate and return load on northern supports
	return calculateNorthernSupportLoad(dish)
}

func PartB(path string) int {
	dish, err := parseInput(path)
	utils.CheckError(err)

	// Perform cycles
	performNCycles(dish, 1000000000)

	// Calculate and return support load
	return calculateNorthernSupportLoad(dish)
}
