package day11

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

// Expands the galaxy by a given expansion factor: an expansion factor of 2x adds 1 blank row/col
// per blank one, an expansion factor of 1,000,000 adds 999,999 blank rows/cols per blank one, etc.
//
// Because the size of the grid gets out of control quickly at larger expansion factors, the result of the
// expansion is simply a slice of galaxy locations, rather than a full Grid.  This is all the data
// we need for today's challenge; though we could return a width/height as well if we needed a full
// grid definition.
func expandGalaxy(galaxyMap *utils.Grid[byte], expansionFactor int) []utils.Point {
	// For each column, calculate how many blank columns precede it
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

		// Increment by our expansion factor - 1, so 2x == add 1, etc.
		if isEmpty {
			previousEmpty += (expansionFactor - 1)
		}
	}

	// For each column, calculate how many blank columns precede it
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

		// Increment by our expansion factor - 1, so 2x == add 1, etc.
		if isEmpty {
			previousEmpty += (expansionFactor - 1)
		}
	}

	// Go through the original grid, and calculate the new positions.  We now know how many columns
	// come before each column, and how many rows come before each row; so we can simply add those
	// values to the existing coordinates to translate original coordinates to new ones.
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

// Finds the sum of the "shortest paths" (actually manhattan distance) between all pairs of galaxies.
//
// Today's challenge asks for the "shortest path" between all pairs of galaxies; but, actual shortest
// path algorithms aren't needed.  The problem says that you are allowed to travel through empty space
// _as well as other galaxies_ on the path between a pair of galaxies; so this implies that nothing
// can "block" your path.  Therefore, the shortest path is always optimal, and is just the distance
// between the two points (manhattan distance, since the problem allows cardinal movements only).
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
