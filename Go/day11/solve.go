package day11

import (
	"aoc/utils"
	"os"
)

func parseInput(path string) (utils.Grid[byte], error) {
	f, err := os.Open(path)
	if err != nil {
		return utils.Grid[byte]{}, err
	}
	defer f.Close()

	return utils.ReadGridFromBytes(f, func(b byte, p utils.Point) (byte, error) { return b, nil })
}

func expandGalaxy(galaxyMap *utils.Grid[byte], expansionFactor int) []utils.Point {
	var xMap []int
	var previousEmpty int
	for x := 0; x < galaxyMap.Width(); x++ {
		isEmpty := true
		for y := 0; y < galaxyMap.Height(); y++ {
			if galaxyMap.GetCopy(utils.Point{X: x, Y: y}) == '#' {
				isEmpty = false
				break
			}
		}

		xMap = append(xMap, previousEmpty)

		if isEmpty {
			previousEmpty += (expansionFactor - 1)
		}
	}

	var yMap []int
	previousEmpty = 0
	for y := 0; y < galaxyMap.Height(); y++ {
		isEmpty := true
		for x := 0; x < galaxyMap.Width(); x++ {
			if galaxyMap.GetCopy(utils.Point{X: x, Y: y}) == '#' {
				isEmpty = false
				break
			}
		}

		yMap = append(yMap, previousEmpty)

		if isEmpty {
			previousEmpty += (expansionFactor - 1)
		}
	}

	var galaxies []utils.Point
	var posIt = galaxyMap.Positions()
	for posIt.Next() {
		curPos := posIt.Current()
		if galaxyMap.GetCopy(curPos) != '#' {
			continue
		}

		galaxies = append(galaxies, curPos.Add(utils.Point{X: xMap[curPos.X], Y: yMap[curPos.Y]}))
	}

	return galaxies
}

func sumGalacticDistances(galaxies []utils.Point) int {
	sum := 0
	for i := 0; i < len(galaxies); i++ {
		for j := i; j < len(galaxies); j++ {
			sum += utils.ManhattanDistance(galaxies[i], galaxies[j])
		}
	}

	return sum
}

func PartA(path string) int {
	galaxyMap, err := parseInput(path)
	utils.CheckError(err)

	galaxies := expandGalaxy(&galaxyMap, 2)
	return sumGalacticDistances(galaxies)
}

func PartB(path string) int {
	galaxyMap, err := parseInput(path)
	utils.CheckError(err)

	galaxies := expandGalaxy(&galaxyMap, 1000000)
	return sumGalacticDistances(galaxies)
}
