package day14

import (
	"aoc/utils"
	"os"
)

// For this one we'll just parse the input into our grid structure.
func parseInput(path string) (utils.Grid[byte], error) {
	f, err := os.Open(path)
	if err != nil {
		return utils.Grid[byte]{}, err
	}
	defer f.Close()

	return utils.ReadGridFromBytes(f, func(b byte, p utils.Point) (byte, error) { return b, nil })
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
	dish, err := parseInput(path)
	utils.CheckError(err)

	// Tilt all rocks north
	rollNorth(dish)

	// Calculate and return support load
	return calculateNorthernSupportLoad(dish)
}

func PartB(path string) int {
	dish, err := parseInput(path)
	utils.CheckError(err)

	// Perform cycles
	for i := 0; i < 1000000000; i++ {
		rollNorth(dish)
		rollWest(dish)
		rollSouth(dish)
		rollEast(dish)
	}

	// Calculate and return support load
	return calculateNorthernSupportLoad(dish)
}
