package day16

import (
	"aoc/utils"
	"os"
)

type Bounce struct {
	Position          utils.Point
	IncomingDirection utils.Point
}

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

func simulateLaser(grid utils.Grid[byte], startingPos utils.Point, startingDir utils.Point, energizedTiles map[utils.Point]bool, bounces map[Bounce]bool) {
	laserPos := startingPos
	laserDir := startingDir

loop:
	for grid.Contains(laserPos) {
		bounce := Bounce{laserPos, laserDir}
		if bounces[bounce] {
			break
		}
		bounces[bounce] = true
		energizedTiles[laserPos] = true

		switch grid.GetCopy(laserPos) {
		case '.': // Same dir
		case '|':
			if laserDir == utils.LEFT || laserDir == utils.RIGHT {
				simulateLaser(grid, laserPos.Add(utils.UP), utils.UP, energizedTiles, bounces)
				simulateLaser(grid, laserPos.Add(utils.DOWN), utils.DOWN, energizedTiles, bounces)
				break loop
			}

			// Otherwise, equiv to '.'
		case '-':
			if laserDir == utils.UP || laserDir == utils.DOWN {
				simulateLaser(grid, laserPos.Add(utils.LEFT), utils.LEFT, energizedTiles, bounces)
				simulateLaser(grid, laserPos.Add(utils.RIGHT), utils.RIGHT, energizedTiles, bounces)
				break loop
			}

			// Otherwise, equiv to '.'

		case '\\':
			switch laserDir {
			case utils.DOWN:
				laserDir = utils.RIGHT
			case utils.LEFT:
				laserDir = utils.UP
			case utils.UP:
				laserDir = utils.LEFT
			case utils.RIGHT:
				laserDir = utils.DOWN
			}

		case '/':
			switch laserDir {
			case utils.DOWN:
				laserDir = utils.LEFT
			case utils.LEFT:
				laserDir = utils.DOWN
			case utils.UP:
				laserDir = utils.RIGHT
			case utils.RIGHT:
				laserDir = utils.UP
			}
		default:
			panic("Unsupported grid value")
		}
		laserPos = laserPos.Add(laserDir)
	}
}

func PartA(path string) int {
	grid, err := parseInput(path)
	utils.CheckError(err)

	energizedTiles := map[utils.Point]bool{}
	simulateLaser(grid, utils.Point{X: 0, Y: 0}, utils.RIGHT, energizedTiles, map[Bounce]bool{})

	return len(energizedTiles)
}

func PartB(path string) int {
	grid, err := parseInput(path)
	utils.CheckError(err)

	maxEnergized := 0
	for x := 0; x < grid.Width(); x++ {
		point := utils.Point{X: x, Y: 0}
		energizedTiles := map[utils.Point]bool{}
		simulateLaser(grid, point, utils.DOWN, energizedTiles, map[Bounce]bool{})
		energized := len(energizedTiles)
		maxEnergized = max(maxEnergized, energized)

		point = utils.Point{X: x, Y: grid.Height() - 1}
		energizedTiles = map[utils.Point]bool{}
		simulateLaser(grid, point, utils.UP, energizedTiles, map[Bounce]bool{})
		energized = len(energizedTiles)
		maxEnergized = max(maxEnergized, energized)
	}

	for y := 0; y < grid.Height(); y++ {
		point := utils.Point{X: 0, Y: y}
		energizedTiles := map[utils.Point]bool{}
		simulateLaser(grid, point, utils.RIGHT, energizedTiles, map[Bounce]bool{})
		energized := len(energizedTiles)
		maxEnergized = max(maxEnergized, energized)

		point = utils.Point{X: grid.Width() - 1, Y: y}
		energizedTiles = map[utils.Point]bool{}
		simulateLaser(grid, point, utils.LEFT, energizedTiles, map[Bounce]bool{})
		energized = len(energizedTiles)
		maxEnergized = max(maxEnergized, energized)
	}

	return maxEnergized
}
