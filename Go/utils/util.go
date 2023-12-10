package utils

import (
	"math"
	"slices"
)

func CheckError(e error) {
	if e != nil {
		panic(e)
	}
}

// Inserts the given element into the slice
func InsertToSlice[T any](a []T, index int, value T) []T {
	if len(a) == index { // nil or empty slice or after last element
		return append(a, value)
	}
	a = append(a[:index+1], a[index:]...) // index < len(a)
	a[index] = value
	return a
}

// Represents a range of numbers (both ends are inclusive)
type Range struct {
	// The start value of the range (inclusive)
	Start int
	// The end value for the range (inclusive)
	End int
}

// Creates a new range based on the given start and length
func NewRange(start int, length int) Range {
	return Range{
		Start: start,
		End:   start + length - 1,
	}
}

// Gets the length of the range.  Exactly this many elements starting from "start" are included
// within the range.
func (r Range) Length() int {
	return r.End - r.Start + 1
}

// Returns whether or not r1 contains the given number.
func (r1 Range) ContainsNum(num int) bool {
	return num >= r1.Start && num <= r1.End
}

// Returns whether or not r1 and r2 overlap.
func (r1 Range) Overlaps(r2 Range) bool {
	return r1.Start <= r2.End && r2.Start <= r1.End
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

func GetOppositeDir(dir Point) Point {
	idx := slices.Index(DIRS_CLOCKWISE, dir)
	return DIRS_CLOCKWISE[(idx+4)%len(DIRS_CLOCKWISE)]
}

// type Rectangle struct {
// 	MinExtent Point
// 	MaxExtent Point
// }

// Builds a map of how often each element occurs in the given slice
func BuildHistogram[T comparable](slice []T) map[T]int {
	hist := map[T]int{}

	for _, v := range slice {
		hist[v] += 1
	}

	return hist
}

// Gets the pair from the map with the maximum value
func MaxFromMap[K comparable](m map[K]int) (K, int) {
	var item K
	maxVal := math.MinInt
	for k, v := range m {
		if v > maxVal {
			maxVal = v
			item = k
		}
	}

	return item, maxVal
}

// Greatest common divisor (GCD) via Euclidean algorithm
func GCD(a, b int) int {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}

// Find Least Common Multiple (LCM) via GCD
func LCM(a, b int, integers ...int) int {
	result := a * b / GCD(a, b)

	for i := 0; i < len(integers); i++ {
		result = LCM(result, integers[i])
	}

	return result
}
