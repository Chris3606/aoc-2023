package day14

import (
	"aoc/utils"
	"fmt"
	"os"
	"strings"
)

// For this one we'll just parse the input into our grid structure.
func parseInput(path string) (utils.Grid[byte], map[utils.Point]bool, error) {
	f, err := os.Open(path)
	if err != nil {
		return utils.Grid[byte]{}, nil, err
	}
	defer f.Close()

	grid, err := utils.ReadGridFromBytes(f, func(b byte, p utils.Point) (byte, error) { return b, nil })
	if err != nil {
		return grid, nil, nil
	}

	rocks := map[utils.Point]bool{}
	var posIt = grid.Positions()
	for posIt.Next() {
		cur := posIt.Current()

		if grid.GetCopy(cur) == 'O' {
			rocks[cur] = true
		}
	}

	return grid, rocks, nil
}

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

// Need to represent rocks in something I can use as a map key.
func rocksToString(grid utils.Grid[byte]) string {
	var s strings.Builder

	posIt := grid.Positions()

	for posIt.Next() {
		cur := posIt.Current()
		if grid.GetCopy(cur) == '.' {
			s.WriteString(fmt.Sprintf("%d", cur.ToIndex(grid.Width())) + " ")
		}
	}
	// for y := 0; y < grid.Height(); y++ {
	// 	for x := 0; x < grid.Width(); x++ {
	// 		cur := utils.Point{X: x, Y: y}
	// 		if grid.GetCopy(cur) == '.' {
	// 			s.WriteString(fmt.Sprintf("%d", cur.ToIndex(grid.Width())) + " ")
	// 		}
	// 	}
	// 	s.WriteByte('\n')
	// }

	return s.String()
}

// func calculateNorthernSupportLoad(rocks map[utils.Point]bool, dishHeight int) int {
// 	sum := 0
// 	for k := range rocks {
// 		sum += dishHeight - k.Y
// 	}

// 	return sum
// }

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

func PartA(path string) int {
	dish, _, err := parseInput(path)
	utils.CheckError(err)

	// Tilt all rocks north
	rollNorth(dish)

	// Calculate and return support load
	return calculateNorthernSupportLoad(dish)
}

func PartB(path string) int {
	const iters = 1000000000
	dish, _, err := parseInput(path)
	utils.CheckError(err)

	stateMap := map[string]int{}
	stateMap[rocksToString(dish)] = 0

	// Perform cycles, and look for a duplicate state
	curIter := 0
	period := 0
	for i := 1; i <= iters; i++ {
		rollNorth(dish)
		rollWest(dish)
		rollSouth(dish)
		rollEast(dish)

		state := rocksToString(dish)

		// We identified a cycle.  Figure out the period.
		if v, ok := stateMap[state]; ok {
			curIter = i
			period = i - v
			break

		} else {
			stateMap[state] = i
		}
	}

	// Found a period, so use that period to project close to the number of iterations we're
	// targeting
	if period != 0 {
		// How many iterations are remaining
		remainingIters := iters - curIter

		// Since we know period iters from this point is a cycle, we'll cut out to the nearest
		// multiple of these iterations
		remainingIters %= period

		// Perform remaining cycles
		for i := 0; i < remainingIters; i++ {
			rollNorth(dish)
			rollWest(dish)
			rollSouth(dish)
			rollEast(dish)
		}
	}

	// Perform cycles

	// // Perform cycles
	// for i := 0; i < iters; i++ {
	// 	rollNorth(dish)
	// 	rollWest(dish)
	// 	rollSouth(dish)
	// 	rollEast(dish)
	// }

	// Calculate and return support load
	return calculateNorthernSupportLoad(dish)
}
