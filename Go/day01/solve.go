package day01

import (
	"aoc/utils"
	"bufio"
	"errors"
	"os"
	"strings"
	"unicode"
)

func parseInput(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)

	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, nil
}

var digitMap = map[string]int{
	"zero":  0,
	"one":   1,
	"two":   2,
	"three": 3,
	"four":  4,
	"five":  5,
	"six":   6,
	"seven": 7,
	"eight": 8,
	"nine":  9,
}

func checkWordDigit(s string) (bool, int) {
	for k, v := range digitMap {
		if strings.Contains(s, k) {
			return true, v
		}
	}

	return false, -1
}

func findFirstDigit(s string, checkWords bool) (int, error) {
	for idx := 0; idx < len(s); idx++ {
		if checkWords {
			found, digit := checkWordDigit(s[:idx])
			if found {
				return digit, nil
			}
		}

		if unicode.IsDigit(rune(s[idx])) {
			return int(s[idx] - '0'), nil
		}
	}

	return -1, errors.New("no first digit found")
}

func findLastDigit(s string, countWords bool) (int, error) {
	for idx := len(s) - 1; idx >= 0; idx-- {
		if countWords {
			found, digit := checkWordDigit(s[idx:])
			if found {
				return digit, nil
			}
		}

		if unicode.IsDigit(rune(s[idx])) {
			return int(s[idx] - '0'), nil
		}
	}

	return -1, errors.New("no last digit found")
}

func findCalibrationValue(line string, countWords bool) (int, error) {
	first, err := findFirstDigit(line, countWords)
	if err != nil {
		return -1, err
	}
	utils.CheckError(err)

	last, err := findLastDigit(line, countWords)
	if err != nil {
		return -1, err
	}

	return 10*first + last, nil
}

func PartA(path string) int {
	lines, err := parseInput(path)
	utils.CheckError(err)

	sum := 0
	for _, line := range lines {
		cv, err := findCalibrationValue(line, false)
		utils.CheckError(err)
		sum += cv
	}

	return sum
}

func PartB(path string) int {
	lines, err := parseInput(path)
	utils.CheckError(err)

	sum := 0
	for _, line := range lines {
		cv, err := findCalibrationValue(line, true)
		utils.CheckError(err)
		sum += cv
	}

	return sum
}
