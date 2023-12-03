package utils

import (
	"bufio"
	"bytes"
	"io"
)

// A split function which can be passed to scanner.Split to split on an arbitrary separator
func ScanDelimiterFunc(separator string) func(data []byte, atEOF bool) (advance int, token []byte, err error) {
	searchBytes := []byte(separator)
	searchLen := len(searchBytes)
	return func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		dataLen := len(data)

		// Return nothing if at end of file and no data passed
		if atEOF && dataLen == 0 {
			return 0, nil, nil
		}

		// Find next separator and return token
		if i := bytes.Index(data, searchBytes); i >= 0 {
			return i + searchLen, data[0:i], nil
		}

		// If we're at EOF, we have a final, non-terminated line. Return it.
		if atEOF {
			return dataLen, data, nil
		}

		// Request more data.
		return 0, nil, nil
	}
}

// Reads grid in the following format
// 012345
// 654568
// 923598
//
// The values can be any arbitrary byte, whether those are characters, actual bytes, etc.
// The parsing function must translate the values to the appropriate result type.
func ReadGridFromBytes[T any](r io.Reader, parser func(byte, Point) (T, error)) (Grid[T], error) {
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanLines)

	var slice []T
	var width int

	var y int
	for scanner.Scan() {
		text := scanner.Text()
		if width == 0 {
			width = len(text)
		}

		for x := range text {
			val, err := parser(text[x], Point{x, y})
			if err != nil {
				return Grid[T]{}, err
			}
			slice = append(slice, val)
		}

		y++
	}

	return GridFromSlice[T](slice, width), nil
}
