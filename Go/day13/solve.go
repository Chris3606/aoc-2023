package day13

import (
	"aoc/utils"
	"bufio"
	"os"
	"strings"
)

const (
	NoLine = iota
	VerticalLine
	HorizontalLine
)

func parseInput(path string) ([]utils.Grid[byte], error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Split(utils.ScanDoubleNewLines)

	var grids []utils.Grid[byte]
	for scanner.Scan() {
		grid, err := utils.ReadGridFromBytes(
			strings.NewReader(scanner.Text()),
			func(b byte, p utils.Point) (byte, error) { return b, nil })
		if err != nil {
			return grids, err
		}

		grids = append(grids, grid)
	}

	return grids, nil
}

// Returns the y value for the mirror line, or -1 if no mirror line is found
func findHorizontalMirrorLine(grid utils.Grid[byte], ignoreVal int) int {
	for y := 1; y < grid.Height(); y++ {
		is_mirror := true
		for y1 := y; y1 < grid.Height(); y1++ {
			mirror_y := y - ((y1 - y) + 1)
			if mirror_y < 0 {
				break
			}

			for x := 0; x < grid.Width(); x++ {
				p := utils.Point{X: x, Y: y1}
				pm := utils.Point{X: x, Y: mirror_y}
				if grid.GetCopy(p) != grid.GetCopy(pm) {
					is_mirror = false
					break
				}
			}

			if !is_mirror {
				break
			}
		}

		if is_mirror && y != ignoreVal {
			return y
		}
	}

	return -1
}

// Returns the x value for the mirror line, or -1 if no mirror line is found
func findVerticalMirrorLine(grid utils.Grid[byte], ignoreVal int) int {
	for x := 1; x < grid.Width(); x++ {
		is_mirror := true
		for x1 := x; x1 < grid.Width(); x1++ {
			mirror_x := x - ((x1 - x) + 1)
			if mirror_x < 0 {
				break
			}

			for y := 0; y < grid.Height(); y++ {
				p := utils.Point{X: x1, Y: y}
				pm := utils.Point{X: mirror_x, Y: y}
				if grid.GetCopy(p) != grid.GetCopy(pm) {
					is_mirror = false
					break
				}
			}

			if !is_mirror {
				break
			}
		}

		if is_mirror && x != ignoreVal {
			return x
		}
	}

	return -1
}

// Finds a mirror line that is NOT the one specified via the "ignore" fields.  To find any
// mirror line, specify NoLine and -1 for ignoreLineType and ignoreLineVal.
func findMirrorLine(grid utils.Grid[byte], ignoreLineType int, ignoreLineVal int) (int, int) {
	vertical := 0
	if ignoreLineType == VerticalLine {
		vertical = findVerticalMirrorLine(grid, ignoreLineVal)
	} else {
		vertical = findVerticalMirrorLine(grid, -1)
	}
	if vertical != -1 {
		return VerticalLine, vertical
	}

	if ignoreLineType == HorizontalLine {
		return HorizontalLine, findHorizontalMirrorLine(grid, ignoreLineVal)
	} else {
		return HorizontalLine, findHorizontalMirrorLine(grid, -1)
	}

}

func PartA(path string) int {
	grids, err := parseInput(path)
	utils.CheckError(err)

	sum := 0
	for _, grid := range grids {
		axis, line := findMirrorLine(grid, NoLine, -1)
		if line == -1 {
			panic("No line of reflection found for a grid.")
		}
		if axis == HorizontalLine {
			line *= 100
		}

		sum += line
	}

	return sum
}

func PartB(path string) int {
	grids, err := parseInput(path)
	utils.CheckError(err)

	sum := 0
	for _, grid := range grids {
		// Find original line
		origAxis, origLine := findMirrorLine(grid, NoLine, -1)
		if origLine == -1 {
			panic("No line of reflection found for a grid.")
		}

		posIt := grid.Positions()

		foundLine := false
		for posIt.Next() {
			cur := posIt.Current()

			// Switch character
			chOrig := grid.GetCopy(cur)
			ch := byte('.')
			if chOrig == '.' {
				ch = '#'
			}
			grid.Set(cur, ch)

			// Check reflection difference
			axis, line := findMirrorLine(grid, origAxis, origLine)

			// Set back for next check
			grid.Set(cur, chOrig)

			if line != -1 {
				if axis == HorizontalLine {
					line *= 100
				}

				sum += line

				foundLine = true
				break
			}
		}

		if !foundLine {
			panic("No new line of reflection found for grid:")
		}

	}

	return sum
}
