package day03

import (
	"aoc/utils"
	"os"
)

type Number struct {
	numDigits int
	value     int
	position  utils.Point
}

func NewNumber() Number {
	return Number{
		numDigits: 0,
		value:     0,
		position:  utils.Point{X: -1, Y: -1},
	}
}

func (num *Number) addDigit(digit byte) {
	num.value *= 10
	num.value += int(digit)
	num.numDigits++
}

// type GridVal = byte

// const (
// 	GridBlank GridVal = iota
// 	GridNumber
// 	GridSymbol
// )

// func parseInput(path string) (map[utils.Point]GridVal, []Number, error) {
// 	f, err := os.Open(path)
// 	if err != nil {
// 		return nil, nil, err
// 	}
// 	defer f.Close()

// 	scanner := bufio.NewScanner(f)
// 	scanner.Split(bufio.ScanLines)

// 	grid := map[utils.Point]GridVal{}
// 	var numbers []Number

// 	width := 0
// 	y := 0
// 	for scanner.Scan() {
// 		text := scanner.Bytes()
// 		var curNum Number

// 		if width == 0 {
// 			width = len(text)
// 		}

// 		for x, ch := range text {
// 			if ch == '.' {
// 				if curNum.value != 0 {
// 					numbers = append(numbers, curNum)
// 					curNum = Number{}
// 				}

// 				continue
// 			}

// 			if ch >= '0' && ch <= '9' {
// 				curNum.addDigit(ch)
// 				grid[utils.Point{X: x, Y: y}] = GridNumber
// 			} else {
// 				grid[utils.Point{X: x, Y: y}] = GridSymbol
// 			}
// 		}
// 		if curNum.value != 0 {
// 			numbers = append(numbers, curNum)
// 		}

// 		y++
// 	}

// 	return grid, numbers, nil
// }

func parseInput(path string) (utils.Grid[byte], []Number, error) {
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
	var numbers []Number
	for y := 0; y < grid.Height(); y++ {
		curNum := NewNumber()
		for x := 0; x < grid.Width(); x++ {
			ch := grid.GetCopy(utils.Point{X: x, Y: y})
			if ch >= '0' && ch <= '9' {
				if curNum.position.X == -1 {
					curNum.position = utils.Point{X: x, Y: y}
				}
				curNum.addDigit(ch - '0')
			} else {
				if curNum.value != 0 {
					numbers = append(numbers, curNum)
					curNum = NewNumber()
				}
			}
		}

		if curNum.value != 0 {
			numbers = append(numbers, curNum)
		}
	}

	return grid, numbers, nil
}

func PartA(path string) int {
	grid, numbers, err := parseInput(path)
	utils.CheckError(err)

	// For each number, check its neighbors for symbols
	sum := 0
	for _, num := range numbers {
		isAdjacent := false

		for dx := 0; dx < num.numDigits; dx++ {
			for _, dir := range utils.DIRS_CLOCKWISE {
				neighbor := num.position.Add(dir)
				neighbor.X += dx

				// Outside grid bounds
				if !grid.Contains(neighbor) {
					continue
				}
				// Neighbor is part of the same number; skip it
				if dir == utils.RIGHT && dx != num.numDigits-1 || dir == utils.LEFT && dx != 0 {
					continue
				}

				val := grid.GetCopy(neighbor)
				if (val < '0' || val > '9') && val != '.' {
					isAdjacent = true
					break
				}
			}
			if isAdjacent {
				break
			}
		}

		if isAdjacent {
			sum += num.value
		}
	}

	return sum
}

func getGearRation(grid *utils.Grid[byte], gearPos utils.Point) int {

}

func PartB(path string) string {
	grid, numbers, err := parseInput(path)
	utils.CheckError(err)

	posIt := grid.Positions()

	for posIt.Next() {
		pos := posIt.Current()

		val := grid.GetCopy(pos)
		if val != '*' {
			break
		}

	}

	return "Not implemented"
}
