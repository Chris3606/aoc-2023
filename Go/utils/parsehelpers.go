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

// ScanDoubleNewLines is a split function for a Scanner that returns a series of strings
// which were separated by a new-line (either \r\n\r\n or \n\n) in the input.
// The last non-empty lines of input will be returned even if it has no
// newline termination.
func ScanDoubleNewLines(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	curData := data[0:]
	idx := 0
	for len(curData) > 0 {
		// Find first new-line
		if i := bytes.IndexByte(curData, '\n'); i >= 0 {
			// Nothing comes after it so this isn't a terminator
			if i+1 >= len(curData) {
				idx += i + 1
				curData = curData[i+1:]
				continue
			}

			// Found a group
			if curData[i+1] == '\n' {
				return idx + i + 2, dropCR(data[0 : idx+i]), nil
			} else if i+2 >= len(curData) || i == 0 { // Need \r\n as the sep but not enough data
				break
			} else if curData[i-1] == '\r' && curData[i+1] == '\r' && curData[i+2] == '\n' { // \r\n\r\n
				return idx + i + 3, dropCR(data[0 : idx+i]), nil
			} else { // Not a separator, only a single new-line
				idx += i + 1
				curData = curData[i+1:]
				continue
			}
		} else {
			break
		}
	}
	// If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		return len(data), dropCR(data), nil
	}
	// Request more data.
	return 0, nil, nil
}

// Gets a string scanner which splits on the given delimiter
func NewStringDelimiterScanner(data string, separator string) *bufio.Scanner {
	s := bufio.NewScanner(strings.NewReader(data))
	s.Split(ScanDelimiterFunc(separator))

	return s
}

// dropCR drops a terminal \r from the data.
func dropCR(data []byte) []byte {
	if len(data) > 0 && data[len(data)-1] == '\r' {
		return data[0 : len(data)-1]
	}
	return data
}

// Reads in groups separated by double-newline.  Each line (NOT each group) is passed through the
// given parser function to produce a value of type T.
//
// If you want to simply get strings representing the group, you can have the parsing function
// simply return the string it is passed.
//
// The primary advantage of using this function over a ScanDelimiterFunc or the like is that it
// automatically handles both unix and windows line-endings; so \n\n or \r\n\r\n work as
// separators.
func ReadGroups[T any](r io.Reader, parser func(string) (T, error)) ([][]T, error) {
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanLines)

	var result [][]T
	var curSlice []T
	for scanner.Scan() {
		if scanner.Text() != "" {
			val, err := parser(scanner.Text())
			if err != nil {
				return result, err
			}

			curSlice = append(curSlice, val)
		} else {
			result = append(result, curSlice)
			curSlice = nil
		}
	}

	if curSlice != nil {
		result = append(result, curSlice)
	}

	return result, scanner.Err()
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
