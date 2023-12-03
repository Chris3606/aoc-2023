package day03

import (
	"aoc/utils"
	"os"
)

type Number struct {
	numDigits int
	value     int
}

func (num *Number) addDigit(digit byte) {
	num.value *= 10
	num.value += int(digit)
	num.numDigits++
}

func parseInput(path string) (utils.Grid[byte], map[utils.Point]*Number, error) {
	f, err := os.Open(path)
	if err != nil {
		return utils.Grid[byte]{}, nil, err
	}
	defer f.Close()

	// Read input into our grid structure
	grid, err := utils.ReadGridFromBytes(f, func(b byte, _ utils.Point) (byte, error) {
		return b, nil
	})
	if err != nil {
		return grid, nil, err
	}

	// Go through the grid, find our numbers
	numbers := map[utils.Point]*Number{}
	for y := 0; y < grid.Height(); y++ {
		var curNum *Number
		for x := 0; x < grid.Width(); x++ {
			ch := grid.GetCopy(utils.Point{X: x, Y: y})
			if ch >= '0' && ch <= '9' {
				if curNum == nil {
					curNum = new(Number)
				}
				curNum.addDigit(ch - '0')
			} else {
				if curNum != nil {
					for dx := 0; dx < curNum.numDigits; dx++ {
						numbers[utils.Point{X: x - 1 - dx, Y: y}] = curNum
					}
					curNum = nil
				}
			}
		}

		if curNum != nil {
			for dx := 0; dx < curNum.numDigits; dx++ {
				numbers[utils.Point{X: grid.Width() - 1 - dx, Y: y}] = curNum
			}
		}
	}

	return grid, numbers, nil
}

func PartA(path string) int {
	grid, numbers, err := parseInput(path)
	utils.CheckError(err)

	// For each number, check its neighbors for symbols
	sum := 0
	numbersAdded := map[*Number]bool{}
	for pos, num := range numbers {
		isAdjacent := false

		for _, dir := range utils.DIRS_CLOCKWISE {
			neighbor := pos.Add(dir)

			// Outside grid bounds
			if !grid.Contains(neighbor) {
				continue
			}

			// Neighbor is part of the same number; skip it
			if numbers[neighbor] == num {
				continue
			}

			val := grid.GetCopy(neighbor)
			if (val < '0' || val > '9') && val != '.' {
				isAdjacent = true
				break
			}
		}

		if isAdjacent && !numbersAdded[num] {
			sum += num.value
			numbersAdded[num] = true
		}
	}

	return sum
}

func getGearRatio(grid *utils.Grid[byte], numbers map[utils.Point]*Number, gearPos utils.Point) int {
	adjacentNumbers := map[*Number]bool{}
	for _, dir := range utils.DIRS_CLOCKWISE {
		neighbor := gearPos.Add(dir)
		if !grid.Contains(neighbor) {
			continue
		}

		// Find adjacent numbers
		number := numbers[neighbor]
		if number != nil {
			adjacentNumbers[number] = true
		}
	}

	ratio := 0
	if len(adjacentNumbers) == 2 {
		ratio = 1
		for k := range adjacentNumbers {
			ratio *= k.value
		}
	}

	return ratio
}

func PartB(path string) int {
	grid, numbers, err := parseInput(path)
	utils.CheckError(err)

	posIt := grid.Positions()

	ratios := 0
	for posIt.Next() {
		pos := posIt.Current()

		val := grid.GetCopy(pos)
		if val != '*' {
			continue
		}

		ratios += getGearRatio(&grid, numbers, pos)
	}

	return ratios
}
