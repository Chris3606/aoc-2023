package day18

import (
	"aoc/utils"
	"bufio"
	"errors"
	"os"
	"strconv"
)

func parseInput(path string) ([]utils.Point, int, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, -1, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)

	vertices := []utils.Point{{X: 0, Y: 0}}
	outerPoints := 0
	for scanner.Scan() {
		parts := utils.NewStringDelimiterScanner(scanner.Text(), " ")

		dir, err := utils.ReadStringFromScanner(parts)
		if err != nil {
			return vertices, outerPoints, err
		}

		val, err := utils.ReadItemFromScanner(parts, strconv.Atoi)
		if err != nil {
			return vertices, outerPoints, err
		}

		var direction utils.Point
		switch dir {
		case "U":
			direction = utils.UP
		case "D":
			direction = utils.DOWN
		case "L":
			direction = utils.LEFT
		case "R":
			direction = utils.RIGHT
		default:
			return vertices, outerPoints, errors.New("invalid direction type")
		}

		last := vertices[len(vertices)-1]
		vertices = append(vertices, utils.Point{X: last.X + (direction.X * val), Y: last.Y + (direction.Y * val)})
		outerPoints += val
	}

	return vertices, outerPoints, nil
}

func parseInput2(path string) ([]utils.Point, int, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, -1, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)

	vertices := []utils.Point{{X: 0, Y: 0}}
	outerPoints := 0
	for scanner.Scan() {
		parts := utils.NewStringDelimiterScanner(scanner.Text(), " ")

		_, err := utils.ReadStringFromScanner(parts)
		if err != nil {
			return vertices, outerPoints, err
		}

		_, err = utils.ReadStringFromScanner(parts)
		if err != nil {
			return vertices, outerPoints, err
		}

		hex, err := utils.ReadStringFromScanner(parts)
		if err != nil {
			return vertices, outerPoints, err
		}
		hex = hex[1 : len(hex)-1] // Strip parens

		v, err := strconv.ParseInt(hex[1:len(hex)-1], 16, 0)
		if err != nil {
			return vertices, outerPoints, err
		}
		val := int(v)

		var direction utils.Point
		switch hex[len(hex)-1] {
		case '0':
			direction = utils.RIGHT
		case '1':
			direction = utils.DOWN
		case '2':
			direction = utils.LEFT
		case '3':
			direction = utils.UP
		default:
			return vertices, outerPoints, errors.New("invalid direction type")
		}

		last := vertices[len(vertices)-1]
		vertices = append(vertices, utils.Point{X: last.X + (direction.X * val), Y: last.Y + (direction.Y * val)})
		outerPoints += val
	}

	return vertices, outerPoints, nil
}

func getArea(vertices []utils.Point, boundaryPoints int) int {
	// Shoelace formula.  This gives us the area of the polygon, but since the integer points are
	// in the center of each square, we're missing 1/2 sq meter around the edge of the outer points.
	area := 0
	for i := 0; i < len(vertices); i++ {
		prev := i - 1
		if prev == -1 {
			prev = len(vertices) - 1
		}
		next := (i + 1) % len(vertices)

		area += vertices[i].X * (vertices[prev].Y - vertices[next].Y)
	}

	area = utils.Abs(area) / 2

	// Pick's theorem area formula (rearranged with some algebra)
	innerPoints := area - boundaryPoints/2 + 1

	return innerPoints + boundaryPoints
}

func PartA(path string) int {
	vertices, boundaryPoints, err := parseInput(path)
	utils.CheckError(err)

	return getArea(vertices, boundaryPoints)
}

func PartB(path string) int {
	vertices, boundaryPoints, err := parseInput2(path)
	utils.CheckError(err)

	return getArea(vertices, boundaryPoints)
}
