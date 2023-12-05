package utils

import (
	"golang.org/x/exp/constraints"
)

type Summable interface {
	constraints.Ordered | constraints.Complex | string
}

func CheckError(e error) {
	if e != nil {
		panic(e)
	}
}

func Sum[T Summable](values []T) T {
	var sum T
	for _, v := range values {
		sum += v
	}

	return sum
}

// Represents a range of numbers (both ends are inclusive)
type Range struct {
	Start int
	End   int
}

// Creates a new range based on the given start and length
func NewRange(start int, length int) Range {
	return Range{
		Start: start,
		End:   start + length - 1,
	}
}

// Returns whether or not r1 contains the given number.
func (r1 Range) ContainsNum(num int) bool {
	return num >= r1.Start && num <= r1.End
}

type Point struct {
	X int
	Y int
}

// Directions
var UP = Point{0, -1}
var UP_RIGHT = Point{1, -1}
var RIGHT = Point{1, 0}
var DOWN_RIGHT = Point{1, 1}
var DOWN = Point{0, 1}
var DOWN_LEFT = Point{-1, 1}
var LEFT = Point{-1, 0}
var UP_LEFT = Point{-1, -1}

var CARDINAL_DIRS_CLOCKWISE = []Point{UP, RIGHT, DOWN, LEFT}
var DIRS_CLOCKWISE = []Point{UP, UP_RIGHT, RIGHT, DOWN_RIGHT, DOWN, DOWN_LEFT, LEFT, UP_LEFT}

func FromIndex(index, width int) Point {
	return Point{index % width, index / width}
}

func (point Point) ToIndex(width int) int {
	return point.Y*width + point.X
}

func (p1 Point) Add(p2 Point) Point {
	return Point{p1.X + p2.X, p1.Y + p2.Y}
}

func (p1 Point) Sub(p2 Point) Point {
	return Point{p1.X - p2.X, p1.Y - p2.Y}
}

type Rectangle struct {
	MinExtent Point
	MaxExtent Point
}
