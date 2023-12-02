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
