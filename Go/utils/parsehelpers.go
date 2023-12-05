package utils

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"strings"
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

// Gets a string scanner which splits on the given delimiter
func NewStringDelimiterScanner(data string, separator string) *bufio.Scanner {
	s := bufio.NewScanner(strings.NewReader(data))
	s.Split(ScanDelimiterFunc(separator))

	return s
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

// Reads in a list of items by parsing the string representation of the scanner.Scan() function
// until Scan returns false.
func ReadItems[T any](scanner *bufio.Scanner, parser func(string) (T, error), ignoreBlank bool) ([]T, error) {
	var results []T
	for scanner.Scan() {
		text := scanner.Text()
		if ignoreBlank && len(text) == 0 {
			continue
		}
		val, err := parser(text)
		if err != nil {
			return nil, err
		}

		results = append(results, val)
	}

	return results, nil
}

// Reads in a list of items by parsing the string representation of the scanner.Scan() function
// until Scan returns false.
func ReadItemsToMap[T comparable](scanner *bufio.Scanner, parser func(string) (T, error), ignoreBlank bool) (map[T]bool, error) {
	results := map[T]bool{}
	for scanner.Scan() {
		text := scanner.Text()
		if ignoreBlank && len(text) == 0 {
			continue
		}
		val, err := parser(text)
		if err != nil {
			return nil, err
		}

		results[val] = true
	}

	return results, nil
}

// Reads a string item from the given scanner
func ReadStringFromScanner(scanner *bufio.Scanner) (string, error) {
	if !scanner.Scan() {
		return "", errors.New("bad data format")
	}
	return scanner.Text(), nil
}

// Reads an item from the scanner by using the given parsing function.
func ReadItemFromScanner[T any](scanner *bufio.Scanner, parser func(string) (T, error)) (T, error) {
	if !scanner.Scan() {
		var noop T
		return noop, errors.New("bad data format")
	}
	return parser(scanner.Text())
}
