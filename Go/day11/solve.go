package day11

import (
	"aoc/utils"
	"os"
)

func parseInput(path string) ([]utils.Point, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	grid, err := utils.ReadGridFromBytes(f, func(b byte, p utils.Point) (byte, error) { return b, nil })
	if err != nil {
		return nil, err
	}

	var xMap []int
	var previousEmpty int
	for x := 0; x < grid.Width(); x++ {
		isEmpty := true
		for y := 0; y < grid.Height(); y++ {
			if grid.GetCopy(utils.Point{X: x, Y: y}) == '#' {
				isEmpty = false
				break
			}
		}

		xMap = append(xMap, previousEmpty)

		if isEmpty {
			previousEmpty++
		}
	}

	var yMap []int
	previousEmpty = 0
	for y := 0; y < grid.Height(); y++ {
		isEmpty := true
		for x := 0; x < grid.Width(); x++ {
			if grid.GetCopy(utils.Point{X: x, Y: y}) == '#' {
				isEmpty = false
				break
			}
		}

		yMap = append(yMap, previousEmpty)

		if isEmpty {
			previousEmpty++
		}
	}

	var galaxies []utils.Point
	var posIt = grid.Positions()
	for posIt.Next() {
		curPos := posIt.Current()
		if grid.GetCopy(curPos) != '#' {
			continue
		}

		galaxies = append(galaxies, curPos.Add(utils.Point{X: xMap[curPos.X], Y: yMap[curPos.Y]}))
	}

	return galaxies, nil
}

func PartA(path string) int {
	galaxies, err := parseInput(path)
	utils.CheckError(err)

	sum := 0
	for i := 0; i < len(galaxies); i++ {
		for j := i; j < len(galaxies); j++ {
			sum += utils.ManhattanDistance(galaxies[i], galaxies[j])
		}
	}

	return sum
}

func PartB(path string) string {
	_, err := parseInput(path)
	utils.CheckError(err)

	return "Not implemented"
}
